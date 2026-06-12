package command

// ListAuditLogsQuery represents paginated audit log list filters.
type ListAuditLogsQuery struct {
	CustomerID string
	Page       int
	PageSize   int
}

// ListAuditLogsResult contains paginated audit log results.
type ListAuditLogsResult struct {
	Entries    []*AuditLogResult
	TotalCount int
	Page       int
	PageSize   int
}

// AuditLogResult contains audit log data.
type AuditLogResult struct {
	ID          string
	CustomerID  string
	Action      string
	PerformedBy string
	Details     string
	CreatedAt   string
}
