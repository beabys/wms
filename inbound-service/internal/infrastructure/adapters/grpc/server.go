package inboundgrpc

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	inboundv1 "github.com/beabys/wms/proto/gen/go/inbound/v1"

	"github.com/beabys/wms/inbound-service/internal/application/inbound/command"
	"github.com/beabys/wms/inbound-service/internal/application/inbound/usecase"
	"github.com/beabys/wms/inbound-service/internal/domain/inbound/model"
)

// InboundServer implements the inbound.v1.InboundService gRPC server.
type InboundServer struct {
	inboundv1.UnimplementedInboundServiceServer
	svc usecase.InboundServiceHandler
}

// NewInboundServer creates a new InboundServer.
func NewInboundServer(svc usecase.InboundServiceHandler) *InboundServer {
	return &InboundServer{svc: svc}
}

// CreateInbound handles CreateInbound RPC.
func (s *InboundServer) CreateInbound(ctx context.Context, req *inboundv1.CreateInboundRequest) (*inboundv1.CreateInboundResponse, error) {
	items := make([]command.SubmitItem, len(req.Items))
	for i, pbItem := range req.Items {
		items[i] = command.SubmitItem{
			SKU:              pbItem.Sku,
			QuantityDeclared: pbItem.QuantityDeclared,
			QuantityReceived: pbItem.QuantityReceived,
			Dimensions:       pbItem.Dimensions,
			Weight:           pbItem.Weight,
		}
	}

	cmd := command.SubmitInboundCommand{
		CustomerID:   req.CustomerId,
		ExpectedDate: req.ExpectedDate,
		Notes:        req.Notes,
		Items:        items,
	}

	result, err := 	s.svc.SubmitInbound(ctx, cmd)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "create inbound: %v", err)
	}

	return &inboundv1.CreateInboundResponse{
		Inbound: domainToProto(result.Inbound),
	}, nil
}

// GetInbound handles GetInbound RPC.
func (s *InboundServer) GetInbound(ctx context.Context, req *inboundv1.GetInboundRequest) (*inboundv1.GetInboundResponse, error) {
	qry := command.GetInboundQuery{InboundID: req.Id}

	result, err := 	s.svc.GetInbound(ctx, qry)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "get inbound: %v", err)
	}

	return &inboundv1.GetInboundResponse{
		Inbound: domainToProto(result.Inbound),
	}, nil
}

// ListInbounds handles ListInbounds RPC.
func (s *InboundServer) ListInbounds(ctx context.Context, req *inboundv1.ListInboundsRequest) (*inboundv1.ListInboundsResponse, error) {
	qry := command.ListInboundsQuery{
		CustomerID: req.CustomerId,
		Status:     req.Status,
	}
	if req.Pagination != nil {
		qry.PageSize = req.Pagination.Limit
		qry.PageToken = fmt.Sprintf("%d", req.Pagination.Page)
	}

	result, err := 	s.svc.ListInbounds(ctx, qry)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list inbounds: %v", err)
	}

	protoInbounds := make([]*inboundv1.Inbound, len(result.Inbounds))
	for i, in := range result.Inbounds {
		protoInbounds[i] = domainToProto(in)
	}

	return &inboundv1.ListInboundsResponse{
		Inbounds: protoInbounds,
	}, nil
}

// ApproveInbound handles ApproveInbound RPC.
// If inspection data is provided AND the inbound is in submitted status,
// it first inspects then approves. If no inspection data and status is inspected,
// it just approves.
func (s *InboundServer) ApproveInbound(ctx context.Context, req *inboundv1.ApproveInboundRequest) (*inboundv1.ApproveInboundResponse, error) {
	// First inspect if inspection data is provided
	if req.Inspection != nil {
		inspectCmd := command.InspectInboundCommand{
			InboundID:   req.Id,
			InspectorID: req.Inspection.InspectorId,
			Notes:       req.Inspection.Notes,
			InspectedAt: time.Now(),
			Passed:      req.Inspection.Passed,
		}

		_, err := 	s.svc.InspectInbound(ctx, inspectCmd)
		if err != nil {
			return nil, status.Errorf(codes.FailedPrecondition, "inspect inbound: %v", err)
		}
	}

	// Then approve
	approveCmd := command.ApproveInboundCommand{InboundID: req.Id}
	result, err := 	s.svc.ApproveInbound(ctx, approveCmd)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "approve inbound: %v", err)
	}

	return &inboundv1.ApproveInboundResponse{
		Inbound: domainToProto(result.Inbound),
	}, nil
}

// FlagInbound handles FlagInbound RPC.
func (s *InboundServer) FlagInbound(ctx context.Context, req *inboundv1.FlagInboundRequest) (*inboundv1.FlagInboundResponse, error) {
	cmd := command.FlagInboundCommand{
		InboundID: req.Id,
		Reason:    req.Reason,
	}

	result, err := 	s.svc.FlagInbound(ctx, cmd)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "flag inbound: %v", err)
	}

	return &inboundv1.FlagInboundResponse{
		Inbound: domainToProto(result.Inbound),
	}, nil
}

// HoldInbound handles HoldInbound RPC.
func (s *InboundServer) HoldInbound(ctx context.Context, req *inboundv1.HoldInboundRequest) (*inboundv1.HoldInboundResponse, error) {
	cmd := command.HoldInboundCommand{
		InboundID: req.Id,
		Reason:    req.Reason,
	}

	result, err := 	s.svc.HoldInbound(ctx, cmd)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "hold inbound: %v", err)
	}

	return &inboundv1.HoldInboundResponse{
		Inbound: domainToProto(result.Inbound),
	}, nil
}

// ReleaseInbound handles ReleaseInbound RPC.
func (s *InboundServer) ReleaseInbound(ctx context.Context, req *inboundv1.ReleaseInboundRequest) (*inboundv1.ReleaseInboundResponse, error) {
	// Extract released_by from auth context
	releasedBy := "unknown"
	if claims, ok := ctx.Value("auth.claims").(map[string]any); ok {
		if sub, ok := claims["sub"].(string); ok {
			releasedBy = sub
		}
	}

	cmd := command.ReleaseInboundCommand{
		InboundID:  req.Id,
		ReleasedBy: releasedBy,
	}

	result, err := 	s.svc.ReleaseInbound(ctx, cmd)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "release inbound: %v", err)
	}

	return &inboundv1.ReleaseInboundResponse{
		Inbound: domainToProto(result.Inbound),
	}, nil
}

// --- mapping helpers ---

func domainToProto(in *model.Inbound) *inboundv1.Inbound {
	if in == nil {
		return nil
	}

	pb := &inboundv1.Inbound{
		Id:           in.ID,
		CustomerId:   in.CustomerID,
		Status:       string(in.Status),
		ExpectedDate: in.ExpectedDate,
		Notes:        in.Notes,
		CreatedAt:    in.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    in.UpdatedAt.Format(time.RFC3339),
	}

	for _, item := range in.Items {
		pb.Items = append(pb.Items, &inboundv1.InboundItem{
			Sku:              item.SKU,
			QuantityDeclared: item.QuantityDeclared,
			QuantityReceived: item.QuantityReceived,
			Dimensions:       item.Dimensions,
			Weight:           item.Weight,
		})
	}

	if in.Inspection != nil {
		pb.Inspection = &inboundv1.Inspection{
			InspectorId: in.Inspection.InspectorID,
			Notes:       in.Inspection.Notes,
			Passed:      in.Inspection.Passed,
		}
		// Photos not in proto Inspection; skip
	}

	if in.HoldRecord != nil {
		pb.Hold = &inboundv1.Hold{
			Reason:    in.HoldRecord.Reason,
			CreatedAt: in.HoldRecord.CreatedAt.Format(time.RFC3339),
		}
		if in.HoldRecord.ReleasedAt != nil {
			pb.Hold.ReleasedAt = in.HoldRecord.ReleasedAt.Format(time.RFC3339)
		}
	}

	return pb
}
