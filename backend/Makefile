include ../.envrc

MIGRATIONS_PATH = ./migrate/migrations

migration:
	@migrate create -seq -ext sql -dir $(MIGRATIONS_PATH) $(filter-out $@,$(MAKECMDGOALS))

migrate-up:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_ADDR) up

migrate-down:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_ADDR) down $(filter-out $@,$(MAKECMDGOALS))

migrate-revert:
	@migrate -path=$(MIGRATIONS_PATH)  -database=$(DB_ADDR)  force 10

generate-proto:
	@echo "Generating Go gRPC code for job.proto..."
	protoc --go_out=. --go-grpc_out=. protos/job.proto
	@echo "✅ Done generating gRPC code."

.PHONY: migration migrate-up migrate-down migrate-revert generate-proto
