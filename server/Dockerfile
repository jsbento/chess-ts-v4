FROM golang:1.24.4-alpine

ENV PROJECT_DIR=/app
ENV CGO_ENABLED=0
ENV GO111MODULE=on

WORKDIR /app

# Install curl for healthchecks
RUN apk add --no-cache curl

COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Install CompileDaemon for hot reloading
RUN go install github.com/githubnemo/CompileDaemon@latest

EXPOSE 8080

# Use CompileDaemon for hot reloading - it will watch for changes and rebuild
ENTRYPOINT ["CompileDaemon", "-build=go build -o /build/api-server ./cmd/api-server", "-command=/build/api-server", "--polling"]