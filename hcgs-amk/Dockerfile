# ---------- STAGE 1: BUILD ----------
FROM golang:1.22-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build binary khusus hcgs
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /out/hcgs-service ./hcgs-amk/cmd/hcgs

# ---------- STAGE 2: RUNTIME ----------
FROM alpine:3.19

WORKDIR /app

COPY --from=builder /out/hcgs-service /app/hcgs-service

EXPOSE 8081

ENTRYPOINT ["/app/hcgs-service"]

