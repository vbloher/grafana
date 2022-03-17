## What

[view api](http://localhost)

This aims to define the unified alerting API as code. It generates OpenAPI definitions from go structs, in the `tooling/definitions` directory.

Routes are defined as `go generate`-style comments, inlined with the struct definitions.


## Running

`make openapi`

## Requires
 - Docker
 - [go-swagger](https://github.com/go-swagger/go-swagger)
 - [goimports](https://pkg.go.dev/golang.org/x/tools/cmd/goimports)
