# syntax=docker/dockerfile:1
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o logger main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/logger .
EXPOSE 8080 9000
ENV UI_PORT=8080
ENV LOG_INGEST_PORT=9000
CMD ["./logger"] 