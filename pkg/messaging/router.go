package messaging

import (
	"context"
	"log/slog"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
)

const (
	defaultMaxRetries = 3
	defaultTimeout    = 30 * time.Second
)

// Router wraps Watermill router with our dependencies.
type Router struct {
	router     *message.Router
	publisher  message.Publisher
	subscriber message.Subscriber
	logger     *slog.Logger
}

// NewRouter creates a new message router with all handlers registered.
func NewRouter(
	publisher message.Publisher,
	subscriber message.Subscriber,
	logger *slog.Logger,
) (*Router, error) {
	router, err := message.NewRouter(message.RouterConfig{}, watermill.NewSlogLogger(logger))
	if err != nil {
		return nil, err
	}

	poisonQueue, err := middleware.PoisonQueue(publisher, "poison_queue")
	if err != nil {
		return nil, err
	}

	router.AddMiddleware(
		middleware.Recoverer,
		middleware.Retry{
			MaxRetries:      defaultMaxRetries,
			InitialInterval: time.Second,
			MaxInterval:     defaultTimeout,
		}.Middleware,
		middleware.CorrelationID,
		poisonQueue,
		middleware.Timeout(defaultTimeout),
		middleware.InstantAck,
	)

	r := &Router{
		router:     router,
		publisher:  publisher,
		subscriber: subscriber,
		logger:     logger,
	}

	return r, nil
}

// Run starts the router (blocking).
func (r *Router) Run(ctx context.Context) error {
	r.logger.Info("Starting message router")
	return r.router.Run(ctx)
}

// Close gracefully shuts down the router.
func (r *Router) Close() error {
	r.logger.Info("Closing message router")
	return r.router.Close()
}

// GetRouter returns the underlying watermill router for direct handler registration
func (r *Router) GetRouter() *message.Router {
	return r.router
}

// RegisterHandler registers a single event handler with the router.
func (r *Router) RegisterHandler(name, topic string, handler message.NoPublishHandlerFunc) {
	r.router.AddConsumerHandler(
		name+"_handler",
		topic,
		r.subscriber,
		handler,
	)
}

// RegisterHandlerFunc is an alias for convenience (accepts func(*message.Message) error)
func (r *Router) RegisterHandlerFunc(name, topic string, handler func(*message.Message) error) {
	r.RegisterHandler(name, topic, message.NoPublishHandlerFunc(handler))
}

// RegisterDomainHandlers is a helper to register all handlers for a domain.
func (r *Router) RegisterDomainHandlers(handlers map[string]message.NoPublishHandlerFunc) {
	for topic, handler := range handlers {
		r.RegisterHandler(topic, topic, handler)
	}
}
