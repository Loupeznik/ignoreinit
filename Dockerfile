FROM golang:1.23-alpine AS builder

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build -o ./ignoreinit ./main.go

FROM alpine:latest

WORKDIR /work

COPY --from=builder /build/ignoreinit /app/ignoreinit

ENTRYPOINT ["/app/ignoreinit"]
