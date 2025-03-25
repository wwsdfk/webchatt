FROM golang:1.23

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

RUN apt-get update && apt-get install -y postgresql-client

RUN go build -o main .

CMD ["/app/main"]