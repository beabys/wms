package command

// AssignCompanyRoleCommand assigns a role (with permissions) to a user.
type AssignCompanyRoleCommand struct {
	CustomerID  string
	UserID      string
	RoleName    string
	Permissions []string
	AssignedBy  string
}

// CompanyRoleResult contains role data.
type CompanyRoleResult struct {
	ID          string
	CustomerID  string
	Name        string
	Permissions []string
	IsDefault   bool
}
