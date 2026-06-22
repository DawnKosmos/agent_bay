# 04 — Feature Implementation Docs

A feature spec is a mini-contract. It tells the implementer exactly what to build, what to check, and how to prove it works. Write one markdown file per feature in `docs/feature/<feature-name>.md`.

## Ponytail pre-flight (answer before writing the spec)

- Does this feature need to exist? If no, write one line saying so and stop.
- Can stdlib, a native platform feature, or an existing dependency cover this? If yes, reference it and stop.
- Is this a one-liner? If yes, write the line and stop.
- Only if none of the above: write the full spec below.

For simple features (one endpoint, no new DB tables, no WS events), collapse the spec into a 5-line summary: goal, endpoint, authz rule, error cases, one check. Don't manufacture sections for structure's sake. YAGNI applies to documentation too.

## Required sections

```markdown
# Feature: <Name>

## Goal
One sentence: what capability this feature adds and which actor uses it.

## UI/UX
Screens, flows, and user-facing copy. Include validations users see.

## API Contract
Endpoints, request/response shapes, status codes, errors. Prefer references to generated types.

## DB Queries
What rows are read or written. Reference sqlc query files or schema locations.

## Events (including WS events)
What events are emitted, who receives them, and what channel/queue they travel on.

## Authorization checks
For every entry point, state: actor, required permission/role, resource scope.

## Error cases
List each failure mode and the returned error/status.

## Tests
Unit, integration, and edge-case tests required.

## Open questions
Anything unresolved. Every open question must have an owner.
```

## Concrete short example: Create a chat message

```markdown
# Feature: Create chat message

## Goal
Authenticated sender posts a message visible to members of the same project chat.

## UI/UX
- A text input below the message list.
- Submit via Enter, disabled when text is empty or >4000 chars.
- Optimistically append to local list; replace with server data on success.

## API Contract
`POST /api/projects/{project_id}/chat/messages`

Request:
```json
{
  "text": "hello team"
}
```

Response `201`:
```json
{
  "messageId": "msg_...",
  "fromId": "usr_...",
  "text": "hello team",
  "createdAt": "2024-..."
}
```

Errors:
- `400` — text empty or >4000 chars (return explicit field error).
- `401` — not authenticated.
- `403` — not a member of the project.
- `404` — project not found (same as 403 for privacy if not member).

## DB Queries
- `GetProjectMember(ctx, projectID, userID)` — verify membership before insert.
- `CreateChatMessage(ctx, arg{fromID, projectID, text})` — insert and return row.

## Events
On success, publish `message.received` via Centrifuge to each project member's personal channel (`personal.{userID}`) so they do not need to poll.

## Authorization checks
- Handler extracts `userID` from JWT.
- Before insert, call the membership repo. Deny if absent.
- The returned message must include only fields visible to the caller.

## Error cases
- Empty text
- Text exceeds limit
- Non-member requests project chat
- Race: project deleted between membership check and insert (return 404)

## Tests
- Handler returns 201 for member
- Handler returns 403 for non-member
- Handler returns 400 for empty text
- Service emits WS event to all members
- Repository row matches inserted text

## Open questions
- Do we support file attachments in V1? Owner: product owner.
```

## Quality bar

A feature spec is ready when an implementer can read it and write the first failing test without asking a question.
