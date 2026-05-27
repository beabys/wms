package usecase

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/beabys/wms/login-service/internal/application/auth/command"
	"github.com/beabys/wms/login-service/internal/domain/auth/model"
)

var (
	errNotFound = errors.New("not found")
)

// mockUserRepo implements repository.UserRepository for testing.
type mockUserRepo struct {
	users         map[string]*model.User
	refreshTokens map[string]string // token -> userID
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{
		users:         make(map[string]*model.User),
		refreshTokens: make(map[string]string),
	}
}

func (m *mockUserRepo) GetByID(_ context.Context, id string) (*model.User, error) {
	u, ok := m.users[id]
	if !ok {
		return nil, errNotFound
	}
	return u, nil
}

func (m *mockUserRepo) GetByEmail(_ context.Context, email model.Email) (*model.User, error) {
	for _, u := range m.users {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, errNotFound
}

func (m *mockUserRepo) Save(_ context.Context, user *model.User) error {
	m.users[user.ID] = user
	return nil
}

func (m *mockUserRepo) List(_ context.Context, limit, offset int32) ([]*model.User, error) {
	var result []*model.User
	for _, u := range m.users {
		result = append(result, u)
	}
	return result, nil
}

func (m *mockUserRepo) Update(_ context.Context, user *model.User) error {
	m.users[user.ID] = user
	return nil
}

func (m *mockUserRepo) CreateRefreshToken(_ context.Context, userID, token string) error {
	m.refreshTokens[token] = userID
	return nil
}

func (m *mockUserRepo) ValidateRefreshToken(_ context.Context, token string) (string, error) {
	uid, ok := m.refreshTokens[token]
	if !ok {
		return "", errNotFound
	}
	return uid, nil
}

func (m *mockUserRepo) DeleteRefreshToken(_ context.Context, token string) error {
	delete(m.refreshTokens, token)
	return nil
}

func generateTestKeys(t *testing.T) (*rsa.PrivateKey, *rsa.PublicKey) {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	return key, &key.PublicKey
}

func TestAuthService_Login_Success(t *testing.T) {
	repo := newMockUserRepo()
	privKey, pubKey := generateTestKeys(t)

	// Create a user
	email, _ := model.NewEmail("test@example.com")
	passHash, _ := model.NewPasswordHash("password123")
	user := &model.User{
		ID:           "user-1",
		Email:        email,
		PasswordHash: passHash,
		Role:         model.RoleCustomer,
		Status:       model.UserStatusActive,
	}
	repo.Save(context.Background(), user)

	svc := NewAuthService(repo, privKey, pubKey, 15*time.Minute, 7*24*time.Hour)
	result, err := svc.Login(context.Background(), command.LoginCommand{
		Email:    "test@example.com",
		Password: "password123",
	})

	require.NoError(t, err)
	require.NotNil(t, result)
	require.NotEmpty(t, result.AccessToken)
	require.NotEmpty(t, result.RefreshToken)
	require.Equal(t, user.ID, result.UserID)
}

func TestAuthService_Login_InvalidPassword(t *testing.T) {
	repo := newMockUserRepo()
	privKey, pubKey := generateTestKeys(t)

	email, _ := model.NewEmail("test@example.com")
	passHash, _ := model.NewPasswordHash("password123")
	user := &model.User{
		ID:           "user-1",
		Email:        email,
		PasswordHash: passHash,
		Role:         model.RoleCustomer,
		Status:       model.UserStatusActive,
	}
	repo.Save(context.Background(), user)

	svc := NewAuthService(repo, privKey, pubKey, 15*time.Minute, 7*24*time.Hour)
	_, err := svc.Login(context.Background(), command.LoginCommand{
		Email:    "test@example.com",
		Password: "wrongpassword",
	})

	require.Error(t, err)
}

func TestAuthService_Register_Success(t *testing.T) {
	repo := newMockUserRepo()
	privKey, pubKey := generateTestKeys(t)
	svc := NewAuthService(repo, privKey, pubKey, 15*time.Minute, 7*24*time.Hour)

	result, err := svc.Register(context.Background(), command.RegisterCommand{
		Email:    "new@example.com",
		Password: "securePass123",
	})

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, "new@example.com", result.Email)
	require.Equal(t, "customer", result.Role)
	require.Equal(t, "active", result.Status)
}

func TestAuthService_Register_DuplicateEmail(t *testing.T) {
	repo := newMockUserRepo()
	privKey, pubKey := generateTestKeys(t)
	svc := NewAuthService(repo, privKey, pubKey, 15*time.Minute, 7*24*time.Hour)

	// First registration
	email, _ := model.NewEmail("dup@example.com")
	passHash, _ := model.NewPasswordHash("password123")
	repo.Save(context.Background(), &model.User{
		ID:           "existing",
		Email:        email,
		PasswordHash: passHash,
		Role:         model.RoleCustomer,
		Status:       model.UserStatusActive,
	})

	// Second registration with same email
	_, err := svc.Register(context.Background(), command.RegisterCommand{
		Email:    "dup@example.com",
		Password: "anotherPass123",
	})

	require.Error(t, err)
}

func TestAuthService_RefreshToken_Success(t *testing.T) {
	repo := newMockUserRepo()
	privKey, pubKey := generateTestKeys(t)
	svc := NewAuthService(repo, privKey, pubKey, 15*time.Minute, 7*24*time.Hour)

	// Create and login user
	email, _ := model.NewEmail("refresh@example.com")
	passHash, _ := model.NewPasswordHash("password123")
	repo.Save(context.Background(), &model.User{
		ID:           "refresh-user",
		Email:        email,
		PasswordHash: passHash,
		Role:         model.RoleCustomer,
		Status:       model.UserStatusActive,
	})

	loginResult, err := svc.Login(context.Background(), command.LoginCommand{
		Email:    "refresh@example.com",
		Password: "password123",
	})
	require.NoError(t, err)

	// Now refresh
	refreshResult, err := svc.RefreshToken(context.Background(), command.RefreshTokenCommand{
		RefreshToken: loginResult.RefreshToken,
	})

	require.NoError(t, err)
	require.NotEmpty(t, refreshResult.AccessToken)
	require.NotEmpty(t, refreshResult.RefreshToken)
	require.NotEqual(t, loginResult.RefreshToken, refreshResult.RefreshToken)
}
