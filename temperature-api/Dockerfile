# Этап сборки
FROM golang:1.21-alpine as builder

WORKDIR /app

# Копируем файлы зависимостей для кэширования
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем исходный код
COPY . /app

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o temperature-api

# Финальный этап
FROM alpine:latest

WORKDIR /app

# Копируем собранное приложение из этапа сборки
COPY --from=builder /app/temperature-api .

# Устанавливаем необходимые пакеты для runtime
RUN apk --no-cache add ca-certificates

# Определяем команду для запуска приложения
CMD ["./temperature-api"]
