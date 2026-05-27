package command

// AdjustReason enumerates reasons for stock adjustments.
type AdjustReason string

const (
	AdjustReasonInboundApproval AdjustReason = "INBOUND_APPROVAL"
	AdjustReasonManual          AdjustReason = "MANUAL"
	AdjustReasonReturn          AdjustReason = "RETURN"
	AdjustReasonDisposal        AdjustReason = "DISPOSAL"
)
