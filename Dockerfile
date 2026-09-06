FROM node:22-alpine AS frontend-build
WORKDIR /opt/app/frontend

COPY frontend/package*.json ./
COPY frontend/vendor/ ./vendor/
RUN --mount=type=cache,target=/root/.npm npm ci --legacy-peer-deps

COPY frontend/ ./
RUN --mount=type=cache,target=/root/.npm npm run start:build

FROM golang:1.27.1-alpine AS backend-build
WORKDIR /opt/app

COPY backend/go.mod backend/go.sum ./backend/
WORKDIR /opt/app/backend
RUN go mod download

WORKDIR /opt/app
COPY backend/ ./backend/

ARG SUB_APP_VERSION
WORKDIR /opt/app/backend
RUN set -eu; \
    version="${SUB_APP_VERSION:-unknown}"; \
    CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w -X main.buildVersion=${version}" -o /opt/app/subscription-page ./cmd/subscription-page

FROM alpine:3.23
WORKDIR /opt/app

LABEL org.opencontainers.image.title="Exodus Subscription Page" \
      org.opencontainers.image.description="Exodus Subscription Page" \
      org.opencontainers.image.url="https://github.com/adam-sizzler/exodus-subscription" \
      org.opencontainers.image.source="https://github.com/adam-sizzler/exodus-subscription" \
      org.opencontainers.image.vendor="Exodus" \
      org.opencontainers.image.licenses="AGPL-3.0" \
      org.opencontainers.image.documentation="https://docs.ex"

RUN apk add --no-cache ca-certificates && mkdir -p /opt/app/ruleset

COPY --from=backend-build /opt/app/subscription-page ./subscription-page
COPY --from=frontend-build /opt/app/frontend/dist/ ./frontend/
COPY docker-entrypoint.sh ./

ENTRYPOINT [ "/bin/sh", "docker-entrypoint.sh" ]
CMD [ "/opt/app/subscription-page" ]
