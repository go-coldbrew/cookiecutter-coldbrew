package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"{{cookiecutter.source_path}}/{{cookiecutter.app_name}}/config"
	"{{cookiecutter.source_path}}/{{cookiecutter.app_name}}/service/metrics"
	mockmetrics "{{cookiecutter.source_path}}/{{cookiecutter.app_name}}/misc/mocks/metrics"
	proto "{{cookiecutter.source_path}}/{{cookiecutter.app_name}}/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TestNew(t *testing.T) {
	s, err := New(config.Get())
	assert.NoError(t, err)
	assert.NotNil(t, s)
}

func TestReadyCheck(t *testing.T) {
	s, err := New(config.Get())
	assert.NoError(t, err)
	assert.NotNil(t, s)

	SetNotReady()
	data, err := s.ReadyCheck(context.Background(), nil)
	assert.Error(t, err)
	assert.Nil(t, data)

	SetReady()
	data, err = s.ReadyCheck(context.Background(), nil)
	assert.NoError(t, err)
	assert.NotNil(t, data)
	assert.NotEmpty(t, data.Data)
}

func TestHealthCheck(t *testing.T) {
	s, err := New(config.Get())
	assert.NoError(t, err)
	assert.NotNil(t, s)

	data, err := s.HealthCheck(context.Background(), nil)
	assert.NoError(t, err)
	assert.NotNil(t, data)
	assert.NotEmpty(t, data.Data)
}

func TestEcho(t *testing.T) {
	const prefix = "testPrefix"
	const msg = "hello"

	m := mockmetrics.NewMetrics(t)
	m.EXPECT().IncEchoTotal(metrics.OutcomeSuccess).Once()
	m.EXPECT().ObserveEchoDuration(metrics.OutcomeSuccess, mock.AnythingOfType("time.Duration")).Once()

	s := &svc{
		Server:     GetHealthCheckServer(),
		monitoring: m,
		prefix:     prefix,
	}

	resp, err := s.Echo(context.Background(), &proto.EchoRequest{Msg: msg})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, prefix+": "+msg, resp.Msg)
}

func TestError(t *testing.T) {
	s, err := New(config.Get())
	assert.NoError(t, err)
	assert.NotNil(t, s)

	resp, err := s.Error(context.Background(), nil)
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestWorkers(t *testing.T) {
	s, err := New(config.Get())
	assert.NoError(t, err)

	w := s.Workers()
	assert.NotEmpty(t, w, "Workers() should return at least one worker")

	var found bool
	for _, worker := range w {
		if worker.GetName() == "cleanup" {
			found = true
			break
		}
	}
	assert.True(t, found, "Workers() should include a worker named 'cleanup'")
}

// fakeStreamEchoServer is a minimal grpc.ServerStreamingServer[proto.EchoToken]
// for unit-testing StreamEcho without a real gRPC connection. Tracks tokens
// sent via Send so tests can assert ordering and payload, and exposes a
// pluggable Context so tests can simulate client disconnect.
type fakeStreamEchoServer struct {
	grpc.ServerStream
	ctx  context.Context
	sent []*proto.EchoToken
}

func (f *fakeStreamEchoServer) Context() context.Context { return f.ctx }

func (f *fakeStreamEchoServer) Send(t *proto.EchoToken) error {
	f.sent = append(f.sent, t)
	return nil
}

func (f *fakeStreamEchoServer) SetHeader(metadata.MD) error  { return nil }
func (f *fakeStreamEchoServer) SendHeader(metadata.MD) error { return nil }
func (f *fakeStreamEchoServer) SetTrailer(metadata.MD)       {}

func TestStreamEcho(t *testing.T) {
	const prefix = "testPrefix"
	const msg = "hello world foo"

	m := mockmetrics.NewMetrics(t)
	m.EXPECT().IncStreamEchoTotal(metrics.OutcomeSuccess).Once()
	m.EXPECT().ObserveStreamEchoDuration(metrics.OutcomeSuccess, mock.AnythingOfType("time.Duration")).Once()
	m.EXPECT().ObserveStreamEchoTTFT(mock.AnythingOfType("time.Duration")).Once()

	s := &svc{
		Server:     GetHealthCheckServer(),
		monitoring: m,
		prefix:     prefix,
	}

	stream := &fakeStreamEchoServer{ctx: context.Background()}
	err := s.StreamEcho(&proto.EchoRequest{Msg: msg}, stream)
	assert.NoError(t, err)

	assert.Len(t, stream.sent, 3, "one frame per whitespace-separated word")
	for i, want := range []string{prefix + ": hello", prefix + ": world", prefix + ": foo"} {
		assert.Equal(t, want, stream.sent[i].GetToken(), "frame %d token", i)
		assert.Equal(t, int32(i), stream.sent[i].GetIndex(), "frame %d index", i)
	}
}

func TestStreamEcho_ContextCanceledMidStream(t *testing.T) {
	const prefix = "testPrefix"
	const msg = "hello world foo bar"

	m := mockmetrics.NewMetrics(t)
	m.EXPECT().IncStreamEchoTotal(metrics.OutcomeCanceled).Once()
	m.EXPECT().ObserveStreamEchoDuration(metrics.OutcomeCanceled, mock.AnythingOfType("time.Duration")).Once()
	m.EXPECT().ObserveStreamEchoTTFT(mock.AnythingOfType("time.Duration")).Maybe()

	s := &svc{
		Server:     GetHealthCheckServer(),
		monitoring: m,
		prefix:     prefix,
	}

	// Cancel the context after a short delay — long enough to emit a frame or
	// two but not enough to finish the four-token stream (each frame waits
	// streamEchoFrameDelay before the next iteration). Asserting the handler
	// stops generating mid-stream is the load-bearing safety property for
	// AI/LLM workloads: client disconnect must halt token production.
	ctx, cancel := context.WithCancel(context.Background())
	time.AfterFunc(streamEchoFrameDelay/2, cancel)

	stream := &fakeStreamEchoServer{ctx: ctx}
	err := s.StreamEcho(&proto.EchoRequest{Msg: msg}, stream)
	assert.Error(t, err, "StreamEcho should return error when context is canceled")
	assert.Less(t, len(stream.sent), 4, "handler should stop emitting after cancel")
}

func BenchmarkEcho(b *testing.B) {
	const prefix = "testPrefix"
	const msg = "hello"
	const expected = prefix + ": " + msg

	cfg := config.Get()
	cfg.Prefix = prefix
	s, err := New(cfg)
	if err != nil {
		b.Fatal(err)
	}

	ctx := context.Background()
	req := &proto.EchoRequest{Msg: msg}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp, err := s.Echo(ctx, req)
		if err != nil {
			b.Fatal(err)
		}
		if resp.Msg != expected {
			b.Fatalf("unexpected response: %s", resp.Msg)
		}
	}
}
