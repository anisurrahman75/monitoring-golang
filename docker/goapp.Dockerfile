FROM golang:1.23.6 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o monitoring-app ./cmd/main.go
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/monitoring-app .
EXPOSE 8080
CMD ["./monitoring-app"]
