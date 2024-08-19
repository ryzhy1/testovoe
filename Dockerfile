FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod tidy
COPY . .
RUN go build -o auth-service ./cmd/main.go
FROM alpine:3.18
RUN apk add --no-cache postgresql-client bash netcat-openbsd
WORKDIR /app
COPY --from=builder /app/auth-service .
COPY config /app/config
EXPOSE 8080
CMD ["./auth-service"]
