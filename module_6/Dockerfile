# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /app/server ./cmd/server/main.go

# Run stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/server .

# Копіюємо .env, якщо ви використовуєте його для конфігурації
# Якщо конфігурація через ENV_VAR, то цей рядок не потрібен
# COPY .env .

EXPOSE 8080 # Або порт, який ви використовуєте

CMD ["./server"]