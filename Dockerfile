# Этап сборки (Builder)
FROM golang:1.22-alpine AS builder

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем go.mod и go.sum для установки зависимостей
COPY go.mod go.sum ./
RUN go mod tidy

# Копируем все остальные файлы проекта
COPY . .

# Собираем приложение
RUN go build -o auth-service ./cmd/main.go

# Этап запуска (Production)
FROM alpine:3.18

# Устанавливаем зависимости для работы с базой данных
RUN apk add --no-cache postgresql-client bash netcat-openbsd

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем бинарный файл из предыдущего этапа
COPY --from=builder /app/auth-service .

# Копируем конфигурационный файл в контейнер
COPY config /app/config

# Указываем порт, на котором будет работать приложение
EXPOSE 8080

# Указываем команду для запуска приложения с ожиданием базы данных
CMD ["./auth-service"]
