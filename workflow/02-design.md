# 02 — Design

`design.md` is the single source of truth for what the system does and why. It is written by the leader model together with the human. It must be complete enough that an implementer can build the right thing without asking clarifying questions every ten minutes.

## Lightweight `design.md` template

```markdown
# <Project> Design

## 1. Problem & Goal
Two paragraphs: what problem this app solves and what success looks like for the first shipped phase.

## 2. Actors & Use Cases
Actors: user, admin, background worker, external service.
Use cases as `<actor> <verb> <object>`.

## 3. Bounded Contexts / Modules
List modules (e.g., `auth`, `project`, `billing`, `notification`). For each, state its public API surface and what it does NOT own.

## 4. Data Model Overview
Entities, relationships, and ownership. No SQL here.

## 5. High-level API/Events
HTTP resources, RPC methods, and WS event families. Include who can call what at this level.

## 6. Deployment & Infra
Runtime, DB, queues, caches, object storage, secrets management.

## 7. Security Model
Authn strategy (OAuth/OIDC/JWT/session) and authz strategy (RBAC, ABAC, resource-scoped, deny-by-default). No implementation details.

## 8. Implementation Phases
Numbered phases with scope boundaries. Each phase ships working, tested, observable code.
```

## What NOT to put here

`design.md` is a contract, not a spec sheet. Do not include:

- Per-field UI layouts or component hierarchy.
- SQL schema, migration names, or indexing decisions.
- Concrete code, file paths, package names, or generated types.
- Internal error message strings exposed to users.

Those belong in feature specs (`feature/<name>.md`), style docs (`style.md`), reference code, and generated contracts.

## Review signal

A good `design.md` lets an implementer list the first three files they will create and the authorization check each one needs.
