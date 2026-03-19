FROM golang:1.24-alpine AS backend-build
WORKDIR /opt/app

COPY backend/ ./backend/

WORKDIR /opt/app/backend

RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /opt/app/subscription-page ./cmd/subscription-page

FROM alpine:3.21
WORKDIR /opt/app

RUN apk add --no-cache ca-certificates

COPY --from=backend-build /opt/app/subscription-page ./subscription-page
COPY frontend/dist/ ./frontend/

ENTRYPOINT ["/opt/app/subscription-page"]
