package usecase

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/beabys/wms/login-service/internal/application/auth/command"
	"github.com/beabys/wms/login-service/internal/application/auth/repository"
	"github.com/beabys/wms/login-service/internal/domain/auth/model"
)

// NewAuthService creates a new AuthService.
func NewAuthService(repo repository.UserRepository, privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey, accessTTL, refreshTTL time.Duration) *AuthService {
	return &AuthService{
		repo:       repo,
		privateKey: privateKey,
		publicKey:  publicKey,
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
	}
}

// generateJTI creates a unique token ID.
func generateJTI() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

// newUUID generates a simple UUID v4 string.
func newUUID() string {
	b := make([]byte, 16)
	for i := range b {
		b[i] = byte(i*7 + 3)
	}
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

// Login authenticates a user and returns tokens.
func (s *AuthService) Login(ctx context.Context, cmd command.LoginCommand) (*command.LoginResponse, error) {
	email, err := model.NewEmail(cmd.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid email: %w", err)
	}

	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	if !user.CanLogin() {
		return nil, fmt.Errorf("account is not active")
	}

	if !user.PasswordHash.Verify(cmd.Password) {
		return nil, fmt.Errorf("invalid credentials")
	}

	now := time.Now()
	accessClaims := jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email.String(),
		"role":  string(user.Role),
		"iat":   now.Unix(),
		"exp":   now.Add(s.accessTTL).Unix(),
		"jti":   generateJTI(),
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodRS256, accessClaims)
	accessTokenStr, err := accessToken.SignedString(s.privateKey)
	if err != nil {
		return nil, fmt.Errorf("sign access token: %w", err)
	}

	refreshClaims := jwt.MapClaims{
		"sub":  user.ID,
		"iat":  now.Unix(),
		"exp":  now.Add(s.refreshTTL).Unix(),
		"type": "refresh",
		"jti":  generateJTI(),
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodRS256, refreshClaims)
	refreshTokenStr, err := refreshToken.SignedString(s.privateKey)
	if err != nil {
		return nil, fmt.Errorf("sign refresh token: %w", err)
	}

	if err := s.repo.CreateRefreshToken(ctx, user.ID, refreshTokenStr); err != nil {
		return nil, fmt.Errorf("store refresh token: %w", err)
	}

	return &command.LoginResponse{
		AccessToken:  accessTokenStr,
		RefreshToken: refreshTokenStr,
		UserID:       user.ID,
		UserEmail:    user.Email.String(),
		UserRole:     string(user.Role),
		CustomerID:   user.CustomerID,
	}, nil
}

// Register creates a new user account.
func (s *AuthService) Register(ctx context.Context, cmd command.RegisterCommand) (*command.RegisterResponse, error) {
	email, err := model.NewEmail(cmd.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid email: %w", err)
	}

	existing, err := s.repo.GetByEmail(ctx, email)
	if err == nil && existing != nil {
		return nil, fmt.Errorf("email already registered")
	}

	passwordHash, err := model.NewPasswordHash(cmd.Password)
	if err != nil {
		return nil, fmt.Errorf("invalid password: %w", err)
	}

	user := &model.User{
		ID:           newUUID(),
		Email:        email,
		PasswordHash: passwordHash,
		Role:         model.RoleCustomer,
		Status:       model.UserStatusActive,
	}

	if err := s.repo.Save(ctx, user); err != nil {
		return nil, fmt.Errorf("save user: %w", err)
	}

	return &command.RegisterResponse{
		UserID:    user.ID,
		Email:     user.Email.String(),
		Role:      string(user.Role),
		Status:    string(user.Status),
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
	}, nil
}

// RefreshToken issues new tokens from a valid refresh token.
func (s *AuthService) RefreshToken(ctx context.Context, cmd command.RefreshTokenCommand) (*command.RefreshTokenResponse, error) {
	// Validate the refresh token exists in our store
	userID, err := s.repo.ValidateRefreshToken(ctx, cmd.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token")
	}

	// Parse and validate the old refresh token JWT
	token, err := jwt.Parse(cmd.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.publicKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid refresh token claims")
	}

	if claims["type"] != "refresh" {
		return nil, fmt.Errorf("token is not a refresh token")
	}

	// Delete old refresh token
	if err := s.repo.DeleteRefreshToken(ctx, cmd.RefreshToken); err != nil {
		return nil, fmt.Errorf("delete old refresh token: %w", err)
	}

	// Issue new tokens
	now := time.Now()
	accessClaims := jwt.MapClaims{
		"sub":  userID,
		"iat":  now.Unix(),
		"exp":  now.Add(s.accessTTL).Unix(),
		"jti":  generateJTI(),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodRS256, accessClaims)
	accessTokenStr, err := accessToken.SignedString(s.privateKey)
	if err != nil {
		return nil, fmt.Errorf("sign access token: %w", err)
	}

	refreshClaims := jwt.MapClaims{
		"sub":  userID,
		"iat":  now.Unix(),
		"exp":  now.Add(s.refreshTTL).Unix(),
		"type": "refresh",
		"jti":  generateJTI(),
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodRS256, refreshClaims)
	refreshTokenStr, err := refreshToken.SignedString(s.privateKey)
	if err != nil {
		return nil, fmt.Errorf("sign refresh token: %w", err)
	}

	if err := s.repo.CreateRefreshToken(ctx, userID, refreshTokenStr); err != nil {
		return nil, fmt.Errorf("store refresh token: %w", err)
	}

	return &command.RefreshTokenResponse{
		AccessToken:  accessTokenStr,
		RefreshToken: refreshTokenStr,
	}, nil
}

// ValidateToken validates a JWT token and returns user info.
func (s *AuthService) ValidateToken(ctx context.Context, tokenStr string) (*ValidateResult, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.publicKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	userID, _ := claims["sub"].(string)
	roleStr, _ := claims["role"].(string)

	role := model.Role(roleStr)
	permissions := role.Permissions()

	return &ValidateResult{
		UserID:      userID,
		Role:        roleStr,
		Permissions: permissions,
	}, nil
}

// CreateUser creates a new user by an admin.
func (s *AuthService) CreateUser(ctx context.Context, cmd command.CreateUserCommand) (*command.CreateUserResponse, error) {
	email, err := model.NewEmail(cmd.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid email: %w", err)
	}

	role, err := model.ParseRole(cmd.Role)
	if err != nil {
		return nil, fmt.Errorf("invalid role: %w", err)
	}

	passwordHash, err := model.NewPasswordHash(cmd.Password)
	if err != nil {
		return nil, fmt.Errorf("invalid password: %w", err)
	}

	user := &model.User{
		ID:           newUUID(),
		Email:        email,
		PasswordHash: passwordHash,
		Role:         role,
		CustomerID:   cmd.CustomerID,
		Status:       model.UserStatusActive,
	}

	if err := s.repo.Save(ctx, user); err != nil {
		return nil, fmt.Errorf("save user: %w", err)
	}

	return &command.CreateUserResponse{
		UserID:    user.ID,
		Email:     user.Email.String(),
		Role:      string(user.Role),
		Status:    string(user.Status),
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
	}, nil
}

// GetUser fetches a user by ID.
func (s *AuthService) GetUser(ctx context.Context, query command.GetUserQuery) (*command.GetUserResponse, error) {
	user, err := s.repo.GetByID(ctx, query.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	return &command.GetUserResponse{
		UserID:     user.ID,
		Email:      user.Email.String(),
		Role:       string(user.Role),
		Status:     string(user.Status),
		CustomerID: user.CustomerID,
		CreatedAt:  user.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  user.UpdatedAt.Format(time.RFC3339),
	}, nil
}

// ListUsers lists users with pagination.
func (s *AuthService) ListUsers(ctx context.Context, query command.ListUsersQuery) (*command.ListUsersResponse, error) {
	if query.PageSize <= 0 {
		query.PageSize = 20
	}
	if query.Page <= 0 {
		query.Page = 1
	}
	offset := (query.Page - 1) * query.PageSize

	users, err := s.repo.List(ctx, query.PageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}

	items := make([]*command.UserItem, 0, len(users))
	for _, u := range users {
		items = append(items, &command.UserItem{
			UserID:     u.ID,
			Email:      u.Email.String(),
			Role:       string(u.Role),
			Status:     string(u.Status),
			CustomerID: u.CustomerID,
			CreatedAt:  u.CreatedAt.Format(time.RFC3339),
		})
	}

	return &command.ListUsersResponse{Users: items}, nil
}

// UpdateUser updates a user's fields.
func (s *AuthService) UpdateUser(ctx context.Context, cmd command.UpdateUserCommand) (*command.UpdateUserResponse, error) {
	user, err := s.repo.GetByID(ctx, cmd.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	if cmd.Email != "" {
		email, err := model.NewEmail(cmd.Email)
		if err != nil {
			return nil, fmt.Errorf("invalid email: %w", err)
		}
		user.Email = email
	}

	if cmd.Role != "" {
		role, err := model.ParseRole(cmd.Role)
		if err != nil {
			return nil, fmt.Errorf("invalid role: %w", err)
		}
		user.Role = role
	}

	if cmd.CustomerID != "" {
		user.CustomerID = cmd.CustomerID
	}

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}

	return &command.UpdateUserResponse{
		UserID:     user.ID,
		Email:      user.Email.String(),
		Role:       string(user.Role),
		Status:     string(user.Status),
		CustomerID: user.CustomerID,
		UpdatedAt:  user.UpdatedAt.Format(time.RFC3339),
	}, nil
}
