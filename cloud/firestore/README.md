# Firestore Database Module

Terraform module for provisioning Google Cloud Firestore databases with IAM bindings for service account access.

## Features

- Supports both Firestore Native mode and Datastore mode
- Configurable location/region for database
- Configurable deletion policy (ABANDON or DELETE)
- IAM bindings for service account access:
  - `roles/datastore.user` for read/write operations
  - `roles/datastore.viewer` for read-only access
- Point-in-time recovery (PITR) configuration
- App Engine integration mode configuration
- Uses GCP provider version 7.11.0
- Applies `lifecycle { create_before_destroy = true }` for zero-downtime updates

## Usage

### Basic Usage (Datastore Mode)

```hcl
module "firestore" {
  source = "../../cloud/firestore"

  name     = "boot-server-db"
  location = "us-central1"

  database_type   = "DATASTORE_MODE"
  deletion_policy = "DELETE"

  service_account_users = [
    "boot-server@my-project.iam.gserviceaccount.com"
  ]
}
```

### Network Boot Server Integration Pattern

This module is designed to store machine-to-image mappings for the network boot server (see [ADR-0005](../../docs/content/r&d/adrs/0005-network-boot-infrastructure-gcp.md)).

```hcl
# Service account for the boot server
module "boot_server_sa" {
  source = "../../cloud/service-account"

  name = "boot-server"
}

# Firestore database for machine-to-image mappings
module "boot_mappings_db" {
  source = "../../cloud/firestore"

  name     = "boot-mappings"
  location = "us-central1"

  database_type                   = "DATASTORE_MODE"
  deletion_policy                 = "DELETE"
  point_in_time_recovery_enabled  = true

  service_account_users = [
    module.boot_server_sa.email
  ]
}

# Boot server Cloud Run service
module "boot_server" {
  source = "../../cloud/rest-api"

  name                  = "boot-server"
  description           = "UEFI HTTP boot server"
  location              = "us-central1"
  service_account_email = module.boot_server_sa.email

  image = {
    artifact_registry_name = "docker-registry"
    name                   = "boot-server"
    tag                    = "latest"
  }

  env = {
    FIRESTORE_DATABASE = module.boot_mappings_db.database_name
  }

  max_instance_count = 2
  max_concurrent_requests = 10
}
```

### Firestore Native Mode

```hcl
module "firestore_native" {
  source = "../../cloud/firestore"

  name     = "app-database"
  location = "us-central1"

  database_type                  = "FIRESTORE_NATIVE"
  concurrency_mode               = "OPTIMISTIC"
  deletion_policy                = "ABANDON"
  point_in_time_recovery_enabled = true

  service_account_users = [
    "app-server@my-project.iam.gserviceaccount.com"
  ]

  service_account_viewers = [
    "read-only-service@my-project.iam.gserviceaccount.com"
  ]
}
```

### Default Database with App Engine Integration

```hcl
module "default_database" {
  source = "../../cloud/firestore"

  name     = "(default)"
  location = "us-central"

  database_type               = "DATASTORE_MODE"
  app_engine_integration_mode = "ENABLED"
  deletion_policy             = "ABANDON"

  service_account_users = [
    "app-engine@my-project.iam.gserviceaccount.com"
  ]
}
```

## Variables

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|----------|
| name | The name of the Firestore database. Use '(default)' for the default database. | `string` | n/a | yes |
| location | The location of the database. See [available locations](https://cloud.google.com/firestore/docs/locations) | `string` | n/a | yes |
| database_type | The type of the database. Possible values: FIRESTORE_NATIVE, DATASTORE_MODE | `string` | `"DATASTORE_MODE"` | no |
| concurrency_mode | The concurrency control mode. Possible values: OPTIMISTIC, PESSIMISTIC | `string` | `"OPTIMISTIC"` | no |
| app_engine_integration_mode | App Engine integration mode. Possible values: ENABLED, DISABLED | `string` | `"DISABLED"` | no |
| deletion_policy | Deletion behavior. Possible values: ABANDON, DELETE | `string` | `"DELETE"` | no |
| point_in_time_recovery_enabled | Whether to enable point-in-time recovery | `bool` | `false` | no |
| service_account_users | Service account emails to grant roles/datastore.user (read/write) | `list(string)` | `[]` | no |
| service_account_viewers | Service account emails to grant roles/datastore.viewer (read-only) | `list(string)` | `[]` | no |

## Outputs

| Name | Description |
|------|-------------|
| database_name | The name of the Firestore database |
| database_id | The ID of the Firestore database |
| database_location | The location of the Firestore database |
| database_type | The type of the Firestore database |
| database_uid | The unique identifier of the Firestore database |

## IAM Bindings

The module creates IAM bindings with conditions to scope access to the specific Firestore database:

- **roles/datastore.user**: Full read/write access to the database (for `service_account_users`)
- **roles/datastore.viewer**: Read-only access to the database (for `service_account_viewers`)

All IAM bindings include a condition to restrict access to only the specified database using the resource name pattern.

## Database Types

### DATASTORE_MODE

Firestore in Datastore mode is optimized for server-side applications and provides:
- Strong consistency for entity group queries
- Compatibility with Datastore APIs
- Suitable for key-value and document storage patterns
- **Recommended for boot server machine-to-image mappings**

### FIRESTORE_NATIVE

Firestore in Native mode provides:
- Real-time synchronization
- Offline support
- Client libraries for web and mobile
- Strong consistency for all queries
- Better for mobile/web applications

## Deletion Policy

- **DELETE**: Database will be deleted when the Terraform resource is destroyed (default)
- **ABANDON**: Database will be retained when the Terraform resource is destroyed

**Warning**: Deleting a Firestore database is irreversible and all data will be lost.

## Point-in-Time Recovery (PITR)

When enabled, PITR allows you to restore the database to any point in time within the retention period (7 days by default). This is useful for:
- Recovering from accidental data deletion
- Rolling back to a known good state
- Compliance and auditing requirements

**Note**: PITR may incur additional storage costs.

## Related Documentation

- [ADR-0005: Network Boot Infrastructure Implementation on Google Cloud](../../docs/content/r&d/adrs/0005-network-boot-infrastructure-gcp.md)
- [Firestore Locations](https://cloud.google.com/firestore/docs/locations)
- [Firestore IAM Roles](https://cloud.google.com/firestore/docs/security/iam)
- [Firestore in Datastore mode](https://cloud.google.com/datastore/docs/firestore-or-datastore)
