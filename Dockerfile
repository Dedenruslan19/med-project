# Stage 1: Build the Go application binary
FROM golang:1.25.3 AS builder
WORKDIR /app
COPY . .
WORKDIR /app/cmd/echo-server
RUN go build -o echo-server main.go

# Stage 2: Create a minimal runtime image and running the application
FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates && \
    rm -rf /var/lib/apt/lists/*
COPY .env .
COPY --from=builder /app/app/echo-server/echo-server /echo-server
CMD ["/echo-server"]