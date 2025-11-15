# ---------- STAGE 1: BUILD ----------
FROM golang:1.22-alpine AS builder

WORKDIR /src

# Copy dependency info dari root monorepo
COPY go.mod go.sum ./
RUN go mod download

# Copy seluruh source code monorepo
COPY . .

# Build binary khusus auth
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /out/auth-service ./auth-amk/cmd/auth

# ---------- STAGE 2: RUNTIME ----------
FROM alpine:3.19

WORKDIR /app

# copy binary dari builder
COPY --from=builder /out/auth-service /app/auth-service

EXPOSE 8080

ENTRYPOINT ["/app/auth-service"]
