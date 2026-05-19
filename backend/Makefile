DB_URL=postgres://postgres:postgres@localhost:5432/go_banking_api?sslmode=disable

run:
	go run ./cmd/api

migrate-up:
	goose -dir migrations postgres "$(DB_URL)" up

migrate-down:
	goose -dir migrations postgres "$(DB_URL)" down

migrate-status:
	goose -dir migrations postgres "$(DB_URL)" status

test:
	go test ./...   

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