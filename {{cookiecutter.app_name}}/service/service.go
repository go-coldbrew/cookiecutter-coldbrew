package service

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"{{cookiecutter.source_path}}/{{cookiecutter.app_name}}/config"
	proto "{{cookiecutter.source_path}}/{{cookiecutter.app_name}}/proto"
	"{{cookiecutter.source_path}}/{{cookiecutter.app_name}}/service/metrics"
	"github.com/go-coldbrew/errors"
	cblog "github.com/go-coldbrew/log"
	"github.com/go-coldbrew/workers"
	"google.golang.org/genproto/googleapis/api/httpbody"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/status"
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

// streamEchoFrameDelay paces frame emission so the streaming behavior is
// visible end-to-end (browser EventSource, curl -N, grpcurl). For real
// workloads — LLM tokens, progress events, etc. — remove the sleep and emit
// frames as the upstream source produces them.
const streamEchoFrameDelay = 50 * time.Millisecond

// StreamEcho streams one EchoToken per whitespace-separated word in the
// request message. Demonstrates server-streaming over both native gRPC and
// the HTTP gateway — the gateway emits newline-delimited JSON by default,
// or Server-Sent Events when the client sends Accept: text/event-stream
// (ColdBrew registers the SSE marshaler by default; set
// DISABLE_SSE_MARSHALER=true on core to opt out).
//
// Context cancellation is the load-bearing piece for AI/LLM workloads:
// client disconnect cancels stream.Context(), and the handler must observe
// it to stop generating (and stop paying for) tokens.
func (s *svc) StreamEcho(req *proto.EchoRequest, stream grpc.ServerStreamingServer[proto.EchoToken]) (err error) {
	ctx := stream.Context()
	start := time.Now()
	outcome := metrics.OutcomeSuccess
	defer func() {
		if err != nil {
			outcome = metrics.OutcomeError
			if ctxErr := ctx.Err(); ctxErr != nil {
				outcome = metrics.OutcomeCanceled
			}
		}
		s.monitoring.IncStreamEchoTotal(outcome)
		s.monitoring.ObserveStreamEchoDuration(outcome, time.Since(start))
	}()

	tokens := strings.Fields(req.GetMsg())
	ctx = cblog.AddAttrsToContext(ctx, slog.Int("stream_echo_tokens", len(tokens)))
	slog.LogAttrs(ctx, slog.LevelInfo, "stream_echo requested")

	for i, token := range tokens {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return status.Error(codes.Canceled, fmt.Sprintf("stream_echo canceled: %v", ctxErr))
		}

		if sendErr := stream.Send(&proto.EchoToken{
			Token: fmt.Sprintf("%s: %s", s.prefix, token),
			Index: int32(i),
		}); sendErr != nil {
			// Preserve the canonical gRPC code when grpc-go already wrapped
			// the error (e.g. Canceled / Unavailable on client disconnect).
			if _, ok := status.FromError(sendErr); ok {
				return sendErr
			}
			return status.Error(codes.Internal, fmt.Sprintf("stream_echo send: %v", sendErr))
		}

		if i == 0 {
			s.monitoring.ObserveStreamEchoTTFT(time.Since(start))
		}

		// ctx-aware pacing: a client disconnect during the artificial delay
		// must stop the handler immediately, not wait for the next iteration.
		select {
		case <-time.After(streamEchoFrameDelay):
		case <-ctx.Done():
			return status.Error(codes.Canceled, fmt.Sprintf("stream_echo canceled: %v", ctx.Err()))
		}
	}
	return nil
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
