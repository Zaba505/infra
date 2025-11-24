---
title: "GET /health/startup"
type: docs
description: "Startup probe endpoint for Cloud Run"
weight: 30
---

Indicates whether the application has completed initialization and is ready to receive traffic.

## Request

**Request Example:**

```http
GET /health/startup HTTP/1.1
Host: boot.example.com
```

## Response

**Response (200 OK):**

Empty response body with HTTP 200 status code.

**Response (503 Service Unavailable):**

Empty response body with HTTP 503 status code.

**Response Headers:**

- `Cache-Control: no-cache, no-store, must-revalidate`

## Startup Check Components

1. **Firestore Connection** - Verifies database connectivity
2. **Cloud Storage Access** - Validates access to boot image buckets

## Cloud Run Configuration

```yaml
startupProbe:
  httpGet:
    path: /health/startup
    port: 8080
  initialDelaySeconds: 0
  timeoutSeconds: 30
  periodSeconds: 10
  failureThreshold: 3
```

## Behavior

- **Success (200)**: Application is fully initialized and ready to serve requests
- **Failure (503)**: Application is still starting up or encountered initialization errors
- **Timeout**: After 30 seconds of no response, Cloud Run considers startup failed

## Observability

**Metrics:**

- `health_check_total{probe="startup",status="ok"}` - Successful startup checks
- `health_check_total{probe="startup",status="error"}` - Failed startup checks
- `health_check_duration_ms{probe="startup"}` - Startup check duration

**Structured Logs:**

```json
{
  "severity": "INFO",
  "timestamp": "2025-11-19T06:00:00Z",
  "message": "Health check completed",
  "probe": "startup",
  "status": "ok",
  "duration_ms": 15
}
```

**Alerts:**

- **Startup Failure**: Alert if startup check fails for > 1 minute

## Testing

### Manual Testing

```bash
curl -v http://localhost:8080/health/startup
```

### Automated Testing

```go
func TestHealthStartup(t *testing.T) {
    resp, err := http.Get("http://localhost:8080/health/startup")
    require.NoError(t, err)
    defer resp.Body.Close()

    assert.Equal(t, http.StatusOK, resp.StatusCode)
}
```

## Troubleshooting

### Startup Check Never Succeeds

**Symptoms:**
- Container restarts repeatedly
- Cloud Run shows "unhealthy" status
- Startup probe returns 503

**Debugging:**
```bash
# Check Cloud Run logs for startup errors
gcloud logging read "resource.type=cloud_run_revision AND labels.service_name=boot-server" \
  --limit 50 --format json | jq '.[] | select(.jsonPayload.probe=="startup")'

# Test locally with debug logging
DEBUG=true go run main.go
```

**Common Causes:**
- Firestore credentials not configured
- Cloud Storage bucket permissions missing
- Network connectivity issues
- Timeout too short for slow dependencies
