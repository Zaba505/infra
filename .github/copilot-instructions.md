# Infrastructure Monorepo - AI Agent Instructions

## Overview

This is an infrastructure-as-code (IaC) monorepo for a personal home lab combining Terraform modules, Go services, and Hugo documentation. The project intentionally over-engineers a home lab setup to explore cloud-native patterns, network booting, and service deployment.

**Key Areas:**
- **Terraform Modules** (`*/main.tf`, `*/variables.tf`): Reusable IaC for GCP and Cloudflare
- **Go Services** (`services/`, `pkg/`): REST APIs and network services built on z5labs/bedrock framework
- **Documentation** (`docs/`): Hugo-based static site with ADRs and infrastructure guides
- **Network Boot** (`rack/`): Fedora CoreOS network booting with iPXE and Ignition configs

## Architecture & Patterns

### Service Architecture (Go)

All Go services follow a **z5labs/bedrock application framework** pattern:

1. **Configuration via embedded YAML**: Services embed a `config.yaml` using `//go:embed`, merged with default config from shared packages (`pkg/rest/default_config.yaml`). Configuration uses Go template syntax with `env` and `default` functions.

2. **Initialization pattern**: Each service has an `Init` function (e.g., `services/machinemgmt/app.Init`) that:
   - Takes `context.Context` and a config struct
   - Returns `[]rest.Endpoint` (for REST services)
   - Initializes dependencies (GCS clients, storage services, etc.)

3. **Service composition**:
   ```
   main.go (embeds config.yaml)
     ↓
   pkg/rest.Run() or pkg/ftp.Run() 
     ↓
   app/service.Init() → []Endpoint
     ↓
   bedrock.Run() → HTTP server
   ```

4. **Shared runtime packages**:
   - `pkg/rest`: HTTP service wrapper with OpenTelemetry, health checks, and bedrock integration
   - `pkg/ftp`: FTP service wrapper (similar pattern to REST)
   - Both use **template-based config** with environment variable substitution

### REST Endpoint Pattern

Endpoints use a functional options pattern from `pkg/rest`:
```go
rest.Get("/path/{param}", handler,
    rest.PathParam("param", "", false),
    rest.StatusCode(http.StatusOK),
)
```

Handlers implement `endpoint.Handler[I, O]` where `I` and `O` are request/response types that satisfy `endpoint.Request` and `endpoint.Response` interfaces.

### Terraform Module Structure

Each infrastructure component is a **reusable Terraform module** in its own directory:
- `artifact-registry/docker/`: GCP Docker registry
- `dns/cloudflare/`: DNS records with optional mTLS
- `https-load-balancer/`: GCP load balancer with Cloud Run backends
- `rest-api/`: Cloud Run service deployment
- `service-account/`: GCP service accounts with IAM

**Conventions:**
- `main.tf`: Resource definitions
- `variables.tf`: Input variables with descriptions and defaults
- Terraform 7.11.0 for GCP provider, consistent across modules
- Use `locals` for computed values and data transformations

### Network Boot Flow

The `rack/` directory contains network booting for HP ProLiant servers:
1. **Butane config** (`bootstrap.bu`) → **Ignition config** (`bootstrap.ign`)
2. **iPXE template** (`bootstrap.ipxe.gotmpl`) references Fedora CoreOS live images
3. Servers PXE boot → fetch iPXE script → download CoreOS kernel/initramfs → apply Ignition config

This enables stateless server provisioning through a dedicated Network Boot VLAN.

## Development Workflows

### Working with Go Services

**Before editing:**
1. Check existing package structure in `services/` or `pkg/`
2. Read the service's `config.yaml` to understand configuration
3. Examine `pkg/rest/default_config.yaml` for shared config defaults

**Adding a new REST endpoint:**
1. Create handler in `services/<name>/` subdirectory (e.g., `bootstrap/bootstrap.go`)
2. Implement `endpoint.Handler[I, O]` interface
3. Return endpoint from package function using `rest.Get/Post/etc`
4. Wire up in `app.Init()` by appending to `[]rest.Endpoint` slice

**Configuration pattern:**
- Services override defaults by embedding their own `config.yaml`
- Use `{{env "VAR_NAME" | default "value"}}` for environment-driven config
- Struct tags: `` config:"field_name" `` for nested config mapping

**Testing:**
- Test files live alongside source (`*_test.go`)
- Use table-driven tests with `t.Run()` for subtests
- Minimal test coverage currently exists (see `services/lb-sink/service/service_test.go`)

### Terraform Development

**Module usage pattern:**
Modules are designed to be consumed by root Terraform configurations (not present in this repo). Each module is self-contained and follows a consistent structure.

**Testing changes:**
```bash
# Format check
terraform fmt -recursive -check

# Validate a specific module
cd <module-directory>
terraform init
terraform validate
```

**Key considerations:**
- GCP resources use lifecycle `create_before_destroy` for zero-downtime updates (e.g., SSL certificates, NEGs)
- Cloudflare authenticated origin pulls require DNS records to exist first (use `depends_on`)
- Use `for_each` with maps for resource collections, not `count`

### Documentation (Hugo)

Documentation lives in `docs/` using the **Docsy theme**:
- `docs/content/`: Markdown content with Hugo front matter
- `docs/content/r&d/adrs/`: Architectural Decision Records (MADR 4.0.0 format)
- ADRs use special front matter: `category`, `status`, `date`, `weight`

**ADR categories:**
- `strategic`: Framework/architecture decisions
- `user-journey`: Feature implementation approaches
- `api-design`: API design trade-offs

**Building docs:**
```bash
cd docs
hugo server  # Development server
hugo         # Production build to public/
```

### Bazel Build System (Legacy/Transitioning)

**Note:** The codebase currently has Bazel build files (`BUILD.bazel`) and symlinks (`bazel-*`), but the branch name suggests a de-Bazeling effort is underway. 

Current Bazel artifacts:
- `pkg/rest/BUILD.bazel`, `pkg/ftp/BUILD.bazel`: Go library definitions
- Embedded resources via `embedsrcs`
- OCI image builds referenced in CI

**When Bazel is removed:**
- Standard `go build` and `go test` workflows will apply
- Docker builds will likely use Dockerfiles or alternative tooling

## Project-Specific Conventions

### Go Code Style

Follow `.github/instructions/go.instructions.md` (comprehensive Go guidelines). Key highlights:
- **Package declarations**: Never duplicate `package` lines; check existing files in directory before creating new ones
- **Error handling**: Wrap errors with context using `fmt.Errorf` with `%w`
- **Interfaces**: Define close to usage, keep small (1-3 methods)
- **Concurrency**: Document goroutine lifecycle; use `sync.WaitGroup.Go()` if go.mod specifies Go 1.25+

### Configuration Management

**Template-based YAML configs:**
- Default config in `pkg/rest/default_config.yaml` and `pkg/ftp/default_config.yaml`
- Service-specific config overrides in `services/<name>/config.yaml`
- Template functions: `env` (lookup env var), `default` (provide fallback)
- Merged at runtime via `config.FromYaml()` with multiple sources

**OpenTelemetry (OTel) configuration:**
All services support OTel tracing, metrics, and logging via shared config structure:
```yaml
otel:
  service_name: {{env "OTEL_SERVICE_NAME"}}
  trace:
    enabled: {{env "OTEL_TRACE_ENABLED" | default false}}
  metric:
    enabled: {{env "OTEL_METRIC_ENABLED" | default false}}
```
Exporters: GCP for trace/metrics, stdout for logs.

### Dependency Management

- **Go modules**: Use `go.mod` at repo root; `docs/go.mod` is separate for Hugo
- **Renovate**: Configured to auto-update dependencies (`.renovate.json`), including Go indirect deps
- **Key dependencies**:
  - `z5labs/bedrock`: Application framework for services
  - `swaggest/openapi-go`: OpenAPI schema generation
  - Google Cloud SDKs: storage, monitoring, trace
  - OpenTelemetry suite for observability

## Integration Points

### GCP Service Deployment

REST services deploy to Cloud Run via `rest-api` Terraform module:
- **Docker images**: Published to GCP Artifact Registry (`artifact-registry/docker`)
- **Service accounts**: Managed by `service-account` module with IAM bindings
- **Environment variables**: Set via `env` map in module, merged with defaults
- **Load balancing**: `https-load-balancer` module routes traffic to Cloud Run NEGs

### Cloudflare Integration

- **DNS**: Managed via `dns/cloudflare` module with proxied records
- **mTLS**: `mtls/cloudflare-gcp` provisions origin certificates and stores in GCP Secret Manager
- Secrets referenced by load balancer for TLS policy

### Storage (GCS)

Services like `machinemgmt` interact with GCS via `cloud.google.com/go/storage`:
- **Pattern**: Pass `*storage.BucketHandle` to service constructors (see `backend.NewStorageService`)
- **Checksums**: Validate CRC32C checksums after reads
- **Buffering**: Use `sync.Pool` for `bytes.Buffer` reuse to minimize allocations

## Common Pitfalls

1. **Reusing `io.Reader` in HTTP requests**: Readers are consumed once. For retries/redirects, buffer the payload and set `req.GetBody` (see `.github/instructions/go.instructions.md` I/O section).

2. **Terraform lifecycle rules**: Forgetting `create_before_destroy` on SSL certificates causes downtime when rotating certs.

3. **Package declarations**: Do not add duplicate `package` statements when editing Go files. Each file has exactly one.

4. **Config merging order**: Default config is read first, then service-specific. Later sources override earlier ones.

5. **Network Boot VLAN isolation**: The Network Boot VLAN is separate from the Homelab VLAN; servers use multiple interfaces (see `docs/content/physical_infrastructure/_index.md`).

## Useful Commands

```bash
# Go
go test ./...                    # Run all tests
go mod tidy                      # Clean dependencies

# Terraform
terraform fmt -recursive -check  # Format check
cd <module> && terraform init    # Initialize module

# Hugo docs
cd docs && hugo server           # Dev server
cd docs && npm install -D postcss postcss-cli autoprefixer  # Install docsy deps

# Bazel (legacy)
bazel build //...                # Build all targets
bazel test //...                 # Run all tests
```

## Additional Context

- **Documentation site**: Hosted at https://zaba505.dev/infra (GitHub Pages)
- **CI/CD**: GitHub Actions for tests, Terraform linting, Hugo builds, Bazel builds/publishes
- **Monorepo rationale**: Single source of truth for IaC, services, and docs simplifies coordination across stack layers
- **Home lab context**: Physical servers (HP ProLiant DL360 Gen 9) network boot via WireGuard tunnel to Kubernetes-hosted Matchbox service in public cloud

---

**When in doubt:** Reference existing implementations in `services/machinemgmt` (most complete service) or `pkg/rest` (canonical REST service framework). For Terraform, `https-load-balancer` demonstrates complex resource dependencies and conditional logic.
