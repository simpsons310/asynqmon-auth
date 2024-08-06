# Build stage
FROM golang:1.21.1-alpine as build

ARG CGO_ENABLED=0

WORKDIR /app

COPY go.* ./
COPY cmd cmd
COPY internal internal

RUN go mod tidy

RUN go build -v -o ./asynqmon_auth cmd/httpserver/main.go

## Runtime stage
FROM alpine:3.16.2

WORKDIR /app

COPY --from=build /app/asynqmon_auth /app/asynqmon_auth

EXPOSE 8080

CMD ["/app/asynqmon_auth", "8080"]