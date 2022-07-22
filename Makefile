build_server:
	docker compose build

start_server:
	docker compose up -d

stop_server:
	docker compose down

build:
	go build -o bin/start_event_router ./cmd/event_router/main.go

test_all:
	# Error codes: Success (0), Failure (1)
	docker compose -f docker-compose.yml -f docker-compose.tests.yml up event_router_integration_test --build --exit-code-from event_router_integration_test
	docker compose -f docker-compose.yml -f docker-compose.tests.yml down

test_all_debug:
	# Error codes: Success (0), Failure (1)
	docker compose -f docker-compose.yml -f docker-compose.tests.yml up event_router_integration_test mongo-express -d --build