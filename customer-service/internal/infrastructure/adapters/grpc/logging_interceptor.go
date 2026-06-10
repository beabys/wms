package grpcdapter

import (
	"context"
	"time"

	"github.com/beabys/wms/pkg/logger"
	"google.golang.org/grpc"
)

// LoggingInterceptor returns a unary server interceptor that logs each gRPC call.
func LoggingInterceptor(log logger.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		duration := time.Since(start)
		if err != nil {
			log.Error("gRPC call failed", err,
				logger.LogField{Key: "method", Value: info.FullMethod},
				logger.LogField{Key: "duration", Value: duration},
			)
		} else {
			log.Info("gRPC call succeeded",
				logger.LogField{Key: "method", Value: info.FullMethod},
				logger.LogField{Key: "duration", Value: duration},
			)
		}
		return resp, err
	}
}
