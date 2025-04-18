FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -o openvpn

FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Copy the binary from builder
COPY --from=builder /app/openvpn /usr/local/bin/openvpn

# Create necessary directories
RUN mkdir -p /etc/openvpn-admin

# Copy configuration files
COPY --from=builder /app/.env /etc/openvpn-admin/.env

WORKDIR /etc/openvpn-admin

# Set the entrypoint
ENTRYPOINT ["openvpn"] 
