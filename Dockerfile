# syntax=docker/dockerfile:1

# Stage 1: Build Assets (Tailwind CSS)
FROM node:22-alpine AS assets
WORKDIR /app
COPY package.json package-lock.json postcss.config.cjs tailwind.config.mjs ./
# Copy source files needed for Tailwind content scanning
COPY app/static/css/main.css ./app/static/css/main.css
COPY app/view ./app/view
# npm ci installs ALL deps (build tools live in devDependencies)
RUN --mount=type=cache,target=/root/.npm \
    npm ci
RUN npm run tailwind-build

# Stage 2: Build Go Binary
FROM golang:1.26-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download
COPY . .
# Overwrite static CSS with the one built in Stage 1
COPY --from=assets /app/app/static/css/output.css ./app/static/css/output.css
# Build static binary
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 go build -ldflags="-s -w" -o server .

# Final Stage: Runtime
FROM alpine:3.21 AS production
WORKDIR /app
RUN apk add --no-cache ca-certificates tzdata \
    && adduser -D -u 10001 appuser
COPY --from=builder /app/server ./server
USER appuser
EXPOSE 8080
ENTRYPOINT ["./server"]
