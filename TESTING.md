# Unit Tests

This project includes comprehensive unit tests for the ePage application.

## Test Coverage

- **58.2%** code coverage across the application

## Test Files

### send_page_test.go
Tests for the Pushover API integration:
- `TestPushoverPayloadStructure` - Validates JSON payload structure and serialization
- `TestPushoverMessageFormat` - Tests message formatting with various input types

### handlers_test.go
Tests for HTTP handlers and template rendering:
- `TestHandleIndexGET` - Tests GET / endpoint returns template
- `TestHandleSendPOST` - Tests POST / with valid and invalid form data
- `TestHandleSendInvalidForm` - Tests error handling for malformed requests
- `TestTemplateCache` - Verifies template caching works correctly
- `TestLoadTemplateNotFound` - Tests handling of missing templates
- `TestLoadTemplateInvalidTemplate` - Tests handling of invalid template syntax
- `TestServerIntegration` - Integration test of full request/response cycle

## Running Tests

Run all tests:
```bash
go test -v ./src
```

Run tests with coverage:
```bash
go test -v -cover ./src
```

Run a specific test:
```bash
go test -v -run TestHandleIndexGET ./src
```

Run tests with detailed coverage report:
```bash
go test -coverprofile=coverage.out ./src
go tool cover -html=coverage.out
```

## Test Design

The tests focus on:
- **Handler logic**: Form validation, routing, and response generation
- **Template handling**: Caching, error handling, and rendering
- **Data structures**: JSON serialization and message formatting
- **Integration**: Full request/response cycles through the HTTP server

The tests use temporary directories and mock templates to avoid file system dependencies and can run in isolation.
