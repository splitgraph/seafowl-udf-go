# Seafowl UDF Go

A simple Seafowl UDF demo, intended to follow in the footsteps of the [Rust example](https://github.com/splitgraph/seafowl-udf-rust) but implemented in Go instead.

# HOWTO

Dependencies: Go, tinygo
| task | command |
|---|---|
| tinygo compile | `tinygo build -o seafowl-udf-go.wasm`
| tests | `go test -v` |
| compile (use tinygo for WASM) | `go build` should output `seafowl-udf-go`, run it via `./seafowl-udf-go 1 2` |
