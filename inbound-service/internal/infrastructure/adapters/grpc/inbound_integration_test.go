// Copyright (c) 2026 WMS.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

//go:build integration

package inboundgrpc

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
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	inboundv1 "github.com/beabys/wms/proto/gen/go/inbound/v1"
)

const (
	addr = "localhost:50003"
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
		"sub": "e2e-test-user",
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

// --- Test 1: Full happy path ---

func TestInboundIntegration_FullHappyPath(t *testing.T) {
	conn := dial(t)
	t.Cleanup(func() { conn.Close() })
	client := inboundv1.NewInboundServiceClient(conn)

	var inboundID string

	t.Run("CreateInbound", func(t *testing.T) {
		req := &inboundv1.CreateInboundRequest{
			CustomerId:   "e2e-happy-customer",
			ExpectedDate: "2026-06-15",
			Notes:        "E2E happy path test",
			Items: []*inboundv1.InboundItem{
				{Sku: "SKU-HAPPY-001", QuantityDeclared: 10},
				{Sku: "SKU-HAPPY-002", QuantityDeclared: 5},
			},
		}
		resp, err := client.CreateInbound(authContext(t), req)
		if err != nil {
			t.Fatalf("CreateInbound: %v", err)
		}
		if resp.Inbound == nil {
			t.Fatal("CreateInbound: response inbound is nil")
		}
		if resp.Inbound.Id == "" {
			t.Fatal("CreateInbound: inbound id is empty")
		}
		if resp.Inbound.CustomerId != "e2e-happy-customer" {
			t.Errorf("CreateInbound: customer_id = %q, want %q", resp.Inbound.CustomerId, "e2e-happy-customer")
		}
		if resp.Inbound.Status != "submitted" {
			t.Errorf("CreateInbound: status = %q, want %q", resp.Inbound.Status, "submitted")
		}
		if len(resp.Inbound.Items) != 2 {
			t.Fatalf("CreateInbound: got %d items, want 2", len(resp.Inbound.Items))
		}
		inboundID = resp.Inbound.Id
		t.Logf("created inbound: id=%s", inboundID)
	})

	if inboundID == "" {
		t.Fatal("inboundID not set, aborting")
	}

	t.Run("GetSubmittedInbound", func(t *testing.T) {
		resp, err := client.GetInbound(authContext(t), &inboundv1.GetInboundRequest{Id: inboundID})
		if err != nil {
			t.Fatalf("GetInbound: %v", err)
		}
		if resp.Inbound == nil {
			t.Fatal("GetInbound: response inbound is nil")
		}
		if resp.Inbound.Id != inboundID {
			t.Errorf("GetInbound: id = %q, want %q", resp.Inbound.Id, inboundID)
		}
		if resp.Inbound.Status != "submitted" {
			t.Errorf("GetInbound: status = %q, want %q", resp.Inbound.Status, "submitted")
		}
	})

	t.Run("ListSubmittedInbounds", func(t *testing.T) {
		req := &inboundv1.ListInboundsRequest{
			Status: "submitted",
		}
		resp, err := client.ListInbounds(authContext(t), req)
		if err != nil {
			t.Fatalf("ListInbounds: %v", err)
		}
		if len(resp.Inbounds) == 0 {
			t.Fatal("ListInbounds: empty results, expected at least 1 inbound")
		}
		found := false
		for _, in := range resp.Inbounds {
			if in.Id == inboundID {
				found = true
				if in.Status != "submitted" {
					t.Errorf("ListInbounds: inbound %s status = %q, want %q", in.Id, in.Status, "submitted")
				}
				break
			}
		}
		if !found {
			t.Errorf("ListInbounds: inbound %s not found in results", inboundID)
		}
	})

	t.Run("InspectInbound", func(t *testing.T) {
		req := &inboundv1.InspectInboundRequest{
			Id: inboundID,
			Inspection: &inboundv1.Inspection{
				InspectorId: "e2e-inspector",
				Notes:       "All items look good",
				Passed:      true,
			},
		}
		resp, err := client.InspectInbound(authContext(t), req)
		if err != nil {
			t.Fatalf("InspectInbound: %v", err)
		}
		if resp.Inbound == nil {
			t.Fatal("InspectInbound: response inbound is nil")
		}
		if resp.Inbound.Status != "inspected" {
			t.Errorf("InspectInbound: status = %q, want %q", resp.Inbound.Status, "inspected")
		}
		if resp.Inbound.Inspection == nil {
			t.Error("InspectInbound: inspection data is nil, expected inspection recorded")
		} else if resp.Inbound.Inspection.InspectorId != "e2e-inspector" {
			t.Errorf("InspectInbound: inspector_id = %q, want %q",
				resp.Inbound.Inspection.InspectorId, "e2e-inspector")
		}
	})

	t.Run("ApproveInbound", func(t *testing.T) {
		req := &inboundv1.ApproveInboundRequest{
			Id: inboundID,
		}
		resp, err := client.ApproveInbound(authContext(t), req)
		if err != nil {
			t.Fatalf("ApproveInbound: %v", err)
		}
		if resp.Inbound == nil {
			t.Fatal("ApproveInbound: response inbound is nil")
		}
		if resp.Inbound.Status != "approved" {
			t.Errorf("ApproveInbound: status = %q, want %q", resp.Inbound.Status, "approved")
		}
	})

	t.Run("GetApprovedInbound", func(t *testing.T) {
		resp, err := client.GetInbound(authContext(t), &inboundv1.GetInboundRequest{Id: inboundID})
		if err != nil {
			t.Fatalf("GetInbound after approve: %v", err)
		}
		if resp.Inbound == nil {
			t.Fatal("GetInbound after approve: response inbound is nil")
		}
		if resp.Inbound.Status != "approved" {
			t.Errorf("GetInbound after approve: status = %q, want %q", resp.Inbound.Status, "approved")
		}
		if resp.Inbound.Inspection == nil {
			t.Error("GetInbound after approve: inspection data is nil")
		}
	})
}

// --- Test 2: Flag flow ---
// State machine: submitted → inspected → flagged.
// FlagInbound on submitted inbound should fail (must inspect first).
// This test verifies the state machine correctly rejects the transition.

func TestInboundIntegration_FlagFlow(t *testing.T) {
	conn := dial(t)
	t.Cleanup(func() { conn.Close() })
	client := inboundv1.NewInboundServiceClient(conn)

	var inboundID string

	t.Run("CreateInboundForFlag", func(t *testing.T) {
		req := &inboundv1.CreateInboundRequest{
			CustomerId:   "e2e-flag-customer",
			ExpectedDate: "2026-06-20",
			Items: []*inboundv1.InboundItem{
				{Sku: "SKU-FLAG-001", QuantityDeclared: 3},
			},
		}
		resp, err := client.CreateInbound(authContext(t), req)
		if err != nil {
			t.Fatalf("CreateInbound: %v", err)
		}
		if resp.Inbound == nil || resp.Inbound.Id == "" {
			t.Fatal("CreateInbound: no id returned")
		}
		inboundID = resp.Inbound.Id
		t.Logf("created inbound: id=%s", inboundID)
	})

	if inboundID == "" {
		t.Fatal("inboundID not set, aborting")
	}

	t.Run("FlagSubmittedInboundRejected", func(t *testing.T) {
		// Flagging a submitted inbound should fail because the state machine
		// requires submitted → inspected → flagged.
		req := &inboundv1.FlagInboundRequest{
			Id:     inboundID,
			Reason: "Items damaged",
		}
		_, err := client.FlagInbound(authContext(t), req)
		if err == nil {
			t.Fatal("FlagInbound on submitted: expected error, got nil")
		}
		st, ok := status.FromError(err)
		if !ok {
			t.Fatalf("error is not a gRPC status: %v", err)
		}
		if st.Code() != codes.FailedPrecondition {
			t.Errorf("status code = %v, want %v", st.Code(), codes.FailedPrecondition)
		}
		if !strings.Contains(st.Message(), "invalid status transition") {
			t.Errorf("unexpected error message: %v", st.Message())
		}
		t.Logf("FlagSubmittedInboundRejected: got expected error: %v", err)
	})

	t.Run("StatusStillSubmitted", func(t *testing.T) {
		// Verify the inbound is still in submitted state
		resp, err := client.GetInbound(authContext(t), &inboundv1.GetInboundRequest{Id: inboundID})
		if err != nil {
			t.Fatalf("GetInbound: %v", err)
		}
		if resp.Inbound == nil {
			t.Fatal("GetInbound: response inbound is nil")
		}
		if resp.Inbound.Status != "submitted" {
			t.Errorf("status = %q, want %q (unchanged)", resp.Inbound.Status, "submitted")
		}
	})
}

// --- Test 3: Hold → Release flow ---
// Hold on submitted inbound should fail (must inspect first).
// Release without hold should also fail.
// This test verifies the state machine enforcement.

func TestInboundIntegration_HoldReleaseFlow(t *testing.T) {
	conn := dial(t)
	t.Cleanup(func() { conn.Close() })
	client := inboundv1.NewInboundServiceClient(conn)

	var inboundID string

	t.Run("CreateInboundForHold", func(t *testing.T) {
		req := &inboundv1.CreateInboundRequest{
			CustomerId:   "e2e-hold-customer",
			ExpectedDate: "2026-06-25",
			Notes:        "E2E hold test",
			Items: []*inboundv1.InboundItem{
				{Sku: "SKU-HOLD-001", QuantityDeclared: 7},
			},
		}
		resp, err := client.CreateInbound(authContext(t), req)
		if err != nil {
			t.Fatalf("CreateInbound: %v", err)
		}
		if resp.Inbound == nil || resp.Inbound.Id == "" {
			t.Fatal("CreateInbound: no id returned")
		}
		inboundID = resp.Inbound.Id
		t.Logf("created inbound: id=%s", inboundID)
	})

	if inboundID == "" {
		t.Fatal("inboundID not set, aborting")
	}

	t.Run("HoldSubmittedRejected", func(t *testing.T) {
		// Hold requires inspected status first
		req := &inboundv1.HoldInboundRequest{
			Id:     inboundID,
			Reason: "Awaiting customs",
		}
		_, err := client.HoldInbound(authContext(t), req)
		if err == nil {
			t.Fatal("HoldInbound on submitted: expected error, got nil")
		}
		st, ok := status.FromError(err)
		if !ok {
			t.Fatalf("error is not a gRPC status: %v", err)
		}
		if st.Code() != codes.FailedPrecondition {
			t.Errorf("status code = %v, want %v", st.Code(), codes.FailedPrecondition)
		}
		t.Logf("HoldSubmittedRejected: got expected error: %v", err)
	})

	t.Run("ReleaseWithoutHoldRejected", func(t *testing.T) {
		req := &inboundv1.ReleaseInboundRequest{
			Id: inboundID,
		}
		_, err := client.ReleaseInbound(authContext(t), req)
		if err == nil {
			t.Fatal("ReleaseInbound without hold: expected error, got nil")
		}
		st, ok := status.FromError(err)
		if !ok {
			t.Fatalf("error is not a gRPC status: %v", err)
		}
		if st.Code() != codes.FailedPrecondition {
			t.Errorf("status code = %v, want %v", st.Code(), codes.FailedPrecondition)
		}
		t.Logf("ReleaseWithoutHoldRejected: got expected error: %v", err)
	})
}

// --- Test 4: Validation errors ---

func TestInboundIntegration_ValidationErrors(t *testing.T) {
	conn := dial(t)
	t.Cleanup(func() { conn.Close() })
	client := inboundv1.NewInboundServiceClient(conn)

	t.Run("EmptyCustomerID", func(t *testing.T) {
		req := &inboundv1.CreateInboundRequest{
			CustomerId:   "",
			ExpectedDate: "2026-06-15",
			Items: []*inboundv1.InboundItem{
				{Sku: "SKU-VALIDATION", QuantityDeclared: 1},
			},
		}
		_, err := client.CreateInbound(authContext(t), req)
		if err == nil {
			t.Fatal("expected error for empty customer_id, got nil")
		}
		st, _ := status.FromError(err)
		t.Logf("EmptyCustomerID: code=%v msg=%v", st.Code(), st.Message())
	})

	t.Run("EmptyItems", func(t *testing.T) {
		req := &inboundv1.CreateInboundRequest{
			CustomerId:   "e2e-validation-customer",
			ExpectedDate: "2026-06-15",
			Items:        []*inboundv1.InboundItem{},
		}
		_, err := client.CreateInbound(authContext(t), req)
		if err == nil {
			t.Fatal("expected error for empty items, got nil")
		}
		st, _ := status.FromError(err)
		t.Logf("EmptyItems: code=%v msg=%v", st.Code(), st.Message())
	})

	t.Run("ZeroQuantityDeclared", func(t *testing.T) {
		req := &inboundv1.CreateInboundRequest{
			CustomerId:   "e2e-validation-customer",
			ExpectedDate: "2026-06-15",
			Items: []*inboundv1.InboundItem{
				{Sku: "SKU-VALIDATION", QuantityDeclared: 0},
			},
		}
		_, err := client.CreateInbound(authContext(t), req)
		if err == nil {
			t.Fatal("expected error for zero quantity_declared, got nil")
		}
		st, _ := status.FromError(err)
		t.Logf("ZeroQuantityDeclared: code=%v msg=%v", st.Code(), st.Message())
	})
}

// --- Test 5: Auth required ---

func TestInboundIntegration_AuthRequired(t *testing.T) {
	conn := dial(t)
	t.Cleanup(func() { conn.Close() })
	client := inboundv1.NewInboundServiceClient(conn)

	t.Run("NoAuthMetadata", func(t *testing.T) {
		req := &inboundv1.CreateInboundRequest{
			CustomerId:   "e2e-no-auth",
			ExpectedDate: "2026-06-15",
			Items: []*inboundv1.InboundItem{
				{Sku: "SKU-NO-AUTH", QuantityDeclared: 1},
			},
		}
		_, err := client.CreateInbound(context.Background(), req)
		if err == nil {
			t.Fatal("expected Unauthenticated error, got nil")
		}
		st, ok := status.FromError(err)
		if !ok {
			t.Fatalf("error is not a gRPC status: %v", err)
		}
		if st.Code() != codes.Unauthenticated {
			t.Errorf("status code = %v, want %v", st.Code(), codes.Unauthenticated)
		}
		t.Logf("NoAuthMetadata: got expected error: %v", err)
	})
}

// --- Test 6: Invalid transitions ---

func TestInboundIntegration_InvalidTransitions(t *testing.T) {
	conn := dial(t)
	t.Cleanup(func() { conn.Close() })
	client := inboundv1.NewInboundServiceClient(conn)

	var inboundID string

	t.Run("CreateInboundForTransitions", func(t *testing.T) {
		req := &inboundv1.CreateInboundRequest{
			CustomerId:   "e2e-transition-customer",
			ExpectedDate: "2026-06-15",
			Items: []*inboundv1.InboundItem{
				{Sku: "SKU-TRANSITION", QuantityDeclared: 1},
			},
		}
		resp, err := client.CreateInbound(authContext(t), req)
		if err != nil {
			t.Fatalf("CreateInbound: %v", err)
		}
		if resp.Inbound == nil || resp.Inbound.Id == "" {
			t.Fatal("CreateInbound: no id returned")
		}
		inboundID = resp.Inbound.Id
		t.Logf("created inbound: id=%s", inboundID)
	})

	if inboundID == "" {
		t.Fatal("inboundID not set, aborting")
	}

	t.Run("ApproveWithoutInspection", func(t *testing.T) {
		// submitted → approved is invalid; must inspect first via InspectInbound
		req := &inboundv1.ApproveInboundRequest{
			Id: inboundID,
		}
		_, err := client.ApproveInbound(authContext(t), req)
		if err == nil {
			t.Fatal("expected error for approve without inspection, got nil")
		}
		st, ok := status.FromError(err)
		if !ok {
			t.Fatalf("error is not a gRPC status: %v", err)
		}
		if st.Code() != codes.FailedPrecondition {
			t.Errorf("status code = %v, want %v", st.Code(), codes.FailedPrecondition)
		}
		t.Logf("ApproveWithoutInspection: got expected error: %v", err)
	})

	t.Run("InspectInbound", func(t *testing.T) {
		req := &inboundv1.InspectInboundRequest{
			Id: inboundID,
			Inspection: &inboundv1.Inspection{
				InspectorId: "e2e-inspector",
				Notes:       "OK",
				Passed:      true,
			},
		}
		resp, err := client.InspectInbound(authContext(t), req)
		if err != nil {
			t.Fatalf("InspectInbound: %v", err)
		}
		if resp.Inbound == nil || resp.Inbound.Status != "inspected" {
			t.Fatalf("InspectInbound: status = %q, want %q",
				resp.Inbound.Status, "inspected")
		}
	})

	t.Run("ApproveInbound", func(t *testing.T) {
		req := &inboundv1.ApproveInboundRequest{
			Id: inboundID,
		}
		resp, err := client.ApproveInbound(authContext(t), req)
		if err != nil {
			t.Fatalf("ApproveInbound: %v", err)
		}
		if resp.Inbound == nil || resp.Inbound.Status != "approved" {
			t.Fatalf("ApproveInbound: status = %q, want %q",
				resp.Inbound.Status, "approved")
		}
	})

	t.Run("ApproveTwice", func(t *testing.T) {
		req := &inboundv1.ApproveInboundRequest{
			Id: inboundID,
		}
		_, err := client.ApproveInbound(authContext(t), req)
		if err == nil {
			t.Fatal("expected error for double approve, got nil")
		}
		st, ok := status.FromError(err)
		if !ok {
			t.Fatalf("error is not a gRPC status: %v", err)
		}
		if st.Code() != codes.FailedPrecondition {
			t.Errorf("status code = %v, want %v", st.Code(), codes.FailedPrecondition)
		}
		t.Logf("ApproveTwice: got expected error: %v", err)
	})

	t.Run("ReleaseWithoutHold", func(t *testing.T) {
		createReq := &inboundv1.CreateInboundRequest{
			CustomerId:   "e2e-release-no-hold",
			ExpectedDate: "2026-06-15",
			Items: []*inboundv1.InboundItem{
				{Sku: "SKU-RELEASE-NO-HOLD", QuantityDeclared: 1},
			},
		}
		createResp, err := client.CreateInbound(authContext(t), createReq)
		if err != nil {
			t.Fatalf("setup CreateInbound: %v", err)
		}
		if createResp.Inbound == nil || createResp.Inbound.Id == "" {
			t.Fatal("setup: no id returned")
		}
		testID := createResp.Inbound.Id

		req := &inboundv1.ReleaseInboundRequest{Id: testID}
		_, err = client.ReleaseInbound(authContext(t), req)
		if err == nil {
			t.Fatal("expected error for release without hold, got nil")
		}
		st, ok := status.FromError(err)
		if !ok {
			t.Fatalf("error is not a gRPC status: %v", err)
		}
		if st.Code() != codes.FailedPrecondition {
			t.Errorf("status code = %v, want %v", st.Code(), codes.FailedPrecondition)
		}
		t.Logf("ReleaseWithoutHold: got expected error: %v", err)
	})
}
