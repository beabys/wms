package command

import "github.com/beabys/wms/inventory-service/internal/domain/inventory/model"

// ListBinLocationsQuery retrieves all bin locations.
type ListBinLocationsQuery struct {
	PageSize  int32
	PageToken string
}

// ListBinLocationsResult is the result of listing bin locations.
type ListBinLocationsResult struct {
	BinLocations  []*model.BinLocation
	NextPageToken string
}
