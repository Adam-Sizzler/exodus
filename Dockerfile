# syntax=docker/dockerfile:1.7

# Multistage build for exodus panel (frontend + backend)
FROM node:20-alpine AS panel-ui
WORKDIR /ui
ENV NODE_OPTIONS=--max-old-space-size=4096
ARG SINGBOX_WASM_URL=https://adam-sizzler.github.io/s-validator/main.wasm
ARG SINGBOX_SCHEMA_URL=https://adam-sizzler.github.io/s-validator/singbox.schema.json
ARG SINGBOX_SCHEMA_CN_URL=https://adam-sizzler.github.io/s-validator/singbox.schema.json
COPY frontend/package*.json ./
RUN --mount=type=cache,target=/root/.npm npm ci --legacy-peer-deps --no-audit --prefer-offline
COPY frontend/ ./
RUN mkdir -p /ui/public/assets
RUN set -eu; \
    fetch_with_retry() { \
      url="$1"; out="$2"; name="$3"; \
      [ -n "$url" ] || return 0; \
      n=1; ok=0; \
      while [ "$n" -le 4 ]; do \
        if wget -qO "$out" "$url"; then \
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
    fetch_with_retry "$SINGBOX_WASM_URL" /ui/public/assets/main.wasm main.wasm; \
    fetch_with_retry "$SINGBOX_SCHEMA_URL" /ui/public/assets/singbox.schema.json singbox.schema.json; \
    fetch_with_retry "$SINGBOX_SCHEMA_CN_URL" /ui/public/assets/singbox.schema.cn.json singbox.schema.cn.json
RUN if [ -f /ui/public/assets/main.wasm ]; then \
      magic="$(od -An -t x1 -N 4 /ui/public/assets/main.wasm | tr -d ' \n')"; \
      [ "$magic" = "0061736d" ] || { echo "Invalid WASM magic for main.wasm: $magic"; exit 1; }; \
      echo "WASM artifact attached: /ui/public/assets/main.wasm"; \
    else \
      echo "WARN: /ui/public/assets/main.wasm is missing. Sing-box validator will be disabled in UI."; \
    fi
RUN --mount=type=cache,target=/root/.npm npm run cb

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

WORKDIR /build

# Copy go mod files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build panel backend binary
RUN CGO_ENABLED=0 \
    GOOS=linux \
    go build \
        -tags "none" \
        -trimpath \
        -buildvcs=true \
        -ldflags="-s -w \
                -X 'exodus/constant.Version=${VERSION}' \
                -X 'exodus/constant.Revision=${REVISION}' \
                -X 'exodus/constant.BuildTags=none' \
                -X 'exodus/constant.CgoEnabled=0'" \
    -o exodus \
    ./backend

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

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/exodus /app/

# Copy built frontend
COPY --from=panel-ui /ui/dist /app/ui
COPY --from=builder /usr/local/go/lib/wasm/wasm_exec.js /app/ui/assets/wasm_exec.js

# Create directories for data and certificates
RUN mkdir -p /app/data /app/certs

EXPOSE 3000

ENTRYPOINT ["/app/exodus"]
