# Multistage build for exodus-node (Go) + sing-box core (with v2ray_api) + supervisord.
FROM golang:1.25-alpine AS builder

RUN apk update && apk add --no-cache git ca-certificates tzdata gcc musl-dev sqlite-dev

ARG VERSION=unknown
ARG REVISION=unknown
ARG BUILD_TAGS=none
ARG CGO_ENABLED_STATUS=enabled

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN set -eu; \
    version="${VERSION}"; \
    revision="${REVISION}"; \
    if [ -z "${version}" ] || [ "${version}" = "unknown" ]; then \
        version="$(git describe --tags --abbrev=0 2>/dev/null || true)"; \
    fi; \
    if [ -z "${version}" ] || [ "${version}" = "unknown" ]; then \
        version="$(git describe --tags --always --dirty 2>/dev/null || true)"; \
    fi; \
    if [ -z "${version}" ]; then \
        version="unknown"; \
    fi; \
    if [ -z "${revision}" ] || [ "${revision}" = "unknown" ]; then \
        revision="$(git rev-parse HEAD 2>/dev/null || true)"; \
    fi; \
    if [ -z "${revision}" ]; then \
        revision="unknown"; \
    fi; \
    CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64 \
    go build \
        -tags "${BUILD_TAGS}" \
        -trimpath \
        -buildvcs=true \
        -ldflags="-s -w \
                -X exodus-node/constant.Version=${version} \
                -X exodus-node/constant.Revision=${revision} \
                -X exodus-node/constant.BuildTags=${BUILD_TAGS} \
                -X exodus-node/constant.CgoEnabled=${CGO_ENABLED_STATUS}" \
        -o /build/exodus-node \
        .

FROM alpine:3.23

ARG SINGBOX_VERSION=v1.13.5
ENV SINGBOX_VERSION=${SINGBOX_VERSION}

RUN apk update && apk add --no-cache ca-certificates tzdata sqlite-libs curl supervisor nftables

RUN curl -fL "https://github.com/Adam-Sizzler/sing-box-v2ray-api/releases/download/${SINGBOX_VERSION}/sing-box-linux-amd64" \
      -o /usr/local/bin/sing-box && \
    chmod +x /usr/local/bin/sing-box

WORKDIR /app

COPY --from=builder /build/exodus-node /app/exodus-node
COPY deploy/supervisord.conf /etc/supervisord.conf
COPY deploy/docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh
COPY deploy/singbox-default.json /app/singbox/config.default.json

RUN chmod +x /usr/local/bin/docker-entrypoint.sh && \
    mkdir -p /run /var/log/supervisor /app/singbox /app/logs /app/certs

ENTRYPOINT ["/usr/local/bin/docker-entrypoint.sh"]
CMD ["/app/exodus-node"]
