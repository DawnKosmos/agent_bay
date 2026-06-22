# 00 — Overview

`agent_bay` treats feature development as an assembly line. Each feature moves through up to eight stations, scaled to feature complexity. The first five are owned by a strong **leader** model together with the human product owner; the last three are owned by cheaper **implementer** models, with review feedback as the quality gate.

## Senior mindset (the Ponytail ladder)

Before building anything — a feature, a module, a phase, a file — stop at the first rung that holds:

1. **Does this need to be built at all?** (YAGNI) — speculative need = skip it, say so in one line.
2. **Does the standard library already do it?** Use it.
3. **Does a native platform feature cover it?** Use it (e.g., `<input type="date">` over a picker lib, CSS over JS, DB constraint over app code).
4. **Does an already-installed dependency solve it?** Use it. Libraries are welcome when they meaningfully reduce code size and complexity — don't reinvent what a good package already does well.
5. **Can it be one line?** Make it one line.
6. **Only then:** write the minimum code that works.

The ladder is a reflex, not a research project. Two rungs work → take the higher one and move on.

Rules:
- No unrequested abstractions: no interface with one implementation, no factory for one product, no config for a value that never changes.
- No boilerplate, no scaffolding "for later." Later can scaffold for itself.
- Deletion over addition. Boring over clever. Fewest files possible.
- Question complex requests: "Did X; Y covers it. Need full X? Say so." Never stall on an answer you can default.
- Pick the edge-case-correct option when two approaches are the same size. Lazy means less code, not the flimsier algorithm.
- Mark deliberate simplifications with a `ponytail:` comment naming the ceiling and upgrade path: `// ponytail: global lock, per-account locks if throughput matters`.

Not lazy about: input validation at trust boundaries, error handling that prevents data loss, security, accessibility, anything explicitly requested.

## Proportionality

Not every feature needs all 8 steps. Simple CRUD features can skip reference implementations (step 5) and collapse design + feature spec (steps 2–4) into one short doc. The leader decides the weight per feature. Manufacturing process steps for a simple feature is itself over-engineering.

## Intensity levels

| Level | What it means |
|-------|---------------|
| **lite** | Full 8-step flow, but each step names the lazier alternative. User picks. |
| **full** (default) | The ladder enforced. Steps scaled to feature complexity. Shortest diff, shortest explanation. |
| **ultra** | YAGNI extremist. Deletion before addition. Ship the one-liner and challenge the rest of the requirement in the same breath. |

## The manufacturer flow

1. **Discovery & techstack** — Understand the domain, actors, features, non-functional requirements, and pick the stack.
2. **Design** — Write a `design.md` that covers the problem, actors, bounded contexts, data model, API/events, security model, and implementation phases.
3. **Style contract** — Write a `style.md` that locks in the senior mindset, error handling, type safety, testing, authorization, security rules, logging, and project conventions.
4. **Feature specs** — For each feature, write a focused `feature/<name>.md` with UI/UX, API contract, DB queries, events, authorization, errors, tests, and open questions.
5. **Reference implementations** — The leader writes 1–2 exemplary handlers/components for non-trivial features. They demonstrate the style and become the training example for implementers.
6. **Codegen** — Generate types and clients at every service boundary that actually exists: Protobuf/gRPC, sqlc, and `libs/gogen` for Go → TypeScript.
7. **Implement** — Implementer models read the docs and produce PR-ready code with tests.
8. **Review** — Use the review checklist (including over-engineering review) before merge. Failed checks loop back to the implementer or escalate to the leader if the spec is wrong.
9. **Debt tracking** — Harvest `ponytail:` comments into a debt ledger so deliberate shortcuts are tracked, not forgotten.

## Two-model split

| Phase | Owner | Deliverable |
|-------|-------|-------------|
| 1–5 | Leader model + human | `docs/design.md`, `docs/style.md`, `docs/feature/<name>.md`, reference code, generated contracts |
| 6 | Leader or tooling | Wired codegen with CI freshness check |
| 7 | Implementer models | PR-ready code + tests per feature spec |
| 8 | Human reviewer (assisted by agents) | Merged code that passes the review checklist |
| 9 | Any | `PONYTAIL-DEBT.md` ledger of deferred shortcuts |

The leader is responsible for clarity. The implementer is responsible for following that clarity exactly. If an implementer hits an ambiguity, the rule is: **stop and ask, do not guess**. Guessing is how authorization checks and error paths get lost.
