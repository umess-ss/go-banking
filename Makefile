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