package usecase

import (
	"context"
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/beabys/wms/auth-service/internal/application/auth/command"
	"github.com/beabys/wms/auth-service/internal/application/auth/repository"
	"github.com/beabys/wms/auth-service/internal/application/auth/validator"
	"github.com/beabys/wms/pkg/logger"
)

// LoginUseCase handles authentication operations.
type LoginUseCase struct {
	logger   logger.Logger
	userRepo repository.UserRepository
	authRepo repository.AuthRepository
	jwtSvc   *Service
}

// NewLoginUseCase creates a new LoginUseCase.
func NewLoginUseCase(
	log logger.Logger,
	userRepo repository.UserRepository,
	authRepo repository.AuthRepository,
	jwtSvc *Service,
) *LoginUseCase {
	return &LoginUseCase{
		logger:   log,
		userRepo: userRepo,
		authRepo: authRepo,
		jwtSvc:   jwtSvc,
	}
}

// Login validates credentials and returns JWT tokens.
func (uc *LoginUseCase) Login(ctx context.Context, cmd command.LoginCommand) (*command.LoginResult, error) {
	if err := validator.Email(cmd.Email); err != nil {
		return nil, err
	}
	if err := validator.Password(cmd.Password); err != nil {
		return nil, err
	}

	user, err := uc.userRepo.GetByEmail(ctx, cmd.Email)
	if err != nil {
		uc.logger.Error("failed to get user by email", err,
			logger.LogField{Key: "email", Value: cmd.Email},
		)
		return nil, fmt.Errorf("invalid credentials")
	}
	if user == nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	if !user.Password.Verify(cmd.Password) {
		return nil, fmt.Errorf("invalid credentials")
	}

	if !user.Active {
		uc.logger.Info("login attempt for inactive user",
			logger.LogField{Key: "user_id", Value: user.ID},
		)
		return nil, fmt.Errorf("account is deactivated")
	}

	accessToken, err := uc.jwtSvc.GenerateAccessToken(user)
	if err != nil {
		uc.logger.Error("failed to generate access token", err)
		return nil, fmt.Errorf("failed to generate token")
	}

	refreshToken, err := uc.jwtSvc.GenerateRefreshToken()
	if err != nil {
		uc.logger.Error("failed to generate refresh token", err)
		return nil, fmt.Errorf("failed to generate token")
	}

	// Store refresh token hash
	tokenHash := hashToken(refreshToken)
	expiresAt := time.Now().Add(7 * 24 * time.Hour) // 7 days
	if err := uc.authRepo.StoreRefreshToken(ctx, user.ID, tokenHash, expiresAt); err != nil {
		uc.logger.Error("failed to store refresh token", err)
		return nil, fmt.Errorf("failed to store token")
	}

	return &command.LoginResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    uc.jwtSvc.GetAccessTokenTTL(),
	}, nil
}

// RefreshToken validates a refresh token and issues new tokens.
func (uc *LoginUseCase) RefreshToken(ctx context.Context, cmd command.RefreshTokenCommand) (*command.LoginResult, error) {
	if err := validator.Token(cmd.RefreshToken); err != nil {
		return nil, err
	}

	tokenHash := hashToken(cmd.RefreshToken)
	userID, err := uc.authRepo.ValidateRefreshToken(ctx, tokenHash)
	if err != nil {
		uc.logger.Error("failed to validate refresh token", err)
		return nil, fmt.Errorf("invalid refresh token")
	}
	if userID == "" {
		return nil, fmt.Errorf("invalid or expired refresh token")
	}

	// Revoke the used token (rotation)
	if err := uc.authRepo.RevokeRefreshToken(ctx, tokenHash); err != nil {
		uc.logger.Error("failed to revoke used refresh token", err)
		return nil, fmt.Errorf("failed to rotate token")
	}

	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		uc.logger.Error("failed to get user by ID", err)
		return nil, fmt.Errorf("invalid user")
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}
	if !user.Active {
		return nil, fmt.Errorf("account is deactivated")
	}

	accessToken, err := uc.jwtSvc.GenerateAccessToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token")
	}

	newRefreshToken, err := uc.jwtSvc.GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token")
	}

	newHash := hashToken(newRefreshToken)
	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	if err := uc.authRepo.StoreRefreshToken(ctx, user.ID, newHash, expiresAt); err != nil {
		return nil, fmt.Errorf("failed to store refresh token")
	}

	return &command.LoginResult{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    uc.jwtSvc.GetAccessTokenTTL(),
	}, nil
}

// ValidateToken validates a JWT and returns its claims.
func (uc *LoginUseCase) ValidateToken(ctx context.Context, cmd command.ValidateTokenCommand) (*command.ValidateTokenResult, error) {
	if err := validator.Token(cmd.Token); err != nil {
		return nil, err
	}

	claims, err := uc.jwtSvc.ValidateToken(cmd.Token)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	return &command.ValidateTokenResult{
		UserID:      claims.UserID,
		Role:        claims.Role,
		CustomerID:  claims.CustomerID,
		Permissions: claims.Permissions,
	}, nil
}

// Logout revokes a refresh token.
func (uc *LoginUseCase) Logout(ctx context.Context, cmd command.LogoutCommand) error {
	if cmd.RefreshToken == "" {
		return nil // no token to revoke, still a success
	}
	tokenHash := hashToken(cmd.RefreshToken)
	if err := uc.authRepo.RevokeRefreshToken(ctx, tokenHash); err != nil {
		uc.logger.Error("failed to revoke refresh token on logout", err)
		return fmt.Errorf("failed to revoke token")
	}
	return nil
}

// GetPublicKey returns the RSA public key PEM.
func (uc *LoginUseCase) GetPublicKey(ctx context.Context) (*command.GetPublicKeyResult, error) {
	pem, err := uc.jwtSvc.GetPublicKeyPEM()
	if err != nil {
		return nil, fmt.Errorf("failed to get public key: %w", err)
	}
	return &command.GetPublicKeyResult{
		PublicKeyPEM: pem,
	}, nil
}

func hashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return fmt.Sprintf("%x", h)
}
