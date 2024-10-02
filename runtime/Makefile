# Variables
EXECUTABLE := modus_runtime
RACE_FLAG := # Empty by default

# Get the OS name
ifneq ($(OS), Windows_NT)
	OS := $(shell uname -s)
endif

# If on Windows, set the executable extension to .exe
ifeq ($(OS), Windows_NT)
	EXECUTABLE := $(EXECUTABLE).exe
endif

# Check if the '-race' argument is passed
ifneq (,$(findstring -race,$(MAKECMDGOALS)))
    RACE_FLAG := -race
    # Remove the '-race' argument from MAKECMDGOALS
    MAKECMDGOALS := $(filter-out -race,$(MAKECMDGOALS))
endif

# PHONY targets to avoid name conflicts with files
.PHONY: all generate build clean test integration

# Default target
all: build

# Go generate before building
generate:
	go generate ./...

# Build the binary
build: generate
	go build -o $(EXECUTABLE)

# Run unit tests with optional race flag
test:
	go test $(RACE_FLAG) ./...

# Run integration tests with the 'integration' build tag and optional race flag
integration:
	go test $(RACE_FLAG) -tags=integration ./...

# Clean the build artifacts
clean:
	go clean
	rm -f $(EXECUTABLE)
