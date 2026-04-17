# Build stage
FROM golang:1.25-alpine AS builder

# Install git and ca-certificates (needed for go modules)
RUN apk update && apk add --no-cache git ca-certificates upx && update-ca-certificates

WORKDIR /build

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary with optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' -a \
    -installsuffix cgo -o main ./cmd/service

# Compress the binary
RUN upx --best --lzma main

# Final stage
FROM scratch

# Copy CA certificates for HTTPS
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the compressed binary
COPY --from=builder /build/main /main

# Expose port
EXPOSE 3000

# Run the binary
ENTRYPOINT ["/main", "app", "start"]