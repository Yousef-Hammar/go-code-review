FROM golang:alpine AS builder

RUN apk add --no-cache git gcc musl-dev

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

WORKDIR /app/cmd/coupon_service
RUN go build -o /app/bin/main .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/bin/main .

EXPOSE 8080

ENTRYPOINT ["./main"]