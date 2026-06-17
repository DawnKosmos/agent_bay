# Design: <Project>

Use this template as the starting `design.md` for a `webdev-go-react-psql` project.

## 1. Problem & Goal

- What problem does the product solve?
- What does Phase 1 ship? What is intentionally out of scope?

## 2. Actors & Use Cases

Actors:

- `end_user` — primary consumer.
- `admin` — manages accounts and global settings.
- `background_job` — scheduled or event-driven work.
- `external_service` — third-party integrations.

Use cases: `<actor> <verb> <object>`.

## 3. Bounded Contexts / Modules

Example modules:

- `auth` — login, logout, token refresh, user identity.
- `project` — project lifecycle and membership.
- `chat` — real-time messages and presence.
- `notification` — push/WS notifications and preferences.

For each module, list what it owns and explicitly what it does NOT own.

## 4. Data Model Overview

Entities, relationships, and ownership. No SQL.

Example:

- `User` — accounts, OAuth identities.
- `Project` — created by a user, has many members.
- `ChatMessage` — belongs to a project and a sender.
- `Notification` — belongs to a recipient.

## 5. High-level API/Events

- HTTP REST under `/api/v1/` (or version as needed).
- WebSocket events via Centrifuge on `personal.{userID}` and project-scoped channels.
- gRPC for any backend-to-mobile or backend-to-backend contract.

List the major endpoints and event families, grouped by module.

## 6. Deployment & Infra

- Go binary built with `go build ./cmd/<app>`.
- PostgreSQL 16 in Docker for local dev; managed Postgres in prod.
- Centrifuge running alongside the Go binary or as a separate service.
- Frontend served as static files by the Go server or a CDN.
- Secrets via environment variables.

## 7. Security Model

- **Authentication**: JWT issued after GitHub OAuth via the Go backend.
- **Authorization**: resource-scoped RBAC; every action checks actor + resource membership/ownership.
- **Deny-by-default**: requests without a valid actor or with insufficient scope are rejected.

## 8. Implementation Phases

1. Auth + user context
2. Core CRUD module (projects, tasks, etc.)
3. Real-time layer (Centrifuge + first WS event)
4. Integrations / advanced features

Each phase ships tested, observable, deployable code.
