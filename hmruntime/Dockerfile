FROM golang:alpine as builder
LABEL maintainer="Hypermode <hello@hypermode.com>"

# build hmruntime binary
WORKDIR /src
COPY go.mod go.sum ./
COPY *.go ./
RUN go build .

FROM ubuntu:20.04

# only update, don't run upgrade
# use cache busting to avoid old versions
# remove /var/lib/apt/lists/* to reduce image size.
# see: https://docs.docker.com/develop/develop-images/dockerfile_best-practices
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    curl \
    iputils-ping \
    jq \
    less \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /src/hmruntime /usr/bin/hmruntime
CMD ["hmruntime"]
