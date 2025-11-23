---
title: "Machine Service"
type: docs
description: "Service for managing machine hardware profiles"
weight: 20
---

The Machine Service is a REST API that manages machine hardware profiles for the network boot infrastructure. It stores machine specifications (CPUs, memory, NICs, drives, accelerators) in Firestore and is queried by the Boot Service during boot operations and by administrators for configuration management.

## Architecture

The service is responsible for:

- **Machine Profile Management**: Creating, listing, retrieving, updating, and deleting machine hardware profiles
- **Hardware Specification Storage**: Storing detailed hardware specifications in Firestore
- **Machine Lookup**: Providing machine profile queries by ID or NIC MAC address

## Components

- **Firestore**: Stores machine hardware profiles
- **REST API**: HTTP endpoints for machine profile management

## Clients

The service is consumed by:

1. **Boot Service**: Queries machine profiles by MAC address during boot operations
2. **Admin Tools**: CLI or web interfaces for managing machine inventory
3. **Monitoring Systems**: Hardware inventory and asset management tools

## Deployment

- **Platform**: GCP Cloud Run
- **Scaling**: Automatic scaling based on request load
- **Availability**: Min instances = 1 for low-latency responses
- **Region**: Same region as Boot Service for minimal latency
