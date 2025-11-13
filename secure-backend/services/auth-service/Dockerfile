FROM golang:1.25.3-alpine

WORKDIR /app

# Copiar TODO el código
COPY . .

# Compilar
RUN go mod tidy && go build -o auth-service ./cmd/server/

EXPOSE 8442

CMD ["./auth-service"]