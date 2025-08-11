build:
	go build -o bin/flashat ./cmd/main.go

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
