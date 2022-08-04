build_server:
	docker compose build

start_server:
	docker compose up -d

stop_server:
	docker compose down

build:
	go build -o bin/start_event_router ./cmd/event_router/main.go

# Error codes: Success (0), Failure (1)
test_all: 
	docker compose -f docker-compose.tests.yml down
	docker compose -f docker-compose.tests.yml up init-redis
	docker compose -f docker-compose.tests.yml up event_router_integration_test --build --exit-code-from event_router_integration_test
	docker compose -f docker-compose.tests.yml down

# Error codes: Success (0), Failure (1)
test_all_debug: 
	docker compose -f docker-compose.tests.yml down
	docker compose -f docker-compose.tests.yml up event_router_integration_test redis-insight-it init-redis -d --build

docker_down:
	docker compose -f docker-compose.yml -f docker-compose.tests.yml down