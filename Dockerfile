FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY internal ./internal
COPY cmd ./cmd

RUN go build -o /app/go_bot ./cmd/

EXPOSE 8080

CMD ["/app/go_bot"]
