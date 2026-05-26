package http

import (
	"errors"
	"net/http"
	"strings"

	v1 "github.com/beabys/wms/customer-service-bff/internal/api/v1"
	"github.com/beabys/wms/customer-service-bff/internal/application"
	"go.uber.org/zap"
)

// extractToken extracts the Bearer token from the Authorization header.
func extractToken(r *http.Request) string {
	return strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
}

// InviteCustomer sends an invitation to a prospective customer.
func (hs *HttpServer) InviteCustomer(w http.ResponseWriter, r *http.Request) {
	body := v1.InviteCustomerJSONBody{}
	if err := decodeJSONBody(r, &body); err != nil {
		hs.Logger.Error("invite customer: invalid request body", zap.Error(err))
		errorResponseJSON(w, http.StatusBadRequest, errors.New("invalid request body"))
		return
	}

	token := extractToken(r)
	inviteToken, err := hs.CustomerClient.InviteCustomer(r.Context(), string(body.Email), token)
	if err != nil {
		hs.Logger.Error("invite customer", zap.Error(err))
		errorResponseJSON(w, http.StatusInternalServerError, err)
		return
	}

	successResponseJSON(w, map[string]interface{}{
		"invitation_token": inviteToken,
	})
}

// RegisterCustomer registers a customer via invitation token.
func (hs *HttpServer) RegisterCustomer(w http.ResponseWriter, r *http.Request) {
	body := v1.RegisterCustomerJSONBody{}
	if err := decodeJSONBody(r, &body); err != nil {
		hs.Logger.Error("register customer: invalid request body", zap.Error(err))
		errorResponseJSON(w, http.StatusBadRequest, errors.New("invalid request body"))
		return
	}

	req := application.CreateCustomerRequest{
		CompanyName: body.CompanyName,
		VATNumber:   body.VatNumber,
		Line1:       body.Line1,
		City:        body.City,
		Country:     body.Country,
	}
	if body.Line2 != nil {
		req.Line2 = *body.Line2
	}
	if body.PostalCode != nil {
		req.PostalCode = *body.PostalCode
	}
	if body.RateCardId != nil {
		req.RateCardID = *body.RateCardId
	}

	token := extractToken(r)
	customer, err := hs.CustomerClient.CreateCustomer(r.Context(), req, token)
	if err != nil {
		hs.Logger.Error("register customer", zap.Error(err))
		errorResponseJSON(w, http.StatusInternalServerError, err)
		return
	}

	successResponseJSON(w, map[string]interface{}{
		"customer": customer,
	})
}

// GetCustomerApprovalQueue lists pending customer approvals (admin).
func (hs *HttpServer) GetCustomerApprovalQueue(w http.ResponseWriter, r *http.Request, params v1.GetCustomerApprovalQueueParams) {
	page := 1
	pageSize := 20
	if params.Page != nil {
		page = *params.Page
	}
	if params.PageSize != nil {
		pageSize = *params.PageSize
	}

	token := extractToken(r)
	result, err := hs.CustomerClient.ListCustomers(r.Context(), "pending", int32(page), int32(pageSize), token)
	if err != nil {
		hs.Logger.Error("get approval queue", zap.Error(err))
		errorResponseJSON(w, http.StatusInternalServerError, err)
		return
	}

	successResponseJSON(w, map[string]interface{}{
		"customers": result.Customers,
		"total":     result.Total,
		"page":      result.Page,
		"page_size": result.PageSize,
	})
}

// ApproveCustomer approves a pending customer (admin).
func (hs *HttpServer) ApproveCustomer(w http.ResponseWriter, r *http.Request, id string) {
	token := extractToken(r)
	customer, err := hs.CustomerClient.ApproveCustomer(r.Context(), id, token)
	if err != nil {
		hs.Logger.Error("approve customer", zap.Error(err))
		errorResponseJSON(w, http.StatusInternalServerError, err)
		return
	}

	successResponseJSON(w, map[string]interface{}{
		"customer": customer,
	})
}

// SuspendCustomer suspends an active customer (admin).
func (hs *HttpServer) SuspendCustomer(w http.ResponseWriter, r *http.Request, id string) {
	body := v1.SuspendCustomerJSONBody{}
	if err := decodeJSONBody(r, &body); err != nil {
		// non-fatal: reason is optional
		body.Reason = nil
	}

	reason := ""
	if body.Reason != nil {
		reason = *body.Reason
	}

	token := extractToken(r)
	customer, err := hs.CustomerClient.SuspendCustomer(r.Context(), id, reason, token)
	if err != nil {
		hs.Logger.Error("suspend customer", zap.Error(err))
		errorResponseJSON(w, http.StatusInternalServerError, err)
		return
	}

	successResponseJSON(w, map[string]interface{}{
		"customer": customer,
	})
}

// GetCustomer retrieves a customer by ID.
func (hs *HttpServer) GetCustomer(w http.ResponseWriter, r *http.Request, id string) {
	token := extractToken(r)
	customer, err := hs.CustomerClient.GetCustomer(r.Context(), id, token)
	if err != nil {
		hs.Logger.Error("get customer", zap.Error(err))
		errorResponseJSON(w, http.StatusNotFound, errors.New("customer not found"))
		return
	}

	successResponseJSON(w, map[string]interface{}{
		"customer": customer,
	})
}

// GetCustomerDashboard returns aggregated customer dashboard data.
func (hs *HttpServer) GetCustomerDashboard(w http.ResponseWriter, r *http.Request, id string) {
	token := extractToken(r)
	customer, err := hs.CustomerClient.GetCustomer(r.Context(), id, token)
	if err != nil {
		hs.Logger.Error("get customer dashboard: get customer", zap.Error(err))
		errorResponseJSON(w, http.StatusNotFound, errors.New("customer not found"))
		return
	}

	// Dashboard aggregation: includes placeholder values for ledger data.
	// In the future, this can call LedgerClient.GetBalance for real balance.
	dashboard := map[string]interface{}{
		"customer":        customer,
		"total_shipments": 0,
		"active_services": 0,
		"balance":         0,
	}

	successResponseJSON(w, dashboard)
}

// decodeJSONBody decodes a JSON request body into the provided target.
func decodeJSONBody(r *http.Request, target interface{}) error {
	defer r.Body.Close()
	return decodeJSON(r.Body, target)
}
