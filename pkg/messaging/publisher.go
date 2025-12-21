package messaging

import (
	"fmt"
	"log/slog"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-nats/v2/pkg/nats"
	"github.com/ThreeDotsLabs/watermill/message"
	nc "github.com/nats-io/nats.go"
)

// PublisherConfig holds essential configuration for NATS JetStream publisher.
type PublisherConfig struct {
	url              string
	name             string
	token            string
	username         string
	password         string
	logger           *slog.Logger
	marshaler        nats.Marshaler
	autoProvision    bool
	maxReconnects    int
	natsOptions      []nc.Option
	disconnectHandler func(*nc.Conn, error)
	reconnectHandler func(*nc.Conn)
}

// PublisherOption is a functional option for configuring the publisher.
type PublisherOption func(*PublisherConfig)

// WithURL sets the NATS server URL.
func WithURL(url string) PublisherOption {
	return func(c *PublisherConfig) {
		c.url = url
	}
}

// WithName sets the client name.
func WithName(name string) PublisherOption {
	return func(c *PublisherConfig) {
		c.name = name
	}
}

// WithToken sets the authentication token.
func WithToken(token string) PublisherOption {
	return func(c *PublisherConfig) {
		c.token = token
	}
}

// WithUserPassword sets username and password authentication.
func WithUserPassword(username, password string) PublisherOption {
	return func(c *PublisherConfig) {
		c.username = username
		c.password = password
	}
}

// WithLogger sets the logger.
func WithLogger(logger *slog.Logger) PublisherOption {
	return func(c *PublisherConfig) {
		c.logger = logger
	}
}

// WithMarshaler sets a custom marshaler (default is GobMarshaler).
func WithMarshaler(marshaler nats.Marshaler) PublisherOption {
	return func(c *PublisherConfig) {
		c.marshaler = marshaler
	}
}

// WithAutoProvision enables/disables JetStream auto-provisioning.
func WithAutoProvision(enable bool) PublisherOption {
	return func(c *PublisherConfig) {
		c.autoProvision = enable
	}
}

// WithMaxReconnects sets the maximum number of reconnection attempts.
// Use -1 for infinite reconnects (default).
func WithMaxReconnects(max int) PublisherOption {
	return func(c *PublisherConfig) {
		c.maxReconnects = max
	}
}

// WithNATSOptions adds custom NATS connection options.
func WithNATSOptions(opts ...nc.Option) PublisherOption {
	return func(c *PublisherConfig) {
		c.natsOptions = append(c.natsOptions, opts...)
	}
}

// WithDisconnectHandler sets a custom disconnect handler.
func WithDisconnectHandler(handler func(*nc.Conn, error)) PublisherOption {
	return func(c *PublisherConfig) {
		c.disconnectHandler = handler
	}
}

// WithReconnectHandler sets a custom reconnect handler.
func WithReconnectHandler(handler func(*nc.Conn)) PublisherOption {
	return func(c *PublisherConfig) {
		c.reconnectHandler = handler
	}
}

// defaultPublisherConfig returns sensible defaults.
func defaultPublisherConfig() *PublisherConfig {
	return &PublisherConfig{
		url:           "nats://localhost:4222",
		name:          "watermill-publisher",
		logger:        slog.Default(),
		marshaler:     &nats.GobMarshaler{},
		autoProvision: true,
		maxReconnects: -1,
		disconnectHandler: func(conn *nc.Conn, err error) {
			if err != nil {
				slog.Error("NATS disconnected", "error", err)
			}
		},
		reconnectHandler: func(conn *nc.Conn) {
			slog.Info("NATS reconnected", "url", conn.ConnectedUrl())
		},
	}
}

// NewPublisher creates a NATS JetStream publisher with production-ready settings.
func NewPublisher(opts ...PublisherOption) (message.Publisher, error) {
	config := defaultPublisherConfig()
	
	// Apply all options
	for _, opt := range opts {
		opt(config)
	}

	if config.url == "" {
		return nil, fmt.Errorf("NATS URL is required")
	}

	natsOpts := buildNATSOptions(config)

	publisher, err := nats.NewPublisher(
		nats.PublisherConfig{
			URL:         config.url,
			NatsOptions: natsOpts,
			Marshaler:   config.marshaler,
			JetStream: nats.JetStreamConfig{
				AutoProvision: config.autoProvision,
			},
		},
		watermill.NewSlogLogger(config.logger),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create publisher: %w", err)
	}

	return publisher, nil
}

// buildNATSOptions constructs NATS connection options.
func buildNATSOptions(config *PublisherConfig) []nc.Option {
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
