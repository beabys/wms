package grpc

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	commonv1 "github.com/beabys/wms/proto/gen/go/common/v1"
	inboundv1 "github.com/beabys/wms/proto/gen/go/inbound/v1"
)

// withToken creates a context with JWT in gRPC metadata.
func withToken(ctx context.Context, token string) context.Context {
	return metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+token)
}

// Client is the gRPC client adapter implementing inbound.InboundService.
type Client struct {
	conn   *grpc.ClientConn
	client inboundv1.InboundServiceClient
}

// NewClient creates a new gRPC client for the inbound service.
func NewClient(addr string) (*Client, error) {
	conn, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("dial grpc: %w", err)
	}

	return &Client{
		conn:   conn,
		client: inboundv1.NewInboundServiceClient(conn),
	}, nil
}

// Close closes the underlying connection.
func (c *Client) Close() error {
	return c.conn.Close()
}

// CreateInbound creates a new inbound.
func (c *Client) CreateInbound(ctx context.Context, customerID, expectedDate, notes string, items []*inboundv1.InboundItem, token string) (*inboundv1.Inbound, error) {
	ctx = withToken(ctx, token)
	resp, err := c.client.CreateInbound(ctx, &inboundv1.CreateInboundRequest{
		CustomerId:   customerID,
		ExpectedDate: expectedDate,
		Notes:        notes,
		Items:        items,
	})
	if err != nil {
		return nil, fmt.Errorf("create inbound: %w", err)
	}
	return resp.Inbound, nil
}

// GetInbound retrieves an inbound by ID.
func (c *Client) GetInbound(ctx context.Context, id string, token string) (*inboundv1.Inbound, error) {
	ctx = withToken(ctx, token)
	resp, err := c.client.GetInbound(ctx, &inboundv1.GetInboundRequest{Id: id})
	if err != nil {
		return nil, fmt.Errorf("get inbound: %w", err)
	}
	return resp.Inbound, nil
}

// ListInbounds lists inbounds with filters.
func (c *Client) ListInbounds(ctx context.Context, customerID, status string, pageSize int32, _ string, token string) ([]*inboundv1.Inbound, error) {
	ctx = withToken(ctx, token)
	resp, err := c.client.ListInbounds(ctx, &inboundv1.ListInboundsRequest{
		CustomerId: customerID,
		Status:     status,
		Pagination: &commonv1.Pagination{
			Limit: pageSize,
			Page:  1,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("list inbounds: %w", err)
	}
	return resp.Inbounds, nil
}

// ApproveInbound approves an inbound.
func (c *Client) ApproveInbound(ctx context.Context, id string, inspection *inboundv1.Inspection, token string) (*inboundv1.Inbound, error) {
	ctx = withToken(ctx, token)
	resp, err := c.client.ApproveInbound(ctx, &inboundv1.ApproveInboundRequest{
		Id:         id,
		Inspection: inspection,
	})
	if err != nil {
		return nil, fmt.Errorf("approve inbound: %w", err)
	}
	return resp.Inbound, nil
}

// FlagInbound flags an inbound.
func (c *Client) FlagInbound(ctx context.Context, id, reason string, token string) (*inboundv1.Inbound, error) {
	ctx = withToken(ctx, token)
	resp, err := c.client.FlagInbound(ctx, &inboundv1.FlagInboundRequest{
		Id:     id,
		Reason: reason,
	})
	if err != nil {
		return nil, fmt.Errorf("flag inbound: %w", err)
	}
	return resp.Inbound, nil
}

// HoldInbound places a hold on an inbound.
func (c *Client) HoldInbound(ctx context.Context, id, reason string, token string) (*inboundv1.Inbound, error) {
	ctx = withToken(ctx, token)
	resp, err := c.client.HoldInbound(ctx, &inboundv1.HoldInboundRequest{
		Id:     id,
		Reason: reason,
	})
	if err != nil {
		return nil, fmt.Errorf("hold inbound: %w", err)
	}
	return resp.Inbound, nil
}

// ReleaseInbound releases a held inbound.
func (c *Client) ReleaseInbound(ctx context.Context, id string, token string) (*inboundv1.Inbound, error) {
	ctx = withToken(ctx, token)
	resp, err := c.client.ReleaseInbound(ctx, &inboundv1.ReleaseInboundRequest{Id: id})
	if err != nil {
		return nil, fmt.Errorf("release inbound: %w", err)
	}
	return resp.Inbound, nil
}
