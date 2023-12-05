# syntax=docker/dockerfile:1

FROM golang:1.21.4

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY /templates/ client.go room.go main.go ./

RUN go build -o godocker

EXPOSE 8080

CMD ["./godocker"]