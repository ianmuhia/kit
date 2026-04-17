package ratelimit

import (
	"errors"
	"testing"
	"time"

	"github.com/mennanov/limiters"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew_defaults(t *testing.T) {
	l := New(nil)
	assert.True(t, l.enabled)
	assert.Equal(t, int64(100), l.capacity)
	assert.Equal(t, time.Minute, l.rate)
	assert.Equal(t, "ratelimit", l.keyPrefix)
	assert.NotNil(t, l.logger)
	assert.NotNil(t, l.registry)
	assert.NotNil(t, l.clock)
}

func TestNew_invalidCapacityFallsBack(t *testing.T) {
	l := New(nil, WithCapacity(0))
	assert.Equal(t, int64(100), l.capacity)
}

func TestNew_invalidRateFallsBack(t *testing.T) {
	l := New(nil, WithRate(-1))
	assert.Equal(t, time.Minute, l.rate)
}

func TestNew_options(t *testing.T) {
	l := New(nil,
		WithEnabled(false),
		WithCapacity(50),
		WithRate(30*time.Second),
		WithKeyPrefix("api"),
	)
	assert.False(t, l.enabled)
	assert.Equal(t, int64(50), l.capacity)
	assert.Equal(t, 30*time.Second, l.rate)
	assert.Equal(t, "api", l.keyPrefix)
}

func TestIsLimitExhausted(t *testing.T) {
	assert.True(t, IsLimitExhausted(limiters.ErrLimitExhausted))
	assert.False(t, IsLimitExhausted(errors.New("other error")))
	assert.False(t, IsLimitExhausted(nil))
}

func TestErrRateLimiterUnavailable_isDistinct(t *testing.T) {
	require.NotNil(t, ErrRateLimiterUnavailable)
	assert.False(t, IsLimitExhausted(ErrRateLimiterUnavailable))
}

func TestLimit_disabledReturnsZero(t *testing.T) {
	l := New(nil, WithEnabled(false))
	wait, err := l.Limit(t.Context(), "user:1")
	require.NoError(t, err)
	assert.Zero(t, wait)
}
