# Multistage build for exodus panel (frontend + backend)
FROM node:20-alpine AS panel-ui
WORKDIR /ui
ENV NODE_OPTIONS=--max-old-space-size=4096
ARG SINGBOX_WASM_URL=https://adam-sizzler.github.io/s-validator/main.wasm
ARG SINGBOX_SCHEMA_URL=https://adam-sizzler.github.io/s-validator/singbox.schema.json
ARG SINGBOX_SCHEMA_CN_URL=https://adam-sizzler.github.io/s-validator/singbox.schema.json
COPY frontend/package*.json ./
RUN --mount=type=cache,target=/root/.npm npm ci --legacy-peer-deps --no-audit --prefer-offline
RUN mkdir -p /tmp/exodus-assets
RUN set -eu; \
    fetch_with_retry() { \
      url="$1"; out="$2"; name="$3"; \
      [ -n "$url" ] || return 0; \
      n=1; ok=0; \
      while [ "$n" -le 4 ]; do \
        if wget -T 12 -t 1 -qO "$out" "$url"; then \
          ok=1; break; \
        fi; \
        echo "WARN: download failed for $name (attempt $n/4): $url"; \
        rm -f "$out"; \
        n=$((n+1)); \
        sleep 2; \
      done; \
      if [ "$ok" -ne 1 ]; then \
        echo "WARN: skip $name after retries: $url"; \
      fi; \
    }; \
    fetch_with_retry "$SINGBOX_WASM_URL" /tmp/exodus-assets/main.wasm main.wasm; \
    fetch_with_retry "$SINGBOX_SCHEMA_URL" /tmp/exodus-assets/singbox.schema.json singbox.schema.json; \
    fetch_with_retry "$SINGBOX_SCHEMA_CN_URL" /tmp/exodus-assets/singbox.schema.cn.json singbox.schema.cn.json
RUN if [ -f /tmp/exodus-assets/main.wasm ]; then \
      magic="$(od -An -t x1 -N 4 /tmp/exodus-assets/main.wasm | tr -d ' \n')"; \
      [ "$magic" = "0061736d" ] || { echo "Invalid WASM magic for main.wasm: $magic"; exit 1; }; \
      echo "WASM artifact attached: /tmp/exodus-assets/main.wasm"; \
    else \
      echo "WARN: /tmp/exodus-assets/main.wasm is missing. Sing-box validator will be disabled in UI."; \
    fi
COPY frontend/ ./
RUN mkdir -p /ui/public/assets && cp -f /tmp/exodus-assets/* /ui/public/assets/ 2>/dev/null || true
RUN --mount=type=cache,target=/root/.npm --mount=type=cache,target=/ui/node_modules/.vite/cache npm run cb

FROM golang:1.25-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

ARG VERSION=latest
ARG REVISION=unknown
ARG FRONTEND_REVISION=unknown
ARG BUILD_BRANCH=unknown
ARG BUILD_TIME=unknown
ARG BUILD_NUMBER=unknown
ARG REPOSITORY_URL=unknown
ARG BUILD_TAGS=none
ARG CGO_ENABLED_STATUS=disabled

WORKDIR /build/backend

# Copy go mod files and download dependencies
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Copy source code
COPY backend/ ./

# Build panel backend binary
RUN CGO_ENABLED=0 \
    GOOS=linux \
    go build \
        -tags "none" \
        -trimpath \
        -buildvcs=false \
        -ldflags="-s -w \
                -X 'exodus/internal/constant.Version=${VERSION}' \
                -X 'exodus/internal/constant.Revision=${REVISION}' \
                -X 'exodus/internal/constant.BuildTags=none' \
                -X 'exodus/internal/constant.CgoEnabled=0'" \
    -o /build/exodus \
    .

# Final stage
FROM alpine:latest

ARG VERSION=latest
ARG REVISION=unknown
ARG FRONTEND_REVISION=unknown
ARG BUILD_BRANCH=unknown
ARG BUILD_TIME=unknown
ARG BUILD_NUMBER=unknown
ARG REPOSITORY_URL=unknown

LABEL org.opencontainers.image.title="exodus" \
      org.opencontainers.image.version="${VERSION}" \
      org.opencontainers.image.revision="${REVISION}" \
      org.opencontainers.image.created="${BUILD_TIME}" \
      org.opencontainers.image.source="${REPOSITORY_URL}"

ENV EXODUS_VERSION="${VERSION}" \
    EXODUS_BACKEND_COMMIT="${REVISION}" \
    EXODUS_FRONTEND_COMMIT="${FRONTEND_REVISION}" \
    EXODUS_GIT_BRANCH="${BUILD_BRANCH}" \
    EXODUS_BUILD_TIME="${BUILD_TIME}" \
    EXODUS_BUILD_NUMBER="${BUILD_NUMBER}" \
    EXODUS_REPOSITORY_URL="${REPOSITORY_URL}"

# Runtime certs/timezone data come from the builder stage to keep the final
# image independent from Alpine repository availability during local rebuilds.
COPY --from=builder /etc/ssl/certs /etc/ssl/certs
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/exodus /app/

# Keep /app/exodus as the server entrypoint, and expose rescue commands
# for `docker exec -it exodus cli` and `docker exec -it exodus exodus`.
RUN printf '%s\n' \
      '#!/bin/sh' \
      'if [ "$#" -eq 0 ]; then' \
      '  exec /app/exodus --rescue' \
      'fi' \
      'exec /app/exodus "$@"' \
    > /usr/local/bin/exodus \
    && printf '%s\n' \
      '#!/bin/sh' \
      'exec /app/exodus --rescue "$@"' \
    > /usr/local/bin/cli \
    && chmod +x /usr/local/bin/exodus /usr/local/bin/cli

# Copy built frontend
COPY --from=panel-ui /ui/dist /app/ui
COPY --from=builder /usr/local/go/lib/wasm/wasm_exec.js /app/ui/assets/wasm_exec.js

# Create directories for data and certificates
RUN mkdir -p /app/data /app/certs

EXPOSE 3000

ENTRYPOINT ["/app/exodus"]
