FROM golang:1.24.0-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/restservice ./cmd/app
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

FROM alpine:3.20
WORKDIR /app
RUN apk add --no-cache postgresql-client
COPY --from=builder /app/restservice /app/restservice
COPY --from=builder /go/bin/goose /usr/local/bin/goose
COPY migrations /app/migrations
COPY Swagger /app/Swagger
EXPOSE 8080
CMD ["/app/restservice"]