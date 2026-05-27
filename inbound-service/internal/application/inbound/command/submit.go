package command

import "github.com/beabys/wms/inbound-service/internal/domain/inbound/model"

// SubmitInboundCommand creates a new inbound.
type SubmitInboundCommand struct {
	CustomerID   string
	ExpectedDate string
	Notes        string
	Items        []SubmitItem
}

// SubmitItem represents an item in the submit command.
type SubmitItem struct {
	SKU              string
	QuantityDeclared int32
	QuantityReceived int32
	Dimensions       string
	Weight           float64
}

// InboundResult is the result returned by handlers.
type InboundResult struct {
	Inbound *model.Inbound
}
