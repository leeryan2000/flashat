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

send:
	$env:GOOS="linux"; $env:GOARCH="amd64"; go build -o flashat-backend main.go
	scp -i leeryan.pem flashat-backend ec2-user@18.140.243.27:/home/ec2-user/flashat/backend/