// Copyright (c) 2026 WMS.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

//go:build integration

package grpc_test

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	commonv1 "github.com/beabys/wms/proto/gen/go/common/v1"
	customerv1 "github.com/beabys/wms/proto/gen/go/customer/v1"
)

const (
	addr = "localhost:50002"
)

// --- helpers ---

func testToken(t *testing.T) string {
	t.Helper()

	key, err := loadPrivateKey()
	if err != nil {
		t.Logf("warning: could not load private key, generating temp key: %v", err)
		key, err = rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			t.Fatalf("generate temp key: %v", err)
		}
	}

	now := time.Now()
	claims := jwt.MapClaims{
		"sub": "e2e-test-ops",
		"iat": now.Unix(),
		"exp": now.Add(1 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signed, err := token.SignedString(key)
	if err != nil {
		t.Fatalf("sign test token: %v", err)
	}
	return signed
}

func loadPrivateKey() (*rsa.PrivateKey, error) {
	candidates := []string{
		"login-service/keys/private.pem",
		"../login-service/keys/private.pem",
		"../../login-service/keys/private.pem",
		"../../../login-service/keys/private.pem",
		"../../../../login-service/keys/private.pem",
		"../../../../../login-service/keys/private.pem",
	}
	var firstErr error
	for _, path := range candidates {
		key, err := readPrivateKey(path)
		if err == nil {
			return key, nil
		}
		if firstErr == nil {
			firstErr = err
		}
	}
	return nil, fmt.Errorf("no private key found (tried %v): %w", candidates, firstErr)
}

func readPrivateKey(path string) (*rsa.PrivateKey, error) {
	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("no PEM block found")
	}
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	rsaKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("private key is not RSA")
	}
	return rsaKey, nil
}

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

func authContext(t *testing.T) context.Context {
	t.Helper()
	return metadata.AppendToOutgoingContext(
		context.Background(),
		"authorization", "Bearer "+testToken(t),
	)
}

func uniqueEmail(prefix string) string {
	return fmt.Sprintf("%s-%d@e2e.wms.test", prefix, time.Now().UnixNano())
}

func uniqueVAT(base string) string {
	return fmt.Sprintf("DE%09d", time.Now().UnixNano()%1000000000)
}

// =============================================================================
// Scenario 1: Invite a customer → success (returns invitation link)
// =============================================================================

func TestCustomerIntegration_InviteCustomer(t *testing.T) {
	conn := dial(t)
	t.Cleanup(func() { conn.Close() })
	client := customerv1.NewCustomerServiceClient(conn)

	req := &customerv1.InviteCustomerRequest{
		Email: uniqueEmail("invite"),
	}
	resp, err := client.InviteCustomer(context.Background(), req)
	if err != nil {
		t.Fatalf("InviteCustomer: %v", err)
	}
	if resp.GetInvitationLink() == "" {
		t.Fatal("InviteCustomer: invitation_link is empty")
	}
	t.Logf("InviteCustomer: link=%s", resp.GetInvitationLink())
}

// =============================================================================
// Scenario 2: Create/register a customer → success (status=pending)
// =============================================================================

func TestCustomerIntegration_CreateCustomer(t *testing.T) {
	conn := dial(t)
	t.Cleanup(func() { conn.Close() })
	client := customerv1.NewCustomerServiceClient(conn)

	req := &customerv1.CreateCustomerRequest{
		CompanyName: "E2E Test GmbH",
		VatNumber:   uniqueVAT("DE123456789"),
		Address: &commonv1.Address{
			Line1:      "Teststr 1",
			City:       "Berlin",
			PostalCode: "10115",
			Country:    "DE",
		},
	}
	resp, err := client.CreateCustomer(context.Background(), req)
	if err != nil {
		t.Fatalf("CreateCustomer: %v", err)
	}
	if resp.GetCustomer() == nil {
		t.Fatal("CreateCustomer: customer is nil")
	}
	if resp.GetCustomer().GetId() == "" {
		t.Fatal("CreateCustomer: customer id is empty")
	}
	if resp.GetCustomer().GetCompanyName() != "E2E Test GmbH" {
		t.Errorf("CreateCustomer: company_name = %q, want %q",
			resp.GetCustomer().GetCompanyName(), "E2E Test GmbH")
	}
	if resp.GetCustomer().GetStatus() != "pending" {
		t.Errorf("CreateCustomer: status = %q, want %q",
			resp.GetCustomer().GetStatus(), "pending")
	}
	t.Logf("CreateCustomer: id=%s status=%s", resp.GetCustomer().GetId(), resp.GetCustomer().GetStatus())
}

// =============================================================================
// Scenario 3: Get customer by ID → correct data
// =============================================================================

func TestCustomerIntegration_GetCustomer(t *testing.T) {
	conn := dial(t)
	t.Cleanup(func() { conn.Close() })
	client := customerv1.NewCustomerServiceClient(conn)

	// Create a customer first
	createResp, err := client.CreateCustomer(context.Background(), &customerv1.CreateCustomerRequest{
		CompanyName: "GetTest Inc",
		VatNumber:   uniqueVAT("DE987654321"),
		Address: &commonv1.Address{
			Line1:      "Getstr 5",
			City:       "Munich",
			PostalCode: "80331",
			Country:    "DE",
		},
	})
	if err != nil {
		t.Fatalf("CreateCustomer (setup): %v", err)
	}
	if createResp.GetCustomer() == nil || createResp.GetCustomer().GetId() == "" {
		t.Fatal("CreateCustomer (setup): no id returned")
	}
	customerID := createResp.GetCustomer().GetId()
	t.Logf("created customer: id=%s", customerID)

	// Get customer by ID
	resp, err := client.GetCustomer(authContext(t), &customerv1.GetCustomerRequest{Id: customerID})
	if err != nil {
		t.Fatalf("GetCustomer: %v", err)
	}
	if resp.GetCustomer() == nil {
		t.Fatal("GetCustomer: customer is nil")
	}
	if resp.GetCustomer().GetId() != customerID {
		t.Errorf("GetCustomer: id = %q, want %q", resp.GetCustomer().GetId(), customerID)
	}
	if resp.GetCustomer().GetCompanyName() != "GetTest Inc" {
		t.Errorf("GetCustomer: company_name = %q, want %q",
			resp.GetCustomer().GetCompanyName(), "GetTest Inc")
	}
	if resp.GetCustomer().GetStatus() == "" {
		t.Error("GetCustomer: status is empty")
	}
	t.Logf("GetCustomer: id=%s company=%s status=%s",
		resp.GetCustomer().GetId(),
		resp.GetCustomer().GetCompanyName(),
		resp.GetCustomer().GetStatus())
}

// =============================================================================
// Scenario 4: Approve a customer → success (status=active)
// =============================================================================

func TestCustomerIntegration_ApproveCustomer(t *testing.T) {
	conn := dial(t)
	t.Cleanup(func() { conn.Close() })
	client := customerv1.NewCustomerServiceClient(conn)

	// Create a pending customer
	createResp, err := client.CreateCustomer(context.Background(), &customerv1.CreateCustomerRequest{
		CompanyName: "ApproveTest AG",
		VatNumber:   uniqueVAT("DE111222333"),
		Address: &commonv1.Address{
			Line1:      "Approvstr 10",
			City:       "Hamburg",
			PostalCode: "20095",
			Country:    "DE",
		},
	})
	if err != nil {
		t.Fatalf("CreateCustomer (setup): %v", err)
	}
	if createResp.GetCustomer() == nil || createResp.GetCustomer().GetId() == "" {
		t.Fatal("CreateCustomer (setup): no id returned")
	}
	customerID := createResp.GetCustomer().GetId()
	t.Logf("created pending customer: id=%s", customerID)

	// Approve the customer
	resp, err := client.ApproveCustomer(authContext(t), &customerv1.ApproveCustomerRequest{
		CustomerId: customerID,
	})
	if err != nil {
		t.Fatalf("ApproveCustomer: %v", err)
	}
	if resp.GetCustomer() == nil {
		t.Fatal("ApproveCustomer: customer is nil")
	}
	if resp.GetCustomer().GetId() != customerID {
		t.Errorf("ApproveCustomer: id = %q, want %q", resp.GetCustomer().GetId(), customerID)
	}
	if resp.GetCustomer().GetStatus() != "active" {
		t.Errorf("ApproveCustomer: status = %q, want %q",
			resp.GetCustomer().GetStatus(), "active")
	}
	t.Logf("ApproveCustomer: id=%s status=%s", resp.GetCustomer().GetId(), resp.GetCustomer().GetStatus())
}

// =============================================================================
// Scenario 5: Suspend an approved customer → success (status=suspended)
// =============================================================================

func TestCustomerIntegration_SuspendCustomer(t *testing.T) {
	conn := dial(t)
	t.Cleanup(func() { conn.Close() })
	client := customerv1.NewCustomerServiceClient(conn)

	// Create a pending customer
	createResp, err := client.CreateCustomer(context.Background(), &customerv1.CreateCustomerRequest{
		CompanyName: "SuspendTest Ltd",
		VatNumber:   uniqueVAT("DE444555666"),
		Address: &commonv1.Address{
			Line1:      "Suspstr 3",
			City:       "Frankfurt",
			PostalCode: "60311",
			Country:    "DE",
		},
	})
	if err != nil {
		t.Fatalf("CreateCustomer (setup): %v", err)
	}
	if createResp.GetCustomer() == nil || createResp.GetCustomer().GetId() == "" {
		t.Fatal("CreateCustomer (setup): no id returned")
	}
	customerID := createResp.GetCustomer().GetId()
	t.Logf("created pending customer: id=%s", customerID)

	// Approve first (required before suspend)
	_, err = client.ApproveCustomer(authContext(t), &customerv1.ApproveCustomerRequest{
		CustomerId: customerID,
	})
	if err != nil {
		t.Fatalf("ApproveCustomer (setup): %v", err)
	}
	t.Log("approved customer")

	// Suspend the approved customer
	resp, err := client.SuspendCustomer(authContext(t), &customerv1.SuspendCustomerRequest{
		CustomerId: customerID,
		Reason:     "e2e test suspension",
	})
	if err != nil {
		t.Fatalf("SuspendCustomer: %v", err)
	}
	if resp.GetCustomer() == nil {
		t.Fatal("SuspendCustomer: customer is nil")
	}
	if resp.GetCustomer().GetId() != customerID {
		t.Errorf("SuspendCustomer: id = %q, want %q", resp.GetCustomer().GetId(), customerID)
	}
	if resp.GetCustomer().GetStatus() != "suspended" {
		t.Errorf("SuspendCustomer: status = %q, want %q",
			resp.GetCustomer().GetStatus(), "suspended")
	}
	t.Logf("SuspendCustomer: id=%s status=%s", resp.GetCustomer().GetId(), resp.GetCustomer().GetStatus())
}

// =============================================================================
// Scenario 6: List customers → returns filtered results
// =============================================================================

func TestCustomerIntegration_ListCustomers(t *testing.T) {
	conn := dial(t)
	t.Cleanup(func() { conn.Close() })
	client := customerv1.NewCustomerServiceClient(conn)

	// Create a customer for listing
	createResp, err := client.CreateCustomer(context.Background(), &customerv1.CreateCustomerRequest{
		CompanyName: "ListTest Corp",
		VatNumber:   uniqueVAT("DE777888999"),
		Address: &commonv1.Address{
			Line1:      "Liststr 20",
			City:       "Cologne",
			PostalCode: "50667",
			Country:    "DE",
		},
	})
	if err != nil {
		t.Fatalf("CreateCustomer (setup): %v", err)
	}
	if createResp.GetCustomer() == nil || createResp.GetCustomer().GetId() == "" {
		t.Fatal("CreateCustomer (setup): no id returned")
	}
	customerID := createResp.GetCustomer().GetId()
	t.Logf("created customer for list: id=%s", customerID)

	// List customers
	resp, err := client.ListCustomers(authContext(t), &customerv1.ListCustomersRequest{
		Pagination: &commonv1.Pagination{
			Page:  1,
			Limit: 100,
		},
	})
	if err != nil {
		t.Fatalf("ListCustomers: %v", err)
	}
	if len(resp.GetCustomers()) == 0 {
		t.Fatal("ListCustomers: empty results, expected at least 1 customer")
	}
	if resp.GetPagination() == nil {
		t.Fatal("ListCustomers: pagination is nil")
	}
	if resp.GetPagination().GetTotal() < 1 {
		t.Errorf("ListCustomers: total = %d, want >= 1", resp.GetPagination().GetTotal())
	}

	// Verify our test customer appears in results
	found := false
	for _, c := range resp.GetCustomers() {
		if c.GetId() == customerID {
			found = true
			if c.GetCompanyName() != "ListTest Corp" {
				t.Errorf("ListCustomers: customer %s company_name = %q, want %q",
					c.GetId(), c.GetCompanyName(), "ListTest Corp")
			}
			break
		}
	}
	if !found {
		t.Errorf("ListCustomers: customer %s not found in results", customerID)
	}
	t.Logf("ListCustomers: total=%d page=%d limit=%d",
		resp.GetPagination().GetTotal(),
		resp.GetPagination().GetPage(),
		resp.GetPagination().GetLimit())
}

// =============================================================================
// Scenario 7: Invite with missing email → error
// =============================================================================

func TestCustomerIntegration_InviteCustomerValidation(t *testing.T) {
	conn := dial(t)
	t.Cleanup(func() { conn.Close() })
	client := customerv1.NewCustomerServiceClient(conn)

	req := &customerv1.InviteCustomerRequest{
		Email: "",
	}
	_, err := client.InviteCustomer(context.Background(), req)
	if err == nil {
		t.Fatal("InviteCustomer with empty email: expected error, got nil")
	}
	st, ok := status.FromError(err)
	if !ok {
		t.Fatalf("error is not a gRPC status: %v", err)
	}
	t.Logf("InviteCustomerValidation: code=%v msg=%v", st.Code(), st.Message())
}

// =============================================================================
// Scenario 8: Auth required — call protected RPC without token → error
// =============================================================================

func TestCustomerIntegration_AuthRequired(t *testing.T) {
	conn := dial(t)
	t.Cleanup(func() { conn.Close() })
	client := customerv1.NewCustomerServiceClient(conn)

	t.Run("GetCustomerNoAuth", func(t *testing.T) {
		_, err := client.GetCustomer(context.Background(), &customerv1.GetCustomerRequest{Id: "nonexistent"})
		if err == nil {
			t.Fatal("GetCustomer without auth: expected error, got nil")
		}
		st, ok := status.FromError(err)
		if !ok {
			t.Fatalf("error is not a gRPC status: %v", err)
		}
		if st.Code() != codes.Unauthenticated {
			t.Errorf("status code = %v, want %v", st.Code(), codes.Unauthenticated)
		}
		t.Logf("GetCustomerNoAuth: got expected Unauthenticated error")
	})

	t.Run("ListCustomersNoAuth", func(t *testing.T) {
		_, err := client.ListCustomers(context.Background(), &customerv1.ListCustomersRequest{
			Pagination: &commonv1.Pagination{Page: 1, Limit: 10},
		})
		if err == nil {
			t.Fatal("ListCustomers without auth: expected error, got nil")
		}
		st, ok := status.FromError(err)
		if !ok {
			t.Fatalf("error is not a gRPC status: %v", err)
		}
		if st.Code() != codes.Unauthenticated {
			t.Errorf("status code = %v, want %v", st.Code(), codes.Unauthenticated)
		}
		t.Logf("ListCustomersNoAuth: got expected Unauthenticated error")
	})

	t.Run("ApproveCustomerNoAuth", func(t *testing.T) {
		_, err := client.ApproveCustomer(context.Background(), &customerv1.ApproveCustomerRequest{
			CustomerId: "nonexistent",
		})
		if err == nil {
			t.Fatal("ApproveCustomer without auth: expected error, got nil")
		}
		st, ok := status.FromError(err)
		if !ok {
			t.Fatalf("error is not a gRPC status: %v", err)
		}
		if st.Code() != codes.Unauthenticated {
			t.Errorf("status code = %v, want %v", st.Code(), codes.Unauthenticated)
		}
		t.Logf("ApproveCustomerNoAuth: got expected Unauthenticated error")
	})
}

// =============================================================================
// Full state machine: Create → Get → Approve → Get → Suspend → Get
// =============================================================================

func TestCustomerIntegration_FullStateMachine(t *testing.T) {
	conn := dial(t)
	t.Cleanup(func() { conn.Close() })
	client := customerv1.NewCustomerServiceClient(conn)

	var customerID string

	// Step 1: Create customer (pending)
	t.Run("CreatePending", func(t *testing.T) {
		resp, err := client.CreateCustomer(context.Background(), &customerv1.CreateCustomerRequest{
			CompanyName: "StateMachine GmbH",
			VatNumber:   uniqueVAT("DE000000001"),
			Address: &commonv1.Address{
				Line1:      "StateStr 99",
				City:       "Stuttgart",
				PostalCode: "70173",
				Country:    "DE",
			},
		})
		if err != nil {
			t.Fatalf("CreateCustomer: %v", err)
		}
		if resp.GetCustomer() == nil || resp.GetCustomer().GetId() == "" {
			t.Fatal("CreateCustomer: no id returned")
		}
		if resp.GetCustomer().GetStatus() != "pending" {
			t.Errorf("CreateCustomer: status = %q, want %q", resp.GetCustomer().GetStatus(), "pending")
		}
		customerID = resp.GetCustomer().GetId()
		t.Logf("step 1 - created pending: id=%s", customerID)
	})

	if customerID == "" {
		t.Fatal("customerID not set after CreatePending, aborting")
	}

	// Step 2: Get customer → verify pending
	t.Run("GetAndVerifyPending", func(t *testing.T) {
		resp, err := client.GetCustomer(authContext(t), &customerv1.GetCustomerRequest{Id: customerID})
		if err != nil {
			t.Fatalf("GetCustomer: %v", err)
		}
		if resp.GetCustomer() == nil {
			t.Fatal("GetCustomer: customer is nil")
		}
		if resp.GetCustomer().GetStatus() != "pending" {
			t.Errorf("GetCustomer: status = %q, want %q", resp.GetCustomer().GetStatus(), "pending")
		}
		t.Log("step 2 - verified pending")
	})

	// Step 3: Approve → active
	t.Run("ApproveToActive", func(t *testing.T) {
		resp, err := client.ApproveCustomer(authContext(t), &customerv1.ApproveCustomerRequest{
			CustomerId: customerID,
		})
		if err != nil {
			t.Fatalf("ApproveCustomer: %v", err)
		}
		if resp.GetCustomer() == nil {
			t.Fatal("ApproveCustomer: customer is nil")
		}
		if resp.GetCustomer().GetStatus() != "active" {
			t.Errorf("ApproveCustomer: status = %q, want %q", resp.GetCustomer().GetStatus(), "active")
		}
		t.Log("step 3 - approved to active")
	})

	// Step 4: Get → verify active
	t.Run("GetAndVerifyActive", func(t *testing.T) {
		resp, err := client.GetCustomer(authContext(t), &customerv1.GetCustomerRequest{Id: customerID})
		if err != nil {
			t.Fatalf("GetCustomer: %v", err)
		}
		if resp.GetCustomer() == nil {
			t.Fatal("GetCustomer: customer is nil")
		}
		if resp.GetCustomer().GetStatus() != "active" {
			t.Errorf("GetCustomer: status = %q, want %q", resp.GetCustomer().GetStatus(), "active")
		}
		t.Log("step 4 - verified active")
	})

	// Step 5: Suspend → suspended
	t.Run("SuspendToSuspended", func(t *testing.T) {
		resp, err := client.SuspendCustomer(authContext(t), &customerv1.SuspendCustomerRequest{
			CustomerId: customerID,
			Reason:     "state machine e2e test",
		})
		if err != nil {
			t.Fatalf("SuspendCustomer: %v", err)
		}
		if resp.GetCustomer() == nil {
			t.Fatal("SuspendCustomer: customer is nil")
		}
		if resp.GetCustomer().GetStatus() != "suspended" {
			t.Errorf("SuspendCustomer: status = %q, want %q", resp.GetCustomer().GetStatus(), "suspended")
		}
		t.Log("step 5 - suspended")
	})

	// Step 6: Get → verify suspended
	t.Run("GetAndVerifySuspended", func(t *testing.T) {
		resp, err := client.GetCustomer(authContext(t), &customerv1.GetCustomerRequest{Id: customerID})
		if err != nil {
			t.Fatalf("GetCustomer: %v", err)
		}
		if resp.GetCustomer() == nil {
			t.Fatal("GetCustomer: customer is nil")
		}
		if resp.GetCustomer().GetStatus() != "suspended" {
			t.Errorf("GetCustomer: status = %q, want %q", resp.GetCustomer().GetStatus(), "suspended")
		}
		t.Log("step 6 - verified suspended")
	})
}
