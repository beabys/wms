package command

import "github.com/beabys/wms/inventory-service/internal/domain/inventory/model"

// CreateBinLocationCommand creates a new bin location.
type CreateBinLocationCommand struct {
	WarehouseZone string
	Aisle         string
	Rack          string
	Shelf         string
}

// BinLocationResult is the result of bin location operations.
type BinLocationResult struct {
	BinLocation *model.BinLocation
}
