# 06 — Codegen

Whenever two services or systems share a contract, generate the contract types instead of writing them by hand. Manual mirroring is the fastest way to silently break HTTP clients, WS consumers, and database queries.

## Rule: codegen at every service boundary

- Backend ↔ backend or backend ↔ mobile: use Protobuf/gRPC.
- Backend ↔ PostgreSQL: use sqlc.
- Backend structs ↔ frontend TypeScript: use `libs/gogen` (a custom Go struct → TypeScript interface generator) or an OpenAPI generator. In this repo, the generator source lives at `libs/gogen/gogen.go` relative to the repository root.
- Frontend network layer: derive from generated types; do not hand-maintain request shapes.

## Tools

### Protobuf/gRPC

For inter-service or backend-to-mobile contracts.

- Define `.proto` files in a shared location (e.g., `libs/proto/<service>/`).
- Generate Go code into `api/gen/<service>/`.
- Generate mobile/frontend stubs if needed.
- Version proto packages explicitly; do not reuse the same package for incompatible changes.

### sqlc

For type-safe SQL.

- Write queries in `.sql` files under `db/query/`.
- Generated Go code lives in `db/query/` or `db/sqlc/` depending on the project.
- Queries are reviewed as code. Parameterized queries are mandatory.

### `libs/gogen`

For Go struct → TypeScript interface generation.

- Keep API/event models in plain Go structs in `ws_api/` or `api/`.
- Tag JSON names explicitly so the generator emits correct TypeScript property names.
- Run `gogen` as a `go generate` step or Makefile target.
- Output to `frontend/src/gen/models.ts`.

## CI rule: generated files must be fresh

Every project must have a CI check that regenerates contracts and fails the build if checked-in generated files differ:

```bash
make generate
git diff --exit-code -- api/gen/ db/query/ frontend/src/gen/
```

Generated files are checked in so the project builds without the generator, but they must never be edited by hand.

## Where generated code lives

```text
api/gen/                  # Protobuf/gRPC generated stubs
db/query/                 # sqlc generated Go code
frontend/src/gen/         # gogen / openapi generated TS code
```

Treat these directories as read-only. If you need a different shape, change the source (proto, sqlc query, Go model) and regenerate.
