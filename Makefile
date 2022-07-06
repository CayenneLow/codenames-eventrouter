start_server:
	go run ./cmd/event_router/main.go

build:
	go build -o bin/start_event_router ./cmd/event_router/main.go