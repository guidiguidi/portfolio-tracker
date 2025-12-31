FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o portfolio-tracker ./cmd/api

FROM alpine:3.20

WORKDIR /app

ENV APP_ENV=local

COPY --from=builder /app/portfolio-tracker /app/portfolio-tracker
COPY config ./config

EXPOSE 8080

CMD ["/app/portfolio-tracker"]
