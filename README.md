# opnsense-go

A Go client library for the OPNsense API.

## Installation

```sh
go get github.com/browningluke/opnsense-go
```

Requires Go 1.23 or later.

## Usage

### Creating a client

The library has two layers: a low-level `api.Client` that handles HTTP and
authentication, and a higher-level `opnsense.Client` that exposes typed
controllers for each OPNsense service.

```go
import (
    "github.com/browningluke/opnsense-go/pkg/api"
    "github.com/browningluke/opnsense-go/pkg/opnsense"
)

apiClient := api.NewClient(api.Options{
    Uri:           "https://<opnsense-host>",
    APIKey:        "<api-key>",
    APISecret:     "<api-secret>",
    AllowInsecure: true, // set to true when the firewall uses a self-signed certificate
})

client := opnsense.NewClient(apiClient)
```

`api.Options` also accepts retry tuning fields:

| Field | Default | Description |
|-------|---------|-------------|
| `MaxRetries` | 4 | Number of times to retry a failed request |
| `MinBackoff` | 1s | Minimum wait between retries |
| `MaxBackoff` | 30s | Maximum wait between retries |
| `Logger` | stdlib default | Custom `*log.Logger` for request/response logging |

### CRUD resources

Most OPNsense resources (host overrides, firewall aliases, IPsec connections,
etc.) follow the same Add/Get/Update/Delete pattern.

```go
import (
    "context"
    "github.com/browningluke/opnsense-go/pkg/unbound"
)

ctx := context.Background()

// Add
id, err := client.Unbound().AddDomainOverride(ctx, &unbound.DomainOverride{
    Enabled:     "1",
    Domain:      "example.internal",
    Server:      "192.168.1.1",
    Description: "internal zone",
})

// Get
override, err := client.Unbound().GetDomainOverride(ctx, id)

// Update
override.Server = "192.168.1.2"
err = client.Unbound().UpdateDomainOverride(ctx, id, override)

// Delete
err = client.Unbound().DeleteDomainOverride(ctx, id)
```

Each mutating call (Add, Update, Delete) automatically triggers the relevant
OPNsense service reconfigure so changes take effect immediately.

### RPC / settings calls

Some controllers expose RPC-style calls for reading or writing global settings
rather than individual records.

```go
// Read the current Unbound settings
settings, err := client.Unbound().SettingsGet(ctx)
if err != nil {
    return err
}

// Modify a field and write back
settings.Unbound.Advanced.HideIdentity = "1"
result, err := client.Unbound().SettingsUpdate(ctx, &settings.Unbound)
```

### SelectedMap and SelectedMapList

OPNsense represents enumerated fields and multi-select fields as maps with a
`selected` key rather than plain strings. This library provides two types that
unmarshal those responses transparently:

- `api.SelectedMap` -- a single-selection field; unmarshals to the selected key
  as a plain `string`-backed type.
- `api.SelectedMapList` -- a multi-selection field; unmarshals to a
  `[]string`-backed type. Marshals back as a comma-separated string.
- `api.SelectedMapListNL` -- same as `SelectedMapList` but marshals with
  newline separators instead of commas.

When constructing a resource to send to the API, assign the key directly:

```go
import "github.com/browningluke/opnsense-go/pkg/api"

alias := &firewall.Alias{
    Type: api.SelectedMap("host"),
}

override := &unbound.HostOverride{
    Type: api.SelectedMap("A"),
}
```

### Error handling

When a GET request targets a resource that does not exist, the library returns
an `*errs.NotFoundError`:

```go
import "github.com/browningluke/opnsense-go/pkg/errs"

override, err := client.Unbound().GetDomainOverride(ctx, id)
if err != nil {
    var notFound *errs.NotFoundError
    if errors.As(err, &notFound) {
        // resource was deleted upstream
    }
    return err
}
```

## Development

### Code generation

The typed controllers and data structs are generated from YAML schema files
under the `schema/` directory. The generator lives in `internal/generate/`.

Regenerate everything:

```sh
make all
```

There are two categories of generated output:

- **Controllers** (`pkg/<service>/controller.go`, `pkg/<service>/*.go`) --
  one per service, generated from `schema/<service>.yml`.
- **Opnsense client** (`pkg/opnsense/client.go`) -- aggregates all controllers
  into a single `Client` interface.

### Adding a new service

1. Create `schema/<service>.yml` describing the endpoints and data types.
   Use an existing schema file as a reference.
2. Create `pkg/<service>/generate.go` with the `//go:generate` directive.
   Copy the file from any existing service package.
3. Run `make all` to generate the controller and update the opnsense client.

### Running tests

Tests require a live OPNsense instance. Set the following environment variables
before running:

```sh
export OPNSENSE_URI="https://<opnsense-host>"
export OPNSENSE_API_KEY="<key>"
export OPNSENSE_API_SECRET="<secret>"
export OPNSENSE_ALLOW_INSECURE="true"  # if using a self-signed certificate
```

Run tests for a specific package:

```sh
go test -v ./pkg/<service>/...
```

## License

[MIT](LICENSE)

