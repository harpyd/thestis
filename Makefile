.PHONY:
.SILENT:

build:
	go build ./...

lint:
	golangci-lint run