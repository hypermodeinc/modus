FROM golang:alpine as builder
LABEL maintainer="Hypermode <hello@hypermode.com>"

# build hmruntime binary
WORKDIR /src
COPY go.mod go.sum ./
COPY *.go ./
RUN go build .

FROM ubuntu:20.04

COPY --from=builder /src/hmruntime /usr/bin/hmruntime
CMD ["hmruntime"]
