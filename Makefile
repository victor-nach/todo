PROTO_DIR=proto
GEN_DIR=$(PROTO_DIR)/gen/go

.PHONY: generate start test start-deps run clean

# Generate proto files
generate:
	if not exist "$(GEN_DIR)" mkdir "$(GEN_DIR)"
	protoc --go_out=$(GEN_DIR) --go-grpc_out=$(GEN_DIR) --proto_path=$(PROTO_DIR) $(PROTO_DIR)/todo.proto

build-run:
	docker-compose up --build

# Stop the services
stop:
	docker-compose down

# Run stopped containers
start:
	docker-compose start

# Start PostgreSQL using docker-compose.dep.yml
start-deps:
	docker-compose -f docker-compose.dep.yml up -d

# Stop PostgreSQL
stop-deps:
	docker-compose -f docker-compose.dep.yml down -v

# Run tests for both services
test:
	@echo "Running tests for api-gateway..."
	cd api-gateway && go test ./...
	@echo "Running tests for internal-service..."
	cd internal-service && go test ./...

