.PHONY:
.SILENT:

export API_V1 = api/openapi/thestis-v1.yml

gen-api-v1:
	oapi-codegen -generate types -o internal/port/http/v1/openapi_type.gen.go -package v1 $$API_V1
	oapi-codegen -generate chi-server -o internal/port/http/v1/openapi_server.gen.go -package v1 $$API_V1

api-v1:
	make gen-api-v1
	cp $$API_V1 swagger/v1/thestis.yml

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
	make run-test-db
	go test -v -race -coverprofile=integration.out ./internal/adapter/... ./internal/config/... || (make stop-test-db && exit 1)
	make stop-test-db


export TEST_DB_URI=mongodb://localhost:27019
export TEST_DB_NAME=test
export TEST_CONTAINER_NAME=test-db

run-test-db:
	docker run --rm -d -p 27019:27017 --name $$TEST_CONTAINER_NAME -e MONGODB_DATABASE=$$TEST_DB_NAME mongo:4.4-bionic

stop-test-db:
	docker stop $$TEST_CONTAINER_NAME