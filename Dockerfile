# syntax=docker/dockerfile:1

FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-w -s" -o /app/watchlist-api ./cmd/main.go

FROM alpine:latest

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/watchlist-api /app/watchlist-api

RUN adduser -D -u 1000 appuser
USER appuser

EXPOSE 8080

CMD ["/app/watchlist-api"]
