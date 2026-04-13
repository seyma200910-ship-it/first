FROM golang:1.26.1 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/app

FROM alpine:3.20

WORKDIR /root/

COPY --from=builder /app/app .

EXPOSE 8080

CMD ["./app"]