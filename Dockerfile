# build the application binary
FROM --platform=$BUILDPLATFORM golang:alpine AS builder
WORKDIR /src

# copy dependencies then switch to the runtime directory
COPY ./lib ./lib
WORKDIR /src/runtime

# copy go.mod and go.sum files separately so that the download step
# is only run when the dependencies change
COPY runtime/go.mod runtime/go.sum ./
RUN go mod download

COPY runtime/ ./
ARG TARGETOS TARGETARCH RUNTIME_RELEASE_VERSION
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o modus_runtime -ldflags "-s -w -X github.com/hypermodeinc/modus/runtime/config.version=$RUNTIME_RELEASE_VERSION" .

# build the container image
FROM ubuntu:22.04
LABEL maintainer="Hypermode Inc. <hello@hypermode.com>"

# add common tools
RUN apt-get update && apt-get install -y --no-install-recommends \
    curl \
    iputils-ping \
    jq \
    less 

# copy runtime binary from the build phase
COPY --from=builder /src/runtime/modus_runtime /usr/bin/modus_runtime

# update certificates every build
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# set the default entrypoint and options
ENTRYPOINT ["modus_runtime", "--jsonlogs"]
