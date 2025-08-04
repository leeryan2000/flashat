build:
	go build -o bin/flashat ./cmd/main.go

run: build
	bin/flashat

test:
	go test -v ./handlers -count=1

setup_db:
	./setup_db.sh
	