FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum .
RUN go mod download
COPY . .
RUN go build -o /bin/run ./cmd/app/main.go

FROM alpine:3.18
WORKDIR /app
COPY --from=builder /bin/run .
ENTRYPOINT ["./run"]
