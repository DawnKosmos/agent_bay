# Codegen

This template relies on generated code at every service boundary. Hand-maintained mirrors are forbidden.

## When to use which tool

| Boundary | Tool | Source of truth | Output |
|----------|------|-----------------|--------|
| Backend ↔ backend/mobile | Protobuf/gRPC | `libs/proto/<service>/*.proto` | `api/gen/<service>/` |
| Backend ↔ PostgreSQL | sqlc | `db/query/*.sql` + schema | `db/query/*.go` |
| Backend structs ↔ frontend TS | `libs/gogen` | `ws_api/*.go`, `api/*.go` | `frontend/src/gen/models.ts` |

## Wiring each tool into the build

### Protobuf/gRPC

Add to `Makefile`:

```make
proto:
	protoc --go_out=paths=source_relative:api/gen \
	       --go-grpc_out=paths=source_relative:api/gen \
	       libs/proto/**/*.proto
```

Put `//go:generate make proto` in a top-level package if you prefer `go generate`.

### sqlc

Add `sqlc.yaml` at project root:

```yaml
version: "2"
sql:
  - engine: "postgresql"
    queries: "db/query/"
    schema: "db/migration/"
    gen:
      go:
        package: "query"
        out: "db/query"
        sql_package: "pgx/v5"
```

Add `make sqlc`:

```make
sqlc:
	sqlc generate
```

### `libs/gogen`

Reference implementation: `libs/gogen/gogen.go` relative to the repository root (in a parent project with `agent_bay` as a submodule, that is `../../libs/gogen/gogen.go` from this file).

Create a small `cmd/gen/main.go` or a `go generate` directive that constructs a `gogen.Generator`:

```go
package main

import (
	"log"
	"github.com/DawnKosmos/carryover/libs/gogen"
)

func main() {
	g := gogen.New(
		[]string{
			"./api",
			"./ws_api",
		},
		"./frontend/src/gen",
	)
	if err := g.Generate(); err != nil {
		log.Fatal(err)
	}
}
```

Expose it in the Makefile:

```make
gen:
	go run ./cmd/gogen
```

## CI freshness check

Add a CI step:

```bash
make generate
git diff --exit-code -- api/gen/ db/query/ frontend/src/gen/
```

This fails if checked-in generated files diverge from source.

## Ownership

- Source files are hand-written and reviewed.
- Generated files are checked in so the project builds without the generator, but they are read-only.
- If a generated file needs editing, change the source and regenerate.
