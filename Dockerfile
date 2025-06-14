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

FROM ubuntu:22.04

# Install runtime dependencies
RUN apt-get update && apt-get install -y \
    ca-certificates \
    tzdata \
    netcat-openbsd \
    && rm -rf /var/lib/apt/lists/*

# Create app directory
WORKDIR /openvpn-admin

# Copy the binary from builder
COPY --from=builder /app/openvpn /openvpn-admin/openvpn

# Copy configuration files
COPY --from=builder /app/.env /openvpn-admin/.env

# Create templates directory and copy source blacklist script and file
RUN mkdir -p /openvpn-admin/templates
COPY scripts/auth-blacklist.sh /openvpn-admin/templates/auth-blacklist.sh
COPY scripts/blacklist.txt /openvpn-admin/templates/blacklist.txt

# Set the entrypoint
ENTRYPOINT ["/openvpn-admin/openvpn"] 
