# opnsense-go
Go library for the OPNsense API

## Development

This library generates the underlying API from custom schema files located in the `schema` directory using `go generate`.  OPNSense appears to have a pretty consistent API for getting and setting parameters, so this approach tends to reduce code duplication.  The code used to generate the API is located in the `internal/generate` directory. 

There are two types of generated objects: individual `controllers` and the opnsense `client` itself. The `controllers` are used to support individual services/components in OPNSense and the `client` represents the service as a whole. 

### Adding Components/Services

Components can be added by creating a `{service}.yml` schema file describing the component/service API under the `schema` directory (for now use the existing schema as a reference) and then adding a `pkg/{service}/generate.go` file. The `generate.go` file can be copied from the existing servicess.  The package can then be regenerated using `make all` in the base directory. 


