# Backend Structure

This is the canonical Go project layout for the `webdev-go-react-psql` template.

## Package responsibilities

```text
cmd/<app>/
  main.go                 # Construct server, db, ws node; wire handlers.

server/
  server.go               # HTTP server, middleware, route registration.
  middleware.go           # Auth, request ID, recovery, logging.

handler/http/
  user.go                 # Thin HTTP handlers around user services.
  project.go              # Thin handlers around project services.
  chat.go                 # Thin handlers around chat services.

service/
  user.go                 # Business logic for users and auth.
  project.go              # Business logic for projects and membership.
  chat.go                 # Business logic for chat and authorization.
  errors.go               # Shared domain errors.

model/
  user.go                 # Domain types shared by service/handlers/db.
  project.go

repo/
  user.go                 # Optional abstraction over db/query.

db/
  query/                  # sqlc generated code (READ-ONLY)
  migration/              # sql-migrate / golang-migrate files

ws/
  ws.go                   # Centrifuge node setup.
  publish.go              # Publishers for personal/project channels.

ws_api/
  ws_api.go               # Hand-written Go structs for WS events.
                          # gogen consumes this package.
```

## Rules

1. **Handlers are thin**. Decode request, call service, map errors, encode response. No business logic.
2. **Services contain behavior**. Authorization checks, validation orchestration, and domain rules live here.
3. **Repositories contain persistence**. Prefer sqlc-generated queries. If you need a custom repo wrapper, keep it small and test it against a real database.
4. **Error mapping is explicit**. Domain errors become HTTP status codes in `handler/http/errors.go` or inline in the handler. Never expose raw DB driver errors to the client.
5. **Generated code is read-only**.Hand-editing `db/query/*.go` or `frontend/src/gen/models.ts` is forbidden.

## Example import direction

```text
handler/http -> service -> model -> db/query
handler/http -> ws? no; service -> ws (events as side effects)
ws_api -> gogen -> frontend/src/gen
```
