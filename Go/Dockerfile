FROM golang:1.23-alpine as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod tidy
COPY . .
RUN go build -o main .
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app .
CMD ["./main", "--seed-user"]
