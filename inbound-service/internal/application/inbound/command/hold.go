package command

// HoldInboundCommand places a hold.
type HoldInboundCommand struct {
	InboundID string
	Reason    string
}
