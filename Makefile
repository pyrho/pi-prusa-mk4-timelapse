build:
	@go build -o bin/timelapse-serial cmd/timelapse-serial/main.go

run: build
	@./bin/timelapse-serial

test:
	@go test -v ./...

default: build

.PHONY: run build test
