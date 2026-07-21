# ==========================================
# STAGE 1: Build Binary
# ==========================================
FROM golang:1.22-alpine AS builder

# Install ca-certificates untuk mendukung HTTPS ke S3/MinIO
RUN apk add --no-cache git ca-certificates

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
FROM alpine:latest

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Copy binary dari stage builder
COPY --from=builder /app/server .

# Expose port (default 8181)
EXPOSE 8181

# Jalankan aplikasi
CMD ["./server"]