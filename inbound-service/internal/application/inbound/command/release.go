package command

// ReleaseInboundCommand releases a hold.
type ReleaseInboundCommand struct {
	InboundID  string
	ReleasedBy string
}
