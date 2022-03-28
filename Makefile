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

run-dev: thestis-build
	docker-compose -f ./deployments/dev/docker-compose.yml --project-directory . up --remove-orphans thestis

stop-dev:
	docker-compose -f ./deployments/dev/docker-compose.yml --project-directory . stop thestis

test-unit:
	go test --short -v -race -coverpkg=./... -coverprofile=unit-all.out ./...
	cat unit-all.out | grep -v .gen.go > unit.out
	rm unit-all.out

test-integration:
	make run-test-db
	make run-test-nats
	MallocNanoZone=0 go test -v -race -coverprofile=integration.out ./internal/adapter/... ./internal/config/... || (make stop-test-nats && make stop-test-db && exit 1)
	make stop-test-db
	make stop-test-nats

test-cover:
	go install github.com/wadey/gocovmerge@latest
	gocovmerge unit.out integration.out > cover.out
	go tool cover -html=cover.out -o cover.html


export TEST_DB_URI=mongodb://localhost:27019
export TEST_DB_NAME=test
export TEST_DB_CONTAINER_NAME=test-db

run-test-db:
	docker run --rm -d -p 27019:27017 --name $$TEST_DB_CONTAINER_NAME -e MONGODB_DATABASE=$$TEST_DB_NAME mongo:4.4-bionic

stop-test-db:
	docker stop $$TEST_DB_CONTAINER_NAME


export TEST_NATS_CONTAINER_NAME=test-nats

run-test-nats:
	docker run --rm -d -p 4222:4222 -ti --name $$TEST_NATS_CONTAINER_NAME nats:latest

stop-test-nats:
	docker stop $$TEST_NATS_CONTAINER_NAME