# ---------- STAGE 1: BUILD ----------
FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build binary khusus hcgs
RUN go build -o hcgs-service ./hcgs-amk/cmd/hcgs

# ---------- STAGE 2: RUNTIME ----------
FROM alpine:3.19

WORKDIR /app

COPY --from=builder /app/hcgs-service .

EXPOSE 8081

CMD ["./hcgs-service"]

