EXECUTABLE := modus_runtime

ifneq ($(OS), Windows_NT)
	OS := $(shell uname -s)
endif

ifeq ($(OS), Windows_NT)
	EXECUTABLE := $(EXECUTABLE).exe
endif

.PHONY: all
all: build

.PHONY: tidy
tidy:
	go mod tidy -v
	go fmt ./...

.PHONY: generate
generate:
	go generate ./...

.PHONY: build
build: generate
	go build -o $(EXECUTABLE)

.PHONY: run
run: generate
	go run .

.PHONY: build-testdata
build-testdata: generate build-testdata-assemblyscript build-testdata-golang

.PHONY: build-testdata-assemblyscript
build-testdata-assemblyscript:
	cd languages/assemblyscript/testdata && npm install && npm run build

.PHONY: build-testdata-golang
build-testdata-golang:
	cd languages/golang/testdata && ./build.sh

.PHONY: test
test: build-testdata
	go test ./...

.PHONY: test-no-cache
test-no-cache: build-testdata
	go test -count=1 ./...

.PHONY: test-race
test-race: build-testdata
	go test -race ./...

.PHONY: test-integration
test-integration:
	go test -race -tags=integration ./integration_tests/...

.PHONY: test-ci
test-race: build-testdata
    go install github.com/jstemmer/go-junit-report/v2@latest
    go test -v -race 2>&1 ./... | go-junit-report -set-exit-code -iocopy -out report.xml

.PHONY: test-integration-ci
test-integration:
	go install github.com/jstemmer/go-junit-report/v2@latest
	go test -tags integration -v -race 2>&1 ./integration_tests/... | go-junit-report -set-exit-code -iocopy -out report.xml

.PHONY: clean
clean:
	go clean
	rm -f $(EXECUTABLE)
	rm -f **/*_generated.go

.PHONY: docker-build
docker-build:
	docker build --build-arg MODUS_BUILD_VERSION="$(shell git describe --tags --always)" -t modus_runtime .
