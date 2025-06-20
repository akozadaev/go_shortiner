# syntax=docker/dockerfile:1

FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o go_shurtiner ./cmd/go_shurtiner/main.go

FROM alpine:latest

RUN adduser -D -u 10001 user

WORKDIR /app

COPY --from=builder /app/go_shurtiner ./shurtiner
COPY ./config/config.docker.yaml ./config/config.local.yaml
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

RUN mkdir -p /app/logs && chown -R user:user /app
RUN chmod +rw /app/*

#USER user

EXPOSE 8080

ENTRYPOINT ["/entrypoint.sh"]
CMD ["./shurtiner"]
