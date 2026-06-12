package grpcdapter

import (
	"errors"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrUnauthenticated  = errors.New("unauthenticated")
	ErrNotFound         = errors.New("not found")
	ErrInvalidArgument  = errors.New("invalid argument")
	ErrPermissionDenied = errors.New("permission denied")
	ErrInternal         = errors.New("internal error")
)

// mapGrpcError converts a gRPC error to an HTTP-friendly error.
// Returns nil if err is nil.
func mapGrpcError(err error) error {
	if err == nil {
		return nil
	}
	st, ok := status.FromError(err)
	if !ok {
		return ErrInternal
	}
	switch st.Code() {
	case codes.Unauthenticated:
		return ErrUnauthenticated
	case codes.NotFound:
		return ErrNotFound
	case codes.InvalidArgument, codes.FailedPrecondition:
		return ErrInvalidArgument
	case codes.PermissionDenied:
		return ErrPermissionDenied
	case codes.Internal, codes.Unavailable, codes.Unknown:
		return ErrInternal
	default:
		return ErrInternal
	}
}

// HTTPStatusFromError returns the HTTP status code for a known error.
func HTTPStatusFromError(err error) int {
	if err == nil {
		return http.StatusOK
	}
	switch {
	case errors.Is(err, ErrUnauthenticated):
		return http.StatusUnauthorized
	case errors.Is(err, ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, ErrInvalidArgument):
		return http.StatusBadRequest
	case errors.Is(err, ErrPermissionDenied):
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}
