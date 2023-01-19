VET_REPORT=vet.report.txt
TEST_REPORT=test.report.txt
GOARCH=amd64
GOOS?=darwin
MYSQL_VER=1.6.2
VERSION?=$(shell date '+%Y%m%d.%H%M%S')
COMMIT=$(shell git rev-parse HEAD)
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)
HOST=$(shell hostname)
APP_PORT=${APP_PORT}

BUILD_DIR=$(shell pwd)
BUILD_DIR_LINK=$(shell readlink ${BUILD_DIR})

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS = -ldflags "-X main.VERSION=${VERSION} -X main.COMMIT=${COMMIT} -X main.BRANCH=${BRANCH} -X main.BUILDHOST=${HOST}"

.PHONY: all
all: clean test fmt build

.PHONY: build-image
build-image: build
	docker build -f Dockerfile -t todo_test .

.PHONY: build
build:
	GO111MODULE=on GOOS=${GOOS} GOARCH=${GOARCH} go build -tags musl ${LDFLAGS} -o tmp/service ./
	@echo "VERSION=${VERSION}" > tmp/version

.PHONY: runlocal
runlocal: openapi
	@GO111MODULE=off go get github.com/cespare/reflex
	APP_ENV=dev go run ./main.go
	 ~/go/bin/reflex -r '\.go' -s -- sh -c "go run ./main.go"

.PHONY: test
test:
	GO111MODULE=on go clean -testcache
	GO111MODULE=on go test $$(go list ./... ) 2>&1

test-coverage:
	GO111MODULE=on go test $$(go list ./... ) -coverprofile tmp/cover.out
	go tool cover -html=tmp/cover.out -o tmp/coverage.html
	open tmp/coverage.html

test-coverage-functions:
	GO111MODULE=on go test $$(go list ./... ) -covermode=atomic -coverprofile tmp/cover.out
	go tool cover -func tmp/cover.out -o tmp/function-coverage.out

.PHONY: godownload
godownload:
	go mod download -x

.PHONY: gotidy
gotidy:
	go mod tidy
	go mod vendor

.PHONY: fmt
fmt:
	GO111MODULE=on go fmt ./...

.PHONY: clean
clean:
	-rm -f ${TEST_REPORT}
	-rm -f ${VET_REPORT}
	-rm -rf tmp

.PHONY: openapi
openapi:
	@GO111MODULE=off go get -v github.com/swaggo/swag/cmd/swag
	rm -rf api/v1/docs
	$(GOPATH)/bin/swag init -g ./cmd/service/main.go -d . --output api/v1/docs

.PHONY: rundocker
rundocker: build-image
	docker run --name mysql -e MYSQL_ROOT_PASSWORD=test -d mysql:$(MYSQL_VER)
	docker run --rm --name todo_test -p $(APP_PORT):$(APP_PORT) todo_test
