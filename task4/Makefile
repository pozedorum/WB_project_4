all: run

build:
	go build -o metrics-server cmd/main.go

run:
	go mod tidy
	go run cmd/main.go

lint:
	golangci-lint run ./...
	go vet ./...

clean:
	rm -f metrics-server