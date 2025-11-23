---
title: "Machine Management Service"
type: docs
description: "Service for managing machine profiles and boot profiles"
weight: 20
---

The Machine Management Service is a REST API that manages machine hardware profiles and boot profiles for the network boot infrastructure. It stores machine hardware specifications and boot profiles in Firestore, manages kernel and initrd blobs in Cloud Storage, and serves both the Boot Service (during boot operations) and administrators (for configuration management).

## Architecture

The service is responsible for:

- **Machine Profile Management**: Creating, listing, retrieving, updating, and deleting machine hardware profiles
- **Boot Profile Management**: Creating, listing, retrieving, updating, and deleting boot profiles
- **Blob Storage**: Uploading and streaming kernel/initrd binaries to/from Cloud Storage
- **Machine-to-Boot Profile Mapping**: Associating machines with boot profiles via machine ID

## Components

- **Firestore**: Stores machine hardware profiles and boot profile metadata
- **Cloud Storage**: Stores kernel and initrd blobs using UUIDv7 identifiers
- **REST API**: HTTP endpoints for machine and boot profile management

## Clients

The service is consumed by:

1. **Boot Service**: Retrieves machine profiles and boot profiles during boot operations
2. **Admin Tools**: CLI or web interfaces for managing machines and profiles
3. **CI/CD Pipelines**: Automated profile deployment workflows

## Deployment

- **Platform**: GCP Cloud Run
- **Scaling**: Automatic scaling based on request load
- **Availability**: Min instances = 1 for low-latency responses
- **Region**: Same region as Boot Service and Cloud Storage for minimal latency
