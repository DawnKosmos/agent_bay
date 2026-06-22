# 03 — Style Contract

`style.md` is the single quality contract. Every implementer must follow it. Code that violates it is rejected in review, even if it works.

## 0. Senior mindset (the Ponytail ladder)

Before writing any code, stop at the first rung that holds:

1. **Does this need to be built at all?** (YAGNI) — speculative need = skip it, say so in one line.
2. **Does the standard library already do this?** Use it.
3. **Does a native platform feature cover it?** Use it (e.g., `<input type="date">` over a picker lib, CSS over JS, DB constraint over app code).
4. **Does an already-installed dependency solve it?** Use it. Libraries are welcome when they meaningfully reduce code size and complexity.
5. **Can it be one line?** Make it one line.
6. **Only then:** write the minimum code that works.

Rules:
- No unrequested abstractions: no interface with one implementation, no factory for one product, no config for a value that never changes.
- No boilerplate, no scaffolding "for later." Later can scaffold for itself.
- Deletion over addition. Boring over clever. Fewest files possible.
- Mark deliberate simplifications with a `ponytail:` comment naming the ceiling and upgrade path.
- Question complex requests: "Did X; Y covers it. Need full X? Say so."
- Pick the edge-case-correct option when two approaches are the same size. Lazy means less code, not the flimsier algorithm.

Not lazy about: input validation at trust boundaries, error handling that prevents data loss, security, accessibility, anything explicitly requested.

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

- **Non-trivial logic** (branches, loops, parsers, money/security paths) leaves ONE runnable check — the smallest thing that fails if the logic breaks. No frameworks, no fixtures unless the project already has them.
- **Authorization boundary test for every entry point** — this is security, not optional. Every handler must have an error branch in its test.
- **Trivial one-liners need no test.** YAGNI applies to tests too.
- **When the project has a test suite**, add to it. When it doesn't, a single `assert`-based self-check or one small test file is enough.
- **Go**: table-driven tests when you do test. Name the fields `name`, `input`, `want`, `wantErr`.
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
- No interface with one implementation. No factory for one product. No wrapper that only delegates. Add abstractions when a second implementation appears, not before.
- Fewest files possible. If two files are always edited together, they should be one file.
- Generated code lives in clearly marked directories (e.g., `api/gen/`, `frontend/src/gen/`, `db/query/`). Generated files are checked in but never hand-edited.
- Naming: `CreateProject` not `Create`; `project_id` in DB maps to `ProjectID` in Go and `projectId` in JSON.
- One package has one reason to change. Handlers are thin; services hold behavior; repositories hold persistence logic.

## 8. Shared libraries

- **Use existing tools from `agent_bay/libs`** — Before writing new utility code, check `agent_bay/libs/go` and `agent_bay/libs/typescript` for reusable packages (e.g., `gogen`, validation helpers, middleware).
- **Contribute new tools** — If you identify a pattern that improves code quality (e.g., a reusable authorization helper, a common error wrapper, a type-safe builder), add it to `agent_bay/libs` instead of duplicating it across projects.
- Libraries in `agent_bay/libs` are versioned as Go modules and can be imported by other projects using Go workspaces, standard module requirements, or git submodules.
