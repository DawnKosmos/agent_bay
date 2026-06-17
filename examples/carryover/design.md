# Design: carryover

This is a draft design document for `carryover` based on the known stack and product direction. It follows the lightweight template from `agent_bay/workflow/02-design.md`.

## 1. Problem & Goal

carryover is a self-hosted issue tracker inspired by Jira, augmented with AI-assisted workflows through an assistant called Levi. Phase 1 ships user identity, project/work tracking, real-time notifications, and a foundation for Levi integration.

## 2. Actors & Use Cases

- `end_user` — creates and manages projects, issues, comments; receives real-time updates.
- `admin` — manages users and global settings (future phase).
- `levi` (external service, via gRPC) — reads and writes project data on behalf of the user when authorized.
- `background_job` — runs digest, search indexing, and async notifications.

Use cases:

- User creates a project.
- User creates an issue inside a project.
- User posts a comment on an issue.
- User receives a real-time notification when mentioned.
- Levi queries project state through the gRPC boundary.

## 3. Bounded Contexts / Modules

- `auth` — GitHub OAuth login, JWT issuance, logout, token refresh.
- `project` — projects, members, roles.
- `issue` — issues, statuses, priorities, assignments.
- `comment` — issue comments and mentions.
- `notification` — user notifications and preferences; publishes WS events.
- `search` — full-text search powered by Typesense.
- `levi_gateway` — gRPC boundary that exposes a controlled view of project/issue data to the Levi service.

## 4. Data Model Overview

- `User` — account, OAuth identity, preferences.
- `Project` — owned by a user, has many members.
- `ProjectMember` — links users to projects with a role.
- `Issue` — belongs to a project, has a reporter, optional assignee, status, priority.
- `Comment` — belongs to an issue and an author; mentions are parsed.
- `Notification` — belongs to a user; delivered via WebSocket and read later.

## 5. High-level API/Events

- HTTP REST under `/api/` (internal versioning to be decided).
- WebSocket events via Centrifuge on `personal.{userID}` for notifications and direct messages.
- gRPC service in `libs/proto/levi/levi.proto` for Levi integration.

## 6. Deployment & Infra

- Go backend in `apps/carryover/`.
- React/TypeScript frontend built with Vite.
- PostgreSQL 16 with sqlc-generated queries.
- Typesense for search.
- Centrifuge for WebSockets.
- Secrets via environment variables; no secrets in the repo.

## 7. Security Model

- **Authentication**: GitHub OAuth callback issues a JWT stored in a secure, httpOnly, same-site cookie.
- **Authorization**: resource-scoped; users can read/write only projects they are members of. Deny-by-default.
- **Levi integration**: every Levi request is authenticated as a user and inherits that user's project membership scope.

## 8. Implementation Phases

1. Auth + user context
2. Project + member management
3. Issue + comment CRUD
4. Real-time notifications via Centrifuge
5. Typesense search integration
6. Levi gRPC gateway

## Notes and open questions

- Authorization matrix (roles: owner, member, viewer) is pending.
- Exact event types beyond `message.received` and `notification.created` are to be defined per feature.
- Public API versioning strategy not decided yet; keep routes simple until it is.
