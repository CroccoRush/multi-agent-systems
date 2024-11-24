# Используем официальный образ Go в качестве базового образа
FROM golang:1.22-alpine

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем файлы зависимостей
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем исходный код в контейнер
COPY . .

# Собираем приложение
RUN go build -o app multi-agent-systems/t1/cmd/experiment

# Запускаем приложение
CMD ["./app"]