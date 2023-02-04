FROM golang:1.17-alpine AS builder
LABEL stage=builder
WORKDIR /app
COPY . .

RUN apk add build-base && go build -o Forum cmd/main.go

FROM alpine:3.6
WORKDIR /app
LABEL authors="Certina01 && AidanaErbolatova" project="forum"
COPY --from=builder /app .
CMD ["/app/Forum"]