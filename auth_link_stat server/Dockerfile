# Первый этап: сборка приложения
FROM golang:1.23.3 AS builder

WORKDIR /app

# Устанавливаем pg_isready и очищаем кеш
RUN apt-get update && apt-get install -y postgresql-client \
    && rm -rf /var/lib/apt/lists/* && apt-get clean

# Устанавливаем переменные окружения
ENV SECRET="your"
ENV CGO_ENABLED=0
# Отключаем CGO для статической компиляции

# Копируем файлы зависимостей
COPY go.mod go.sum ./

# Скачиваем зависимости
RUN go mod download && go mod verify

# Копируем весь код
COPY . .

# Компилируем бинарник
RUN go build -o /app/auth_link_stat ./cmd/main.go



# Второй этап: финальный образ (без лишних инструментов)
FROM alpine:latest

WORKDIR /app

# Устанавливаем postgresql-client
RUN apk add --no-cache postgresql-client

COPY .env /app/.env

# Копируем бинарный файл из предыдущего этапа
COPY --from=builder /app/auth_link_stat /app/auth_link_stat

# Копируем wait-for-db.sh и делаем исполняемым
COPY --from=builder /app/wait-for-db.sh /app/wait-for-db.sh
RUN chmod +x /app/wait-for-db.sh

# Запуск приложения
CMD ["/app/auth_link_stat"]
