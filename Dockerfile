# set up the image to build the explorer app
FROM --platform=$BUILDPLATFORM node:alpine AS node-builder
WORKDIR /src

# install dependencies first
COPY ./runtime/explorer/content/package*.json ./
RUN npm ci

# copy and build the rest of the application
COPY ./runtime/explorer/content .
RUN npm run build

# set up the image to build the runtime
FROM --platform=$BUILDPLATFORM golang:alpine AS builder
WORKDIR /src

# copy lib dependencies
COPY ./lib ./lib

# Copy and modify go.work file
COPY ./go.work ./
RUN sed -i '/^[[:space:]]*\.\/sdk\//d' ./go.work
RUN sed -i '/^[[:space:]]*\.\/.*\/testdata/d' ./go.work

# switch to the runtime directory
WORKDIR /src/runtime

# copy go.mod and go.sum files separately so that the download step
# is only run when the dependencies change
COPY runtime/go.mod runtime/go.sum ./
RUN go mod download

# copy the rest of the runtime source and the compiled explorer app
COPY runtime/ ./
COPY --from=node-builder /src/dist ./explorer/content/dist

# build the runtime binary
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
