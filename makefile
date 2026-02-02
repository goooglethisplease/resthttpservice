APP_NAME := restservice
MAIN_PKG := ./cmd/app
BIN_DIR := bin
SWAGGER_DIR := docs

.PHONY: build run tidy swag test docker-up docker-down

build:
	go build -o $(BIN_DIR)/$(APP_NAME) $(MAIN_PKG)

run:
	go run $(MAIN_PKG)

tidy:
	go mod tidy

swag:
	swag init -g cmd/app/main.go -o $(SWAGGER_DIR)

test:
	go test ./...

docker-up:
	docker compose up --build

docker-down:
	docker compose down -v