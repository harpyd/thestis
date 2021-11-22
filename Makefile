.PHONY:
.SILENT:

thestis-validate-build:
	go mod download && CGO_ENABLES=0 go build -o ./.bin/thestis-validate ./cmd/thestis-validate

lint:
	golangci-lint run