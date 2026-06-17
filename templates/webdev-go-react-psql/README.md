# webdev-go-react-psql

Default web template for `agent_bay`. Use this when you are building a browser-first web application with a Go backend, a React/TypeScript frontend, PostgreSQL for persistence, and Centrifuge for live updates.

## When to use this template

- Backend: Go with either `net/http`, Echo, Gin, or Fiber.
- Database: PostgreSQL (>= 16 recommended), accessed via sqlc.
- Real-time: Centrifuge WebSockets, with server-side publishers and generated TS event types.
- Frontend: React 19, Vite, strict TypeScript, generated API/WS types.
- Auth: JWT issued by backend, populated from GitHub OAuth/OIDC.

## Quickstart

1. Copy this template into your target project as `docs/stack-template/`.
2. Fill `design.template.md` with the project's actors, bounded contexts, and security model.
3. Fill `style.template.md` with exact package locations (e.g., `cmd/carryover/`, `server/`, `handler/`, `db/query/`, `ws/`, `ws_api/`, `frontend/src/gen/`).
4. For the first 3 features, fill one `feature.template.md` per feature.
5. Scaffold the repo layout, generate types (`make generate`), and run tests.
6. Implement features one at a time; commit each feature end-to-end before starting the next.

## Folder contents

- `design.template.md` — lightweight design doc with stack hints.
- `style.template.md` — stack-specific quality contract.
- `feature.template.md` — per-feature spec template.
- `backend/structure.md` — Go package layout.
- `backend/reference-handler/` — a canonical HTTP handler example.
- `frontend/structure.md` — React/TS package layout.
- `frontend/reference-component/` — a canonical component example.
- `db/schema-patterns.md` — Postgres conventions.
- `ws/centrifuge.md` — type-safe WebSocket architecture.
- `codegen/README.md` — codegen tools and CI wiring.
