# Docker Support

This document describes how to use Docker with the faceit-cli application.

## Quick Start

### Build and Run

```bash
# Build Docker image
make docker-build

# Run with environment file
make docker-run
```

## Docker Image Details

### Multi-stage Build

The Dockerfile uses a multi-stage build approach:

1. **Builder stage**: Uses `golang:1.23-alpine` for compilation
2. **Final stage**: Uses `gcr.io/distroless/static-debian11` for minimal runtime

### Image Size

- **Builder stage**: ~300MB (includes Go toolchain)
- **Final stage**: ~20MB (distroless + binary only)

### Security Features

- **Distroless base image**: No shell, package manager, or unnecessary tools
- **Non-root user**: Runs as non-root for security
- **Minimal attack surface**: Only the application binary and required libraries

## Environment Variables

The Docker container supports all the same environment variables as the native application:

```bash
# Required
FACEIT_API_KEY=your_api_key

# Optional
LOG_LEVEL=info
KAFKA_ENABLED=false
PRODUCTION_MODE=false
LOG_TO_STDOUT=true
CACHE_ENABLED=false
MATCHES_PER_PAGE=10
MAX_MATCHES_TO_LOAD=100
COMPARISON_MATCHES=20
```

## Docker Compose Integration

You can run the application with Kafka using Docker Compose:

```yaml
version: '3.8'
services:
  faceit-cli:
    build: .
    environment:
      - FACEIT_API_KEY=${FACEIT_API_KEY}
      - KAFKA_ENABLED=true
      - KAFKA_BROKERS=kafka:9092
      - LOG_LEVEL=info
    depends_on:
      - kafka
    stdin_open: true
    tty: true

  kafka:
    image: confluentinc/cp-kafka:latest
    environment:
      KAFKA_BROKER_ID: 1
      KAFKA_ZOOKEEPER_CONNECT: zookeeper:2181
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://kafka:9092
      KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR: 1
    ports:
      - "9092:9092"
```

## GitHub Container Registry

The application automatically builds and pushes Docker images to GitHub Container Registry:

- **Registry**: `ghcr.io/armitageee/faceit-cli`
- **Tags**: 
  - `latest` (main branch)
  - `develop` (develop branch)
  - `v1.0.0` (release tags)
  - `v1.0` (major.minor tags)

### Pull and Run

```bash
# Pull latest image
docker pull ghcr.io/armitageee/faceit-cli:latest

# Run with environment file
docker run --rm -it --env-file .env ghcr.io/armitageee/faceit-cli:latest

# Run with specific variables
docker run --rm -it -e FACEIT_API_KEY=your_key ghcr.io/armitageee/faceit-cli:latest
```

## Development

### Local Development with Docker

```bash
# Build and test
make docker-build
make test-docker

# Run with local changes
docker run --rm -it --env-file .env -v $(pwd):/app faceit-cli:latest
```

### Debugging

```bash
# Run with shell access (for debugging)
docker run --rm -it --entrypoint /bin/sh faceit-cli:latest

# Check image layers
docker history faceit-cli:latest

# Inspect image
docker inspect faceit-cli:latest
```

## Troubleshooting

### Common Issues

1. **Permission denied**: Ensure the binary is executable
2. **Environment variables not loaded**: Check `.env` file format
3. **Network issues**: Verify Docker network configuration
4. **Build failures**: Check Go version compatibility

### Debug Commands

```bash
# Check image size
docker images faceit-cli

# View container logs
docker logs <container_id>

# Execute commands in running container
docker exec -it <container_id> /bin/sh
```

## Performance

### Resource Usage

- **CPU**: Minimal (TUI application)
- **Memory**: ~10-20MB base usage
- **Network**: Only for API calls and Kafka (if enabled)

### Optimization Tips

1. Use `distroless` base image for minimal size
2. Enable caching for reduced API calls
3. Use production mode for reduced logging
4. Set appropriate timeouts for API calls

## Security Considerations

1. **API Keys**: Never hardcode in Dockerfile
2. **Secrets**: Use Docker secrets or environment variables
3. **Base Image**: Use distroless for minimal attack surface
4. **Updates**: Regularly update base images and dependencies
5. **Scanning**: Use tools like `docker scan` to check for vulnerabilities
