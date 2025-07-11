FROM golang:1.24.5-alpine3.21 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o bin/main cmd/server/main.go

FROM alpine:latest

WORKDIR /root

COPY --from=builder /app/bin/main .
COPY --from=builder /app/.env .

CMD ["./main"]