package service

import (
	"context"
	"fmt"

	"{{cookiecutter.source_path}}/{{cookiecutter.app_name}}/config"
	proto "{{cookiecutter.source_path}}/{{cookiecutter.app_name}}/proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/genproto/googleapis/api/httpbody"
	"github.com/go-coldbrew/errors"
	"github.com/go-coldbrew/log"
	"google.golang.org/grpc/health"
)

// svc should implement the service interface defined in the proto file
var _ proto.{{cookiecutter.service_name}}Server = (*svc)(nil)

// Service interface for the service
type svc struct {
	// health server for the service
	*health.Server
	// TODO: remove this, since this is just to demonstrate how to use config
	// prefix to be added to the message in the response
	prefix string
}

// ReadinessProbe for the service
// This is called by the kubernetes readiness probe
func (s *svc) ReadyCheck(ctx context.Context, _ *emptypb.Empty) (*httpbody.HttpBody, error) {
	return GetReadyState(ctx)
}

// LivenessProbe for the service
// This is called by the kubernetes liveness probe
func (s *svc) HealthCheck(ctx context.Context, _ *emptypb.Empty) (*httpbody.HttpBody, error) {
	return GetHealthCheck(ctx), nil
}

// Echo returns the message with the prefix added
// TODO: remove this, since this is just to demonstrate how to use endpoints and config
func (s *svc) Echo(_ context.Context, req *proto.EchoRequest) (*proto.EchoResponse, error) {
	return &proto.EchoResponse{
		Msg: fmt.Sprintf("%s: %s", s.prefix, req.GetMsg()),
	}, nil
}

// Error returns an error to the client
// TODO: remove this, since this is just to demonstrate how to use endpoints and config
func (s *svc) Error(ctx context.Context, req *proto.EchoRequest) (*proto.EchoResponse, error) {
	err := errors.New("This is an Error")
	log.Info(ctx, "error requested")
	return nil, errors.Wrap(err, "endpoint error")
}

// Creates a new Service instance and returns it
func New(cfg config.Config) (*svc, error) {
	// TODO: Application should validate the config here and return an error if it is invalid or missing
	s := &svc{
		// This is the health server for the service that is used for grpc
		Server: GetHealthCheckServer(),
		// TODO: remove this, since this is just to demonstrate how to use config
		prefix: cfg.Prefix,
	}
	// TODO: Application should initialize the service here and return an error if it fails to initialize

	// we call SetReady() here to indicate that the service is ready to serve traffic
	// service will not receive any traffic until this is called
	SetReady()
	return s, nil
}
