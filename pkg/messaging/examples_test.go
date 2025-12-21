package messaging_test

import (
	"log/slog"

	"github.com/ThreeDotsLabs/watermill-nats/v2/pkg/nats"
	"github.com/ianmuhia/kit/pkg/messaging"
)

func ExampleNewPublisher_basic() {
	// Create a publisher with default settings
	publisher, err := messaging.NewPublisher()
	if err != nil {
		panic(err)
	}
	defer publisher.Close()

	// Use the publisher...
}

func ExampleNewPublisher_withURL() {
	// Create a publisher with custom URL
	publisher, err := messaging.NewPublisher(
		messaging.WithURL("nats://production.example.com:4222"),
	)
	if err != nil {
		panic(err)
	}
	defer publisher.Close()
}

func ExampleNewPublisher_withAuth() {
	// Create a publisher with token authentication
	publisher, err := messaging.NewPublisher(
		messaging.WithURL("nats://secure.example.com:4222"),
		messaging.WithToken("my-secret-token"),
		messaging.WithName("my-service-publisher"),
	)
	if err != nil {
		panic(err)
	}
	defer publisher.Close()
}

func ExampleNewPublisher_withUserPassword() {
	// Create a publisher with username/password auth
	publisher, err := messaging.NewPublisher(
		messaging.WithURL("nats://secure.example.com:4222"),
		messaging.WithUserPassword("username", "password"),
	)
	if err != nil {
		panic(err)
	}
	defer publisher.Close()
}

func ExampleNewPublisher_withLogger() {
	// Create a publisher with custom logger
	logger := slog.Default()

	publisher, err := messaging.NewPublisher(
		messaging.WithURL("nats://localhost:4222"),
		messaging.WithLogger(logger),
	)
	if err != nil {
		panic(err)
	}
	defer publisher.Close()
}

func ExampleNewPublisher_withJSONMarshaler() {
	// Create a publisher with JSON marshaler instead of Gob
	publisher, err := messaging.NewPublisher(
		messaging.WithURL("nats://localhost:4222"),
		messaging.WithMarshaler(&nats.JSONMarshaler{}),
	)
	if err != nil {
		panic(err)
	}
	defer publisher.Close()
}

func ExampleNewPublisher_full() {
	// Create a fully configured publisher
	logger := slog.Default()

	publisher, err := messaging.NewPublisher(
		messaging.WithURL("nats://production.example.com:4222"),
		messaging.WithName("payment-service-publisher"),
		messaging.WithToken("prod-token"),
		messaging.WithLogger(logger),
		messaging.WithMarshaler(&nats.JSONMarshaler{}),
		messaging.WithAutoProvision(true),
		messaging.WithMaxReconnects(-1), // infinite

	)
	if err != nil {
		panic(err)
	}
	defer publisher.Close()
}

func ExampleNewPublisher_development() {
	// Simple setup for local development
	publisher, err := messaging.NewPublisher(
		messaging.WithName("dev-publisher"),
	)
	if err != nil {
		panic(err)
	}
	defer publisher.Close()
}

func ExampleNewPublisher_production() {
	// Production-ready setup with auth and custom settings
	logger := slog.Default()

	publisher, err := messaging.NewPublisher(
		messaging.WithURL("nats://nats-cluster.prod:4222"),
		messaging.WithName("order-service-publisher"),
		messaging.WithToken("secure-prod-token"),
		messaging.WithLogger(logger),
		messaging.WithMarshaler(&nats.JSONMarshaler{}),
		messaging.WithMaxReconnects(-1),
		messaging.WithAutoProvision(true),
	)
	if err != nil {
		panic(err)
	}
	defer publisher.Close()
}

// Subscriber Examples

func ExampleNewSubscriber_basic() {
	// Create a subscriber with default settings
	subscriber, err := messaging.NewSubscriber()
	if err != nil {
		panic(err)
	}
	defer subscriber.Close()

	// Use the subscriber...
}

func ExampleNewSubscriber_withURL() {
	// Create a subscriber with custom URL and durable prefix
	subscriber, err := messaging.NewSubscriber(
		messaging.WithSubscriberURL("nats://production.example.com:4222"),
		messaging.WithDurablePrefix("payment-service"),
	)
	if err != nil {
		panic(err)
	}
	defer subscriber.Close()
}

func ExampleNewSubscriber_withAuth() {
	// Create a subscriber with token authentication
	subscriber, err := messaging.NewSubscriber(
		messaging.WithSubscriberURL("nats://secure.example.com:4222"),
		messaging.WithSubscriberToken("my-secret-token"),
		messaging.WithSubscriberName("my-service-subscriber"),
		messaging.WithDurablePrefix("my-service"),
	)
	if err != nil {
		panic(err)
	}
	defer subscriber.Close()
}

func ExampleNewSubscriber_withJSONMarshaler() {
	// Create a subscriber with JSON unmarshaler instead of Gob
	subscriber, err := messaging.NewSubscriber(
		messaging.WithSubscriberURL("nats://localhost:4222"),
		messaging.WithUnmarshaler(&nats.JSONMarshaler{}),
		messaging.WithDurablePrefix("json-service"),
	)
	if err != nil {
		panic(err)
	}
	defer subscriber.Close()
}

func ExampleNewSubscriber_full() {
	// Create a fully configured subscriber
	logger := slog.Default()

	subscriber, err := messaging.NewSubscriber(
		messaging.WithSubscriberURL("nats://production.example.com:4222"),
		messaging.WithSubscriberName("order-service-subscriber"),
		messaging.WithDurablePrefix("order-service"),
		messaging.WithSubscriberToken("prod-token"),
		messaging.WithSubscriberLogger(logger),
		messaging.WithUnmarshaler(&nats.JSONMarshaler{}),
		messaging.WithSubscriberAutoProvision(true),
		messaging.WithSubscriberMaxReconnects(-1), // infinite

	)
	if err != nil {
		panic(err)
	}
	defer subscriber.Close()
}

func ExampleNewSubscriber_development() {
	// Simple setup for local development
	subscriber, err := messaging.NewSubscriber(
		messaging.WithSubscriberName("dev-subscriber"),
		messaging.WithDurablePrefix("dev"),
	)
	if err != nil {
		panic(err)
	}
	defer subscriber.Close()
}

func ExampleNewSubscriber_production() {
	// Production-ready setup with auth and custom settings
	logger := slog.Default()

	subscriber, err := messaging.NewSubscriber(
		messaging.WithSubscriberURL("nats://nats-cluster.prod:4222"),
		messaging.WithSubscriberName("notification-service-subscriber"),
		messaging.WithDurablePrefix("notification-service"),
		messaging.WithSubscriberToken("secure-prod-token"),
		messaging.WithSubscriberLogger(logger),
		messaging.WithUnmarshaler(&nats.JSONMarshaler{}),
		messaging.WithSubscriberMaxReconnects(-1),
		messaging.WithSubscriberAutoProvision(true),
	)
	if err != nil {
		panic(err)
	}
	defer subscriber.Close()
}
