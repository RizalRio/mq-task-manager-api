# Sesuaikan password dan dbname dengan yang ada di .env kamu
DB_URL=postgres://postgres:asdqwe123@localhost:5432/gk_task_manager?sslmode=disable

migrate-up:
	migrate -path db/migrations -database "$(DB_URL)" -verbose up

migrate-down:
	migrate -path db/migrations -database "$(DB_URL)" -verbose down

.PHONY: migrate-up migrate-down