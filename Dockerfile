
FROM golang:1.24 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/pet-shop ./cmd/app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/migrator ./cmd/migrator

FROM alpine:3.20
RUN adduser -D appuser
WORKDIR /app

COPY --from=builder /out/pet-shop /app/pet-shop
COPY --from=builder /out/migrator /app/migrator
COPY config ./config
COPY migrations ./migrations

EXPOSE 3000
USER appuser

CMD ["/app/pet-shop"]