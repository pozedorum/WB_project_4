FROM golang:1.24.4-alpine

WORKDIR /app

# Установка зависимостей
RUN apk add --no-cache git postgresql-client

# Копируем сначала только mod файлы для кэширования зависимостей
COPY go.mod go.sum ./

# Скачиваем зависимости
RUN go mod download

# Копируем остальные файлы проекта
COPY . .

# Собираем приложение-trimpath -gcflags="all=-N -l"
#RUN go build -o url-shortener ./cmd/main.go
RUN go -trimpath -gcflags="all=-N -l" build -o url-shortener ./cmd/main.go
EXPOSE 8080 6060 5432

CMD ["./url-shortener"]