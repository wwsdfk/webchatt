#!/bin/sh
# Ждем, пока БД будет доступна
until pg_isready -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER"; do
  echo "Жду БД..."
  sleep 2
done

# Запускаем миграции (если они есть)
if [ -f "migrations/init.sql" ]; then
  echo "Запускаю миграции..."
  psql -h "$DB_HOST" -U "$DB_USER" -d "$DB_NAME" -f migrations/init.sql
fi

# Запускаем бэкенд
exec go run main.go