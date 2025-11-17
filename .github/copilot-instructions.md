# Infrastructure Monorepo - AI Agent Instructions

This is a personal home lab infrastructure monorepo containing IaC, services, and documentation. The project is intentionally over-engineered for learning purposes.

## Architecture Overview

**Three-tier architecture**: Internet → Cloudflare (mTLS proxy + DDoS protection) → Home Lab ↔ GCP (via Wireguard)

- **`cloud/`**: Reusable Terraform modules for GCP infrastructure (not root modules - these are components)
- **`services/`**: Go microservices deployed to GCP Cloud Run
- **`docs/`**: Hugo documentation site (Docsy theme) deployed to GitHub Pages at `https://zaba505.github.io/infra/`

### Cloud Modules Pattern

Each directory in `cloud/` is a **reusable Terraform module** (e.g., `cloud/rest-api/`, `cloud/https-load-balancer/`):
- Contains `main.tf` and `variables.tf` 
- Not root modules - these are imported by actual infrastructure deployments elsewhere
- `cloud/rest-api/`: Provisions GCP Cloud Run services with standard health checks (`/health/startup`, `/health/liveness`)
- `cloud/https-load-balancer/`: GCP HTTPS load balancer with mTLS support using Cloudflare-to-GCP trust anchors
- All modules use GCP provider version `7.11.0`

### Go Services Pattern

Services follow the `z5labs/humus` framework pattern (see `services/lb-sink/`):
- `main.go`: Embeds `config.yaml` and calls `rest.Run()` with app initializer
- `app/app.go`: Implements `Init(ctx, cfg) (*rest.Api, error)` to wire up endpoints and health checks
- `endpoint/`: OpenAPI-first handlers implementing `RequestBody()`, `Responses()`, and `ServeHTTP()`
- Config file: `config.yaml` with minimal `openapi` section (title, version)
- **Current Go version: 1.24.0** - use new `sync.WaitGroup.Go()` method instead of classic `Add`/`Done` pattern

## Key Dependencies

- **Humus framework** (`github.com/z5labs/humus`): HTTP service framework with built-in OpenAPI support
- **OpenTelemetry**: Full observability stack pre-configured (metrics, traces, logs)
- Uses `swaggest/openapi-go` for OpenAPI 3 schema generation

## Development Workflows

### Terraform
- **Linting**: `terraform fmt -recursive -check` (enforced in CI on `**.tf` changes)
- **No apply automation**: Changes trigger CI checks only; manual apply required
- Store secrets in GCP Secret Manager, reference in Terraform via `google_secret_manager_secret_version_access`

### Documentation
- **Local preview**: `cd docs && hugo server` (requires Hugo 0.147.0+ extended, Dart Sass, Node for PostCSS)
- **Build**: `hugo --gc --minify --baseURL "https://zaba505.github.io/infra/"`
- **Deploy**: Automatic to GitHub Pages on push to `main` affecting `docs/**`
- **PR previews**: Enabled via `peaceiris/actions-gh-pages@v4` with preview deployment

### Go Services
- **No Makefile/Dockerfile present**: Services likely built/deployed via external tooling or GCP buildpacks
- Use `go mod tidy` for dependency management (Renovate configured for automatic updates)
- Renovate runs before 4am, updates indirect dependencies, auto-tidies go.mod

## Project Conventions

### Health Checks
All Cloud Run services expose:
- `/health/startup` (startup probe: 30s timeout, 10s period, 3 failures)
- `/health/liveness` (liveness probe: 30s timeout, 10s period, 3 failures)

### mTLS Configuration
Load balancers use GCP Certificate Manager trust configs with Cloudflare CA certificates stored in Secret Manager. The `mtls/cloudflare-gcp/` module provisions this trust relationship.

### Resource Naming
- Terraform resources use `create_before_destroy` lifecycle for zero-downtime updates (SSL certs, NEGs)
- Backend services named after Cloud Run service names
- Network endpoint groups: `{service}-{region}-neg` pattern

### Code Organization
- Services use package-per-component: `app/`, `endpoint/` subdirectories
- Embedded config: Use `//go:embed config.yaml` pattern in `main.go`
- OpenAPI-first: Handlers define schemas via struct methods before implementation

## Common Pitfalls

- **Terraform module vs root confusion**: `cloud/` modules are imported elsewhere; don't expect `terraform.tfstate` here
- **Hugo baseURL**: Always use `https://zaba505.github.io/infra/` for production builds (trailing slash matters)
- **Go 1.24+ concurrency**: Don't use old `wg.Add(1); go func() { defer wg.Done() }()` - use `wg.Go(task)`
- **Health check ports**: Must match `HTTP_PORT` env var in Cloud Run config (default 8080)

## External References

- Full documentation: https://zaba505.dev/infra (redirects to GitHub Pages site)
- Architecture decisions: `docs/content/r&d/adrs/` (uses MADR format)
- Physical infrastructure notes: `docs/content/physical_infrastructure/`
