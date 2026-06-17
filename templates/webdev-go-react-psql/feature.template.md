# Feature: <FeatureName>

## Goal

One sentence describing the capability and the primary actor.

## UI/UX

- Screens / components involved.
- User-visible validation rules and copy.
- Loading and empty states.

## API Contract

`METHOD /api/v1/...`

Request shape (reference generated types if available):

```json
{
  "field": "value"
}
```

Response `200/201`:

```json
{
  "id": "..."
}
```

Error responses with status codes.

## DB Queries

- Name of sqlc query files and what each query does.
- Any transaction boundaries.

## Events (including WS events)

- Event type string.
- Payload shape.
- Channel or queue.
- Which actors receive it.

## Authorization checks

For each entry point:

- Actor source (JWT claim, context value).
- Required permission / role.
- Resource scope check.
- Failure mode.

## Error cases

| Error | Trigger | Returned status/message |
|-------|---------|-------------------------|
| ...   | ...     | ...                     |

## Tests

- Unit tests.
- Integration tests.
- WS/event tests.

## Open questions

- Question, owner, deadline.
