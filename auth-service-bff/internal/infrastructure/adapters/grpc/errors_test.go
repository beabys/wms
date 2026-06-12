package grpcdapter

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestMapGrpcError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		wantErr  error
		wantNil  bool
	}{
		{
			name:    "nil input returns nil",
			err:     nil,
			wantNil: true,
		},
		{
			name:    "unauthenticated",
			err:     status.Error(codes.Unauthenticated, "bad token"),
			wantErr: ErrUnauthenticated,
		},
		{
			name:    "not found",
			err:     status.Error(codes.NotFound, "user not found"),
			wantErr: ErrNotFound,
		},
		{
			name:    "invalid argument",
			err:     status.Error(codes.InvalidArgument, "bad request"),
			wantErr: ErrInvalidArgument,
		},
		{
			name:    "failed precondition",
			err:     status.Error(codes.FailedPrecondition, "precondition failed"),
			wantErr: ErrInvalidArgument,
		},
		{
			name:    "permission denied",
			err:     status.Error(codes.PermissionDenied, "no access"),
			wantErr: ErrPermissionDenied,
		},
		{
			name:    "internal",
			err:     status.Error(codes.Internal, "server error"),
			wantErr: ErrInternal,
		},
		{
			name:    "unavailable",
			err:     status.Error(codes.Unavailable, "service down"),
			wantErr: ErrInternal,
		},
		{
			name:    "unknown",
			err:     status.Error(codes.Unknown, "something"),
			wantErr: ErrInternal,
		},
		{
			name:    "non-grpc error",
			err:     errors.New("plain error"),
			wantErr: ErrInternal,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mapGrpcError(tt.err)
			if tt.wantNil {
				assert.NoError(t, got)
				return
			}
			assert.ErrorIs(t, got, tt.wantErr)
		})
	}
}

func TestHTTPStatusFromError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want int
	}{
		{"nil", nil, http.StatusOK},
		{"unauthenticated", ErrUnauthenticated, http.StatusUnauthorized},
		{"not found", ErrNotFound, http.StatusNotFound},
		{"invalid argument", ErrInvalidArgument, http.StatusBadRequest},
		{"permission denied", ErrPermissionDenied, http.StatusForbidden},
		{"internal", ErrInternal, http.StatusInternalServerError},
		{"unknown error", errors.New("unknown"), http.StatusInternalServerError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := HTTPStatusFromError(tt.err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSentinelErrors(t *testing.T) {
	assert.Error(t, ErrUnauthenticated)
	assert.Error(t, ErrNotFound)
	assert.Error(t, ErrInvalidArgument)
	assert.Error(t, ErrPermissionDenied)
	assert.Error(t, ErrInternal)

	assert.NotEqual(t, ErrUnauthenticated, ErrNotFound)
	assert.NotEqual(t, ErrNotFound, ErrInternal)
}
