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

0. **Senior mindset (the Ponytail ladder)** — the ladder, anti-abstraction rules, `ponytail:` comment convention.
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
- Testing: one runnable check for non-trivial logic, authorization boundary test for every entry point, trivial one-liners need no test. YAGNI applies to tests.
- No interface with one implementation. No factory for one product. No wrapper that only delegates. Add abstractions when a second implementation appears.
- Fewest files possible. If two files are always edited together, they should be one file.
- Project-specific conventions must include package layout, generated-code locations, and naming rules.

## Forbidden shortcuts

- Do not write generic advice such as "handle errors gracefully." State the exact rule, e.g., "return errors as explicit values; no `panic` in request paths."
- Do not leave sections as bullet placeholders; fill them with project-specific examples.
- Do not omit authorization. Deny-by-default is mandatory.
- Do not mandate full test suites for trivial code.
- Do not add abstraction rules that enable single-implementation interfaces.

Return only the `style.md` content inside a markdown code block.
