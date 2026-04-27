package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"{{cookiecutter.source_path}}/{{cookiecutter.app_name}}/config"
	proto "{{cookiecutter.source_path}}/{{cookiecutter.app_name}}/proto"
	"{{cookiecutter.source_path}}/{{cookiecutter.app_name}}/service/metrics"
	"github.com/go-coldbrew/errors"
	cblog "github.com/go-coldbrew/log"
	"github.com/go-coldbrew/workers"
	"google.golang.org/genproto/googleapis/api/httpbody"
	"google.golang.org/grpc/health"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Compile-time interface checks — remove or adjust as you customize your service.
var (
	_ proto.{{cookiecutter.service_name}}Server         = (*svc)(nil)
	_ interface{ Stop() }                      = (*svc)(nil)
	_ interface{ Workers() []*workers.Worker } = (*svc)(nil)
)

// Service interface for the service
type svc struct {
	// health server for the service
	*health.Server
	// application metrics
	monitoring metrics.Metrics
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
// TODO: remove this, since this is just to demonstrate how to use endpoints, config, and logging
func (s *svc) Echo(ctx context.Context, req *proto.EchoRequest) (resp *proto.EchoResponse, err error) {
	start := time.Now()
	outcome := metrics.OutcomeSuccess
	defer func() {
		if err != nil {
			outcome = metrics.OutcomeError
		}
		s.monitoring.IncEchoTotal(outcome)
		s.monitoring.ObserveEchoDuration(outcome, time.Since(start))
	}()

	// Add typed context fields — these appear in all logs for this request.
	// ColdBrew interceptors already add trace_id and grpcMethod automatically.
	ctx = cblog.AddAttrsToContext(ctx, slog.Int("echo_msg_len", len(req.GetMsg())))

	slog.LogAttrs(ctx, slog.LevelInfo, "echo requested")

	return &proto.EchoResponse{
		Msg: fmt.Sprintf("%s: %s", s.prefix, req.GetMsg()),
	}, nil
}

// Error returns an error to the client
// TODO: remove this, since this is just to demonstrate how to use endpoints and config
func (s *svc) Error(ctx context.Context, req *proto.EchoRequest) (*proto.EchoResponse, error) {
	err := errors.New("This is an Error")
	slog.LogAttrs(ctx, slog.LevelInfo, "error requested")
	return nil, errors.Wrap(err, "endpoint error")
}

func (s *svc) Stop() {
	// Close database connections, flush buffers, etc.
}

// Workers returns background workers managed by ColdBrew via CBWorkerProvider.
// Workers are started alongside gRPC/HTTP servers with automatic panic recovery
// and configurable restart. Add your periodic tasks and long-running consumers here.
func (s *svc) Workers() []*workers.Worker {
	return []*workers.Worker{
		workers.NewWorker("cleanup").
			HandlerFunc(s.cleanup).
			Every(5 * time.Minute).
			WithJitter(10),
		// Uncomment to add a queue consumer:
		// workers.NewWorker("queue-consumer").HandlerFunc(s.consumeMessages),
	}
}

func (s *svc) cleanup(ctx context.Context, info *workers.WorkerInfo) error {
	slog.LogAttrs(ctx, slog.LevelInfo, "running periodic cleanup")
	// TODO: Add your cleanup logic here (e.g., purge expired sessions, compact data)
	return nil
}

// New creates a new Service instance and returns it
func New(cfg config.Config) (*svc, error) {
	// TODO: Application should validate the config here and return an error if it is invalid or missing
	s := &svc{
		// This is the health server for the service that is used for grpc
		Server: GetHealthCheckServer(),
		// application metrics
		monitoring: metrics.New(),
		// TODO: remove this, since this is just to demonstrate how to use config
		prefix: cfg.Prefix,
	}
	// TODO: Application should initialize the service here and return an error if it fails to initialize

	// we call SetReady() here to indicate that the service is ready to serve traffic
	// service will not receive any traffic until this is called
	SetReady()
	return s, nil
}
