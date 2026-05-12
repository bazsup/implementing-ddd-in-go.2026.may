# Implementing Domain-Driven Design in Go

A hands-on workshop project for learning and applying tactical Domain-Driven Design patterns in Go. The project evolves session by session — starting from a minimal HTTP server and growing into a full DDD application with aggregates, value objects, domain events, and repositories.

## Requirements

- Go 1.26+

## Installation

```bash
git clone <repository-url>
cd implementing-ddd-in-go
go mod tidy
```

## Environment Variables

Copy the example and adjust as needed:

```bash
cp .env.example .env
```

| Variable | Default | Description          |
|----------|---------|----------------------|
| `PORT`   | `8080`  | HTTP server port     |

## Development

```bash
go run ./cmd
```

## Build & Run

```bash
go build -o server ./cmd
./server
```

## Tests

```bash
go test ./...
```
