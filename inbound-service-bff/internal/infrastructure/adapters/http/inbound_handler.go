package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"go.uber.org/zap"

	v1 "github.com/beabys/wms/inbound-service-bff/internal/api/v1"
	inboundv1 "github.com/beabys/wms/proto/gen/go/inbound/v1"
)

// extractToken extracts the Bearer token from the Authorization header.
func extractToken(r *http.Request) string {
	return strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
}

// SubmitInbound handles POST /v1/inbounds.
func (hs *HttpServer) SubmitInbound(w http.ResponseWriter, r *http.Request) {
	var req v1.SubmitInboundJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		hs.Logger.Error("decode submit inbound body", zap.Error(err))
		errorResponseJSON(w, http.StatusBadRequest, err)
		return
	}

	items := make([]*inboundv1.InboundItem, len(req.Items))
	for i, item := range req.Items {
		it := &inboundv1.InboundItem{
			Sku:              item.SKU,
			QuantityDeclared: int32(item.QuantityDeclared),
		}
		if item.QuantityReceived != nil {
			it.QuantityReceived = int32(*item.QuantityReceived)
		}
		if item.Dimensions != nil {
			it.Dimensions = *item.Dimensions
		}
		if item.Weight != nil {
			it.Weight = *item.Weight
		}
		items[i] = it
	}

	notes := ""
	if req.Notes != nil {
		notes = *req.Notes
	}

	token := extractToken(r)
	inbound, err := hs.InboundClient.CreateInbound(r.Context(), req.CustomerId, req.ExpectedDate.String(), notes, items, token)
	if err != nil {
		hs.Logger.Error("create inbound", zap.Error(err))
		errorResponseJSON(w, http.StatusInternalServerError, err)
		return
	}

	successResponseJSON(w, inbound)
}

// GetInboundQueue handles GET /v1/inbounds/queue.
func (hs *HttpServer) GetInboundQueue(w http.ResponseWriter, r *http.Request, params v1.GetInboundQueueParams) {
	pageSize := int32(0)
	if params.PageSize != nil {
		pageSize = int32(*params.PageSize)
	}

	status := ""
	if params.Status != nil {
		status = *params.Status
	}

	token := extractToken(r)
	inbounds, err := hs.InboundClient.ListInbounds(r.Context(), "", status, pageSize, "", token)
	if err != nil {
		hs.Logger.Error("list inbounds queue", zap.Error(err))
		errorResponseJSON(w, http.StatusInternalServerError, err)
		return
	}

	successResponseJSON(w, inbounds)
}

// GetInbound handles GET /v1/inbounds/{id}.
func (hs *HttpServer) GetInbound(w http.ResponseWriter, r *http.Request, id string) {
	token := extractToken(r)
	inbound, err := hs.InboundClient.GetInbound(r.Context(), id, token)
	if err != nil {
		hs.Logger.Error("get inbound", zap.Error(err), zap.String("id", id))
		errorResponseJSON(w, http.StatusNotFound, err)
		return
	}

	successResponseJSON(w, inbound)
}

// InspectInbound handles POST /v1/inbounds/{id}/inspect.
func (hs *HttpServer) InspectInbound(w http.ResponseWriter, r *http.Request, id string) {
	var req v1.InspectInboundJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		hs.Logger.Error("decode inspect inbound body", zap.Error(err))
		errorResponseJSON(w, http.StatusBadRequest, err)
		return
	}

	inspection := &inboundv1.Inspection{
		InspectorId: req.InspectorId,
	}
	if req.Notes != nil {
		inspection.Notes = *req.Notes
	}
	if req.Passed != nil {
		inspection.Passed = *req.Passed
	}

	token := extractToken(r)
	inbound, err := hs.InboundClient.ApproveInbound(r.Context(), id, inspection, token)
	if err != nil {
		hs.Logger.Error("inspect inbound", zap.Error(err), zap.String("id", id))
		errorResponseJSON(w, http.StatusInternalServerError, err)
		return
	}

	successResponseJSON(w, inbound)
}

// ApproveInbound handles POST /v1/inbounds/{id}/approve.
func (hs *HttpServer) ApproveInbound(w http.ResponseWriter, r *http.Request, id string) {
	token := extractToken(r)
	inbound, err := hs.InboundClient.ApproveInbound(r.Context(), id, nil, token)
	if err != nil {
		hs.Logger.Error("approve inbound", zap.Error(err), zap.String("id", id))
		errorResponseJSON(w, http.StatusInternalServerError, err)
		return
	}

	successResponseJSON(w, inbound)
}

// FlagInbound handles POST /v1/inbounds/{id}/flag.
func (hs *HttpServer) FlagInbound(w http.ResponseWriter, r *http.Request, id string) {
	var req v1.FlagInboundJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		hs.Logger.Error("decode flag inbound body", zap.Error(err))
		errorResponseJSON(w, http.StatusBadRequest, err)
		return
	}

	token := extractToken(r)
	inbound, err := hs.InboundClient.FlagInbound(r.Context(), id, req.Reason, token)
	if err != nil {
		hs.Logger.Error("flag inbound", zap.Error(err), zap.String("id", id))
		errorResponseJSON(w, http.StatusInternalServerError, err)
		return
	}

	successResponseJSON(w, inbound)
}

// HoldInbound handles POST /v1/inbounds/{id}/hold.
func (hs *HttpServer) HoldInbound(w http.ResponseWriter, r *http.Request, id string) {
	var req v1.HoldInboundJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		hs.Logger.Error("decode hold inbound body", zap.Error(err))
		errorResponseJSON(w, http.StatusBadRequest, err)
		return
	}

	token := extractToken(r)
	inbound, err := hs.InboundClient.HoldInbound(r.Context(), id, req.Reason, token)
	if err != nil {
		hs.Logger.Error("hold inbound", zap.Error(err), zap.String("id", id))
		errorResponseJSON(w, http.StatusInternalServerError, err)
		return
	}

	successResponseJSON(w, inbound)
}

// ReleaseInbound handles POST /v1/inbounds/{id}/release.
func (hs *HttpServer) ReleaseInbound(w http.ResponseWriter, r *http.Request, id string) {
	token := extractToken(r)
	inbound, err := hs.InboundClient.ReleaseInbound(r.Context(), id, token)
	if err != nil {
		hs.Logger.Error("release inbound", zap.Error(err), zap.String("id", id))
		errorResponseJSON(w, http.StatusInternalServerError, err)
		return
	}

	successResponseJSON(w, inbound)
}

// GetCustomerInbounds handles GET /v1/inbounds/customer/{id}.
func (hs *HttpServer) GetCustomerInbounds(w http.ResponseWriter, r *http.Request, id string, params v1.GetCustomerInboundsParams) {
	pageSize := int32(0)
	if params.PageSize != nil {
		pageSize = int32(*params.PageSize)
	}

	token := extractToken(r)
	inbounds, err := hs.InboundClient.ListInbounds(r.Context(), id, "", pageSize, "", token)
	if err != nil {
		hs.Logger.Error("list customer inbounds", zap.Error(err), zap.String("customer_id", id))
		errorResponseJSON(w, http.StatusInternalServerError, err)
		return
	}

	successResponseJSON(w, inbounds)
}
