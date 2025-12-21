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

// Config holds rate limiter configuration.
type Config struct {
	// Enabled controls whether rate limiting is active.
	Enabled bool
	// Capacity is the maximum number of requests allowed in the rate window.
	Capacity int64
	// Rate is the time window for rate limiting.
	Rate time.Duration
	// KeyPrefix is the Redis key prefix for rate limit data.
	KeyPrefix string
}

// DefaultConfig returns sensible defaults for rate limiting.
func DefaultConfig() Config {
	return Config{
		Enabled:   true,
		Capacity:  100,
		Rate:      time.Minute,
		KeyPrefix: "itara:ratelimit",
	}
}

// Limiter provides per-key rate limiting using Redis backend.
type Limiter struct {
	redisClient *redis.Client
	config      Config
	logger      *slog.Logger
	registry    *limiters.Registry
	clock       limiters.Clock
	mu          sync.RWMutex
}

// New creates a new rate limiter with Redis backend.
func New(redisClient *redis.Client, cfg Config, logger *slog.Logger) *Limiter {
	if logger == nil {
		logger = slog.Default()
	}

	return &Limiter{
		redisClient: redisClient,
		config:      cfg,
		logger:      logger,
		registry:    limiters.NewRegistry(),
		clock:       limiters.NewSystemClock(),
	}
}

// Limit applies rate limiting for a specific key (e.g., IP address, user ID).
// Returns the wait duration if rate limited, or zero if the request is allowed.
func (l *Limiter) Limit(ctx context.Context, key string) (time.Duration, error) {
	if !l.config.Enabled {
		return 0, nil
	}

	limiter := l.registry.GetOrCreate(
		key,
		func() any {
			return limiters.NewTokenBucket(
				l.config.Capacity,
				l.config.Rate,
				limiters.NewLockNoop(),
				limiters.NewTokenBucketRedis(
					l.redisClient,
					l.config.KeyPrefix+":"+key,
					l.config.Rate,
					false,
				),
				l.clock,
				&slogLogger{logger: l.logger},
			)
		},
		l.config.Rate*2, // TTL for inactive limiters
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
