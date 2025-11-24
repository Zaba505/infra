---
title: "GET /health/liveness"
type: docs
description: "Liveness probe endpoint for Cloud Run"
weight: 31
---

Indicates whether the application is alive and healthy. Used by Cloud Run to detect and restart unhealthy instances.

## Request

**Request Example:**

```http
GET /health/liveness HTTP/1.1
Host: boot.example.com
```

## Response

**Response (200 OK):**

Empty response body with HTTP 200 status code.

**Response (503 Service Unavailable):**

Empty response body with HTTP 503 status code.

**Response Headers:**

- `Cache-Control: no-cache, no-store, must-revalidate`

## Liveness Check Components

1. **HTTP Server Health** - Verifies the HTTP server is responsive
2. **Basic health validation** - Ensures the application can handle requests

## Cloud Run Configuration

```yaml
livenessProbe:
  httpGet:
    path: /health/liveness
    port: 8080
  initialDelaySeconds: 0
  timeoutSeconds: 30
  periodSeconds: 10
  failureThreshold: 3
```

## Behavior

- **Success (200)**: Application is healthy and functioning normally
- **Failure (503)**: Application is unhealthy and should be restarted
- **Consecutive Failures**: After 3 consecutive failures (30 seconds), Cloud Run restarts the instance

## Graceful Degradation

The health check is designed with graceful degradation in mind:

- **Critical Failures**: Return 503 and trigger restart (e.g., database connection lost)
- **Non-Critical Failures**: Log warnings but return 200 (e.g., temporary Cloud Storage timeout)
- **Transient Errors**: Retry internally before reporting failure

## Observability

**Metrics:**

- `health_check_total{probe="liveness",status="ok"}` - Successful liveness checks
- `health_check_total{probe="liveness",status="error"}` - Failed liveness checks
- `health_check_duration_ms{probe="liveness"}` - Liveness check duration

**Structured Logs:**

```json
{
  "severity": "INFO",
  "timestamp": "2025-11-19T06:00:00Z",
  "message": "Health check completed",
  "probe": "liveness",
  "status": "ok",
  "duration_ms": 15
}
```

**Alerts:**

- **Liveness Failure**: Alert if liveness check fails 3+ times consecutively
- **High Restart Rate**: Alert if container restarts > 3 times in 5 minutes

## Testing

### Manual Testing

```bash
curl -v http://localhost:8080/health/liveness
```

### Load Testing

Health check endpoints should handle high request rates without degrading application performance:

- **Target**: 100 requests/second sustained
- **Timeout**: < 10ms average response time
- **Resource Impact**: < 1% CPU, < 10MB memory overhead

## Troubleshooting

### Liveness Check Intermittent Failures

**Symptoms:**
- Occasional container restarts
- Liveness probe returns 503 sporadically
- High request latency

**Debugging:**
```bash
# Check error rate in last 5 minutes
gcloud monitoring time-series list \
  --filter='metric.type="custom.googleapis.com/health_check_total" AND metric.labels.status="error"' \
  --interval-start-time="5 minutes ago"

# Check for resource exhaustion (Cloud Run)
gcloud run services describe boot-server --region=<region> --format=json | jq '.status'
```

**Common Causes:**
- Database connection pool exhausted
- Memory pressure triggering GC pauses
- High request volume overwhelming server
- Dependency timeouts

## Security Considerations

### Unauthenticated Access

Health check endpoints are **intentionally unauthenticated** to allow Cloud Run infrastructure to probe without credentials. This is safe because:

1. Endpoints return only HTTP status codes (no response body)
2. No sensitive data is returned
3. Rate limiting prevents abuse
4. Endpoints are read-only

### Information Disclosure

Health checks return only HTTP status codes with no response body, ensuring:

- No internal IP addresses disclosed
- No error messages or stack traces exposed
- No database connection strings revealed
- No API keys or secrets leaked

Detailed diagnostics are logged internally (not returned in response):

```json
{
  "severity": "ERROR",
  "message": "Firestore connection failed",
  "error": "rpc error: code = PermissionDenied desc = Missing or insufficient permissions"
}
```
