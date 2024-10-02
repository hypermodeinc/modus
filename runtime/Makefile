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

.PHONY: test
test:
	go test ./...

.PHONY: test-race
test-race:
	go test -race ./...

.PHONY: test-integration
test-integration:
	go test -race -tags=integration ./integration_tests/...

.PHONY: clean
clean:
	go clean
	rm -f $(EXECUTABLE)
