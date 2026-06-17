# 08 — Review Checklist

Use this checklist before merging any PR. A single unchecked item blocks merge unless the reviewer explicitly documents why it is acceptable.

## Code quality

- [ ] All tests pass locally and in CI.
- [ ] New code has tests: happy path, error path, and at least one authorization boundary.
- [ ] No `any` or `unknown` types without a justification comment.
- [ ] No `panic` in Go request paths.
- [ ] No swallowed errors; every deferred/ignored return is intentional and commented.
- [ ] No dead code, commented-out experiments, or unused imports.

## Authorization

- [ ] Every entry point checks actor identity.
- [ ] Every action checks resource scope or role/permission matrix.
- [ ] Deny-by-default is enforced; missing auth state returns 403, not 200.

## Generated code

- [ ] Generated files are up to date (`make generate` produces no diff).
- [ ] Generated files are not hand-edited.
- [ ] New or changed service boundaries use codegen (proto, sqlc, gogen).

## Security

- [ ] Input validation happens at the boundary.
- [ ] SQL queries use sqlc or parameterized statements.
- [ ] Secrets are not hard-coded or logged.
- [ ] CORS, CSP, and rate limiting are configured per `style.md`.

## Data & migrations

- [ ] DB migrations are reversible where possible.
- [ ] Migrations do not destroy data unless explicitly approved.
- [ ] Indexes and foreign keys are added for new relationships.

## Observability

- [ ] Structured logs include request IDs and action context.
- [ ] New errors are logged once with enough detail to debug.
- [ ] Metrics or tracing added for new operations where appropriate.

## Documentation

- [ ] Feature spec matches the implementation.
- [ ] Reference code is updated if the pattern changed.
- [ ] PR description explains what changed and why.

## Merge rule

If the PR changes authorization, error handling, or generated contracts, it requires a second reviewer.
