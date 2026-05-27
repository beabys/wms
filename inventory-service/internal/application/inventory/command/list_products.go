package command

import "github.com/beabys/wms/inventory-service/internal/domain/inventory/model"

// ListProductsQuery filters products.
type ListProductsQuery struct {
	PageSize   int32
	PageToken  string
	Category   string
	Search     string
}

// ListProductsResult is the result of listing products.
type ListProductsResult struct {
	Products      []*model.Product
	NextPageToken string
}
