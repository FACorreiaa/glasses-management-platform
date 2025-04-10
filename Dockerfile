FROM node:alpine-slim as assets
WORKDIR /app
COPY package.json package-lock.json postcss.config.cjs fonts.css tailwind.css ./
RUN mkdir -p app/static/css app/static/fonts
RUN npm install --only=production --ci
RUN npm run fonts && npm run tailwind-build

# Define the "base" stage
FROM golang:alpine as base
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . .

# Define the final stage
FROM base as production
WORKDIR /app
COPY --from=base /app/app/static/css/output.css ./controller/static/css/
COPY --from=base /app/app/static/fonts/* ./controller/static/fonts/
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /app/server
#RUN upx --best --lzma /app/server
EXPOSE 6968
ENTRYPOINT ["/app/server"]
