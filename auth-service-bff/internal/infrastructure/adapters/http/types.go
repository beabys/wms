package httpadapter

import (
	"net/http"

	"github.com/beabys/wms/pkg/logger"
)

// HTTPServer wraps an http.Server with lifecycle management.
type HTTPServer struct {
	Server *http.Server
	Logger logger.Logger
}
