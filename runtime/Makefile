EXECUTABLE := modus_runtime
VERSION := $(shell git describe --tags --always --match 'runtime/*' | sed 's/^runtime\///')
LDFLAGS := -X github.com/hypermodeinc/modus/runtime/app.version=$(VERSION)

ifneq ($(OS), Windows_NT)
	OS := $(shell uname -s)
endif

ifeq ($(OS), Windows_NT)
	EXECUTABLE := $(EXECUTABLE).exe
	LDFLAGS := $(LDFLAGS) -checklinkname=0
endif

.PHONY: all
all: build

.PHONY: tidy
tidy:
	go mod tidy -v
	go fmt ./...

.PHONY: copy-db-migrations
copy-db-migrations:
	cp db/migrations/*.sql ../tools/modus-dbmigrate/migrations/

.PHONY: build-explorer
build-explorer:
	cd explorer/content && npm ci && npm run build

.PHONY: build
build: build-explorer copy-db-migrations
	go build -o $(EXECUTABLE) -ldflags "$(LDFLAGS)" .

.PHONY: run
run: build-explorer
	@ARGS="$(filter-out $@,$(MAKECMDGOALS))" && \
	MODUS_ENV=dev go run . $$ARGS

.PHONY: runapp
runapp:
	@ARGS="$(filter-out $@,$(MAKECMDGOALS))" && \
	MODUS_ENV=dev go run . -appPath $$ARGS

.PHONY: build-testdata
build-testdata: build-testdata-assemblyscript build-testdata-golang

.PHONY: build-testdata-assemblyscript
build-testdata-assemblyscript:
	cd languages/assemblyscript/testdata && npm ci && npm run build

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

.PHONY: build-integration-testdata
build-integration-testdata:
	cd ../sdk/assemblyscript/src && npm ci
	cd ../sdk/assemblyscript/examples/postgresql && npm ci && npm run build
	cp ../sdk/assemblyscript/examples/postgresql/build/postgresql-example.wasm ./integration_tests/testdata

.PHONY: test-integration
test-integration: build-integration-testdata
	go test -race -tags=integration -count=1 ./integration_tests/...

.PHONY: test-ci
test-ci:
	go install github.com/jstemmer/go-junit-report/v2@latest
	go test -v -race 2>&1 ./... | go-junit-report -set-exit-code -iocopy -out report.xml

.PHONY: test-integration-ci
test-integration-ci: build-integration-testdata
	go install github.com/jstemmer/go-junit-report/v2@latest
	go test -tags integration -v -race 2>&1 ./integration_tests/... | go-junit-report -set-exit-code -iocopy -out report.xml

.PHONY: clean
clean:
	go clean
	rm -f $(EXECUTABLE)
	rm -f **/*_generated.go

.PHONY: docker-build
docker-build:
	cd .. && docker build --build-arg RUNTIME_RELEASE_VERSION="$(VERSION)" -t modus_runtime .

# Catch all rule to prevent make errors for non-target arguments
%:
	@true
