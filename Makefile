build:
	@go build -o bin/timelapse-serial cmd/timelapse-serial/main.go

watch:
	@ gow -e=css,js,go,mod,html run ./cmd/timelapse-serial/main.go --configPath ~/.timelapse-serial

install: build
	@cp ./configs/config.toml /usr/local/etc/timelapse-serial.toml
	@cp ./bin/timelapse-serial /usr/local/bin/.

test:
	@go test -v ./...

default: build

.PHONY: run build test
