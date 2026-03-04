# Этап 1: Сборка приложения
FROM golang:1.25.1-alpine AS builder

WORKDIR /application

# Копируем go.mod и go.sum для кэширования зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем остальной код
COPY . .

# Компилируем приложение
RUN go build -o ./app ./cmd/app

# Этап 2: Создание минимального итогового образа
FROM alpine:latest AS final

WORKDIR /application

# Устанавливаем переменную окружения для порта
ARG APP_PORT=8080
ENV PORT=${APP_PORT}

# Копируем бинарник, миграции и конфигурационный файл с сервисами
COPY --from=builder /application/app /usr/local/bin/app
COPY --from=builder /application/internal/database/migrations ./internal/database/migrations
COPY --from=builder /application/services.yml ./services.yml
COPY --from=builder /application/web ./web

# Указываем порт
EXPOSE ${APP_PORT}

# Запускаем приложение
CMD ["app"]
