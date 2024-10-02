# build the application binary
FROM --platform=$BUILDPLATFORM golang:alpine AS builder
WORKDIR /src

# copy go.mod and go.sum files separately so that the download step
# is only run when the dependencies change
COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./
ARG TARGETOS TARGETARCH MODUS_BUILD_VERSION
RUN MODUS_BUILD_VERSION=$MODUS_BUILD_VERSION go generate ./...
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o modus_runtime

# build the container image
FROM ubuntu:22.04
LABEL maintainer="Hypermode, Inc. <hello@hypermode.com>"

# add common tools
RUN apt-get update && apt-get install -y --no-install-recommends \
    curl \
    iputils-ping \
    jq \
    less 

# copy runtime binary from the build phase
COPY --from=builder /src/modus_runtime /usr/bin/modus_runtime

# create a link to the old name(s) for backwards compatibility
RUN ln -s /usr/bin/modus_runtime /usr/bin/hmruntime
RUN ln -s /usr/bin/modus_runtime /usr/bin/hypruntime

# update certificates every build
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# set the default entrypoint and options
ENTRYPOINT ["modus_runtime", "--jsonlogs"]
