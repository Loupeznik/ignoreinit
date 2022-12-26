FROM golang:1.19-alpine as builder

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
