package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReady(t *testing.T) {
	ctx := context.Background()
	_, err := GetReadyState(ctx)
	assert.Error(t, err, "initial state should be not ready")

	SetReady()
	ready, err := GetReadyState(ctx)
	assert.NoError(t, err, "should be ready")
	assert.NotNil(t, ready, "should return response")
	assert.NotEmpty(t, ready.Data, "should return response")

	SetNotReady()
	ready, err = GetReadyState(ctx)
	assert.Error(t, err, "should not be ready")
	assert.Nil(t, ready, "should not return response")
}
