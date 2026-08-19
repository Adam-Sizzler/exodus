# Multistage build for exodus-node (Go) + sing-box core (with v2ray_api) + s6-overlay.
FROM golang:1.25.12-alpine AS builder

RUN apk update && apk add --no-cache git ca-certificates tzdata gcc musl-dev sqlite-dev

ARG VERSION=unknown

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN set -eu; \
    version="${VERSION}"; \
    if [ -z "${version}" ] || [ "${version}" = "unknown" ]; then \
        version="$(git describe --tags --abbrev=0 2>/dev/null || true)"; \
    fi; \
    if [ -z "${version}" ] || [ "${version}" = "unknown" ]; then \
        version="$(git describe --tags --always --dirty 2>/dev/null || true)"; \
    fi; \
    if [ -z "${version}" ]; then \
        version="unknown"; \
    fi; \
    CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64 \
    go build \
        -trimpath \
        -buildvcs=true \
        -ldflags="-s -w \
                -X exodus-node/constant.Version=${version}" \
        -o /build/exodus-node-app \
        .

FROM alpine:3.23 AS singbox

ARG SINGBOX_VERSION=v1.13.13

RUN apk update && apk add --no-cache ca-certificates curl \
    && curl -fL "https://github.com/Adam-Sizzler/sing-box-v2ray-api/releases/download/${SINGBOX_VERSION}/sing-box-linux-amd64" \
        -o /usr/local/bin/sing-box \
    && chmod +x /usr/local/bin/sing-box

FROM alpine:3.23

ARG S6_OVERLAY_VERSION=3.2.0.2

LABEL org.opencontainers.image.title="Exodus Node" \
      org.opencontainers.image.description="Exodus Node with built-in Sing-box Core" \
      org.opencontainers.image.url="https://github.com/exodus/node" \
      org.opencontainers.image.source="https://github.com/exodus/node" \
      org.opencontainers.image.vendor="Exodus" \
      org.opencontainers.image.licenses="AGPL-3.0" \
      org.opencontainers.image.documentation="https://docs.exodus.dev"

RUN apk update && apk add --no-cache ca-certificates tzdata sqlite-libs curl xz

# Install s6-overlay
RUN S6_ARCH="$(uname -m)" && \
    curl -sSL -o /tmp/s6-noarch.tar.xz "https://github.com/just-containers/s6-overlay/releases/download/v${S6_OVERLAY_VERSION}/s6-overlay-noarch.tar.xz" && \
    curl -sSL -o /tmp/s6-arch.tar.xz "https://github.com/just-containers/s6-overlay/releases/download/v${S6_OVERLAY_VERSION}/s6-overlay-${S6_ARCH}.tar.xz" && \
    xz -dc /tmp/s6-noarch.tar.xz | tar -C / -xpf - && \
    xz -dc /tmp/s6-arch.tar.xz | tar -C / -xpf - && \
    rm -f /tmp/s6-noarch.tar.xz /tmp/s6-arch.tar.xz

WORKDIR /opt/app

COPY --from=builder /build/exodus-node-app /opt/app/exodus-node
COPY --from=singbox /usr/local/bin/sing-box /usr/local/bin/sing-box

COPY deploy/s6-overlay/etc/s6-overlay /etc/s6-overlay

RUN chmod -R +x /etc/s6-overlay && \
    chmod +x /opt/app/exodus-node && \
    mkdir -p /run /opt/app/singbox /var/log/singbox /usr/local/share/asn && \
    printf '#!/bin/sh\ntail -n +1 -f /var/log/singbox/current\n' > /usr/local/bin/slogs && \
    printf '#!/bin/sh\ntail -n +1 -f /var/log/singbox/current | grep -E "(^| )(WARN|ERROR|FATAL)( |\\[)|panic:"\n' > /usr/local/bin/serrors && \
    chmod +x /usr/local/bin/slogs /usr/local/bin/serrors

ENV S6_VERBOSITY=1

ENTRYPOINT ["/init"]

CMD ["/command/with-contenv", "/opt/app/exodus-node"]
