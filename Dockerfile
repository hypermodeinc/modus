# build hmruntime binary
FROM --platform=$BUILDPLATFORM golang:alpine as builder
WORKDIR /src
COPY ./ ./
ARG TARGETOS TARGETARCH
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH go build .

# build runtime image
FROM ubuntu:20.04
LABEL maintainer="Hypermode <hello@hypermode.com>"
COPY --from=builder /src/hmruntime /usr/bin/hmruntime

ENTRYPOINT ["hmruntime"]
