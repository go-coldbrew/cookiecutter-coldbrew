package service

import (
	"context"
	"encoding/json"
	"sync/atomic"

	"google.golang.org/genproto/googleapis/api/httpbody"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"{{cookiecutter.source_path}}/{{cookiecutter.app_name}}/version"
)

var (
	// hc is the health check response for the service in JSON format
	hc []byte
	// ready is the readiness state of the service
	ready int32
)

// readiness states
const (
	statusReady    = 1
	statusNotReady = 0
)

func init() {
	// initialize the health check response for the service
	hc, _ = json.Marshal(version.Get())
}

// GetHealthCheck returns the health check response for the service
// This is used by the Kubernetes liveness probe to check the health of the service
func GetHealthCheck(context.Context) *httpbody.HttpBody {
	return &httpbody.HttpBody{
		ContentType: "application/json",
		Data:        hc,
	}
}

// GetReadyState returns the readiness state of the service and an error if the service is not ready
// This is used by the Kubernetes readiness probe to check the readiness of the service
func GetReadyState(ctx context.Context) (*httpbody.HttpBody, error) {
	st := atomic.LoadInt32(&ready)
	if st == statusReady {
		return GetHealthCheck(ctx), nil
	}
	return nil, status.Error(codes.Internal, "Not Ready to server traffic")
}

// SetNotReady sets the readiness state of the service to not ready
func SetNotReady() {
	atomic.StoreInt32(&ready, statusNotReady)
}

// SetReady sets the readiness state of the service to ready
func SetReady() {
	atomic.StoreInt32(&ready, statusReady)
}
