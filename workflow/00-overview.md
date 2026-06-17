# 00 — Overview

`agent_bay` treats feature development as an assembly line. Each feature moves through eight stations. The first five are owned by a strong **leader** model together with the human product owner; the last three are owned by cheaper **implementer** models, with review feedback as the quality gate.

## The 8-step manufacturer flow

1. **Discovery & techstack** — Understand the domain, actors, features, non-functional requirements, and pick the stack.
2. **Design** — Write a `design.md` that covers the problem, actors, bounded contexts, data model, API/events, security model, and implementation phases.
3. **Style contract** — Write a `style.md` that locks in error handling, type safety, testing, authorization, security rules, logging, and project conventions.
4. **Feature specs** — For each feature, write a focused `feature/<name>.md` with UI/UX, API contract, DB queries, events, authorization, errors, tests, and open questions.
5. **Reference implementations** — The leader writes 1–2 exemplary handlers/components. They demonstrate the style and become the training example for implementers.
6. **Codegen** — Generate types and clients at every service boundary: Protobuf/gRPC, sqlc, and `libs/gogen` for Go → TypeScript.
7. **Implement** — Implementer models read the docs and produce PR-ready code with tests.
8. **Review** — Use the review checklist before merge. Failed checks loop back to the implementer or escalate to the leader if the spec is wrong.

## Two-model split

| Phase | Owner | Deliverable |
|-------|-------|-------------|
| 1–5 | Leader model + human | `docs/design.md`, `docs/style.md`, `docs/feature/<name>.md`, reference code, generated contracts |
| 6 | Leader or tooling | Wired codegen with CI freshness check |
| 7 | Implementer models | PR-ready code + tests per feature spec |
| 8 | Human reviewer (assisted by agents) | Merged code that passes the review checklist |

The leader is responsible for clarity. The implementer is responsible for following that clarity exactly. If an implementer hits an ambiguity, the rule is: **stop and ask, do not guess**. Guessing is how authorization checks and error paths get lost.
