# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Repository Overview

Personal home lab infrastructure monorepo containing IaC, services, and documentation. Intentionally over-engineered for learning purposes.

**Architecture**: Internet → Cloudflare (mTLS proxy + DDoS) → Home Lab ↔ GCP (Wireguard)

## Key Commands

### Documentation (Hugo)
```bash
cd docs && hugo server                    # Local preview at localhost:1313
hugo --gc --minify --baseURL "https://zaba505.github.io/infra/"  # Production build
```
Requires: Hugo 0.147.0+ extended, Dart Sass, Node.js (PostCSS)

### Terraform
```bash
terraform fmt -recursive -check          # Lint (enforced in CI)
```
No automated apply; manual deployment required. Modules in `cloud/` are reusable components, not root modules.

### Go Services
```bash
go mod tidy                               # Dependency management
```
Current version: Go 1.24.0 (module requires 1.25 for `sync.WaitGroup.Go()`)

## Architecture Patterns

### Go Services (Humus Framework)

Services follow `z5labs/humus` pattern in `services/`:

```go
// main.go - Embed config and bootstrap
//go:embed config.yaml
var configBytes []byte

func main() {
    rest.Run(bytes.NewReader(configBytes), app.Init)
}
```

```go
// app/app.go - Wire up API
func Init(ctx context.Context, cfg Config) (*rest.Api, error) {
    return rest.NewApi(
        cfg.OpenApi.Title,
        cfg.OpenApi.Version,
        rest.Liveness(handler),
        rest.Readiness(handler),
        endpoint.Handlers...,
    ), nil
}
```

```go
// endpoint/ - OpenAPI-first handlers
type Handler struct{}
func (h *Handler) RequestBody() openapi3.RequestBodyOrRef { }
func (h *Handler) Responses() openapi3.Responses { }
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) { }
```

Standard health checks: `/health/startup`, `/health/liveness` (30s timeout, 10s period, 3 failures)

### Terraform Modules

Each `cloud/*` directory is a reusable module (not a root module):
- `cloud/rest-api/` - Cloud Run service provisioner
- `cloud/https-load-balancer/` - HTTPS LB with mTLS support
- `cloud/mtls/cloudflare-gcp/` - Cloudflare-to-GCP trust anchors

All use GCP provider v7.11.0 with `create_before_destroy` lifecycle for zero-downtime.

Resource naming:
- Network endpoint groups: `{service}-{region}-neg`
- Backend services: Named after Cloud Run service names

## Critical Go Patterns

**Package declarations**: Each `.go` file must have exactly ONE `package` line. Never duplicate. Check existing files in the directory before creating new ones.

**Go 1.24+ WaitGroup**: Use new method, not classic pattern:
```go
// Correct (Go 1.25+)
var wg sync.WaitGroup
wg.Go(task)
wg.Wait()

// Avoid
wg.Add(1)
go func() { defer wg.Done(); task() }()
```

**Error handling**: Return early, wrap with context using `%w`, keep messages lowercase without punctuation.

**Interfaces**: Accept interfaces, return concrete types. Keep small (1-3 methods). Define close to usage.

## CI/CD Workflows

- **terraform.yml** - Lints `**.tf` on PR/push to main
- **docs.yaml** - Auto-deploys Hugo to GitHub Pages on `docs/**` changes
- **codeql.yaml** - Security analysis on `.go` changes
- **Renovate** - Runs before 4am, auto-tidies go.mod

## Important Files

- `.github/copilot-instructions.md` - Comprehensive AI agent instructions
- `.github/instructions/go.instructions.md` - Detailed Go coding standards
- `docs/content/r&d/adrs/` - Architecture Decision Records (MADR 4.0.0 format)

## Common Pitfalls

- **Terraform module confusion**: `cloud/` contains reusable modules imported elsewhere; no `terraform.tfstate` here
- **Hugo baseURL**: Always include trailing slash for production builds
- **Health check ports**: Must match `HTTP_PORT` env var (default 8080)
- **Secrets**: Store in GCP Secret Manager, reference via `google_secret_manager_secret_version_access`
