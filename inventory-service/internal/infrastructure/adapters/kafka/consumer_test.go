package kafka

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/beabys/wms/inventory-service/internal/application/inventory/usecase"
	"github.com/beabys/wms/inventory-service/internal/domain/inventory/model"
)

// ---------------------------------------------------------------------------
// In-memory repositories for testing
// ---------------------------------------------------------------------------

type inMemoryProductRepo struct {
	products map[string]*model.Product // keyed by ID (SKU)
}

func newInMemoryProductRepo() *inMemoryProductRepo {
	return &inMemoryProductRepo{products: make(map[string]*model.Product)}
}

func (r *inMemoryProductRepo) GetByID(_ context.Context, id string) (*model.Product, error) {
	p, ok := r.products[id]
	if !ok {
		return nil, assert.AnError // simulate not-found
	}
	return p, nil
}

func (r *inMemoryProductRepo) Save(_ context.Context, product *model.Product) error {
	r.products[product.ID] = product
	return nil
}

func (r *inMemoryProductRepo) Update(_ context.Context, product *model.Product) error {
	r.products[product.ID] = product
	return nil
}

func (r *inMemoryProductRepo) List(_ context.Context, _, _ string, _ int32, _ string) ([]*model.Product, string, error) {
	var list []*model.Product
	for _, p := range r.products {
		list = append(list, p)
	}
	return list, "", nil
}

func (r *inMemoryProductRepo) Archive(_ context.Context, id string) error {
	if p, ok := r.products[id]; ok {
		p.Archive()
	}
	return nil
}

type inMemoryStockRepo struct {
	entries []*model.StockEntry
}

func newInMemoryStockRepo() *inMemoryStockRepo {
	return &inMemoryStockRepo{}
}

func (r *inMemoryStockRepo) GetByID(_ context.Context, _ string) (*model.StockEntry, error) {
	return nil, assert.AnError
}

func (r *inMemoryStockRepo) Save(_ context.Context, entry *model.StockEntry) error {
	r.entries = append(r.entries, entry)
	return nil
}

func (r *inMemoryStockRepo) UpdateQuantity(_ context.Context, _ *model.StockEntry) error {
	return nil
}

func (r *inMemoryStockRepo) List(_ context.Context, _, _, _ string, _ int32, _ string) ([]*model.StockEntry, string, error) {
	return r.entries, "", nil
}

// noopBinLocationRepo is a minimal stub that returns nothing.
type noopBinLocationRepo struct{}

func (n *noopBinLocationRepo) GetByID(_ context.Context, _ string) (*model.BinLocation, error) {
	return nil, assert.AnError
}

func (n *noopBinLocationRepo) Save(_ context.Context, _ *model.BinLocation) error {
	return nil
}

func (n *noopBinLocationRepo) List(_ context.Context, _ int32, _ string) ([]*model.BinLocation, string, error) {
	return nil, "", nil
}

// ---------------------------------------------------------------------------
// Mock message reader
// ---------------------------------------------------------------------------

type mockMessageReader struct {
	messages []kafka.Message
	idx      int
	err      error
}

func (m *mockMessageReader) ReadMessage(ctx context.Context) (kafka.Message, error) {
	if m.err != nil {
		return kafka.Message{}, m.err
	}
	if m.idx >= len(m.messages) {
		// block until context cancels (simulating no more messages)
		<-ctx.Done()
		return kafka.Message{}, ctx.Err()
	}
	msg := m.messages[m.idx]
	m.idx++
	return msg, nil
}

func (m *mockMessageReader) Close() error {
	return nil
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func marshalEvent(evt InboundApprovedEvent) []byte {
	data, err := json.Marshal(evt)
	if err != nil {
		panic(err)
	}
	return data
}

// findProductBySKU returns the first product with the given SKU from the in-memory repo.
func findProductBySKU(repo *inMemoryProductRepo, sku string) *model.Product {
	for _, p := range repo.products {
		if p.SKU == sku {
			return p
		}
	}
	return nil
}

func newTestConsumer(evts ...InboundApprovedEvent) (*InboundEventConsumer, *inMemoryProductRepo, *inMemoryStockRepo) {
	msgs := make([]kafka.Message, len(evts))
	for i, evt := range evts {
		msgs[i] = kafka.Message{
			Key:   []byte(evt.InboundID),
			Value: marshalEvent(evt),
		}
	}
	reader := &mockMessageReader{messages: msgs}

	productRepo := newInMemoryProductRepo()
	stockRepo := newInMemoryStockRepo()
	binLocationRepo := &noopBinLocationRepo{}
	noopPub := usecase.EventPublisher(&noopEventPublisher{})

	svc := usecase.NewInventoryService(productRepo, stockRepo, binLocationRepo, noopPub)

	logger, _ := zap.NewDevelopment()

	return &InboundEventConsumer{
		reader: reader,
		invSvc: svc,
		logger: logger,
	}, productRepo, stockRepo
}

// noopEventPublisher is a stub for EventPublisher used in tests.
type noopEventPublisher struct{}

func (n *noopEventPublisher) PublishStockAdjusted(_ context.Context, _ *model.StockEntry, _ string) error {
	return nil
}

func (n *noopEventPublisher) PublishLowStockAlert(_ context.Context, _ *model.StockEntry, _ float64) error {
	return nil
}

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

func TestNewInboundEventConsumer(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	svc := usecase.NewInventoryService(
		newInMemoryProductRepo(),
		newInMemoryStockRepo(),
		&noopBinLocationRepo{},
		&noopEventPublisher{},
	)
	c := NewInboundEventConsumer([]string{"localhost:9092"}, svc, logger)
	require.NotNil(t, c)
	assert.NotNil(t, c.reader)
	defer c.Close()
}

func TestConsumer_ProcessMessage_ValidEvent_CreatesProductAndStock(t *testing.T) {
	evt := InboundApprovedEvent{
		InboundID:  "inb_001",
		CustomerID: "cust_001",
		Items: []InboundItem{
			{
				SKU:              "SKU-001",
				QuantityDeclared: 100,
				QuantityReceived: 95,
				Weight:           1.5,
			},
		},
		Timestamp: time.Now(),
	}

	consumer, productRepo, stockRepo := newTestConsumer(evt)

	// Manually process — simulate what Run would do
	msg := kafka.Message{
		Key:   []byte("inb_001"),
		Value: marshalEvent(evt),
	}
	consumer.processMessage(context.Background(), msg)

	// Verify product was created and find its actual ID
	product := findProductBySKU(productRepo, "SKU-001")
	require.NotNil(t, product)
	assert.Equal(t, "SKU-001", product.SKU)
	assert.Equal(t, "Product SKU-001", product.Name)
	assert.Equal(t, "unit", product.Unit)

	// Verify stock was added with the correct product ID
	require.Len(t, stockRepo.entries, 1)
	assert.Equal(t, product.ID, stockRepo.entries[0].ProductID)
	assert.Equal(t, float64(95), stockRepo.entries[0].Quantity)
}

func TestConsumer_ProcessMessage_ExistingProduct(t *testing.T) {
	// Pre-seed a product
	productRepo := newInMemoryProductRepo()
	now := time.Now()
	product, err := model.RegisterProduct("SKU-002", "SKU-002", "Existing Product", "desc", "cat", "kg", 2.0, 10)
	require.NoError(t, err)
	product.CreatedAt = now
	product.UpdatedAt = now
	err = productRepo.Save(context.Background(), product)
	require.NoError(t, err)

	stockRepo := newInMemoryStockRepo()
	binLocationRepo := &noopBinLocationRepo{}
	noopPub := &noopEventPublisher{}
	svc := usecase.NewInventoryService(productRepo, stockRepo, binLocationRepo, noopPub)

	logger, _ := zap.NewDevelopment()
	consumer := &InboundEventConsumer{
		reader: &mockMessageReader{},
		invSvc: svc,
		logger: logger,
	}

	evt := InboundApprovedEvent{
		InboundID:  "inb_002",
		CustomerID: "cust_001",
		Items: []InboundItem{
			{
				SKU:              "SKU-002",
				QuantityDeclared: 50,
				QuantityReceived: 50,
				Weight:           1.0,
			},
		},
		Timestamp: time.Now(),
	}

	consumer.processMessage(context.Background(), kafka.Message{
		Key:   []byte("inb_002"),
		Value: marshalEvent(evt),
	})

	// Verify product still has same name (not overwritten)
	p := findProductBySKU(productRepo, "SKU-002")
	require.NotNil(t, p)
	assert.Equal(t, "Existing Product", p.Name)

	// Verify stock was added with the correct product ID
	require.Len(t, stockRepo.entries, 1)
	assert.Equal(t, p.ID, stockRepo.entries[0].ProductID)
	assert.Equal(t, float64(50), stockRepo.entries[0].Quantity)
}

func TestConsumer_ProcessMessage_InvalidJSON(t *testing.T) {
	consumer, _, stockRepo := newTestConsumer()

	msg := kafka.Message{
		Key:   []byte("bad-key"),
		Value: []byte("{invalid json"),
	}
	// Should not panic
	consumer.processMessage(context.Background(), msg)

	// No stock should be added
	assert.Empty(t, stockRepo.entries)
}

func TestConsumer_Run_ContextCancellation(t *testing.T) {
	// Consumer with no messages — will block on ReadMessage until ctx cancels.
	reader := &mockMessageReader{
		messages: nil,
	}
	logger, _ := zap.NewDevelopment()
	svc := usecase.NewInventoryService(
		newInMemoryProductRepo(),
		newInMemoryStockRepo(),
		&noopBinLocationRepo{},
		&noopEventPublisher{},
	)
	consumer := &InboundEventConsumer{
		reader: reader,
		invSvc: svc,
		logger: logger,
	}

	ctx, cancel := context.WithCancel(context.Background())
	wg, ctx := errgroup.WithContext(ctx)

	err := consumer.Run(ctx, wg)
	require.NoError(t, err)

	// Cancel after a short delay
	go func() {
		time.Sleep(10 * time.Millisecond)
		cancel()
	}()

	// Wait for goroutine to exit
	err = wg.Wait()
	require.NoError(t, err)
}
