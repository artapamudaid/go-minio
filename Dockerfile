# ==========================================
# STAGE 1: Build Binary
# ==========================================
FROM golang:1.25-bookworm AS builder

RUN apt-get update && apt-get install -y --no-install-recommends git ca-certificates && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy dependency files lebih dulu untuk memanfaatkan Docker Caching
COPY go.mod go.sum ./
RUN go mod download

# Copy seluruh source code
COPY . .

# Build binary Go dengan optimasi ukuran (-ldflags="-w -s")
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o server ./cmd/server/main.go

# ==========================================
# STAGE 2: Lightweight Runner
# ==========================================
FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates tzdata && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy binary dari stage builder
COPY --from=builder /app/server .

# Expose port (default 8080)
EXPOSE 8080

# Jalankan aplikasi
CMD ["./server"]