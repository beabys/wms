package inboundgrpc

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	inboundv1 "github.com/beabys/wms/proto/gen/go/inbound/v1"

	"github.com/beabys/wms/inbound-service/internal/application/inbound/usecase"
	"github.com/beabys/wms/inbound-service/internal/domain/inbound/model"
)

// mockGRPCRepo implements repository.InboundRepository.
type mockGRPCRepo struct {
	inbounds map[string]*model.Inbound
}

func newMockGRPCRepo() *mockGRPCRepo {
	return &mockGRPCRepo{inbounds: make(map[string]*model.Inbound)}
}

func (m *mockGRPCRepo) GetByID(_ context.Context, id string) (*model.Inbound, error) {
	in, ok := m.inbounds[id]
	if !ok {
		return nil, assert.AnError
	}
	return in, nil
}

func (m *mockGRPCRepo) Save(_ context.Context, inbound *model.Inbound) error {
	m.inbounds[inbound.ID] = inbound
	return nil
}

func (m *mockGRPCRepo) UpdateStatus(_ context.Context, inbound *model.Inbound) error {
	m.inbounds[inbound.ID] = inbound
	return nil
}

func (m *mockGRPCRepo) List(_ context.Context, customerID, status string, pageSize int32, pageToken string) ([]*model.Inbound, string, error) {
	var result []*model.Inbound
	for _, in := range m.inbounds {
		if customerID != "" && in.CustomerID != customerID {
			continue
		}
		if status != "" && string(in.Status) != status {
			continue
		}
		result = append(result, in)
	}
	return result, "", nil
}

type mockGRPCEvents struct{}

func (m *mockGRPCEvents) PublishSubmitted(_ context.Context, _ *model.Inbound) error { return nil }
func (m *mockGRPCEvents) PublishApproved(_ context.Context, _ *model.Inbound) error  { return nil }
func (m *mockGRPCEvents) PublishFlagged(_ context.Context, _ *model.Inbound) error   { return nil }

func TestInboundServer_CreateAndGet(t *testing.T) {
	// Setup server
	repo := newMockGRPCRepo()
	events := &mockGRPCEvents{}
	svc := usecase.NewInboundService(repo, events)
	srv := NewInboundServer(svc)

	// Start in-process gRPC
	server := grpc.NewServer()
	inboundv1.RegisterInboundServiceServer(server, srv)

	lis, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)

	go func() {
		_ = server.Serve(lis)
	}()
	defer server.Stop()

	// Client
	conn, err := grpc.NewClient(lis.Addr().String(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)
	defer conn.Close()

	client := inboundv1.NewInboundServiceClient(conn)
	ctx := context.Background()

	// Create inbound
	createResp, err := client.CreateInbound(ctx, &inboundv1.CreateInboundRequest{
		CustomerId:   "cust_1",
		ExpectedDate: "2026-06-01",
		Notes:        "test",
		Items: []*inboundv1.InboundItem{
			{Sku: "SKU001", QuantityDeclared: 10},
		},
	})
	require.NoError(t, err)
	require.NotNil(t, createResp.Inbound)
	assert.Equal(t, "cust_1", createResp.Inbound.CustomerId)
	assert.Equal(t, "submitted", createResp.Inbound.Status)
	assert.Len(t, createResp.Inbound.Items, 1)

	// Get inbound
	getResp, err := client.GetInbound(ctx, &inboundv1.GetInboundRequest{
		Id: createResp.Inbound.Id,
	})
	require.NoError(t, err)
	assert.Equal(t, createResp.Inbound.Id, getResp.Inbound.Id)
}

func TestInboundServer_ApproveInbound(t *testing.T) {
	repo := newMockGRPCRepo()
	events := &mockGRPCEvents{}
	svc := usecase.NewInboundService(repo, events)
	srv := NewInboundServer(svc)

	// Pre-save a submitted inbound
	in, _ := model.SubmitInbound("inb_test_approve", "cust_1", "2026-06-01", "",
		[]model.InboundItem{{SKU: "SKU001", QuantityDeclared: 10}})
	_ = repo.Save(context.Background(), in)

	// Start in-process gRPC
	server := grpc.NewServer()
	inboundv1.RegisterInboundServiceServer(server, srv)

	lis, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	go func() { _ = server.Serve(lis) }()
	defer server.Stop()

	conn, err := grpc.NewClient(lis.Addr().String(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)
	defer conn.Close()

	client := inboundv1.NewInboundServiceClient(conn)

	// Approve with inspection (inspect + approve in one call)
	resp, err := client.ApproveInbound(context.Background(), &inboundv1.ApproveInboundRequest{
		Id: "inb_test_approve",
		Inspection: &inboundv1.Inspection{
			InspectorId: "insp_1",
			Notes:       "all good",
			Passed:      true,
		},
	})
	require.NoError(t, err)
	assert.Equal(t, "approved", resp.Inbound.Status)
}

func TestInboundServer_HoldAndRelease(t *testing.T) {
	repo := newMockGRPCRepo()
	events := &mockGRPCEvents{}
	svc := usecase.NewInboundService(repo, events)
	srv := NewInboundServer(svc)

	in, _ := model.SubmitInbound("inb_test_hr", "cust_1", "2026-06-01", "",
		[]model.InboundItem{{SKU: "SKU001", QuantityDeclared: 10}})
	in.Inspect(model.Inspection{InspectorID: "insp_1", InspectedAt: time.Now()})
	repo.inbounds[in.ID] = in
	in.Status = model.StatusInspected

	server := grpc.NewServer()
	inboundv1.RegisterInboundServiceServer(server, srv)
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	go func() { _ = server.Serve(lis) }()
	defer server.Stop()

	conn, err := grpc.NewClient(lis.Addr().String(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)
	defer conn.Close()

	client := inboundv1.NewInboundServiceClient(conn)
	ctx := context.Background()

	// Hold
	holdResp, err := client.HoldInbound(ctx, &inboundv1.HoldInboundRequest{
		Id:     "inb_test_hr",
		Reason: "docs missing",
	})
	require.NoError(t, err)
	assert.Equal(t, "held", holdResp.Inbound.Status)

	// Release
	relResp, err := client.ReleaseInbound(ctx, &inboundv1.ReleaseInboundRequest{
		Id: "inb_test_hr",
	})
	require.NoError(t, err)
	assert.Equal(t, "released", relResp.Inbound.Status)
}
