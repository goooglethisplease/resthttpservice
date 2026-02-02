# Build stage
FROM golang:1.24.0-alpine AS builder

WORKDIR /app

# Copy dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/restservice ./cmd/app

# Install goose for migrations
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Runtime stage
FROM alpine:3.20

WORKDIR /app

# Install required tools
RUN apk add --no-cache postgresql-client

# Copy application binary
COPY --from=builder /app/restservice /app/restservice

# Copy goose binary
COPY --from=builder /go/bin/goose /usr/local/bin/goose

# Copy migrations
COPY migrations /app/migrations

# Expose port
EXPOSE 8080

# Default command
CMD ["/app/restservice"]