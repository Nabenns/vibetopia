# Login Server Dockerfile
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum* ./
RUN go mod download 2>/dev/null || true
COPY . .
RUN CGO_ENABLED=0 go build -o /login-server ./cmd/login-server

FROM alpine:3.20
RUN apk add --no-cache ca-certificates
COPY --from=builder /login-server /login-server
EXPOSE 8080
CMD ["/login-server"]
