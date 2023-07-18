# Seafowl UDF Go

A simple Seafowl UDF demo, intended to follow in the footsteps of the [Rust example](https://github.com/splitgraph/seafowl-udf-rust) but implemented in Go instead.

# HOWTO

Dependencies: Go, tinygo
| task | command |
|---|---|
| tinygo compile with wasi[^1] | `tinygo build -o seafowl-udf-go.wasm -target=wasi`
| tinygo compile | `tinygo build -o seafowl-udf-go.wasm`
| tests | `go test -v` |
| compile (use tinygo for WASM) | `go build` should output `seafowl-udf-go`, run it via `./seafowl-udf-go 1 2` |

---

# How to run in Seafowl

It's basically identical to the Rust UDF [docs](https://github.com/splitgraph/seafowl-udf-rust#loading-the-wasm-module-into-seafowl-as-a-udf)

1. Compile into a WASM using `tinygo build -o seafowl-udf-go.wasm -target=wasi`
2. Install the module using [`./create_udf.sh`](./create_udf.sh)
3. Run Seafowl (e.g. `seafowl` or the Docker image, whatever)
4. Invoke the function using [`./query_udf.sh`](./query_udf.sh)

   You can invoke it directly via SQL, the shell script is just a convenience e.g.

   ```sql
   SELECT addi64(1, 2) AS RESULT
   ```

   using e.g. [curl](https://seafowl.io/docs/guides/querying-http) or [psql](https://seafowl.io/docs/guides/querying-postgresql)

   And you should get back

   ```
   $ ./query_udf.sh
   SELECT addi64(1, 2) AS RESULT;
   HTTP/1.1 200 OK
   content-type: application/json; arrow-schema=%7B%22fields%22%3A%5B%7B%22children%22%3A%5B%5D%2C%22name%22%3A%22result%22%2C%22nullable%22%3Atrue%2C%22type%22%3A%7B%22bitWidth%22%3A64%2C%22isSigned%22%3Atrue%2C%22name%22%3A%22int%22%7D%7D%5D%2C%22metadata%22%3A%7B%7D%7D
   x-seafowl-query-time: 426
   vary: Authorization, Content-Type, Origin, X-Seafowl-Query
   transfer-encoding: chunked
   date: Tue, 18 Jul 2023 06:34:14 GMT

   {"result":3}
   ```

[^1]: If you don't need concurrency/goroutines, the binary size can be squeezed even smaller via `tinygo build -scheduler=none -o seafowl-udf-go.wasm -target=wasi`
