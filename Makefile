.PHONY:
.SILENT:

openapi:
	oapi-codegen -generate types -o internal/port/http/v1/openapi_type.gen.go -package v1 api/openapi/thestis.yml
	oapi-codegen -generate chi-server -o internal/port/http/v1/openapi_server.gen.go -package v1 api/openapi/thestis.yml

thestis-validate-build:
	go mod download && CGO_ENABLES=0 go build -o ./.bin/thestis-validate ./cmd/thestis-validate

thestis-build:
	go mod download && CGO_ENABLED=0 GOOS=linux go build -o ./.bin/thestis ./cmd/thestis

lint:
	golangci-lint run

dev: thestis-build
	docker-compose -f ./deployments/dev/docker-compose.yml --project-directory . up --remove-orphans thestis

test-unit:
	go test --short -v -race -coverpkg=./... -coverprofile=unit-all.out ./...
	cat unit-all.out | grep -v .gen.go > unit.out
	rm unit-all.out

test-integration:
	go test -v -race -coverprofile=integration.out ./internal/adapter/...