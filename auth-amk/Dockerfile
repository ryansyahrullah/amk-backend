# ---------- STAGE 1: BUILD ----------
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy dependency info dari root monorepo
COPY go.mod go.sum ./
RUN go mod download

# Copy seluruh source code monorepo
COPY . .

# Build binary khusus auth
RUN go build -o auth-service ./auth-amk/cmd/auth

# ---------- STAGE 2: RUNTIME ----------
FROM alpine:3.19

WORKDIR /app

# copy binary dari builder
COPY --from=builder /app/auth-service .

EXPOSE 8080

CMD ["./auth-service"]
