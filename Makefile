.PHONY:
.SILENT:

thestis-validate-build:
	go mod download && CGO_ENABLES=0 go build -o ./.bin/thestis-validate ./cmd/thestis-validate

lint:
	golangci-lint run

test-unit:
	go test --short -v -race -coverpkg=./... -coverprofile=unit-all.out ./...
	cat unit-all.out | grep -v .gen.go > unit.out
	rm unit-all.out