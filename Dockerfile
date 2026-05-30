# Stage 1: Build Assets
FROM node:alpine-slim AS assets
WORKDIR /app
COPY package.json package-lock.json postcss.config.cjs tailwind.config.mjs ./
# Copy source files needed for Tailwind (templ, html, etc)
COPY app/static/css/main.css ./app/static/css/main.css
COPY app/view ./app/view
RUN npm install --only=production --ci
RUN npm run tailwind-build

# Stage 2: Build Go Binary
FROM golang:1.24-alpine AS base
WORKDIR /app

# Install build dependencies if needed
RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

# Copy source and assets
COPY . .
# Overwrite static assets with those built in Stage 1
COPY --from=assets /app/app/static/css/output.css ./app/static/css/output.css

# Define the final stage for production
FROM base AS production
WORKDIR /app
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /app/server
EXPOSE 8080
ENTRYPOINT ["/app/server"]
