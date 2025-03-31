# Используем официальный образ Go
FROM golang:1.23.6 AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы проекта
COPY . .

# Скачиваем зависимости
RUN go mod download

# Собираем бинарник
RUN go build -o main .

# Используем минимальный образ для финального контейнера
FROM alpine:latest
WORKDIR /root/

RUN apt-get update && apt-get install -y postgresql-client

# Копируем скомпилированное приложение
COPY --from=builder /app/main .

# Открываем нужный порт
EXPOSE 8080

# Запускаем сервер
CMD ["./main"]