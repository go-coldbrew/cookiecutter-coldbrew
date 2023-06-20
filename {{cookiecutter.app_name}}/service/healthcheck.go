package service

import (
	"context"
	"encoding/json"

	"google.golang.org/genproto/googleapis/api/httpbody"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
	"{{cookiecutter.source_path}}/{{cookiecutter.app_name}}/version"
)

var (
	// hc is the health check response for the service in JSON format
	hc []byte
	// hcServer is the health check server for the service
	hcServer *health.Server
)

const (
	serviceName = "{{cookiecutter.grpc_package}}.{{cookiecutter.service_name}}"
)

func init() {
	// initialize the health check response for the service
	hc, _ = json.Marshal(version.Get())

	// initialize the health check server for the service
	hcServer = health.NewServer()

	// register the health check server for the service
	hcServer.SetServingStatus(serviceName, healthpb.HealthCheckResponse_NOT_SERVING)
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
	resp, err := hcServer.Check(context.Background(), &healthpb.HealthCheckRequest{
		Service: serviceName,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if resp.Status == healthpb.HealthCheckResponse_SERVING {
		return GetHealthCheck(ctx), nil
	}
	return nil, status.Error(codes.Internal, "Not Ready to server traffic")
}

// SetNotReady sets the readiness state of the service to not ready
func SetNotReady() {
	hcServer.Shutdown()
}

// SetReady sets the readiness state of the service to ready
func SetReady() {
	hcServer.Resume()
}

// GetHealthCheckServer returns the health check server for the service
func GetHealthCheckServer() *health.Server {
	return hcServer
}
