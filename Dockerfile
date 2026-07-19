FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /todo main.go

FROM ubuntu:latest

WORKDIR /app

COPY --from=builder /todo /app/todo
COPY --from=builder /app/web /app/web

ENV TODO_PORT=7540
ENV TODO_DBFILE=/app/scheduler.db

EXPOSE 7540

CMD ["/app/todo"]
