package messaging

import (
	"fmt"
	"log/slog"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-nats/v2/pkg/nats"
	"github.com/ThreeDotsLabs/watermill/message"
	nc "github.com/nats-io/nats.go"
)

// SubscriberConfig holds essential configuration for NATS JetStream subscriber.
type SubscriberConfig struct {
	url               string
	name              string
	durablePrefix     string
	token             string
	username          string
	password          string
	logger            *slog.Logger
	unmarshaler       nats.Unmarshaler
	autoProvision     bool
	maxReconnects     int
	natsOptions       []nc.Option
	disconnectHandler func(*nc.Conn, error)
	reconnectHandler  func(*nc.Conn)
}

// SubscriberOption is a functional option for configuring the subscriber.
type SubscriberOption func(*SubscriberConfig)

// WithSubscriberURL sets the NATS server URL.
func WithSubscriberURL(url string) SubscriberOption {
	return func(c *SubscriberConfig) {
		c.url = url
	}
}

// WithSubscriberName sets the client name.
func WithSubscriberName(name string) SubscriberOption {
	return func(c *SubscriberConfig) {
		c.name = name
	}
}

// WithDurablePrefix sets the durable prefix for JetStream consumers.
func WithDurablePrefix(prefix string) SubscriberOption {
	return func(c *SubscriberConfig) {
		c.durablePrefix = prefix
	}
}

// WithSubscriberToken sets the authentication token.
func WithSubscriberToken(token string) SubscriberOption {
	return func(c *SubscriberConfig) {
		c.token = token
	}
}

// WithSubscriberUserPassword sets username and password authentication.
func WithSubscriberUserPassword(username, password string) SubscriberOption {
	return func(c *SubscriberConfig) {
		c.username = username
		c.password = password
	}
}

// WithSubscriberLogger sets the logger.
func WithSubscriberLogger(logger *slog.Logger) SubscriberOption {
	return func(c *SubscriberConfig) {
		c.logger = logger
	}
}

// WithUnmarshaler sets a custom unmarshaler (default is GobMarshaler).
func WithUnmarshaler(unmarshaler nats.Unmarshaler) SubscriberOption {
	return func(c *SubscriberConfig) {
		c.unmarshaler = unmarshaler
	}
}

// WithSubscriberAutoProvision enables/disables JetStream auto-provisioning.
func WithSubscriberAutoProvision(enable bool) SubscriberOption {
	return func(c *SubscriberConfig) {
		c.autoProvision = enable
	}
}

// WithSubscriberMaxReconnects sets the maximum number of reconnection attempts.
// Use -1 for infinite reconnects (default).
func WithSubscriberMaxReconnects(max int) SubscriberOption {
	return func(c *SubscriberConfig) {
		c.maxReconnects = max
	}
}

// WithSubscriberNATSOptions adds custom NATS connection options.
func WithSubscriberNATSOptions(opts ...nc.Option) SubscriberOption {
	return func(c *SubscriberConfig) {
		c.natsOptions = append(c.natsOptions, opts...)
	}
}

// WithSubscriberDisconnectHandler sets a custom disconnect handler.
func WithSubscriberDisconnectHandler(handler func(*nc.Conn, error)) SubscriberOption {
	return func(c *SubscriberConfig) {
		c.disconnectHandler = handler
	}
}

// WithSubscriberReconnectHandler sets a custom reconnect handler.
func WithSubscriberReconnectHandler(handler func(*nc.Conn)) SubscriberOption {
	return func(c *SubscriberConfig) {
		c.reconnectHandler = handler
	}
}

// defaultSubscriberConfig returns sensible defaults.
func defaultSubscriberConfig() *SubscriberConfig {
	return &SubscriberConfig{
		url:           "nats://localhost:4222",
		name:          "watermill-subscriber",
		durablePrefix: "service",
		logger:        slog.Default(),
		unmarshaler:   &nats.GobMarshaler{},
		autoProvision: true,
		maxReconnects: -1,
		disconnectHandler: func(conn *nc.Conn, err error) {
			if err != nil {
				slog.Error("NATS subscriber disconnected", "error", err)
			}
		},
		reconnectHandler: func(conn *nc.Conn) {
			slog.Info("NATS subscriber reconnected", "url", conn.ConnectedUrl())
		},
	}
}

// NewSubscriber creates a NATS JetStream subscriber with production-ready settings.
func NewSubscriber(opts ...SubscriberOption) (message.Subscriber, error) {
	config := defaultSubscriberConfig()

	// Apply all options
	for _, opt := range opts {
		opt(config)
	}

	if config.url == "" {
		return nil, fmt.Errorf("NATS URL is required")
	}

	natsOpts := buildSubscriberNATSOptions(config)

	subscriber, err := nats.NewSubscriber(
		nats.SubscriberConfig{
			URL:         config.url,
			NatsOptions: natsOpts,
			Unmarshaler: config.unmarshaler,
			JetStream: nats.JetStreamConfig{
				AutoProvision: config.autoProvision,
				DurablePrefix: config.durablePrefix,
			},
		},
		watermill.NewSlogLogger(config.logger),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create subscriber: %w", err)
	}

	return subscriber, nil
}

// buildSubscriberNATSOptions constructs NATS connection options for subscriber.
func buildSubscriberNATSOptions(config *SubscriberConfig) []nc.Option {
	opts := []nc.Option{
		nc.Name(config.name),
		nc.MaxReconnects(config.maxReconnects),
		nc.RetryOnFailedConnect(true),
		nc.DisconnectErrHandler(config.disconnectHandler),
		nc.ReconnectHandler(config.reconnectHandler),
	}

	// Add custom NATS options
	opts = append(opts, config.natsOptions...)

	// Authentication
	if config.token != "" {
		opts = append(opts, nc.Token(config.token))
	} else if config.username != "" && config.password != "" {
		opts = append(opts, nc.UserInfo(config.username, config.password))
	}

	return opts
}
