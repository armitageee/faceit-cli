# Logging and Kafka Integration

This document describes the logging system and Kafka integration in faceit-cli.

## Overview

The application uses a centralized logging system with the following features:

- **Structured JSON logging** with timestamps and service identification
- **Configurable log levels**: debug, info, warn, error
- **Kafka integration** for log aggregation (optional)
- **Fallback logging** to stdout when Kafka is unavailable

## Configuration

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `LOG_LEVEL` | `info` | Log level: debug, info, warn, error |
| `KAFKA_ENABLED` | `false` | Enable Kafka logging |
| `KAFKA_BROKERS` | `localhost:9092` | Comma-separated list of Kafka brokers |
| `KAFKA_TOPIC` | `faceit-cli-logs` | Kafka topic for log messages |
| `PRODUCTION_MODE` | `false` | Enable production mode |
| `LOG_TO_STDOUT` | `true` | Log to stdout (set to false in production) |

### Example Configuration

```bash
# Development mode
LOG_LEVEL=debug
LOG_TO_STDOUT=true
KAFKA_ENABLED=false

# Production mode with Kafka
PRODUCTION_MODE=true
LOG_TO_STDOUT=false
KAFKA_ENABLED=true
KAFKA_BROKERS=localhost:9092,localhost:9093
KAFKA_TOPIC=faceit-cli-logs
LOG_LEVEL=info
```

## Docker Compose Setup

### Starting Kafka Infrastructure

```bash
# Start Kafka KRaft and Kafka UI (no Zookeeper needed!)
make kafka-up

# Or manually
docker-compose up -d
```

### Accessing Kafka UI

Once started, Kafka UI is available at: http://localhost:8080

### Stopping Infrastructure

```bash
# Stop all services
make kafka-down

# Or manually
docker-compose down
```

## Log Structure

### Console Logs

Logs are written to stdout in JSON format:

```json
{
  "timestamp": "2024-01-15T10:30:45Z",
  "level": "info",
  "message": "Player profile loaded successfully",
  "service": "faceit-cli",
  "nickname": "player123",
  "player_id": "abc-123-def"
}
```

### Kafka Logs

When Kafka is enabled, logs are also sent to the configured Kafka topic with the same structure.

## Usage Examples

### Running with Debug Logging

```bash
LOG_LEVEL=debug go run main.go
```

### Running with Kafka Integration

```bash
# Start Kafka first
make kafka-up

# Run with Kafka logging
make run-kafka

# Or manually
KAFKA_ENABLED=true LOG_LEVEL=debug go run main.go
```

### Monitoring Logs

```bash
# View Kafka infrastructure logs
make kafka-logs

# View application logs with Kafka
KAFKA_ENABLED=true LOG_LEVEL=debug go run main.go 2>&1 | jq .
```

## Log Levels

- **debug**: Detailed information for debugging
- **info**: General information about application flow
- **warn**: Warning messages for potential issues
- **error**: Error messages for failures

## Kafka Topic Management

### Automatic Topic Creation

The application has **three ways** to create the Kafka topic:

1. **Automatic creation** - Topic is created automatically when first message is sent (enabled by `KAFKA_AUTO_CREATE_TOPICS_ENABLE=true`)
2. **Docker Compose init** - Topic is created during infrastructure startup by the `kafka-init` service (KRaft mode)
3. **Manual creation** - You can create the topic manually using Makefile commands

### Topic Creation Methods

```bash
# Method 1: Automatic (default) - Topic created on first log message
KAFKA_ENABLED=true go run main.go

# Method 2: Docker Compose init - Topic created during startup
make kafka-up  # Creates topic automatically

# Method 3: Manual creation
make kafka-create-topic

# Method 4: Using Kafka UI (recommended for management)
# Go to http://localhost:8080 and create topic "faceit-cli-logs"

# Method 5: Direct kafka-topics command
docker exec -it faceit-cli-kafka kafka-topics --create \
  --topic faceit-cli-logs \
  --bootstrap-server localhost:9092 \
  --partitions 3 \
  --replication-factor 1
```

### Topic Management Commands

```bash
# List all topics
make kafka-topics

# Create topic manually
make kafka-create-topic

# View topic details in Kafka UI
make kafka-ui  # Opens http://localhost:8080
```

## Troubleshooting

### Kafka Connection Issues

1. Ensure Kafka is running: `docker-compose ps`
2. Check Kafka logs: `make kafka-logs`
3. Verify broker configuration: `KAFKA_BROKERS=localhost:9092`

### Log Level Issues

- Ensure log level is one of: debug, info, warn, error
- Case insensitive: `LOG_LEVEL=DEBUG` works the same as `LOG_LEVEL=debug`

### Performance Considerations

- Kafka logging is asynchronous and won't block the application
- If Kafka is unavailable, logs are still written to stdout
- Consider log volume when using debug level in production

## Integration with Monitoring Systems

The structured JSON logs can be easily integrated with:

- **ELK Stack** (Elasticsearch, Logstash, Kibana)
- **Grafana Loki**
- **Fluentd**
- **Prometheus** (with log-based metrics)

Example log parsing for monitoring:

```bash
# Count errors by service
jq 'select(.level == "error") | .service' logs.json | sort | uniq -c

# Monitor API response times
jq 'select(.message | contains("API call")) | .duration' logs.json
```

## Production Mode

Production mode is designed for deployment environments where you don't want logs cluttering stdout:

### Features
- **No stdout output** - Logs are not printed to console
- **Kafka-only logging** - All logs go to Kafka for centralized collection
- **Optimized performance** - Reduced I/O overhead
- **Clean console** - Application output is not mixed with logs

### Usage
```bash
# Using Makefile commands
make run-production      # Production mode without Kafka
make run-prod-kafka      # Production mode with Kafka

# Manual configuration
PRODUCTION_MODE=true LOG_TO_STDOUT=false KAFKA_ENABLED=true go run main.go
```

### Production Deployment
For production deployments, consider:
- Set `LOG_LEVEL=info` or `LOG_LEVEL=warn` to reduce log volume
- Use `LOG_TO_STDOUT=false` to keep console clean
- Enable `KAFKA_ENABLED=true` for centralized logging
- Configure proper Kafka brokers for your infrastructure
