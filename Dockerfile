FROM golang:1.25-alpine AS backend-build
WORKDIR /opt/app

COPY backend/go.mod backend/go.sum ./backend/
WORKDIR /opt/app/backend
RUN go mod download

WORKDIR /opt/app
COPY backend/ ./backend/

ARG SUB_APP_VERSION=unknown
WORKDIR /opt/app/backend
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w -X main.buildVersion=${SUB_APP_VERSION}" -o /opt/app/subscription-page ./cmd/subscription-page

FROM alpine:3.23
WORKDIR /opt/app

RUN apk add --no-cache ca-certificates

COPY --from=backend-build /opt/app/subscription-page ./subscription-page
COPY frontend/dist/ ./frontend/

ENTRYPOINT ["/opt/app/subscription-page"]
