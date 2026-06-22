# Leader Prompt: Project Setup

Copy-paste the context below into a strong leader model chat. The output must be a concrete setup plan, including repo structure, initial commands, and codegen wiring.

---

You are setting up a new project repository. I will give you the project name, a summary, and the chosen stack. Produce a setup checklist and initial file list.

## Default stack

- Go backend (Go 1.22+) with Fiber or standard `net/http`.
- React 19 + Vite + strict TypeScript frontend.
- PostgreSQL 16 + sqlc.
- Centrifuge WebSockets.
- Protobuf/gRPC for any backend-to-mobile or inter-service contracts.
- `libs/gogen` for Go → TypeScript model generation.

## Output format

1. **Repo root structure** — list directories and key files.
2. **Initial commands** — exact `go mod init`, `npm create vite@latest`, sqlc init, etc.
3. **Makefile targets** — `make generate`, `make test`, `make dev`.
4. **Codegen wiring** — how proto, sqlc, and gogen fit into `go generate` or Makefile.
5. **CI starter** — fresh generated-code check and test run.
6. **First three features to spec** — in priority order.
7. **Open questions** — anything that blocks setup.

## Constraints

- Do not propose custom frameworks or tools outside the default stack unless I explicitly asked.
- Every target in the Makefile must be concrete and runnable.
- Include a `.env.example` listing.
- Apply the Ponytail ladder to the stack: question whether each tool is needed. A single-service app may not need Protobuf/gRPC. sqlc is always worth it. gogen is worth it when frontend consumes backend types.
- Fewest dependencies that satisfy real requirements. Libraries are welcome when they reduce code size and complexity. Every dependency must justify itself by making the code meaningfully smaller or safer.
- No boilerplate, no scaffolding "for later."

## Forbidden shortcuts

- Do not leave steps as "configure X"; give the file content or exact command.
- Do not omit testing or codegen from CI.
- Do not skip authorization wiring (JWT middleware, user context).
- Do not add codegen for boundaries that don't exist.
- Do not create abstraction layers before there's a second implementation.

Return the plan as markdown with code blocks for commands and file contents.
