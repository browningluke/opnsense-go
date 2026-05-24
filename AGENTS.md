# opnsense-go

A Go client library for the OPNsense API. Requires Go 1.23+.

## Commands

```sh
# Regenerate all code from schema files (run after any schema change)
make all

# Run all acceptance tests (requires a live OPNsense instance)
make testacc

# Run tests for a specific service package
make testacc PKG=unbound

# Run a single test
make testacc PKG=unbound TEST=TestUnboundDomainOverride

# Acceptance test environment variables (required)
export OPNSENSE_URI="https://<opnsense-host>"
export OPNSENSE_API_KEY="<key>"
export OPNSENSE_API_SECRET="<secret>"
export OPNSENSE_ALLOW_INSECURE="true"
```

## Architecture

The library has two layers:

- **`pkg/api`** — low-level HTTP client (`api.Client`). Handles authentication, retries, mutex locking, and the raw CRUD/RPC helpers (`api.Add`, `api.Get`, `api.Update`, `api.Delete`, `api.Call`).
- **`pkg/opnsense`** — high-level typed client (`opnsense.Client`). Aggregates all service controllers behind a single interface. `pkg/opnsense/client.go` is fully generated.
- **`pkg/<service>/`** — one package per OPNsense service (e.g. `unbound`, `firewall`, `ipsec`). Each package contains a generated `controller.go` and one generated file per resource or RPC group.

### Code generation

All files marked `// Code generated … DO NOT EDIT.` are produced by the generator:

```
schema/<service>.yml         ← source of truth for a service
internal/generate/schema/    ← schema structs (parsed by generator)
internal/generate/api/       ← generator entry-point + Go templates
internal/generate/generator/ ← template renderer (gofmt applied after render)
```

Running `make all` invokes `go generate` in each `pkg/<service>/` directory (via `generate.go`), which calls `internal/generate/api/main.go -controller <service>` and produces:
- `pkg/<service>/controller.go` — Controller struct + Client() accessor
- `pkg/<service>/<filename>.go` — one file per resource or RPC group
- `pkg/opnsense/client.go` — aggregated Client interface

The CI `make-generate-check` workflow fails PRs if committed generated files are out of sync with schemas.

### Mutex on writes

All mutating operations (Add, Update, Delete) acquire a global mutex (`GlobalMutexKV` with key `"OPNSENSE"`) before writing. This prevents concurrent writes from interleaving with the mandatory service reconfigure that follows every mutation.

## Key Conventions

### Never edit generated files
Edit `schema/<service>.yml`, then run `make all`. The `generate.go` in every service package must contain **only** the `//go:generate` directive and `package` declaration — nothing else.

### Adding a new service
1. Create `schema/<service>.yml` (use an existing file as reference).
2. Create `pkg/<service>/generate.go` with the `//go:generate` directive (copy from any existing service).
3. Run `make all`.

### Schema YAML structure
- `resources` block → CRUD resources. Each resource becomes `Add<Name>`, `Get<Name>`, `Update<Name>`, `Delete<Name>` methods.
  - `readOnly: true` — omits Add/Update/Delete.
  - `getByFilter: true` — uses `api.GetFilter` (lookup by key in a flat map) instead of `api.Get` (lookup by UUID).
  - `getAll: true` — adds a `Get<Name>All` method.
  - `reconfigure: "null"` on an endpoint — suppresses the post-mutation reconfigure call for that resource.
- `rpc` block → non-CRUD calls. Each entry in `rpc_calls` becomes a typed method on the controller.
- `reconfigureEndpoint` at the top level sets the default for all resources in the controller; individual resources can override it.

### SelectedMap types
OPNsense returns enumerated fields as a map of `{key: {selected, value}}` objects. Three custom types handle these transparently:
- `api.SelectedMap` — single selection; marshals/unmarshals as a plain string key.
- `api.SelectedMapList` — multi-selection; marshals as comma-separated string.
- `api.SelectedMapListNL` — multi-selection; marshals with newline separators.

When constructing a resource, assign the key directly: `Type: api.SelectedMap("host")`.

### Error handling
`api.Get` returns `*errs.NotFoundError` when the resource doesn't exist (OPNsense returns an unmarshalable response for missing UUIDs). Always check with `errors.As(err, &notFound)`.

### Monad wrapping
OPNsense's API wraps request/response bodies in a top-level key (the "monad"). The schema `monad` field sets this key. `api.Add`/`api.Update` wrap outgoing structs; `api.Get` unwraps incoming responses automatically.

## Keeping AGENTS.md Up-to-Date

Update this file whenever you make changes that affect how someone would work in this repo:

- **New service added** — add it to the architecture overview if it introduces a new pattern.
- **Generator or schema structure changes** — update the Code generation section and schema YAML structure reference.
- **New `api` types or helpers** — document them in Key Conventions (e.g. a new `SelectedMap` variant).
- **Build/test workflow changes** — update the Commands section (targets, env vars, flags).
- **New conventions established** — add a subsection under Key Conventions.

The goal is that this file always reflects the current state of the repo so that any contributor (human or AI) can orient themselves from it alone.
