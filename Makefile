start_server:
	docker compose up -d --build

build:
	go build -o bin/start_event_router ./cmd/event_router/main.go