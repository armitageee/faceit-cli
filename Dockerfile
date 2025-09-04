# Multi-stage build for faceit-cli
FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o faceit-cli main.go

# Final stage - distroless image
FROM gcr.io/distroless/static-debian11

# Copy the binary
COPY --from=builder /app/faceit-cli /faceit-cli

# Copy timezone data
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Set timezone
ENV TZ=UTC

# Expose port (if needed for health checks)
EXPOSE 8080

# Run the application
ENTRYPOINT ["/faceit-cli"]
