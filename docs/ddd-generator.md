# DDD Generator Documentation

A comprehensive Domain-Driven Design (DDD) code generator for Go projects following hexagonal architecture principles.

## Table of Contents

- [Quick Start](#quick-start)
- [Installation](#installation)
- [Usage](#usage)
- [Generated Structure](#generated-structure)
- [Examples](#examples)
- [Best Practices](#best-practices)

## Quick Start

```bash
# Install the generator
go install github.com/ianmuhia/kit/cmd/ddd-gen@latest

# Generate a basic domain
ddd-gen --domain=user

# Generate with all features
ddd-gen --domain=order --all
```

## Installation

### From Source

```bash
git clone https://github.com/ianmuhia/kit.git
cd kit
make install-ddd-gen
```

### Using Go Install

```bash
go install github.com/ianmuhia/kit/cmd/ddd-gen@latest
```

## Usage

### Command-Line Flags

| Flag | Alias | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--domain` | `-d` | string | *required* | Domain name (e.g., `booking`, `user`, `order`) |
| `--output` | `-o` | string | `./internal` | Output directory for generated code |
| `--with-tests` | `-t` | bool | `false` | Generate test files |
| `--with-messaging` | `-m` | bool | `false` | Generate messaging/pub-sub adapter |
| `--with-river` | `-r` | bool | `false` | Generate River job queue adapter |
| `--with-cqrs` | `-c` | bool | `false` | Generate CQRS components (Watermill) |
| `--with-workflows` | `-w` | bool | `false` | Generate Temporal workflow adapter |
| `--with-decorators` | | bool | `false` | Generate service decorators |
| `--all` | | bool | `false` | Generate all optional components |

### Basic Examples

```bash
# Minimal domain (entity + repository + service + HTTP + Postgres)
ddd-gen -d booking

# Domain with tests
ddd-gen -d booking -t

# Domain with CQRS
ddd-gen -d order -c

# Domain with everything
ddd-gen -d payment --all
```

## Generated Structure

### Minimal Generation

```
internal/booking/
├── booking.go              # Domain entity (aggregate root)
├── repository.go           # Repository interface (port)
├── errors.go              # Domain-specific errors
├── events.go              # Domain events
├── validation.go          # Domain validation rules
├── app/
│   └── service.go         # Application service
└── adapters/
    ├── booking_http.go    # HTTP handlers (adapter)
    └── booking_postgres.go # PostgreSQL repository (adapter)
```

### Full Generation (--all)

```
internal/booking/
├── booking.go
├── repository.go
├── errors.go
├── events.go
├── validation.go
├── app/
│   ├── service.go
│   ├── service_test.go         # Unit tests
│   ├── decorators.go           # Service decorators
│   └── wiring_example.go       # Dependency injection example
├── adapters/
│   ├── booking_http.go
│   ├── booking_postgres.go
│   ├── booking_messaging.go    # Pub/sub event handlers
│   ├── booking_river.go        # River job queue integration
│   └── booking_temporal.go     # Temporal workflows
└── cqrs/
    ├── commands.go             # Command definitions
    ├── command_handlers.go     # Command handlers
    ├── events.go              # CQRS events
    ├── event_handlers.go      # Event handlers
    └── wiring.go              # Watermill CQRS configuration
```

## Examples

### Example 1: E-commerce Order Domain

```bash
ddd-gen -d order --with-cqrs --with-river --with-workflows
```

**Use Case:** Order processing with:

- CQRS for command/query separation
- River for background job processing (order fulfillment)
- Temporal for order workflow orchestration

**Generated Components:**

- Order entity with domain logic
- CQRS commands (CreateOrder, UpdateOrder, CancelOrder)
- Event handlers for order state changes
- River jobs for email notifications, inventory updates
- Temporal workflow for order fulfillment saga

### Example 2: User Authentication Domain

```bash
ddd-gen -d user --with-decorators --with-tests
```

**Use Case:** User management with:

- Service decorators for audit logging and metrics
- Comprehensive test coverage

**Generated Components:**

- User entity with validation
- Service with decorators (logging, metrics, caching)
- HTTP endpoints for user CRUD
- PostgreSQL repository
- Unit tests for service layer

### Example 3: Notification Domain

```bash
ddd-gen -d notification --with-messaging --with-river
```

**Use Case:** Async notification system with:

- Event-driven architecture via messaging
- Background job processing

**Generated Components:**

- Notification entity
- Pub/sub event handlers (Watermill)
- River jobs for email/SMS delivery
- Retry logic and dead letter handling

### Example 4: Payment Domain

```bash
ddd-gen -d payment --all
```

**Use Case:** Complete payment processing system with all features.

## Best Practices

### 1. Domain Layer

**DO:**

- Keep domain logic pure and framework-agnostic
- Use domain events for cross-aggregate communication
- Implement validation in the domain layer
- Make entities immutable where possible

**DON'T:**

- Put infrastructure concerns in domain layer
- Depend on external packages in domain entities
- Use database-specific types in domain

### 2. Application Layer

**DO:**

- Coordinate between domain and infrastructure
- Handle transactions in service methods
- Use decorators for cross-cutting concerns
- Return domain errors from services

**DON'T:**

- Put business logic in services (it belongs in domain)
- Access database directly (use repositories)

### 3. Adapter Layer

**DO:**

- Implement repository interfaces defined in domain
- Handle infrastructure-specific concerns
- Map between domain models and DTOs
- Validate input at HTTP boundaries

**DON'T:**

- Pass DTOs to domain layer
- Put business logic in adapters

### 4. CQRS (Optional)

**DO:**

- Use commands for writes, queries for reads
- Keep command handlers focused and simple
- Use events for eventual consistency
- Separate read and write models when complexity justifies it

**DON'T:**

- Over-engineer simple CRUD operations
- Create unnecessary complexity with CQRS

### 5. Testing

**DO:**

- Test domain logic thoroughly (unit tests)
- Use table-driven tests for multiple scenarios
- Mock external dependencies in tests
- Test error cases and edge cases

**DON'T:**

- Test implementation details
- Skip integration tests for adapters

## Customization

After generation, you should customize:

1. **Entity Logic**: Add your business rules in `{domain}.go`
2. **Repository Methods**: Extend repository interface with domain-specific queries
3. **Service Methods**: Implement use cases in `app/service.go`
4. **HTTP Routes**: Wire up HTTP handlers in your main application
5. **Validation**: Add domain-specific validation rules
6. **Events**: Define domain events for state changes

## Dependencies

The generated code may require:

```bash
# Core dependencies (always needed)
go get github.com/lib/pq  # PostgreSQL driver

# Optional dependencies
go get github.com/ThreeDotsLabs/watermill  # Messaging/CQRS
go get github.com/riverqueue/river         # Job queue
go get go.temporal.io/sdk                   # Workflows
```

## Troubleshooting

### Issue: Import errors after generation

**Solution:** Run `go mod tidy` in your project root.

### Issue: Templates not found

**Solution:** Ensure you're using the installed version, not running from source without embedded templates.

### Issue: Generated code doesn't compile

**Solution:** Check that all required dependencies are installed and run `go mod tidy`.

## Additional Resources

- [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/)
- [Domain-Driven Design](https://martinfowler.com/bliki/DomainDrivenDesign.html)
- [CQRS Pattern](https://martinfowler.com/bliki/CQRS.html)
- [Event Sourcing](https://martinfowler.com/eaaDev/EventSourcing.html)

## Support

For issues and questions:

- GitHub Issues: [https://github.com/ianmuhia/kit/issues](https://github.com/ianmuhia/kit/issues)
- Discussions: [https://github.com/ianmuhia/kit/discussions](https://github.com/ianmuhia/kit/discussions)
