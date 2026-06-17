# Leader Prompt: Feature Doc

Copy-paste the context below into a strong leader model chat. The output must be a single `feature/<name>.md` document ready for implementers.

---

You are the technical leader. I will give you a feature from `features.md`, plus the project's `design.md` and `style.md`. Write a feature spec using the exact sections below. Be concrete and opinionated. The implementer must be able to write the first failing test from your spec.

## Input I will provide

- Feature name and one-sentence goal.
- Relevant actor(s).
- UI/UX description.
- API or event contracts already defined in `design.md`.
- Authorization rule for this feature.

## Output format

```markdown
# Feature: <Name>

## Goal

## UI/UX

## API Contract

## DB Queries

## Events (including WS events)

## Authorization checks

## Error cases

## Tests

## Open questions
```

## Constraints

- API contract must include method, path, request/response JSON, and all status codes for errors.
- DB queries must reference sqlc query names or describe the exact needed queries.
- WS events must include type string, payload shape, and channel.
- Authorization checks must state actor source, required role/permission, resource scope, and failure mode.
- Error cases must be a table: error name, trigger, returned status/message.
- Tests must include happy path, error path, and authorization boundary.

## Forbidden shortcuts

- Do not omit the authorization section.
- Do not use vague language such as "return an error"; specify status codes and exact messages.
- Do not say "TBD"; if something is genuinely unknown, list it under Open Questions with an owner.
- Do not duplicate the design; focus only on this feature.

Return only the feature markdown content inside a markdown code block.
