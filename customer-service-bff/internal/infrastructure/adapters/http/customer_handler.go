package httpadapter

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/beabys/wms/customer-service-bff/internal/api/v1"
	grpcdapter "github.com/beabys/wms/customer-service-bff/internal/infrastructure/adapters/grpc"
	"github.com/beabys/wms/customer-service-bff/internal/domain/model"
	"github.com/beabys/wms/customer-service-bff/internal/infrastructure/adapters/http/context"
	"github.com/beabys/wms/pkg/logger"
)

// RegisterCustomer registers a new customer (public endpoint).
func (s *Server) RegisterCustomer(w http.ResponseWriter, r *http.Request) {
	var req v1.RegisterCustomerJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	modelReq := &model.RegisterCustomerRequest{
		Token:       req.Token,
		CompanyName: req.CompanyName,
		Email:       string(req.Email),
		Password:    req.Password,
		Phone:       stringPtrValue(req.Phone),
		VatNumber:   stringPtrValue(req.VatNumber),
		Address:     stringPtrValue(req.Address),
		City:        stringPtrValue(req.City),
		PostalCode:  stringPtrValue(req.PostalCode),
		Country:     stringPtrValue(req.Country),
	}

	result, err := s.customerSvc.RegisterCustomer(r.Context(), modelReq)
	if err != nil {
		mapAndWriteError(s, w, err)
		return
	}

	s.logger.Info("customer registered",
		logger.LogField{Key: "company_name", Value: req.CompanyName},
	)
	s.writeCreated(w, result)
}

// ListCustomers lists customers with optional filters.
func (s *Server) ListCustomers(w http.ResponseWriter, r *http.Request, params v1.ListCustomersParams) {
	var (
		status         *string
		page, pageSize *int
	)
	if params.Status != nil && *params.Status != "" {
		status = params.Status
	}
	if params.Page != nil && *params.Page > 0 {
		page = params.Page
	}
	if params.PageSize != nil && *params.PageSize > 0 {
		pageSize = params.PageSize
	}

	result, err := s.customerSvc.ListCustomers(r.Context(), status, page, pageSize)
	if err != nil {
		mapAndWriteError(s, w, err)
		return
	}

	s.writeSuccess(w, result)
}

// GetMyCustomer retrieves the customer profile for the authenticated user.
func (s *Server) GetMyCustomer(w http.ResponseWriter, r *http.Request) {
	adminUserID := extractUserIDFromJWT(r)
	if adminUserID == "" {
		s.writeError(w, http.StatusUnauthorized, "unauthenticated")
		return
	}

	result, err := s.customerSvc.GetMyCustomer(r.Context(), adminUserID)
	if err != nil {
		mapAndWriteError(s, w, err)
		return
	}

	s.writeSuccess(w, result)
}

// GetCustomer retrieves a customer by ID.
func (s *Server) GetCustomer(w http.ResponseWriter, r *http.Request, id string) {
	result, err := s.customerSvc.GetCustomer(r.Context(), id)
	if err != nil {
		mapAndWriteError(s, w, err)
		return
	}

	s.writeSuccess(w, result)
}

// UpdateCustomer updates a customer's fields.
func (s *Server) UpdateCustomer(w http.ResponseWriter, r *http.Request, id string) {
	var req v1.UpdateCustomerJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	modelReq := &model.UpdateCustomerRequest{
		Phone:      stringPtrValue(req.Phone),
		Address:    stringPtrValue(req.Address),
		City:       stringPtrValue(req.City),
		PostalCode: stringPtrValue(req.PostalCode),
		Country:    stringPtrValue(req.Country),
	}

	result, err := s.customerSvc.UpdateCustomer(r.Context(), id, modelReq)
	if err != nil {
		mapAndWriteError(s, w, err)
		return
	}

	s.writeSuccess(w, result)
}

// ApproveCustomer approves a customer.
func (s *Server) ApproveCustomer(w http.ResponseWriter, r *http.Request, id string) {
	approvedBy := extractUserIDFromJWT(r)

	result, err := s.customerSvc.ApproveCustomer(r.Context(), id, approvedBy)
	if err != nil {
		mapAndWriteError(s, w, err)
		return
	}

	s.logger.Info("customer approved",
		logger.LogField{Key: "customer_id", Value: id},
	)
	s.writeSuccess(w, result)
}

// RejectCustomer rejects a customer.
func (s *Server) RejectCustomer(w http.ResponseWriter, r *http.Request, id string) {
	var reqBody struct {
		Reason string `json:"reason"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		reqBody.Reason = ""
	}

	err := s.customerSvc.RejectCustomer(r.Context(), id, reqBody.Reason)
	if err != nil {
		mapAndWriteError(s, w, err)
		return
	}

	s.logger.Info("customer rejected",
		logger.LogField{Key: "customer_id", Value: id},
	)
	s.writeSuccess(w, map[string]string{"message": "customer rejected"})
}

// SuspendCustomer suspends a customer.
func (s *Server) SuspendCustomer(w http.ResponseWriter, r *http.Request, id string) {
	var reqBody struct {
		Reason string `json:"reason"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		reqBody.Reason = ""
	}

	err := s.customerSvc.SuspendCustomer(r.Context(), id, reqBody.Reason)
	if err != nil {
		mapAndWriteError(s, w, err)
		return
	}

	s.logger.Info("customer suspended",
		logger.LogField{Key: "customer_id", Value: id},
	)
	s.writeSuccess(w, map[string]string{"message": "customer suspended"})
}

// RestoreCustomer restores a suspended customer.
func (s *Server) RestoreCustomer(w http.ResponseWriter, r *http.Request, id string) {
	restoredBy := extractUserIDFromJWT(r)
	result, err := s.customerSvc.RestoreCustomer(r.Context(), id, restoredBy)
	if err != nil {
		mapAndWriteError(s, w, err)
		return
	}

	s.logger.Info("customer restored",
		logger.LogField{Key: "customer_id", Value: id},
	)
	s.writeSuccess(w, result)
}

// AssignCompanyRole assigns a company role to a user.
func (s *Server) AssignCompanyRole(w http.ResponseWriter, r *http.Request) {
	var req v1.AssignCompanyRoleJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	modelReq := &model.AssignCompanyRoleRequest{
		CustomerID:  req.CustomerId,
		UserID:      req.UserId,
		RoleName:    req.RoleName,
		Permissions: slicePtrValue(req.Permissions),
		AssignedBy:  extractUserIDFromJWT(r),
	}

	err := s.customerSvc.AssignCompanyRole(r.Context(), modelReq)
	if err != nil {
		mapAndWriteError(s, w, err)
		return
	}

	s.writeSuccess(w, map[string]string{"message": "company role assigned"})
}

// ListCompanyRoles lists company roles.
func (s *Server) ListCompanyRoles(w http.ResponseWriter, r *http.Request, params v1.ListCompanyRolesParams) {
	if params.CustomerId == "" {
		s.writeError(w, http.StatusBadRequest, "customer_id is required")
		return
	}

	result, err := s.customerSvc.ListCompanyRoles(r.Context(), params.CustomerId)
	if err != nil {
		mapAndWriteError(s, w, err)
		return
	}

	s.writeSuccess(w, result)
}

// ListAuditLogs retrieves paginated audit logs for a customer.
func (s *Server) ListAuditLogs(w http.ResponseWriter, r *http.Request, id string, params v1.ListAuditLogsParams) {
	page := 1
	pageSize := 20
	if params.Page != nil && *params.Page > 0 {
		page = *params.Page
	}
	if params.PageSize != nil && *params.PageSize > 0 {
		pageSize = *params.PageSize
	}

	result, err := s.customerSvc.ListAuditLogs(r.Context(), id, page, pageSize)
	if err != nil {
		mapAndWriteError(s, w, err)
		return
	}

	s.writeSuccess(w, result)
}

// GetUserPermissions gets permissions for a user within a company.
func (s *Server) GetUserPermissions(w http.ResponseWriter, r *http.Request, params v1.GetUserPermissionsParams) {
	if params.UserId == "" || params.CustomerId == "" {
		s.writeError(w, http.StatusBadRequest, "user_id and customer_id are required")
		return
	}

	result, err := s.customerSvc.GetUserPermissions(r.Context(), params.UserId, params.CustomerId)
	if err != nil {
		mapAndWriteError(s, w, err)
		return
	}

	s.writeSuccess(w, result)
}

// ---- helpers ----

func mapAndWriteError(s *Server, w http.ResponseWriter, err error) {
	if err == nil {
		return
	}
	status := grpcdapter.HTTPStatusFromError(err)
	msg := err.Error()
	if errors.Is(err, grpcdapter.ErrInternal) {
		status = http.StatusInternalServerError
		msg = "internal server error"
	}
	s.writeError(w, status, msg)
}

// extractUserIDFromJWT returns the raw JWT token from the request context.
func extractUserIDFromJWT(r *http.Request) string {
	token, _ := r.Context().Value(httpctx.ContextKeyJWT).(string)
	return token
}

// stringPtrValue safely dereferences a string pointer, returning empty string if nil.
func stringPtrValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// slicePtrValue safely dereferences a slice pointer, returning nil if nil.
func slicePtrValue(s *[]string) []string {
	if s == nil {
		return nil
	}
	return *s
}
