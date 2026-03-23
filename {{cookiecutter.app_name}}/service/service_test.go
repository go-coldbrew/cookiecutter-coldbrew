package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"{{cookiecutter.source_path}}/{{cookiecutter.app_name}}/config"
	proto "{{cookiecutter.source_path}}/{{cookiecutter.app_name}}/proto"
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

	s, err := New(config.Get())
	assert.NoError(t, err)
	assert.NotNil(t, s)

	// override the prefix
	s.prefix = prefix

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

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp, err := s.Echo(context.Background(), &proto.EchoRequest{Msg: msg})
		if err != nil {
			b.Fatal(err)
		}
		if resp.Msg != expected {
			b.Fatalf("unexpected response: %s", resp.Msg)
		}
	}
}
