# 01 — Discovery & Techstack

Before writing code, the human + leader interview the idea into a small set of written artifacts. This is not negotiation; it is translation from human intent into constraints that implementer models can execute.

## Discovery interview

Answer the following in one or two paragraphs each:

- **Domain** — What problem does the app solve? What is the core noun (project, chat, order, invoice)?
- **Actors** — Who uses the system? Humans, background jobs, external services, admin tools.
- **Features** — A feature is a user-visible capability. List them as `<actor> <verb> <object>`: `user creates a project`, `admin suspends an account`, `system emits a digest email`.
- **Non-functional requirements** — Latency, throughput, offline support, multi-tenancy, self-hosting, mobile-first, real-time, audit trails.
- **Deployment target** — Docker on a single VPS, Kubernetes, serverless functions, edge, on-premise.
- **Compliance & security** — Authentication source (GitHub OAuth, OIDC, SAML), data residency, PII, retention, audit logging, rate limiting, public vs private API.

If an answer is unknown, write it down as an open question with an owner, do not make it up.

## Techstack decision guide

The default stack is:

- **Backend** — Go (Fiber or standard `net/http`), PostgreSQL.
- **Frontend** — React 19 + Vite + strict TypeScript.
- **Real-time** — Centrifuge WebSockets.
- **Data access** — sqlc for type-safe SQL queries.
- **Service contracts** — Protobuf/gRPC between backend services or backend-to-mobile; `libs/gogen` for Go structs → TypeScript interfaces.

### Ponytail stack check

Before committing to the default stack, ask: does the platform or an existing tool already cover what we need? Use PostgreSQL features (JSONB, full-text search, pub/sub via LISTEN/NOTIFY) instead of adding Redis/Elasticsearch when the workload fits. Question every non-functional requirement — is it real or speculative? "We might need multi-tenancy" is not a requirement, it's a maybe.

Libraries and packages are welcome when they reduce code size and complexity. Don't reinvent what a good package does well. For simple things, a few lines of your own code may be enough. For complex things (routing, validation, WebSocket management), a well-maintained library that shrinks your codebase is the lazy choice. Every dependency must justify itself by making the code meaningfully smaller or safer — not by novelty.

### When to deviate

Consider a different backend language or framework only for **technical** or **team** reasons, not novelty:

- Team already ships production code in another typed language (Rust, C#, Java).
- Strict latency requirements that favor a runtime with no GC pauses and the team can operate it.
- Existing shared library ecosystem (e.g., heavy ML Python integration).

Consider a different database only when:

- The access pattern is overwhelmingly document/key-value and relational joins are rare.
- A specialized index (graph, full-text, time-series) is needed and the operational cost is accepted; still keep PostgreSQL as the source of truth unless proven otherwise.

Consider a different frontend only when:

- The team owns a deployed design system in another framework.
- Mobile native performance is a first-class requirement from day one.

## Output artifacts

Create a `docs/` directory in the target project and commit the following before design begins:

```text
docs/
├── features.md      # Actor/feature list and non-functional requirements
└── techstack.md     # Rationale, versions, deployment target, auth source
```

These files are living documents. Update them when requirements change, but every feature spec must trace back to an entry in `features.md`.
