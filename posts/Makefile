SHELL := bash.exe

build:
	go build -o bin/flashat-posts main.go

run: tidy build
	bin/flashat-posts

test:
	go test ./... -race -count=1

tidy:
	go mod tidy
	go fmt ./...

migrate_up:
	./migrate.sh up

migrate_down:
	./migrate.sh down -all

migrate:
	migrate -path db/migrations -database "postgres://dev:dev@localhost:5432/flashat_posts?sslmode=disable" version
