FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o auth-service ./cmd/app

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/auth-service .


COPY local.yaml ./config.yaml

EXPOSE 8765

CMD ["./auth-service"]