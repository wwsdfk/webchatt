# Строим React-приложение
FROM node:18-alpine AS build

# Устанавливаем рабочую директорию для проекта
WORKDIR /app

# Копируем package.json и устанавливаем зависимости
COPY package*.json ./
RUN npm install

# Копируем все остальные файлы и строим приложение
COPY . .
RUN npm run build

# Запускаем через Nginx
FROM nginx:alpine

# Копируем скомпилированные файлы из предыдущего этапа
COPY --from=build /app/build /usr/share/nginx/html

# Экспонируем порт 80 для фронтенда
EXPOSE 80

# Запуск Nginx
CMD ["nginx", "-g", "daemon off;"]
