package command

// ListCustomersQuery lists customers, optionally filtered by status.
type ListCustomersQuery struct {
	Status   string
	Page     int32
	PageSize int32
}
