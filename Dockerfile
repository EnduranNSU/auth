FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o auth ./cmd/auth

FROM alpine:3.20

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/auth /app/auth

COPY config/config.yaml /app/config/config.yaml

ENV GIN_MODE=release

EXPOSE 8080

ENTRYPOINT ["/app/auth"]
