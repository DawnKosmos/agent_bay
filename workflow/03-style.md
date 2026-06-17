# 03 — Style Contract

`style.md` is the single quality contract. Every implementer must follow it. Code that violates it is rejected in review, even if it works.

## 1. Error handling

- **Explicit errors as values**. In Go, return `(T, error)` from every function that can fail. No `panic` in request paths; `panic` is reserved for programmer errors detected at process startup.
- **Do not swallow errors**. If you cannot handle an error, return it with context using `fmt.Errorf("...: %w", err)`. Log the error only once, at the outermost layer that owns the request context.
- **TypeScript**: return `Result<T, E>` or use `try/catch` and report to the caller. Do not catch and return `undefined`. Async errors must reach a boundary where they are logged and surfaced to the user.
- **Map errors at the boundary**. A repository error becomes a domain error; a domain error becomes an HTTP status + structured response. Every handler must have an error branch in its test.

## 2. Type safety

- Go: avoid `interface{}` where a concrete type or generic will do.
- TypeScript: `strict: true`, no implicit `any`.
- **Generated types everywhere**. Use sqlc for DB models, `gogen` or OpenAPI generators for API shapes. Do not hand-write types that mirror generated ones unless the generator is broken.
- `any` is allowed only for true unknown external payloads, and must be immediately parsed into a typed value using `zod`, `io-ts`, or an explicit decoder.

## 3. Testing

- **Unit tests** for pure logic: validation, authorization decisions, domain calculations.
- **Integration tests** for handlers and repos: spin up a real test DB for repositories; use a real HTTP router (or a thin test server) for handlers.
- **Go**: table-driven tests. Name the fields `name`, `input`, `want`, `wantErr`.
- **TypeScript**: prefer pure function tests plus narrow UI integration tests. Use MSW for network mocking or a test DB when testing data hooks.
- Every test must assert on errors, not just success paths.

## 4. Authorization

- **Every action checks actor + resource scope**. No action is public unless explicitly documented as public.
- Deny-by-default. If authorization state is missing or ambiguous, return `403` or abort.
- Maintain a role/permission matrix in `design.md` and reference it in every feature spec.
- Use the same authorization path for HTTP, RPC, and WS handlers.

## 5. Security rules

- **Input validation** at the boundary. Sanitize, never trust client-provided IDs or filter values.
- **SQL injection prevention** through sqlc or parameterized queries. No string concatenation into SQL.
- **XSS**: escape output, set a strict Content Security Policy, and avoid `dangerouslySetInnerHTML`.
- **Secrets** live in environment variables only. Never log secrets, tokens, or raw passwords.
- **CORS** is explicit and narrow; do not reflect the `Origin` header.
- **Rate limiting** on auth endpoints and expensive operations.

## 6. Logging & observability

- Structured logs (JSON in production).
- Request IDs propagated through context.
- Trace context on outgoing calls.
- Log the actor, the action, the resource, and the outcome. Do not log PII or secrets.

## 7. Project-specific conventions

- Module boundaries match bounded contexts from `design.md`.
- Generated code lives in clearly marked directories (e.g., `api/gen/`, `frontend/src/gen/`, `db/query/`). Generated files are checked in but never hand-edited.
- Naming: `CreateProject` not `Create`; `project_id` in DB maps to `ProjectID` in Go and `projectId` in JSON.
- One package has one reason to change. Handlers are thin; services hold behavior; repositories hold persistence logic.
