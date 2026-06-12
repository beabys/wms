package grpcdapter

import (
	"context"
	"errors"
	"testing"

	"github.com/beabys/wms/pkg/logger"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

type captureLogger struct {
	infoMsgs  []string
	errorMsgs []string
}

func (l *captureLogger) GetLogger() any                                    { return nil }
func (l *captureLogger) Debug(msg string, fields ...logger.LogField)       {}
func (l *captureLogger) Info(msg string, fields ...logger.LogField)        { l.infoMsgs = append(l.infoMsgs, msg) }
func (l *captureLogger) Warn(msg string, fields ...logger.LogField)        {}
func (l *captureLogger) Error(msg string, err error, fields ...logger.LogField) { l.errorMsgs = append(l.errorMsgs, msg) }
func (l *captureLogger) Fatal(msg string, fields ...logger.LogField)       {}

func TestLoggingInterceptor_Success(t *testing.T) {
	log := &captureLogger{}
	interceptor := LoggingInterceptor(log)
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "ok", nil
	}

	resp, err := interceptor(context.Background(), "req", &grpc.UnaryServerInfo{
		FullMethod: "/test.Service/Method",
	}, handler)

	assert.NoError(t, err)
	assert.Equal(t, "ok", resp)
	assert.Len(t, log.infoMsgs, 1)
	assert.Contains(t, log.infoMsgs[0], "gRPC call succeeded")
	assert.Len(t, log.errorMsgs, 0)
}

func TestLoggingInterceptor_Error(t *testing.T) {
	log := &captureLogger{}
	interceptor := LoggingInterceptor(log)
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, errors.New("test error")
	}

	resp, err := interceptor(context.Background(), "req", &grpc.UnaryServerInfo{
		FullMethod: "/test.Service/Method",
	}, handler)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Len(t, log.errorMsgs, 1)
	assert.Contains(t, log.errorMsgs[0], "gRPC call failed")
	assert.Len(t, log.infoMsgs, 0)
}

func TestLoggingInterceptor_Panic(t *testing.T) {
	log := &captureLogger{}
	interceptor := LoggingInterceptor(log)
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		panic("test panic")
	}

	assert.Panics(t, func() {
		interceptor(context.Background(), "req", &grpc.UnaryServerInfo{
			FullMethod: "/test.Service/Method",
		}, handler)
	})
}
