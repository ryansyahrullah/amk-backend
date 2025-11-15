# ---------- STAGE 1: BUILD ----------
FROM golang:1.22-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build binary khusus fat
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /out/fat-service ./fat-amk/cmd/fat

# ---------- STAGE 2: RUNTIME ----------
FROM alpine:3.19

WORKDIR /app

COPY --from=builder /out/fat-service /app/fat-service

EXPOSE 8082

ENTRYPOINT ["/app/fat-service"]
