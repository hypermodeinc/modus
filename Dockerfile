# build hmruntime binary
FROM --platform=$BUILDPLATFORM golang:alpine as builder
WORKDIR /src
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

ENTRYPOINT ["hmruntime"]
