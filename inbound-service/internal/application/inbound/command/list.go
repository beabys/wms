package command

import "github.com/beabys/wms/inbound-service/internal/domain/inbound/model"

// ListInboundsQuery lists inbounds with filters.
type ListInboundsQuery struct {
	CustomerID string
	Status     string
	PageSize   int32
	PageToken  string
}

// ListInboundsResult is a paginated result.
type ListInboundsResult struct {
	Inbounds      []*model.Inbound
	NextPageToken string
}
