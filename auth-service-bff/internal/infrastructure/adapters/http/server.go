package httpadapter

import (
	"fmt"
	"net/http"

	"github.com/beabys/wms/auth-service-bff/internal/api/v1"
	"github.com/beabys/wms/pkg/logger"
)

// Server implements v1.ServerInterface for the login BFF HTTP API.
type Server struct {
	authSvc AuthHandlerService
	logger  logger.Logger
}

// NewServer creates a new Server with the given auth service.
func NewServer(authSvc AuthHandlerService, log logger.Logger) *Server {
	return &Server{
		authSvc: authSvc,
		logger:  log,
	}
}

// compile-time check
var _ v1.ServerInterface = (*Server)(nil)

// writeSuccess writes a success JSON response with status 200.
func (s *Server) writeSuccess(w http.ResponseWriter, data interface{}) {
	successResponseJSON(w, http.StatusOK, data)
}

// writeError writes an error JSON response with the given status and message.
func (s *Server) writeError(w http.ResponseWriter, status int, msg string) {
	s.logger.Warn("request failed",
		logger.LogField{Key: "status", Value: status},
		logger.LogField{Key: "error", Value: msg},
	)
	errorResponseJSON(w, status, fmt.Errorf("%s", msg))
}
