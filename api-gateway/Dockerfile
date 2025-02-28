# Build stage
FROM golang:1.23.1-alpine AS builder

# Set the working directory
WORKDIR /build

# Copy the entire repository into the build context
COPY . .

# Switch to the api-gateway directory
WORKDIR /build/api-gateway

# Download dependencies
RUN go mod download

# Build the binary statically
RUN CGO_ENABLED=0 GOOS=linux go build -o api-gateway ./cmd

# Final stage
FROM alpine:latest

# Install CA certificates
RUN apk --no-cache add ca-certificates

# Set the working directory
WORKDIR /root/

# Copy the built binary from the builder stage
COPY --from=builder /build/api-gateway/api-gateway .

# Expose the port
EXPOSE 8080

# Run the service
CMD ["./api-gateway"]