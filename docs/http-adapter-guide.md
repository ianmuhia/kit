# Production-Ready HTTP Adapter with Huma

This guide covers the production-ready HTTP adapter template that uses Huma v2 for building robust REST APIs.

## Features

### 🚀 Production-Ready Features

1. **Clean API Handler**
   - No middleware conflicts - middleware is configured at root level
   - Focused on API endpoint registration and business logic
   - Structured logging with slog
   - Functional options pattern for configuration

2. **Advanced Validation**
   - Built-in Huma validation (min/max length, patterns, ranges)
   - Custom validation logic support
   - Detailed error responses with field-level errors
   - Pattern descriptions for user-friendly error messages

3. **Flexible API Design**
   - RESTful endpoints (GET, POST, PUT, PATCH, DELETE)
   - Versioned API paths (`/api/v1`)
   - Customizable path prefixes with `RegisterWithPrefix()`
   - Health check endpoint

4. **Rich Request/Response Models**
   - Comprehensive pagination with HATEOAS links
   - Field filtering support
   - Soft delete support
   - Optimistic locking with versioning
   - Metadata fields for extensibility

5. **Error Handling**
   - Domain-specific error mapping to HTTP status codes
   - Structured error responses
   - Detailed logging without exposing internals to clients
   - Operation-specific error handling

6. **Documentation**
   - Auto-generated OpenAPI 3.1 specification
   - Interactive API documentation
   - Comprehensive field documentation
   - Examples for all endpoints

## Architecture

The adapter is designed to be **middleware-agnostic**. All middleware (CORS, rate limiting, authentication, logging, etc.) should be configured at the application root level to:

- Avoid duplication and conflicts
- Provide consistent behavior across all APIs
- Allow centralized middleware management
- Make testing easier

## Usage

### Basic Setup with Root-Level Middleware

```go
package main

import (
    "log/slog"
    "net/http"
    "os"
    
    "github.com/danielgtaylor/huma/v2"
    "github.com/danielgtaylor/huma/v2/adapters/humachi"
    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
    
    "yourapp/internal/user/adapters"
    "yourapp/internal/user/app"
)

func main() {
    // Setup logger
    logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
    
    // Create router with middleware (configured once at root level)
    router := chi.NewMux()
    router.Use(middleware.RequestID)
    router.Use(middleware.RealIP)
    router.Use(middleware.Logger)
    router.Use(middleware.Recoverer)
    router.Use(middleware.Compress(5))
    
    // Optional: Add your custom middleware here
    // router.Use(corsMiddleware)
    // router.Use(authMiddleware)
    // router.Use(rateLimitMiddleware)
    
    // Create Huma API
    config := huma.DefaultConfig("My API", "1.0.0")
    humaAPI := humachi.New(router, config)
    
    // Initialize service
    service := app.NewService(/* dependencies */)
    
    // Create and register API handler
    api := adapters.NewUserAPI(
        service,
        adapters.WithLogger(logger),
    )
    
    api.Register(humaAPI)              // Uses /api/v1 prefix
    api.RegisterHealthCheck(humaAPI)   // Adds /health endpoint
    
    // Start server
    http.ListenAndServe(":8080", router)
}
```

## API Endpoints

### Resource Endpoints

| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/api/v1/{domain}s` | Create a new resource |
| `GET` | `/api/v1/{domain}s/{id}` | Get resource by ID |
| `GET` | `/api/v1/{domain}s` | List resources with pagination |
| `PUT` | `/api/v1/{domain}s/{id}` | Update resource (full replacement) |
| `PATCH` | `/api/v1/{domain}s/{id}` | Partially update resource |
| `DELETE` | `/api/v1/{domain}s/{id}` | Delete resource (soft delete by default) |

### System Endpoints

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/health` | Health check endpoint |
| `GET` | `/docs` | Interactive API documentation |
| `GET` | `/openapi.json` | OpenAPI 3.1 specification |

## Request Examples

### Create Resource

```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "description": "A test user",
    "active": true,
    "metadata": "{\"role\": \"admin\"}"
  }'
```

### List Resources with Filtering and Pagination

```bash
curl "http://localhost:8080/api/v1/users?page=1&page_size=20&active=true&sort_by=created_at&sort_order=desc&search=john"
```

### Partial Update (PATCH)

```bash
curl -X PATCH http://localhost:8080/api/v1/users/123 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Jane Doe"
  }'
```

### Get with Field Selection

```bash
curl "http://localhost:8080/api/v1/users/123?fields=id,name,created_at"
```

## Response Examples

### Single Resource Response

```json
{
  "id": 123,
  "name": "John Doe",
  "description": "A test user",
  "active": true,
  "metadata": "{\"role\": \"admin\"}",
  "created_at": "2024-01-01T12:00:00Z",
  "updated_at": "2024-01-01T12:00:00Z",
  "deleted_at": null,
  "version": 1
}
```

### List Response with Pagination

```json
{
  "items": [
    {
      "id": 123,
      "name": "John Doe",
      "description": "A test user",
      "active": true,
      "created_at": "2024-01-01T12:00:00Z",
      "updated_at": "2024-01-01T12:00:00Z"
    }
  ],
  "pagination": {
    "total": 100,
    "page": 1,
    "page_size": 20,
    "total_pages": 5,
    "has_next": true,
    "has_previous": false,
    "next_page": 2,
    "_links": {
      "self": "/api/v1/users?page=1&page_size=20",
      "first": "/api/v1/users?page=1&page_size=20",
      "last": "/api/v1/users?page=5&page_size=20",
      "next": "/api/v1/users?page=2&page_size=20"
    }
  }
}
```

### Error Response

```json
{
  "status": 400,
  "title": "Bad Request",
  "detail": "Validation failed",
  "errors": [
    {
      "location": "body.name",
      "message": "name must be at least 3 characters",
      "value": "ab"
    }
  ]
}
```

## Configuration Options

### APIConfig

```go
type APIConfig struct {
    // EnableMetrics enables request metrics collection
    EnableMetrics bool
    
    // RateLimitPerMinute sets the rate limit for API requests (0 = unlimited)
    RateLimitPerMinute int
    
    // RequestTimeout sets the maximum duration for a request
    RequestTimeout time.Duration
    
    // EnableCORS enables CORS middleware
    EnableCORS bool
    
    // AllowedOrigins lists allowed CORS origins (empty = all)
    AllowedOrigins []string
}
```

## Middleware

The template includes several built-in middleware:

1. **Request ID**: Adds unique request ID to each request
2. **Real IP**: Extracts real client IP from headers
3. **Logging**: Structured logging for each request
4. **Recovery**: Panic recovery with error logging
5. **CORS**: Configurable CORS support
6. **Timeout**: Request timeout handling
7. **Rate Limiting**: Basic rate limiting (upgrade to distributed for production)

### Adding Custom Middleware

```go
router := api.SetupRouter()

// Add custom middleware
router.Use(yourCustomMiddleware)

// Or add to specific routes
router.Group(func(r chi.Router) {
    r.Use(authenticationMiddleware)
    
    humaAPI := humachi.New(r, config)
    api.Register(humaAPI)
})
```

## Validation

### Built-in Validation Tags

The template uses Huma's validation tags:

```go
type CreateInput struct {
    Body struct {
        Name string `json:"name" minLength:"3" maxLength:"100" pattern:"^[a-zA-Z0-9\\s\\-_]+$"`
        Email string `json:"email" format:"email"`
        Age int `json:"age" minimum:"18" maximum:"120"`
    }
}
```

### Custom Validation

Implement custom validation in the handler:

```go
func (api *UserAPI) validateCreateInput(input *CreateUserInput) error {
    // Custom validation logic
    if isReservedName(input.Body.Name) {
        return fmt.Errorf("name '%s' is reserved", input.Body.Name)
    }
    return nil
}
```

## Error Handling

### Domain Error Mapping

The template maps domain errors to HTTP status codes:

```go
func (api *API) handleError(err error, operation string) error {
    switch {
    case err == domain.ErrNotFound:
        return huma.Error404NotFound("Resource not found")
    case err == domain.ErrAlreadyExists:
        return huma.Error409Conflict("Resource already exists", err)
    case err == domain.ErrValidationFailed:
        return huma.Error422UnprocessableEntity("Validation failed", err)
    default:
        return huma.Error500InternalServerError("Internal error")
    }
}
```

## Deployment Checklist

- [ ] Configure appropriate rate limits
- [ ] Set up distributed rate limiting (Redis, etc.)
- [ ] Configure CORS origins for your domain
- [ ] Set appropriate request timeouts
- [ ] Enable production logging (JSON format)
- [ ] Set up health checks in load balancer
- [ ] Configure TLS/HTTPS
- [ ] Set up request ID propagation to downstream services
- [ ] Configure connection pooling for database
- [ ] Set up monitoring and alerting
- [ ] Enable request/response compression
- [ ] Configure appropriate max header/body sizes

## Best Practices

1. **Use Structured Logging**: Always use slog with appropriate context
2. **Handle Errors Gracefully**: Map domain errors, don't expose internal details
3. **Version Your API**: Use path versioning (`/api/v1`)
4. **Implement Health Checks**: Use for load balancer and monitoring
5. **Add Request IDs**: Essential for distributed tracing
6. **Use PATCH for Partial Updates**: Don't require full objects for updates
7. **Implement Pagination**: Always paginate list endpoints
8. **Add HATEOAS Links**: Help clients navigate the API
9. **Document Everything**: Use Huma's doc tags extensively
10. **Validate Input**: Use both struct tags and custom validation

## Testing

```go
func TestCreateUser(t *testing.T) {
    // Setup
    service := &mockService{}
    api := adapters.NewUserAPI(service)
    router := api.SetupRouter()
    humaAPI := humachi.New(router, huma.DefaultConfig("Test", "1.0.0"))
    api.Register(humaAPI)
    
    // Test
    w := httptest.NewRecorder()
    body := `{"name":"John Doe","description":"Test user"}`
    req := httptest.NewRequest("POST", "/api/v1/users", strings.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    
    router.ServeHTTP(w, req)
    
    // Assert
    assert.Equal(t, http.StatusCreated, w.Code)
}
```

## Further Reading

- [Huma Documentation](https://huma.rocks/)
- [Chi Router Documentation](https://go-chi.io/)
- [OpenAPI 3.1 Specification](https://spec.openapis.org/oas/v3.1.0)
- [REST API Best Practices](https://restfulapi.net/)
- [HTTP Status Codes](https://httpstatuses.com/)
