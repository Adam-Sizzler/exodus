# Multistage build for exodus panel (frontend + backend)
FROM node:22-alpine AS panel-ui
WORKDIR /ui
ENV NODE_OPTIONS=--max-old-space-size=4096

ARG SINGBOX_ASSETS_URL=https://adam-sizzler.github.io/s-validator

COPY frontend/package*.json ./
RUN --mount=type=cache,target=/root/.npm npm install --legacy-peer-deps --no-audit --prefer-offline

COPY frontend/ ./

RUN if [ ! -f public/assets/main.wasm ] || [ ! -f public/assets/wasm_exec.js ] || [ ! -f public/assets/singbox.schema.json ]; then \
      apk add --no-cache curl \
      && mkdir -p public/assets \
      && curl -L ${SINGBOX_ASSETS_URL}/wasm_exec.js -o public/assets/wasm_exec.js \
      && curl -L ${SINGBOX_ASSETS_URL}/singbox.schema.json -o public/assets/singbox.schema.json \
      && curl -L ${SINGBOX_ASSETS_URL}/main.wasm -o public/assets/main.wasm; \
    else \
      echo "Assets already present locally, skipping download"; \
    fi

RUN --mount=type=cache,target=/root/.npm --mount=type=cache,target=/ui/node_modules/.vite/cache \
    if [ -d dist ]; then \
      echo "Using existing dist directory from host"; \
    else \
      npm run cb; \
    fi

FROM golang:1.25.12-alpine AS builder

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
