# Build stage
FROM golang:1.23.1-alpine AS builder

# Set the working directory
WORKDIR /build

# Copy the entire repository into the build context
COPY . .

# Switch to the internal-service directory
WORKDIR /build/internal-service

# Download dependencies
RUN go mod download

# Build the binary from the cmd directory
RUN CGO_ENABLED=0 GOOS=linux go build -o internal-service ./cmd

# Final stage
FROM alpine:latest

# Install CA certificates
RUN apk --no-cache add ca-certificates

# Set the working directory
WORKDIR /root/

# Copy the built binary from the builder stage
COPY --from=builder /build/internal-service/internal-service .

# Expose the port
EXPOSE 50051

# Run the service
CMD ["./internal-service"]