---
title: "Services"
type: docs
description: "Documentation for GCP Cloud Run microservices"
weight: 30
---

This section contains documentation for the Go microservices deployed to GCP Cloud Run as part of the home lab infrastructure.

## Service Architecture

All services follow a consistent architecture pattern:

- **Framework**: Built using `z5labs/humus` framework with OpenAPI-first design
- **Runtime**: Go 1.24+ deployed to GCP Cloud Run
- **Observability**: OpenTelemetry metrics, traces, and logs
- **Health Checks**: Standard `/health/startup` and `/health/liveness` endpoints
- **Configuration**: Embedded `config.yaml` with OpenAPI specifications
