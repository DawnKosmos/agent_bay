# Leader Prompt: style.md

Copy-paste the context below into a strong leader model chat. The output must be a `style.md` file that implementers can enforce during code review.

---

You are the technical leader for a software project. Write a `style.md` quality contract using the exact sections below. The contract must be specific enough that an implementer can read a PR and know whether it passes without asking questions.

## Input I will provide

- The `design.md` for the project.
- Stack: Go backend + React/TypeScript + PostgreSQL + Centrifuge WebSockets (unless I stated otherwise).
- Any project-specific package paths (e.g., `cmd/carryover/`, `server/`, `handler/http/`, `service/`, `db/query/`, `ws/`, `ws_api/`, `frontend/src/gen/`).

## Output format

Write `style.md` with these exact sections:

1. **Error handling**
2. **Type safety**
3. **Testing**
4. **Authorization**
5. **Security rules**
6. **Logging & observability**
7. **Project-specific conventions**

## Constraints

- Go: no `panic` in request paths; errors returned explicitly.
- TypeScript: `strict: true`; no implicit `any`; generated types are source of truth.
- Every handler or component must check authorization.
- Tests must cover happy path, error path, and at least one authorization boundary.
- Project-specific conventions must include package layout, generated-code locations, and naming rules.

## Forbidden shortcuts

- Do not write generic advice such as "handle errors gracefully." State the exact rule, e.g., "return errors as explicit values; no `panic` in request paths."
- Do not leave sections as bullet placeholders; fill them with project-specific examples.
- Do not omit authorization. Deny-by-default is mandatory.

Return only the `style.md` content inside a markdown code block.
