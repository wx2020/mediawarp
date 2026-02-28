FROM golang:1.24 AS builder
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.io,direct \
    CGO_ENABLED=0

WORKDIR /builder
COPY . .
ARG MEDIAWARP_VERSION=dev
RUN go mod download && \
    go build -ldflags="-s -w -X MediaWarp/internal/config.appVersion=${MEDIAWARP_VERSION}" \
    -o MediaWarp

FROM alpine:latest
COPY --from=builder /builder/MediaWarp /MediaWarp

RUN chmod +x /MediaWarp

EXPOSE 9000
VOLUME ["/etc/localtime", "/etc/timezone", "/config", "/logs", "/custom"]
ENTRYPOINT ["/MediaWarp"]
