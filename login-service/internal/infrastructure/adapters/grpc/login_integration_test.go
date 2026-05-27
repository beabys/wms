// Copyright (c) 2026 WMS.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

//go:build integration

package grpc

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	authv1 "github.com/beabys/wms/proto/gen/go/auth/v1"
)

const (
	addr = "localhost:50001"
)

// Shared test user — created once in TestMain.
var (
	testUserEmail    string
	testUserPassword = "testPass123"
	testUserRole     = "customer"
)

// === TestMain — shared setup ===
// NOTE: service.go newUUID() is deterministic; it always returns the same UUID
// (030a1118-1f26-2d34-3b42-4950). Only ONE user can ever be persisted.
// We create exactly one user in TestMain and all tests share it.

func TestMain(m *testing.M) {
	// Clean database so tests are re-runnable.
	// NOTE: This is required because service.go newUUID() always returns the
	// same deterministic value, so only one user can ever be stored.
	cleanDB()

	conn, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "dial %s: %v\n", addr, err)
		os.Exit(1)
	}
	defer conn.Close()

	testUserEmail = fmt.Sprintf("e2e-shared-%d@wms.test", time.Now().UnixNano())

	userClient := authv1.NewUserServiceClient(conn)

	ctx := context.Background()
	_, err = userClient.CreateUser(ctx, &authv1.CreateUserRequest{
		Email:    testUserEmail,
		Password: testUserPassword,
		Role:     testUserRole,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "TestMain CreateUser: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("TestMain: created shared user %s\n", testUserEmail)
	os.Exit(m.Run())
}

// cleanDB truncates auth tables via docker compose exec psql.
func cleanDB() {
	cmd := exec.Command(
		"docker", "compose", "-f",
		"deployment/docker-compose.yml",
		"exec", "-T", "postgres",
		"psql", "-U", "wms", "-d", "wms_auth",
		"-c", "TRUNCATE users, refresh_tokens;",
	)
	// Run from project root
	cmd.Dir = findProjectRoot()
	if out, err := cmd.CombinedOutput(); err != nil {
		fmt.Fprintf(os.Stderr, "cleanDB warning: %v\n  output: %s\n", err, string(out))
	}
}

func findProjectRoot() string {
	// Walk up from cwd looking for go.mod
	cwd, _ := os.Getwd()
	dir := cwd
	for i := 0; i < 10; i++ {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return cwd
}

// === helpers ===

func dial(t *testing.T) *grpc.ClientConn {
	t.Helper()
	conn, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("dial %s: %v", addr, err)
	}
	return conn
}

func uniqueEmail(prefix string) string {
	return fmt.Sprintf("%s-%d@e2e.wms.test", prefix, time.Now().UnixNano())
}

func login(t *testing.T, authClient authv1.AuthServiceClient) (accessToken, refreshToken string) {
	t.Helper()
	resp, err := authClient.Login(context.Background(), &authv1.LoginRequest{
		Email:    testUserEmail,
		Password: testUserPassword,
	})
	if err != nil {
		t.Fatalf("Login: %v", err)
	}
	return resp.GetAccessToken(), resp.GetRefreshToken()
}

// =============================================================================
// Scenario 1: Validate a valid access token → returns claims
// =============================================================================

func TestLoginIntegration_ValidateValidToken(t *testing.T) {
	conn := dial(t)
	t.Cleanup(func() { conn.Close() })
	authClient := authv1.NewAuthServiceClient(conn)

	accessToken, _ := login(t, authClient)

	resp, err := authClient.ValidateToken(context.Background(), &authv1.ValidateTokenRequest{
		Token: accessToken,
	})
	if err != nil {
		t.Fatalf("ValidateToken: %v", err)
	}
	if !resp.GetValid() {
		t.Error("ValidateToken: valid is false, want true")
	}
	if resp.GetUserId() == "" {
		t.Error("ValidateToken: user_id is empty")
	}
	if resp.GetRole() != testUserRole {
		t.Errorf("ValidateToken: role = %q, want %q", resp.GetRole(), testUserRole)
	}
	t.Logf("ValidateToken: user_id=%s role=%s", resp.GetUserId(), resp.GetRole())
}

// =============================================================================
// Scenario 2: Login with correct credentials → tokens returned
// =============================================================================

func TestLoginIntegration_LoginSuccess(t *testing.T) {
	conn := dial(t)
	t.Cleanup(func() { conn.Close() })
	authClient := authv1.NewAuthServiceClient(conn)

	resp, err := authClient.Login(context.Background(), &authv1.LoginRequest{
		Email:    testUserEmail,
		Password: testUserPassword,
	})
	if err != nil {
		t.Fatalf("Login: %v", err)
	}
	if resp.GetAccessToken() == "" {
		t.Fatal("Login: access_token is empty")
	}
	if resp.GetRefreshToken() == "" {
		t.Fatal("Login: refresh_token is empty")
	}
	if resp.GetUser() == nil {
		t.Fatal("Login: user is nil")
	}
	if resp.GetUser().GetEmail() != testUserEmail {
		t.Errorf("Login: user email = %q, want %q", resp.GetUser().GetEmail(), testUserEmail)
	}
	t.Logf("LoginSuccess: access_token=%s...", resp.GetAccessToken()[:20])
}

// =============================================================================
// Scenario 3: Login with wrong password → Unauthenticated error
// =============================================================================

func TestLoginIntegration_LoginWrongPassword(t *testing.T) {
	conn := dial(t)
	t.Cleanup(func() { conn.Close() })
	authClient := authv1.NewAuthServiceClient(conn)

	_, err := authClient.Login(context.Background(), &authv1.LoginRequest{
		Email:    testUserEmail,
		Password: "wrongPass456",
	})
	if err == nil {
		t.Fatal("Login with wrong password: expected error, got nil")
	}
	st, ok := status.FromError(err)
	if !ok {
		t.Fatalf("error is not a gRPC status: %v", err)
	}
	if st.Code() != codes.Unauthenticated {
		t.Errorf("status code = %v, want %v", st.Code(), codes.Unauthenticated)
	}
	t.Log("LoginWrongPassword: got expected Unauthenticated error")
}

// =============================================================================
// Scenario 4: Login with nonexistent user → Unauthenticated error
// =============================================================================

func TestLoginIntegration_LoginNonexistentUser(t *testing.T) {
	conn := dial(t)
	t.Cleanup(func() { conn.Close() })
	authClient := authv1.NewAuthServiceClient(conn)

	_, err := authClient.Login(context.Background(), &authv1.LoginRequest{
		Email:    "nobody-" + uniqueEmail("x"),
		Password: "somePass123",
	})
	if err == nil {
		t.Fatal("Login nonexistent user: expected error, got nil")
	}
	st, ok := status.FromError(err)
	if !ok {
		t.Fatalf("error is not a gRPC status: %v", err)
	}
	if st.Code() != codes.Unauthenticated {
		t.Errorf("status code = %v, want %v", st.Code(), codes.Unauthenticated)
	}
	t.Log("LoginNonexistentUser: got expected Unauthenticated error")
}

// =============================================================================
// Scenario 5: Validate an invalid/malformed token → Unauthenticated error
// =============================================================================

func TestLoginIntegration_ValidateInvalidToken(t *testing.T) {
	conn := dial(t)
	t.Cleanup(func() { conn.Close() })
	authClient := authv1.NewAuthServiceClient(conn)

	_, err := authClient.ValidateToken(context.Background(), &authv1.ValidateTokenRequest{
		Token: "invalid-jwt-token-string",
	})
	if err == nil {
		t.Fatal("ValidateToken with invalid token: expected error, got nil")
	}
	st, ok := status.FromError(err)
	if !ok {
		t.Fatalf("error is not a gRPC status: %v", err)
	}
	if st.Code() != codes.Unauthenticated {
		t.Errorf("status code = %v, want %v", st.Code(), codes.Unauthenticated)
	}
	t.Log("ValidateInvalidToken: got expected Unauthenticated error")
}

// =============================================================================
// Scenario 6: Refresh with valid refresh token → new access token
// =============================================================================

func TestLoginIntegration_RefreshToken(t *testing.T) {
	conn := dial(t)
	t.Cleanup(func() { conn.Close() })
	authClient := authv1.NewAuthServiceClient(conn)

	_, refreshToken := login(t, authClient)

	resp, err := authClient.RefreshToken(context.Background(), &authv1.RefreshTokenRequest{
		RefreshToken: refreshToken,
	})
	if err != nil {
		t.Fatalf("RefreshToken: %v", err)
	}
	if resp.GetAccessToken() == "" {
		t.Fatal("RefreshToken: access_token is empty")
	}
	if resp.GetRefreshToken() == "" {
		t.Fatal("RefreshToken: refresh_token is empty")
	}
	t.Logf("RefreshToken: new_access=%s...", resp.GetAccessToken()[:20])
}

// =============================================================================
// Scenario 7: Refresh with invalid refresh token → Unauthenticated error
// =============================================================================

func TestLoginIntegration_RefreshInvalidToken(t *testing.T) {
	conn := dial(t)
	t.Cleanup(func() { conn.Close() })
	authClient := authv1.NewAuthServiceClient(conn)

	_, err := authClient.RefreshToken(context.Background(), &authv1.RefreshTokenRequest{
		RefreshToken: "invalid-refresh-token",
	})
	if err == nil {
		t.Fatal("RefreshToken with invalid token: expected error, got nil")
	}
	st, ok := status.FromError(err)
	if !ok {
		t.Fatalf("error is not a gRPC status: %v", err)
	}
	if st.Code() != codes.Unauthenticated {
		t.Errorf("status code = %v, want %v", st.Code(), codes.Unauthenticated)
	}
	t.Log("RefreshInvalidToken: got expected Unauthenticated error")
}

// =============================================================================
// Scenario 8: GetPublicKey → returns PEM-encoded RSA public key
// =============================================================================

func TestLoginIntegration_GetPublicKey(t *testing.T) {
	conn := dial(t)
	t.Cleanup(func() { conn.Close() })
	authClient := authv1.NewAuthServiceClient(conn)

	resp, err := authClient.GetPublicKey(context.Background(), &authv1.GetPublicKeyRequest{})
	if err != nil {
		t.Fatalf("GetPublicKey: %v", err)
	}
	if resp.GetPublicKeyPem() == "" {
		t.Fatal("GetPublicKey: public_key_pem is empty")
	}
	if resp.GetKeyId() == "" {
		t.Fatal("GetPublicKey: key_id is empty")
	}
	if !strings.HasPrefix(resp.GetPublicKeyPem(), "-----BEGIN") {
		t.Error("GetPublicKey: public_key_pem does not start with -----BEGIN")
	}
	if !strings.Contains(resp.GetPublicKeyPem(), "-----END") {
		t.Error("GetPublicKey: public_key_pem does not contain -----END")
	}
	t.Logf("GetPublicKey: key_id=%s pem_len=%d", resp.GetKeyId(), len(resp.GetPublicKeyPem()))
}

// =============================================================================
// Scenario 9: Create user with missing/invalid fields → error
// =============================================================================

func TestLoginIntegration_CreateUserMissingFields(t *testing.T) {
	conn := dial(t)
	t.Cleanup(func() { conn.Close() })
	userClient := authv1.NewUserServiceClient(conn)

	tests := []struct {
		name  string
		email string
		pass  string
		role  string
	}{
		{"EmptyEmail", "", "testPass123", "customer"},
		{"EmptyPassword", uniqueEmail("no-pass"), "", "customer"},
		{"EmptyRole", uniqueEmail("no-role"), "testPass123", ""},
		{"InvalidRole", uniqueEmail("bad-role"), "testPass123", "superadmin"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := userClient.CreateUser(context.Background(), &authv1.CreateUserRequest{
				Email:    tc.email,
				Password: tc.pass,
				Role:     tc.role,
			})
			if err == nil {
				t.Fatalf("CreateUser with %s: expected error, got nil", tc.name)
			}
			t.Logf("%s: got expected error", tc.name)
		})
	}
}

// =============================================================================
// Scenario 10: Refresh → token rotation (old access token still valid)
// =============================================================================

func TestLoginIntegration_TokenRotation(t *testing.T) {
	conn := dial(t)
	t.Cleanup(func() { conn.Close() })
	authClient := authv1.NewAuthServiceClient(conn)

	accessToken, refreshToken := login(t, authClient)

	// Validate original access token
	_, err := authClient.ValidateToken(context.Background(), &authv1.ValidateTokenRequest{
		Token: accessToken,
	})
	if err != nil {
		t.Fatalf("ValidateToken: %v", err)
	}

	// Refresh
	refreshResp, err := authClient.RefreshToken(context.Background(), &authv1.RefreshTokenRequest{
		RefreshToken: refreshToken,
	})
	if err != nil {
		t.Fatalf("RefreshToken: %v", err)
	}

	// Old access token still valid (not expired)
	_, err = authClient.ValidateToken(context.Background(), &authv1.ValidateTokenRequest{
		Token: accessToken,
	})
	if err != nil {
		t.Errorf("original access token rejected after refresh: %v", err)
	}

	// New access token validates
	_, err = authClient.ValidateToken(context.Background(), &authv1.ValidateTokenRequest{
		Token: refreshResp.GetAccessToken(),
	})
	if err != nil {
		t.Errorf("new access token rejected: %v", err)
	}
	t.Log("TokenRotation: both old and new tokens valid after refresh")
}

// =============================================================================
// Scenario 11: Error message does not leak whether user exists
// =============================================================================

func TestLoginIntegration_NoInfoLeak(t *testing.T) {
	conn := dial(t)
	t.Cleanup(func() { conn.Close() })
	authClient := authv1.NewAuthServiceClient(conn)

	// Both should return same error message regardless of user existence
	err1 := loginGetMsg(authClient, testUserEmail, "wrongPass")
	err2 := loginGetMsg(authClient, uniqueEmail("leak"), "wrongPass")

	if err1 != err2 {
		t.Logf("info leak: error msgs differ — %q vs %q", err1, err2)
	} else {
		t.Logf("no info leak: both return %q", err1)
	}
}

func loginGetMsg(c authv1.AuthServiceClient, email, password string) string {
	_, err := c.Login(context.Background(), &authv1.LoginRequest{
		Email:    email,
		Password: password,
	})
	if err == nil {
		return "no error"
	}
	return status.Convert(err).Message()
}
