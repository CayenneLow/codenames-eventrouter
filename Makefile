build_server:
	docker compose build

start_server:
	docker compose up -d

stop_server:
	docker compose down

build:
	go build -o bin/start_event_router ./cmd/event_router/main.go

test_all: start_server
	go test -timeout 5s -v -cover ./...
	docker compose down