# syntax=docker/dockerfile:1
FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o shortiner ./cmd/go_shurtiner/main.go

FROM debian:bookworm-slim

WORKDIR /app

RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/shortiner .

 COPY ./config/config.local.yaml ./config/

ENV DB_DSN="host=localhost user=short password=short dbname=short port=5433 sslmode=disable TimeZone=Europe/Moscow"

EXPOSE 2025

CMD ["./shortiner"]
