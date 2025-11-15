# ---------- STAGE 1: BUILD ----------
FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build binary khusus fat
RUN go build -o fat-service ./fat-amk/cmd/fat

# ---------- STAGE 2: RUNTIME ----------
FROM alpine:3.19

WORKDIR /app

COPY --from=builder /app/fat-service .

EXPOSE 8082

CMD ["./fat-service"]
