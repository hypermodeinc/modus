# build hmruntime binary
FROM --platform=$BUILDPLATFORM golang:alpine as runtime-builder
WORKDIR /src
COPY hmruntime/ ./
ARG TARGETOS TARGETARCH
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH go build .

# build example plugin
FROM node:20-alpine as plugin-builder
WORKDIR /src
COPY plugins/as/ ./
WORKDIR /src/hmplugin1
RUN npm install
RUN npm run build:release

# build runtime image
FROM ubuntu:20.04
LABEL maintainer="Hypermode <hello@hypermode.com>"
COPY --from=runtime-builder /src/hmruntime /usr/bin/hmruntime
COPY --from=plugin-builder /src/hmplugin1/build/release.wasm /plugins/hmplugin1.wasm

ENTRYPOINT ["hmruntime", "--plugins=/plugins"]
