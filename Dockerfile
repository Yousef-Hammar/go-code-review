FROM golang:1.22.7-alpine3.20 AS builder

RUN apk add --no-cache git gcc musl-dev

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

WORKDIR /app/cmd/coupon_service
RUN go build -o /app/bin/main .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/bin/main .

EXPOSE 8080

ENTRYPOINT ["./main"]