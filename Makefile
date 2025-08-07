build:
	go build -o bin/flashat ./cmd/main.go

run: setup tidy build 
	bin/flashat

test:
	go test -v ./handlers -count=1

setup:
	.\setup_db.sh

tidy:
	go mod tidy
	go fmt ./...

migrate_up: 
	./migrate.sh up

migrate_down:
	./migrate.sh down -all
