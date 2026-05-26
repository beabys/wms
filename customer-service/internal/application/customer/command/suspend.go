package command

// SuspendCustomerCommand suspends an active customer.
type SuspendCustomerCommand struct {
	CustomerID string
	Reason     string
}
