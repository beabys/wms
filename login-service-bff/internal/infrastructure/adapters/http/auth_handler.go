package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"go.uber.org/zap"

	v1 "github.com/beabys/wms/login-service-bff/internal/api/v1"
)

// Postv1AuthLogin handles POST /v1/auth/login.
func (hs *HttpServer) Postv1AuthLogin(w http.ResponseWriter, r *http.Request) {
	var req v1.Postv1AuthLoginJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		hs.Logger.Warn("login: invalid request body", zap.Error(err))
		errorResponseJSON(w, http.StatusBadRequest, err)
		return
	}

	result, err := hs.GRPCClient.Login(r.Context(), string(req.Email), req.Password)
	if err != nil {
		hs.Logger.Warn("login: gRPC failed", zap.Error(err))
		errorResponseJSON(w, http.StatusUnauthorized, err)
		return
	}

	user := result.GetUser()
	successResponseJSON(w, map[string]interface{}{
		"access_token":  result.GetAccessToken(),
		"refresh_token": result.GetRefreshToken(),
		"user": map[string]interface{}{
			"id":          user.GetId(),
			"email":       user.GetEmail(),
			"role":        user.GetRole(),
			"customer_id": user.GetCustomerId(),
			"created_at":  user.GetCreatedAt(),
		},
	})
}

// Postv1AuthRegister handles POST /v1/auth/register.
func (hs *HttpServer) Postv1AuthRegister(w http.ResponseWriter, r *http.Request) {
	var req v1.Postv1AuthRegisterJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		hs.Logger.Warn("register: invalid request body", zap.Error(err))
		errorResponseJSON(w, http.StatusBadRequest, err)
		return
	}

	companyName := ""
	if req.CompanyName != nil {
		companyName = *req.CompanyName
	}

	result, err := hs.GRPCClient.Register(r.Context(), string(req.Email), req.Password, companyName)
	if err != nil {
		hs.Logger.Warn("register: gRPC failed", zap.Error(err))
		errorResponseJSON(w, http.StatusBadRequest, err)
		return
	}

	user := result.GetUser()
	successResponseJSON(w, map[string]interface{}{
		"user": map[string]interface{}{
			"id":          user.GetId(),
			"email":       user.GetEmail(),
			"role":        user.GetRole(),
			"customer_id": user.GetCustomerId(),
			"created_at":  user.GetCreatedAt(),
		},
	})
}

// Postv1AuthRefresh handles POST /v1/auth/refresh.
func (hs *HttpServer) Postv1AuthRefresh(w http.ResponseWriter, r *http.Request) {
	var req v1.Postv1AuthRefreshJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		hs.Logger.Warn("refresh: invalid request body", zap.Error(err))
		errorResponseJSON(w, http.StatusBadRequest, err)
		return
	}

	result, err := hs.GRPCClient.RefreshToken(r.Context(), req.RefreshToken)
	if err != nil {
		hs.Logger.Warn("refresh: gRPC failed", zap.Error(err))
		errorResponseJSON(w, http.StatusUnauthorized, err)
		return
	}

	successResponseJSON(w, map[string]interface{}{
		"access_token":  result.GetAccessToken(),
		"refresh_token": result.GetRefreshToken(),
	})
}

// Postv1AuthLogout handles POST /v1/auth/logout.
func (hs *HttpServer) Postv1AuthLogout(w http.ResponseWriter, r *http.Request) {
	var req v1.Postv1AuthLogoutJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		hs.Logger.Warn("logout: invalid request body", zap.Error(err))
		errorResponseJSON(w, http.StatusBadRequest, err)
		return
	}

	if err := hs.GRPCClient.DeleteRefreshToken(r.Context(), req.RefreshToken); err != nil {
		hs.Logger.Warn("logout: failed", zap.Error(err))
		errorResponseJSON(w, http.StatusInternalServerError, err)
		return
	}

	successResponseJSON(w, map[string]interface{}{})
}

// Getv1AuthMe handles GET /v1/auth/me.
func (hs *HttpServer) Getv1AuthMe(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")

	validateResp, err := hs.GRPCClient.ValidateToken(r.Context(), token)
	if err != nil {
		hs.Logger.Warn("me: token validation failed", zap.Error(err))
		errorResponseJSON(w, http.StatusUnauthorized, err)
		return
	}

	userResp, err := hs.GRPCClient.GetUser(r.Context(), validateResp.GetUserId(), token)
	if err != nil {
		hs.Logger.Warn("me: get user failed", zap.Error(err))
		errorResponseJSON(w, http.StatusInternalServerError, err)
		return
	}

	user := userResp.GetUser()
	successResponseJSON(w, map[string]interface{}{
		"user": map[string]interface{}{
			"id":          user.GetId(),
			"email":       user.GetEmail(),
			"role":        user.GetRole(),
			"customer_id": user.GetCustomerId(),
			"created_at":  user.GetCreatedAt(),
		},
	})
}
