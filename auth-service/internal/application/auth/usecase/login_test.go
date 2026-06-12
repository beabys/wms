package usecase

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/beabys/wms/auth-service/internal/application/auth/command"
	"github.com/beabys/wms/auth-service/internal/domain/auth/model"
	repomocks "github.com/beabys/wms/auth-service/mocks/application/auth/repository"
	"github.com/beabys/wms/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setupLoginTest(t *testing.T) (*LoginUseCase, *repomocks.UserRepository, *repomocks.AuthRepository) {
	t.Helper()
	log := logger.NewSlogLogger(slog.LevelDebug)

	userRepo := repomocks.NewUserRepository(t)
	authRepo := repomocks.NewAuthRepository(t)

	privPEM, pubPEM, err := GenerateKeyPair()
	require.NoError(t, err)

	jwtSvc, err := NewService(privPEM, pubPEM, 15*time.Minute, 7*24*time.Hour)
	require.NoError(t, err)

	uc := NewLoginUseCase(log, userRepo, authRepo, jwtSvc)
	return uc, userRepo, authRepo
}

func TestLoginUseCase_Login(t *testing.T) {
	t.Run("successful login returns tokens", func(t *testing.T) {
		uc, userRepo, authRepo := setupLoginTest(t)

		user, err := model.NewUser("test@example.com", "password123", "Test", "admin", nil)
		require.NoError(t, err)

		userRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(user, nil)
		authRepo.On("StoreRefreshToken", mock.Anything, user.ID, mock.Anything, mock.Anything).Return(nil)

		result, err := uc.Login(context.Background(), command.LoginCommand{
			Email:    "test@example.com",
			Password: "password123",
		})
		require.NoError(t, err)
		assert.NotEmpty(t, result.AccessToken)
		assert.NotEmpty(t, result.RefreshToken)
		assert.Greater(t, result.ExpiresIn, int64(0))
	})

	t.Run("wrong password returns error", func(t *testing.T) {
		uc, userRepo, _ := setupLoginTest(t)

		user, _ := model.NewUser("test@example.com", "password123", "Test", "admin", nil)

		userRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(user, nil)

		_, err := uc.Login(context.Background(), command.LoginCommand{
			Email:    "test@example.com",
			Password: "wrongpassword",
		})
		assert.Error(t, err)
	})

	t.Run("non-existent email returns error", func(t *testing.T) {
		uc, userRepo, _ := setupLoginTest(t)

		userRepo.On("GetByEmail", mock.Anything, "nonexistent@example.com").Return(nil, nil)

		_, err := uc.Login(context.Background(), command.LoginCommand{
			Email:    "nonexistent@example.com",
			Password: "password123",
		})
		assert.Error(t, err)
	})

	t.Run("inactive user returns error", func(t *testing.T) {
		uc, userRepo, _ := setupLoginTest(t)

		user, _ := model.NewUser("inactive@example.com", "password123", "Inactive", "admin", nil)
		user.Deactivate()

		userRepo.On("GetByEmail", mock.Anything, "inactive@example.com").Return(user, nil)

		_, err := uc.Login(context.Background(), command.LoginCommand{
			Email:    "inactive@example.com",
			Password: "password123",
		})
		assert.Error(t, err)
	})

	t.Run("empty email returns error", func(t *testing.T) {
		uc, _, _ := setupLoginTest(t)

		_, err := uc.Login(context.Background(), command.LoginCommand{
			Email:    "",
			Password: "password123",
		})
		assert.Error(t, err)
	})

	t.Run("empty password returns error", func(t *testing.T) {
		uc, _, _ := setupLoginTest(t)

		_, err := uc.Login(context.Background(), command.LoginCommand{
			Email:    "test@example.com",
			Password: "",
		})
		assert.Error(t, err)
	})
}

func TestLoginUseCase_RefreshToken(t *testing.T) {
	t.Run("valid refresh token returns new tokens", func(t *testing.T) {
		uc, userRepo, authRepo := setupLoginTest(t)

		user, _ := model.NewUser("test@example.com", "password123", "Test", "admin", nil)

		// Login expectations
		userRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(user, nil)
		authRepo.On("StoreRefreshToken", mock.Anything, user.ID, mock.Anything, mock.Anything).Return(nil)

		loginResult, err := uc.Login(context.Background(), command.LoginCommand{
			Email:    "test@example.com",
			Password: "password123",
		})
		require.NoError(t, err)

		// RefreshToken expectations
		authRepo.On("ValidateRefreshToken", mock.Anything, mock.Anything).Return(user.ID, nil)
		authRepo.On("RevokeRefreshToken", mock.Anything, mock.Anything).Return(nil)
		userRepo.On("GetByID", mock.Anything, user.ID).Return(user, nil)
		authRepo.On("StoreRefreshToken", mock.Anything, user.ID, mock.Anything, mock.Anything).Return(nil)

		result, err := uc.RefreshToken(context.Background(), command.RefreshTokenCommand{
			RefreshToken: loginResult.RefreshToken,
		})
		require.NoError(t, err)
		assert.NotEmpty(t, result.AccessToken)
		assert.NotEmpty(t, result.RefreshToken)
		assert.NotEqual(t, loginResult.RefreshToken, result.RefreshToken)

		// Verify revoke was called
		authRepo.AssertCalled(t, "RevokeRefreshToken", mock.Anything, mock.Anything)
	})

	t.Run("invalid refresh token returns error", func(t *testing.T) {
		uc, _, authRepo := setupLoginTest(t)

		authRepo.On("ValidateRefreshToken", mock.Anything, mock.Anything).Return("", nil)

		_, err := uc.RefreshToken(context.Background(), command.RefreshTokenCommand{
			RefreshToken: "invalid-token",
		})
		assert.Error(t, err)
	})

	t.Run("empty refresh token returns error", func(t *testing.T) {
		uc, _, _ := setupLoginTest(t)

		_, err := uc.RefreshToken(context.Background(), command.RefreshTokenCommand{
			RefreshToken: "",
		})
		assert.Error(t, err)
	})
}

func TestLoginUseCase_ValidateToken(t *testing.T) {
	t.Run("valid token returns claims", func(t *testing.T) {
		uc, userRepo, authRepo := setupLoginTest(t)

		user, _ := model.NewUser("test@example.com", "password123", "Test", "admin", nil)

		userRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(user, nil)
		authRepo.On("StoreRefreshToken", mock.Anything, user.ID, mock.Anything, mock.Anything).Return(nil)

		result, err := uc.Login(context.Background(), command.LoginCommand{
			Email:    "test@example.com",
			Password: "password123",
		})
		require.NoError(t, err)

		validateResult, err := uc.ValidateToken(context.Background(), command.ValidateTokenCommand{
			Token: result.AccessToken,
		})
		require.NoError(t, err)
		assert.Equal(t, user.ID, validateResult.UserID)
		assert.Equal(t, "admin", validateResult.Role)
	})

	t.Run("invalid token returns error", func(t *testing.T) {
		uc, _, _ := setupLoginTest(t)

		_, err := uc.ValidateToken(context.Background(), command.ValidateTokenCommand{
			Token: "invalid.jwt.token",
		})
		assert.Error(t, err)
	})

	t.Run("empty token returns error", func(t *testing.T) {
		uc, _, _ := setupLoginTest(t)

		_, err := uc.ValidateToken(context.Background(), command.ValidateTokenCommand{
			Token: "",
		})
		assert.Error(t, err)
	})
}

func TestLoginUseCase_GetPublicKey(t *testing.T) {
	t.Run("returns PEM public key", func(t *testing.T) {
		uc, _, _ := setupLoginTest(t)

		result, err := uc.GetPublicKey(context.Background())
		require.NoError(t, err)
		assert.Contains(t, result.PublicKeyPEM, "-----BEGIN PUBLIC KEY-----")
	})
}

func TestHashToken(t *testing.T) {
	t.Run("produces consistent hash", func(t *testing.T) {
		h1 := hashToken("test-token")
		h2 := hashToken("test-token")
		assert.Equal(t, h1, h2)
	})

	t.Run("different tokens produce different hashes", func(t *testing.T) {
		h1 := hashToken("token-1")
		h2 := hashToken("token-2")
		assert.NotEqual(t, h1, h2)
	})
}
