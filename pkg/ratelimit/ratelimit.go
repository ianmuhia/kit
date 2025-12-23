// Package ratelimit provides distributed rate limiting using Redis.
package ratelimit

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/mennanov/limiters"
	"github.com/redis/go-redis/v9"
)

// Limiter provides per-key rate limiting using Redis backend.
type Limiter struct {
	redisClient *redis.Client
	enabled     bool
	capacity    int64
	rate        time.Duration
	keyPrefix   string
	logger      *slog.Logger
	registry    *limiters.Registry
	clock       limiters.Clock
	mu          sync.RWMutex
}

// Option is a functional option for configuring the Limiter.
type Option func(*Limiter)

// WithEnabled controls whether rate limiting is active.
func WithEnabled(enabled bool) Option {
	return func(l *Limiter) {
		l.enabled = enabled
	}
}

// WithCapacity sets the maximum number of requests allowed in the rate window.
func WithCapacity(capacity int64) Option {
	return func(l *Limiter) {
		l.capacity = capacity
	}
}

// WithRate sets the time window for rate limiting.
func WithRate(rate time.Duration) Option {
	return func(l *Limiter) {
		l.rate = rate
	}
}

// WithKeyPrefix sets the Redis key prefix for rate limit data.
func WithKeyPrefix(prefix string) Option {
	return func(l *Limiter) {
		l.keyPrefix = prefix
	}
}

// WithLogger sets a custom logger.
func WithLogger(logger *slog.Logger) Option {
	return func(l *Limiter) {
		l.logger = logger
	}
}

// WithClock sets a custom clock (useful for testing).
func WithClock(clock limiters.Clock) Option {
	return func(l *Limiter) {
		l.clock = clock
	}
}

// New creates a new rate limiter with Redis backend.
// Default configuration:
//   - Enabled: true
//   - Capacity: 100 requests
//   - Rate: 1 minute
//   - KeyPrefix: "ratelimit"
func New(redisClient *redis.Client, opts ...Option) *Limiter {
	l := &Limiter{
		redisClient: redisClient,
		enabled:     true,
		capacity:    100,
		rate:        time.Minute,
		keyPrefix:   "ratelimit",
		logger:      slog.Default(),
		registry:    limiters.NewRegistry(),
		clock:       limiters.NewSystemClock(),
	}

	for _, opt := range opts {
		opt(l)
	}

	return l
}

// Limit applies rate limiting for a specific key (e.g., IP address, user ID).
// Returns the wait duration if rate limited, or zero if the request is allowed.
func (l *Limiter) Limit(ctx context.Context, key string) (time.Duration, error) {
	if !l.enabled {
		return 0, nil
	}

	limiter := l.registry.GetOrCreate(
		key,
		func() any {
			return limiters.NewTokenBucket(
				l.capacity,
				l.rate,
				limiters.NewLockNoop(),
				limiters.NewTokenBucketRedis(
					l.redisClient,
					l.keyPrefix+":"+key,
					l.rate,
					false,
				),
				l.clock,
				&slogLogger{logger: l.logger},
			)
		},
		l.rate*2, // TTL for inactive limiters
		l.clock.Now(),
	)

	return limiter.(*limiters.TokenBucket).Limit(ctx)
}

// IsLimitExhausted checks if the error indicates rate limit exceeded.
func IsLimitExhausted(err error) bool {
	return err == limiters.ErrLimitExhausted
}

// slogLogger adapts slog.Logger to limiters.Logger interface.
type slogLogger struct {
	logger *slog.Logger
}

func (s *slogLogger) Log(v ...interface{}) {
	s.logger.Debug("rate limiter", "message", v)
}
