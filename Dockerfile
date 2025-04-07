FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o Go_wallet .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/Go_wallet .
COPY config.env .

EXPOSE 8080

CMD ["./Go_wallet"]