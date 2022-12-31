# syntax=docker/dockerfile:1
FROM golang:1.19-alpine

WORKDIR /app/go-mongo-auth

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o ./cmd/build ./cmd/go-mongo-auth/main.go

RUN chmod +x ./docker-entrypoint.sh

EXPOSE 8091

ENTRYPOINT [ "./docker-entrypoint.sh" ]
