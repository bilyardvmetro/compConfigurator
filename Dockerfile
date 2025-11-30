FROM golang:1.23-alpine AS builder
LABEL authors="bilyardvmetro"

WORKDIR /app
COPY go.mod ./
RUN go mod dowload

COPY . .
RUN go build -o server ./cmd/app

FROM alpine:3.20

WORKDIR /app
COPY --fro=builder /app/server .

EXPOSE 8080

ENTRYPOINT ["top", "-b"]