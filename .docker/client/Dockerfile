FROM golang:1.24

WORKDIR /app

# Копируем весь проект, включая go.mod, internal и cmd
COPY . .

# Применяем replace и скачиваем зависимости
RUN go mod tidy

# Собираем клиент
WORKDIR /app/cmd/client
RUN CGO_ENABLED=1 GOOS=linux go build -o /app/client

# Заглушка — можно заменить на запуск клиента
CMD ["tail", "-f", "/dev/null"]