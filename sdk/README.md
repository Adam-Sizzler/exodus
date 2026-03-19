# Core SDK (cerberus-node)

This package provides a unified core API client (similar to cerberus `XtlsApi` approach):

- one entrypoint: `sdk.New(Config)`
- one logical service: `API.Stats`
- sing-box implementation hidden behind `StatsService`

## Structure

- `sdk.New(...)` creates a `singbox` provider
- provider uses local generated protobuf stubs in `sdk/singboxapi`

## Why this exists

Before this SDK layer, core wiring logic was hardcoded in handler code (`api/grpc_response.go`).
Now business code talks only to `sdk.API`, while transport details stay inside this package.
