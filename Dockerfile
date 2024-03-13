# build hmruntime binary
FROM --platform=$BUILDPLATFORM golang:alpine as builder
WORKDIR /src

# copy go.mod and go.sum files separately so that
# it is only run when the dependencies change
COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./
ARG TARGETOS TARGETARCH
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH go build .

# build runtime image
FROM ubuntu:22.04
LABEL maintainer="Hypermode <hello@hypermode.com>"
COPY --from=builder /src/hmruntime /usr/bin/hmruntime

# add common tools
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    curl \
    iputils-ping \
    jq \
    less \
    && rm -rf /var/lib/apt/lists/*

ENTRYPOINT ["hmruntime", "--jsonlogs"]
