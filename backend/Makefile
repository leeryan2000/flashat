SHELL := bash.exe

build:
	go build -o bin/flashat main.go

run: tidy build 
	 bin/flashat

test:
	go test -v ./handlers -count=1

setup:
	./migrate.sh goto 1
	./migrate.sh up

tidy:
	go mod tidy
	go fmt ./...

migrate_up: 
	./migrate.sh up

migrate_down:
	./migrate.sh down -all

migrate:
	migrate -path db/migrations -database "postgres://dev:dev@localhost:5432/dev?sslmode=disable" version