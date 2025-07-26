# Используем официальный образ Go в качестве базового образа для сборки
FROM golang:1.22-alpine AS builder

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем go.mod и go.sum для загрузки зависимостей
COPY go.mod go.sum ./

# Загружаем зависимости Go
RUN go mod download

# Копируем исходный код приложения
COPY . .

# Собираем приложение
# -o /app/anarchy-core: указывает имя выходного файла
# ./cmd/anarchy-core: указывает путь к main.go
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o /app/anarchy-core ./cmd/anarchy-core

# Используем облегченный образ Alpine Linux для финального образа
FROM alpine:latest

# Устанавливаем ca-certificates для работы с HTTPS
RUN apk --no-cache add ca-certificates

# Устанавливаем рабочую директорию
WORKDIR /root/

# Копируем скомпилированный бинарник из образа builder
COPY --from=builder /app/anarchy-core .

# Копируем .env файл. В продакшене лучше использовать переменные окружения напрямую.
# Если вы используете .env для локальной разработки, убедитесь, что он не содержит чувствительных данных в продакшене.
COPY .env .

# Открываем порт, на котором будет работать приложение
EXPOSE 8080

# Запускаем приложение
CMD ["./anarchy-core"]

