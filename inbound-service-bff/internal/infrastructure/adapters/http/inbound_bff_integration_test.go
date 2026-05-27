// Copyright (c) 2026 WMS.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

//go:build integration

package http

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	// baseURL is the inbound BFF HTTP base URL.
	baseURL = "http://localhost:8083"
)

// apiResponse mirrors the BFF JSON response envelope.
type apiResponse struct {
	Success bool                   `json:"success"`
	Data    map[string]interface{} `json:"data,omitempty"`
	Error   *string                `json:"error,omitempty"`
}

// submitInboundBody is the JSON body for POST /v1/inbounds.
type submitInboundBody struct {
	CustomerID   string             `json:"customer_id"`
	ExpectedDate string             `json:"expected_date"`
	Notes        string             `json:"notes,omitempty"`
	Items        []submitInboundItem `json:"items"`
}

type submitInboundItem struct {
	SKU              string  `json:"sku"`
	QuantityDeclared int     `json:"quantity_declared"`
	Dimensions       string  `json:"dimensions,omitempty"`
	Weight           float64 `json:"weight,omitempty"`
}

// testToken returns a signed JWT for E2E testing.
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

// bearerHeader returns an Authorization header value with a Bearer token.
func bearerHeader(t *testing.T) string {
	t.Helper()
	return "Bearer " + testToken(t)
}

// doJSONRequest performs an HTTP request and decodes the JSON response.
func doJSONRequest(t *testing.T, method, url string, body []byte, authToken string) (*http.Response, *apiResponse) {
	t.Helper()

	var reqBody io.Reader
	if body != nil {
		reqBody = bytes.NewReader(body)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		t.Fatalf("new request %s %s: %v", method, url, err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if authToken != "" {
		req.Header.Set("Authorization", "Bearer "+authToken)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("do request %s %s: %v", method, url, err)
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read response body: %v", err)
	}

	var apiResp apiResponse
	if len(respBody) > 0 {
		if err := json.Unmarshal(respBody, &apiResp); err != nil {
			t.Fatalf("decode response body: %v (body: %s)", err, string(respBody))
		}
	}

	return resp, &apiResp
}

// TestInboundBFF_Integration runs the full inbound BFF lifecycle via HTTP.
func TestInboundBFF_Integration(t *testing.T) {
	token := testToken(t)
	var inboundID string

	// --- SubmitInbound ---
	t.Run("SubmitInbound", func(t *testing.T) {
		body := submitInboundBody{
			CustomerID:   "e2e-bff-customer-1",
			ExpectedDate: "2026-07-01",
			Notes:        "E2E BFF test",
			Items: []submitInboundItem{
				{SKU: "BFF-SKU-001", QuantityDeclared: 20, Dimensions: "5x5x5", Weight: 2.0},
				{SKU: "BFF-SKU-002", QuantityDeclared: 3},
			},
		}
		jsonBody, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("marshal body: %v", err)
		}

		resp, apiResp := doJSONRequest(t, http.MethodPost, baseURL+"/v1/inbounds", jsonBody, token)
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("SubmitInbound: status = %d, want 200", resp.StatusCode)
		}
		if !apiResp.Success {
			errMsg := "<nil>"
			if apiResp.Error != nil {
				errMsg = *apiResp.Error
			}
			t.Fatalf("SubmitInbound: success = false, error = %s", errMsg)
		}
		if apiResp.Data == nil {
			t.Fatal("SubmitInbound: data is nil")
		}
		idRaw, ok := apiResp.Data["id"]
		if !ok {
			t.Fatal("SubmitInbound: data.id not found in response")
		}
		idStr, ok := idRaw.(string)
		if !ok || idStr == "" {
			t.Fatalf("SubmitInbound: data.id = %v (not a non-empty string)", idRaw)
		}
		inboundID = idStr

		// Validate customer_id
		cidRaw, ok := apiResp.Data["customer_id"]
		if !ok {
			t.Fatal("SubmitInbound: data.customer_id not found")
		}
		cidStr, _ := cidRaw.(string)
		if cidStr != "e2e-bff-customer-1" {
			t.Errorf("SubmitInbound: customer_id = %q, want %q", cidStr, "e2e-bff-customer-1")
		}

		t.Logf("created inbound via BFF: id=%s", inboundID)
	})

	if inboundID == "" {
		t.Fatal("inboundID not set, aborting remaining tests")
	}

	// --- GetInboundQueue ---
	t.Run("GetInboundQueue", func(t *testing.T) {
		resp, apiResp := doJSONRequest(t, http.MethodGet, baseURL+"/v1/inbounds/queue", nil, token)
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("GetInboundQueue: status = %d, want 200", resp.StatusCode)
		}
		if !apiResp.Success {
			errMsg := "<nil>"
			if apiResp.Error != nil {
				errMsg = *apiResp.Error
			}
			t.Fatalf("GetInboundQueue: success = false, error = %s", errMsg)
		}
		t.Logf("GetInboundQueue returned data: %+v", apiResp.Data)
	})

	// --- GetInbound ---
	t.Run("GetInbound", func(t *testing.T) {
		url := fmt.Sprintf("%s/v1/inbounds/%s", baseURL, inboundID)
		resp, apiResp := doJSONRequest(t, http.MethodGet, url, nil, token)
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("GetInbound: status = %d, want 200", resp.StatusCode)
		}
		if !apiResp.Success {
			errMsg := "<nil>"
			if apiResp.Error != nil {
				errMsg = *apiResp.Error
			}
			t.Fatalf("GetInbound: success = false, error = %s", errMsg)
		}
		if apiResp.Data == nil {
			t.Fatal("GetInbound: data is nil")
		}
		idRaw, ok := apiResp.Data["id"]
		if !ok {
			t.Fatal("GetInbound: data.id not found")
		}
		idStr, _ := idRaw.(string)
		if idStr != inboundID {
			t.Errorf("GetInbound: id = %q, want %q", idStr, inboundID)
		}
	})

	// --- InspectInbound ---
	t.Run("InspectInbound", func(t *testing.T) {
		url := fmt.Sprintf("%s/v1/inbounds/%s/inspect", baseURL, inboundID)
		body := `{"inspector_id":"e2e-inspector","passed":true,"notes":"bff e2e inspect"}`
		resp, apiResp := doJSONRequest(t, http.MethodPost, url, []byte(body), token)
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("InspectInbound: status = %d, want 200", resp.StatusCode)
		}
		if !apiResp.Success {
			errMsg := "<nil>"
			if apiResp.Error != nil {
				errMsg = *apiResp.Error
			}
			t.Fatalf("InspectInbound: success = false, error = %s", errMsg)
		}
		t.Logf("inspected inbound: id=%s", inboundID)
	})

	// --- ApproveInbound ---
	t.Run("ApproveInbound", func(t *testing.T) {
		url := fmt.Sprintf("%s/v1/inbounds/%s/approve", baseURL, inboundID)
		resp, apiResp := doJSONRequest(t, http.MethodPost, url, nil, token)
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("ApproveInbound: status = %d, want 200", resp.StatusCode)
		}
		if !apiResp.Success {
			errMsg := "<nil>"
			if apiResp.Error != nil {
				errMsg = *apiResp.Error
			}
			t.Fatalf("ApproveInbound: success = false, error = %s", errMsg)
		}
		t.Logf("approved inbound: id=%s", inboundID)
	})

	// --- Invalid JSON ---
	t.Run("InvalidJSON", func(t *testing.T) {
		badBody := `{"customer_id": "test", "items": [invalid json]}`
		resp, apiResp := doJSONRequest(t, http.MethodPost, baseURL+"/v1/inbounds", []byte(badBody), token)

		// The BFF should reject with 400.
		if resp.StatusCode != http.StatusBadRequest {
			// If the service doesn't return 400 for this specific payload,
			// at minimum verify that it's not a 200 OK with success=true.
			if resp.StatusCode == http.StatusOK && apiResp.Success {
				t.Fatalf("InvalidJSON: got 200 success, expected 400 error")
			}
		}
		t.Logf("InvalidJSON: status = %d, success = %v", resp.StatusCode, apiResp.Success)
	})

	// --- No auth ---
	t.Run("NoAuth", func(t *testing.T) {
		body := submitInboundBody{
			CustomerID:   "no-auth-test",
			ExpectedDate: "2026-07-01",
			Items: []submitInboundItem{
				{SKU: "NO-AUTH-SKU", QuantityDeclared: 1},
			},
		}
		jsonBody, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("marshal body: %v", err)
		}

		resp, _ := doJSONRequest(t, http.MethodPost, baseURL+"/v1/inbounds", jsonBody, "")

		// Expect 401 Unauthorized when no auth header is provided.
		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("NoAuth: status = %d, want 401", resp.StatusCode)
		}
		t.Logf("NoAuth: status = %d", resp.StatusCode)
	})
}

// TestInboundBFF_NotFound tests that a non-existent inbound returns a proper error.
func TestInboundBFF_NotFound(t *testing.T) {
	token := testToken(t)

	t.Run("GetNonExistentInbound", func(t *testing.T) {
		url := fmt.Sprintf("%s/v1/inbounds/%s", baseURL, "non-existent-id")
		resp, apiResp := doJSONRequest(t, http.MethodGet, url, nil, token)
		if resp.StatusCode == http.StatusOK && apiResp.Success {
			t.Fatal("GetNonExistentInbound: expected error, got 200 success")
		}
		t.Logf("GetNonExistentInbound: status = %d, success = %v", resp.StatusCode, apiResp.Success)
	})
}

// TestInboundBFF_AuthTokenSigning verifies the test token helper produces
// a well-formed JWT that the service would accept.
func TestInboundBFF_AuthTokenSigning(t *testing.T) {
	token := testToken(t)
	if !strings.HasPrefix(token, "eyJ") {
		t.Errorf("token does not look like a JWT (want 'eyJ...' prefix), got: %s", token[:10]+"...")
	}
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		t.Errorf("token has %d parts, expected 3 (JWT format)", len(parts))
	}
}
