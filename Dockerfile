# Multistage build for exodus-node (Go) + sing-box core (with v2ray_api) + s6-overlay.
FROM golang:1.25.12-alpine AS builder

RUN apk update && apk add --no-cache git ca-certificates tzdata gcc musl-dev sqlite-dev

ARG VERSION=unknown
ARG REVISION=unknown
ARG BUILD_TAGS=none
ARG CGO_ENABLED_STATUS=enabled

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

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
        -o /build/exodus-node-app \
        .

FROM alpine:3.23 AS singbox

ARG SINGBOX_VERSION=v1.13.13
ARG ASN_LMDB_URL=https://github.com/Adam-Sizzler/lmdb-go/releases/download/latest/asn-prefixes-lmdb.tar.gz

RUN mkdir -p /usr/local/bin /usr/local/share/asn \
    && wget -qO /usr/local/bin/sing-box "https://github.com/Adam-Sizzler/sing-box-v2ray-api/releases/download/${SINGBOX_VERSION}/sing-box-linux-amd64" \
    && chmod +x /usr/local/bin/sing-box \
    && wget -qO /tmp/asn-prefixes-lmdb.tar.gz "${ASN_LMDB_URL}" \
    && tar -xzf /tmp/asn-prefixes-lmdb.tar.gz -C /usr/local/share/asn \
    && rm -f /tmp/asn-prefixes-lmdb.tar.gz

FROM alpine:3.23

ARG S6_OVERLAY_VERSION=3.2.0.2

LABEL org.opencontainers.image.title="Exodus Node" \
      org.opencontainers.image.description="Exodus Node with built-in Sing-box Core" \
      org.opencontainers.image.url="https://github.com/exodus/node" \
      org.opencontainers.image.source="https://github.com/exodus/node" \
      org.opencontainers.image.vendor="Exodus" \
      org.opencontainers.image.licenses="AGPL-3.0" \
      org.opencontainers.image.documentation="https://docs.exodus.dev"

WORKDIR /app

COPY --from=builder /build/exodus-node-app /app/exodus-node
COPY --from=singbox /usr/local/bin/sing-box /usr/local/bin/sing-box
COPY --from=singbox /usr/local/share/asn /usr/local/share/asn

COPY deploy/s6-overlay/etc/s6-overlay /etc/s6-overlay

RUN S6_ARCH="$(uname -m)" \
    && wget -qO /tmp/s6-noarch.tar.xz "https://github.com/just-containers/s6-overlay/releases/download/v${S6_OVERLAY_VERSION}/s6-overlay-noarch.tar.xz" \
    && wget -qO /tmp/s6-arch.tar.xz "https://github.com/just-containers/s6-overlay/releases/download/v${S6_OVERLAY_VERSION}/s6-overlay-${S6_ARCH}.tar.xz" \
    && xz -dc /tmp/s6-noarch.tar.xz | tar -C / -xpf - \
    && xz -dc /tmp/s6-arch.tar.xz | tar -C / -xpf - \
    && rm -f /tmp/s6-noarch.tar.xz /tmp/s6-arch.tar.xz \
    && chmod -R +x /etc/s6-overlay \
    && chmod +x /app/exodus-node \
    && mkdir -p /run /app/singbox /app/logs

ENTRYPOINT ["/init"]

CMD ["/command/with-contenv", "/app/exodus-node"]
