build:
	go build -o bin/flashat

run: build
	bin/flashat

test:
	go test -v ./... -count=1

setup_db:
	./setup_db.sh
	