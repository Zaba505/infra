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
go test ./...                             # Run all tests
go test ./services/machine/...            # Test specific service
go mod tidy                               # Dependency management
```
**Go version**: 1.26.0 (use `sync.WaitGroup.Go()` — the classic `Add`/`Done` pattern is obsolete)

**Module**: `github.com/Zaba505/infra`

## Architecture Patterns

### Go Services

Two services live under `services/`:
- `machine/` — full chi-based HTTP service backed by Firestore; the reference shape for new services
- `lb-sink/` — placeholder stub (a `main.go` that prints hello); exists to anchor a load-balancer backend

Services do **not** use the humus framework. The current pattern (see `services/machine/`) is:

```
services/{service-name}/
├── main.go          # Calls os.Exit(app.Main(context.Background()))
├── app/
│   └── app.go      # ConfigFromEnv(), Main(ctx) int — wires chi router + HTTP server
├── endpoint/        # HTTP handlers; each registers routes on *chi.Mux
│   └── endpointpb/  # Protobuf-generated types for request/response
└── service/         # Backend service clients (e.g. Firestore)
```

**Key dependencies** (in `go.mod`):
- `github.com/go-chi/chi/v5` — HTTP routing
- `github.com/z5labs/bedrock` — config from environment (`config.Env`, `config.Default`, `config.Must`)
- `github.com/sourcegraph/conc/pool` — structured concurrency
- `google.golang.org/protobuf` — request/response serialization
- `go.opentelemetry.io/otel` — tracing

**Config** comes from environment variables via `ConfigFromEnv(ctx)`, not embedded YAML. (Note: `services/machine/config.yaml` is a vestigial humus-era artifact — the live config path is env vars.)
```go
func ConfigFromEnv(ctx context.Context) Config {
    return Config{
        HTTP: HTTPConfig{
            Port: config.Must(ctx, config.Default(8080, config.IntFromString(config.Env("HTTP_PORT")))),
        },
        Firestore: FirestoreConfig{
            ProjectID: config.Must(ctx, config.Env("GCP_PROJECT_ID")),
        },
    }
}
```

**HTTP server** starts with self-signed TLS, listens on `HTTP_PORT` (default 8080):
```go
func Main(ctx context.Context) int {
    sigCtx, cancel := signal.NotifyContext(ctx)
    defer cancel()
    // ... init deps, build chi mux, start TLS server via conc pool
}
```

**Endpoint registration** — each handler file exposes a registration function:
```go
func RegisterMachines(mux *chi.Mux, firestoreClient FirestoreClient) {
    handler := &registerMachinesHandler{...}
    mux.Method(http.MethodPost, "/api/v1/machines", handler)
}
```

**Request/response**: protobuf over HTTP. Content-Type is `application/x-protobuf`; errors use `application/problem+protobuf`.

**Error handling**: use `pkg/errorpb` helpers — `NewInternalError`, `NewValidationError`, `NewConflictError` — then call `errorHandler(ctx, w, err)` which dispatches to `WriteHttpResponse`.

**Service layer interfaces**: define a small interface in `endpoint/` (not in `service/`) for the backend client, then pass the concrete `service.*Client` at registration time. This keeps the endpoint package testable.

### pkg/errorpb

Shared error types generated from protobuf. Problems serialize to protobuf and are written as `application/problem+protobuf` responses. Three concrete types: `*Problem` (generic/internal), `*ValidationProblem` (400), `*ConflictProblem` (409).

### Terraform Modules

Each `cloud/*` directory is a reusable Terraform module (not a root module). The set spans ~14 GCP building blocks covering compute (`compute-engine`, `rest-api`), networking (`vpc-network`, `https-load-balancer`, `internal-application-load-balancer`, `network-load-balancer`, `ip`, `dns`), storage/data (`storage-bucket`, `firestore`, `artifact-registry`), IAM (`service-account`), and mTLS trust (`mtls/cloudflare-gcp`).

All use GCP provider v7.11.0 with `create_before_destroy` lifecycle for zero-downtime.

Resource naming:
- Network endpoint groups: `{service}-{region}-neg`
- Backend services: named after Cloud Run service names

## Critical Go Patterns

**Package declarations**: each `.go` file must have exactly ONE `package` line. Check existing files before creating new ones.

**WaitGroup**: Go 1.26 — use `wg.Go(task)`, not `wg.Add(1); go func() { defer wg.Done() }()`.

**Structured concurrency**: use `pool.New().WithErrors().WithContext(ctx)` from `github.com/sourcegraph/conc/pool`.

**Error handling**: return early, wrap with `%w`, messages lowercase without trailing punctuation.

**Interfaces**: define in the package that uses them (not where implemented), keep to 1–3 methods.

## CI/CD Workflows

- **terraform.yml** — Lints `**.tf` on PR/push to main
- **docs.yaml** — Auto-deploys Hugo to GitHub Pages on `docs/**` changes to main
- **docs-preview.yaml** — PR preview sites for `docs/**` changes
- **codeql.yaml** — Security analysis on `.go` changes
- **Renovate** — Runs before 4am, auto-tidies go.mod, updates indirect dependencies

## Commit & Branch Conventions

- Branch: `story/issue-{number}/{description}` or `fix/issue-{number}/{description}`
- Commit: `feat(issue-123): description` or `fix(issue-123): description`

## Important Files

- `.github/copilot-instructions.md` — Comprehensive AI agent instructions
- `.github/instructions/go.instructions.md` — Detailed Go coding standards
- `.github/agents/` — Agent definitions (e.g. `review-capability.agent.md`)
- `.github/prompts/` — Reusable prompt templates (e.g. `new-adr.prompt.md`)
- `docs/content/r&d/adrs/` — Architecture Decision Records (MADR 4.0.0 format)

## Common Pitfalls

- **Terraform module confusion**: `cloud/` contains reusable modules; no `terraform.tfstate` here
- **Hugo baseURL**: always include trailing slash for production builds
- **Health check ports**: must match `HTTP_PORT` env var (default 8080)
- **Secrets**: store in GCP Secret Manager, reference via `google_secret_manager_secret_version_access`
- **humus framework**: not used in this repo — ignore any docs or examples referencing `rest.Run`, `rpc.Producer`, `rpc.Handler`, or embedded `config.yaml`
