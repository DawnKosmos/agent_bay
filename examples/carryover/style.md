# Style: carryover

Carryover-specific style rules. Generic rules live in `agent_bay/workflow/03-style.md`.

## Stack and tools

- **Backend**: Go 1.22+, Fiber framework, JWT auth after GitHub OAuth.
- **Database**: PostgreSQL with sqlc (`pgx/v5`).
- **Search**: Typesense via `libs/ts_util`.
- **Real-time**: Centrifuge WebSockets.
- **Go → TS types**: `libs/gogen` reading from `ws_api/`.
- **Backend ↔ Levi**: gRPC via `libs/proto/levi`.

## Package layout

```text
apps/carryover/
├── server/              # Fiber setup, middleware
├── handler/http/        # HTTP handlers (thin)
├── service/             # Business logic and authz
├── model/               # Domain types shared between layers
├── db/
│   ├── query/           # sqlc generated code (READ-ONLY)
│   └── migration/       # Up/down migrations
├── ws/                  # Centrifuge publishers
├── ws_api/              # WS contract structs; consumed by gogen
└── levi_gateway/        # gRPC handlers for Levi

frontend/
├── src/api/             # API client + react-query hooks
├── src/components/      # UI components
├── src/hooks/           # Auth, WS, feature hooks
├── src/gen/             # gogen output (READ-ONLY)
└── src/ws/              # Centrifuge subscription + dispatch

libs/
├── gogen/               # Generator source
├── proto/levi/          # Protobuf contracts
└── ts_util/             # Typesense helpers
```

## Codegen

- `ws_api/` structs are hand-written and reviewed.
- Run `go run ./cmd/gogen` (or Makefile target) to produce `frontend/src/gen/models.ts`.
- Run `sqlc generate` to refresh `db/query/`.
- Run `make proto` to regenerate gRPC stubs from `libs/proto/levi/levi.proto`.
- CI must fail if generated files are stale.

## Error handling

- Fiber handlers return `error` and map it via a centralized error handler.
- No `panic` in request paths.
- Errors crossing the gRPC boundary are converted to gRPC status codes with detailed messages.

## Type safety

- Go: avoid `interface{}`; use generics where useful.
- TypeScript: `strict: true`.
- Frontend network and WS payloads use types from `src/gen/models.ts`.

## Authorization

- Extract user identity from JWT middleware into Fiber context locals.
- Every handler must verify actor + project membership.
- Levi gRPC methods verify the same actor context and enforce project scope.

## Testing

- Go: table-driven tests.
- Handler tests use the real Fiber app with injected services.
- Repository tests use a temporary Postgres container.
- TS: vitest + MSW for API mocking.
