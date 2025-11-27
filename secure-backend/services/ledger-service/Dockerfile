FROM golang:1.25.3-alpine

WORKDIR /app

COPY . .

RUN go mod tidy && go build -o ledger-service ./cmd/server/

EXPOSE 8443

CMD ["./ledger-service"]