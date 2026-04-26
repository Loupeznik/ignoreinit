FROM golang:1.26.2-alpine3.23 AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/ignoreinit .

FROM alpine:3.23.4

RUN apk add --no-cache ca-certificates && addgroup -S ignoreinit && adduser -S -G ignoreinit ignoreinit && mkdir -p /work && chown ignoreinit:ignoreinit /work

WORKDIR /work

COPY --from=builder /out/ignoreinit /usr/local/bin/ignoreinit

USER ignoreinit:ignoreinit

ENTRYPOINT ["/usr/local/bin/ignoreinit"]
