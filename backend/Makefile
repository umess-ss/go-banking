DB_URL_LOCAL=postgres://postgres:postgres@localhost:5432/go_banking?sslmode=disable
DB_URL_DOCKER=postgres://postgres:postgres@localhost:5433/go_banking?sslmode=disable

run:
	go run ./cmd/api

test:
	go test ./...

migrate-up:
	migrate -path migrations -database "$(DB_URL_LOCAL)" up

migrate-down:
	migrate -path migrations -database "$(DB_URL_LOCAL)" down 1

migrate-force:
	migrate -path migrations -database "$(DB_URL_LOCAL)" force $(VERSION)

migrate-version:
	migrate -path migrations -database "$(DB_URL_LOCAL)" version

migrate-up-docker:
	migrate -path migrations -database "$(DB_URL_DOCKER)" up

migrate-down-docker:
	migrate -path migrations -database "$(DB_URL_DOCKER)" down 1

migrate-version-docker:
	migrate -path migrations -database "$(DB_URL_DOCKER)" version

docker-up:
	docker compose up --build

docker-down:
	docker compose down

docker-reset:
	docker compose down -v
	docker compose up --build

docker-logs:
	docker compose logs -f backend

docker-ps:
	docker compose ps