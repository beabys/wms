package transformer

import (
	"time"

	"github.com/beabys/wms/auth-service/internal/application/auth/command"
	"github.com/beabys/wms/auth-service/internal/domain/auth/model"
	authv1 "github.com/beabys/wms/proto/gen/go/auth/v1"
	commonv1 "github.com/beabys/wms/proto/gen/go/common/v1"
)

// UserToProto converts a domain User to a proto User.
func UserToProto(u *model.User) *authv1.User {
	if u == nil {
		return nil
	}
	return &authv1.User{
		Id:         u.ID,
		Email:      u.Email,
		Name:       u.Name,
		Role:       u.Role,
		CustomerId: customerIDString(u.CustomerID),
		Active:     u.Active,
		CreatedAt:  timeToProto(u.CreatedAt),
		UpdatedAt:  timeToProto(u.UpdatedAt),
	}
}

// UserFromProto converts a proto User to a domain User.
func UserFromProto(u *authv1.User) *model.User {
	if u == nil {
		return nil
	}
	var customerID *string
	if cid := u.GetCustomerId(); cid != "" {
		customerID = &cid
	}
	return &model.User{
		ID:         u.GetId(),
		Email:      u.GetEmail(),
		Name:       u.GetName(),
		Role:       u.GetRole(),
		CustomerID: customerID,
		Active:     u.GetActive(),
		CreatedAt:  protoToTime(u.GetCreatedAt()),
		UpdatedAt:  protoToTime(u.GetUpdatedAt()),
	}
}

// UserResultToProto converts a command.UserResult to a proto User.
func UserResultToProto(u *command.UserResult) *authv1.User {
	if u == nil {
		return nil
	}
	return &authv1.User{
		Id:         u.ID,
		Email:      u.Email,
		Name:       u.Name,
		Role:       u.Role,
		CustomerId: customerIDString(u.CustomerID),
		Active:     u.Active,
		CreatedAt:  parseTimestamp(u.CreatedAt),
		UpdatedAt:  parseTimestamp(u.UpdatedAt),
	}
}

// InviteToProto converts a domain InviteToken to a proto InviteEntry.
func InviteToProto(i *model.InviteToken) *authv1.InviteEntry {
	if i == nil {
		return nil
	}
	return &authv1.InviteEntry{
		Id:        i.ID,
		Email:     i.Email,
		Token:     i.Token,
		InvitedBy: i.InvitedBy,
		Status:    i.Status,
		ExpiresAt: i.ExpiresAt.Unix(),
		CreatedAt: i.CreatedAt.Unix(),
	}
}

// InviteFromProto converts a proto InviteEntry to a domain InviteToken.
func InviteFromProto(i *authv1.InviteEntry) *model.InviteToken {
	if i == nil {
		return nil
	}
	return &model.InviteToken{
		ID:        i.GetId(),
		Email:     i.GetEmail(),
		Token:     i.GetToken(),
		InvitedBy: i.GetInvitedBy(),
		Status:    i.GetStatus(),
		ExpiresAt: time.Unix(i.GetExpiresAt(), 0),
		CreatedAt: time.Unix(i.GetCreatedAt(), 0),
	}
}

// RoleToProto converts a domain Role to a proto Role.
func RoleToProto(r *model.Role) *authv1.Role {
	if r == nil {
		return nil
	}
	perms := make([]*authv1.Permission, 0, len(r.Permissions))
	for _, p := range r.Permissions {
		perms = append(perms, PermissionToProto(&p))
	}
	return &authv1.Role{
		Id:          r.ID,
		Name:        r.Name,
		Description: r.Description,
		Permissions: perms,
		IsSystem:    r.IsSystem,
	}
}

// PermissionToProto converts a domain Permission to a proto Permission.
func PermissionToProto(p *model.Permission) *authv1.Permission {
	if p == nil {
		return nil
	}
	return &authv1.Permission{
		Service:  string(p.Service),
		Action:   string(p.Action),
		Resource: string(p.Resource),
	}
}

// PermissionFromProto converts a proto Permission to a domain Permission.
// Unknown/invalid string values are preserved as-is (validation happens at use sites).
func PermissionFromProto(p *authv1.Permission) *model.Permission {
	if p == nil {
		return nil
	}
	return &model.Permission{
		Service:  model.Service(p.GetService()),
		Action:   model.Action(p.GetAction()),
		Resource: model.Resource(p.GetResource()),
	}
}

// customerIDString returns the customer ID string or empty string if nil.
func customerIDString(cid *string) string {
	if cid == nil {
		return ""
	}
	return *cid
}

// timeToProto converts a time.Time to a proto Timestamp.
func timeToProto(t time.Time) *commonv1.Timestamp {
	return &commonv1.Timestamp{
		Seconds: t.Unix(),
		Nanos:   int32(t.Nanosecond()),
	}
}

// protoToTime converts a proto Timestamp to a time.Time.
func protoToTime(ts *commonv1.Timestamp) time.Time {
	if ts == nil {
		return time.Time{}
	}
	return time.Unix(ts.GetSeconds(), int64(ts.GetNanos()))
}

// parseTimestamp parses an RFC3339 timestamp string to a proto Timestamp.
func parseTimestamp(ts string) *commonv1.Timestamp {
	t, err := time.Parse("2006-01-02T15:04:05Z", ts)
	if err != nil {
		return nil
	}
	return &commonv1.Timestamp{
		Seconds: t.Unix(),
		Nanos:   int32(t.Nanosecond()),
	}
}
