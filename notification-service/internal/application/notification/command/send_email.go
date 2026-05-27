package command

// SendEmailCommand represents a request to send an email.
type SendEmailCommand struct {
	To           string
	Subject      string
	Body         string
	TemplateID   string
	TemplateData map[string]string
}

// SendEmailResult contains the result of a send email operation.
type SendEmailResult struct {
	MessageID string
	Status    string
}
