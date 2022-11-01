include .env
export
APP_DSN := mysql://${DB_USER}:${DB_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/${DB_DATABASE}?charset=utf8mb4&parseTime=True&loc=Local
MIGRATE := docker run --rm -v $(shell pwd)/migrations:/migrations --name consume-fcm_migrate --network consume-fcm_backend migrate/migrate -path=/migrations/ -database "$(APP_DSN)"

run:
	@echo "Running go..."
	@docker-compose up

down:
	@echo "Shutdown docker container"
	@docker-compose down

.PHONY: migrate
migrate: ## run all new database migrations
	@echo "Running all new database migrations..."
	@$(MIGRATE) up

.PHONY: migrate-down
migrate-down: ## revert database to the last migration step
	@echo "Reverting database to the last migration step..."
	@$(MIGRATE) down 1

.PHONY: migrate-new
migrate-new: ## create a new database migration
	@read -p "Enter the name of the new migration: " name; \
	$(MIGRATE) create -ext sql -dir /migrations/ $${name// /_}
