FROM golang:1.23.5

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o cred-service .

EXPOSE 4000

CMD ["./cred-service"]
