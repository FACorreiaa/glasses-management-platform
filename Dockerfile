# Stage 1: Build Assets
FROM node:alpine-slim AS assets
WORKDIR /app
COPY package.json package-lock.json postcss.config.cjs tailwind.config.mjs ./
# Copy source files needed for Tailwind
COPY app/static/css/main.css ./app/static/css/main.css
COPY app/view ./app/view
RUN npm install --only=production --ci
RUN npm run tailwind-build

# Stage 2: Build Go Binary
FROM golang:1.24-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# Overwrite static CSS with the one built in Stage 1
COPY --from=assets /app/app/static/css/output.css ./app/static/css/output.css
# Build the binary
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o server .

# Final Stage: Runtime
FROM alpine:latest
WORKDIR /app
RUN apk add --no-cache ca-certificates tzdata
# Copy only the binary from the builder
COPY --from=builder /app/server ./server
EXPOSE 8080
ENTRYPOINT ["./server"]
