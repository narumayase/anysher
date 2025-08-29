# anysher - Kafka and HTTP integration lib

este proyecto provee una librería para enviar mensajitos de manera configurable, puede enviar mensajes json a kafka o tmb por post http.

## Features

- 

### Prerequisites

- Go 1.21 or higher
- Kafka (optional, for Kafka integration)

## 🚀 Installation

1. Install dependencies:

```bash
go mod tidy
```

2. Configure environment variables:

```bash
cp env.example .env
# Edit .env with the values described below.
```

3. Run the application:

```bash
go run main.go
```

## 🔧 Configuration

### Environment Variables

Create a `.env` file based on `env.example`:

- `API_ENDPOINT`: API endpoint to send payload (default: https://api.groq.com/openai/v1/responses)
- `PORT`: Server port (default: 8080)
- `LOG_LEVEL`: Log level (debug, info, warn, error, fatal, panic - default: info)
- `KAFKA_BROKER`: Comma-separated list of Kafka brokers
- `KAFKA_TOPIC`: Kafka topic to send events to

### Usage

//TODO 

## 🎗️ Architecture

This project follows Clean Architecture principles:

- **Domain**: Entities, repository interfaces, and use cases
- **Application**: Implementation of use cases
- **Infrastructure**: OpenAI and Groq repository implementations
- **Interfaces**: HTTP controllers and routers

## 📁 Project Structure

```
anysher/
├── internal/             # Project-specific code
│   ├── infrastructure/   # Repository implementations
├── go.mod                # Go dependencies
├── README_ES.md          # README in spanish
└── README.md             # This file
```

## 🧪 Testing

### Running Tests

To run all tests:

```bash
go test ./...
```

### Test Coverage

To check test coverage (excluding mocks):

```bash
# Generate coverage report
go test -coverprofile=coverage.out ./...

# View coverage report in terminal
go tool cover -func=coverage.out

# Generate HTML coverage report
go tool cover -html=coverage.out -o coverage.html

# View coverage excluding mocks
go test -coverprofile=coverage.out ./... && \
go tool cover -func=coverage.out | grep -v "mocks"
```

### Running Benchmarks

```bash
go test -bench=. ./...
```

## BackLog

- [x] Unit Tests
