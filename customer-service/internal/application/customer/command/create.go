package command

import "github.com/beabys/wms/customer-service/internal/domain/customer/model"

// CreateCustomerCommand is sent by an admin or BFF to create a new customer.
type CreateCustomerCommand struct {
	CompanyName string
	VATNumber   string
	Address     model.Address
	RateCardID  string
}
