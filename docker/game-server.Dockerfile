# Game Server Dockerfile
FROM golang:1.22-alpine AS builder
RUN apk add --no-cache gcc musl-dev git
WORKDIR /app
COPY go.mod go.sum* ./
RUN go mod download 2>/dev/null || true
COPY . .
RUN CGO_ENABLED=1 go build -o /game-server ./cmd/game-server

FROM alpine:3.20
RUN apk add --no-cache ca-certificates
COPY --from=builder /game-server /game-server
COPY items.dat /items.dat
EXPOSE 17091/udp
CMD ["/game-server"]
