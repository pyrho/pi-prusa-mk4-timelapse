build:
	go build -o bin/timelapse-serial cmd/timelapse-serial/main.go

run: build
	./bin/timelapse-serial

.PHONY: run build
