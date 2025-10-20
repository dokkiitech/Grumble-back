# syntax=docker/dockerfile:1

FROM golang:1.23-alpine AS dev

WORKDIR /app

# 開発で必要になりがちなツールを追加。
RUN apk add --no-cache git curl

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

CMD ["go", "run", "./cmd/api"]
