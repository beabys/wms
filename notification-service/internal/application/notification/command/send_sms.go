package command

// SendSMSCommand represents a request to send an SMS.
type SendSMSCommand struct {
	To      string
	Message string
}

// SendSMSResult contains the result of a send SMS operation.
type SendSMSResult struct {
	MessageID string
	Status    string
}
