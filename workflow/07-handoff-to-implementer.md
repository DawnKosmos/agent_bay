# 07 — Handoff to Implementer

The handoff is the boundary between design and construction. The implementer model's job is to read the docs and produce PR-ready code that follows them exactly. The leader's job is to make that possible.

Before handing off, the leader applies the Ponytail ladder to the feature: does this need to be built, or can stdlib/native/existing dep cover it? If the answer is "don't build it," the handoff is one line, not a spec.

The implementer is also empowered to push back: if the implementer sees a stdlib/native/existing-dep solution that the spec missed, they implement the lazy version and note it in the PR description.

## Required input for an implementer

Before starting, an implementer must have:

1. `docs/design.md` — what the system does and why.
2. `docs/style.md` — the quality contract.
3. `docs/feature/<name>.md` — the specific feature to build.
4. Reference handler/component files showing the pattern.
5. Access to generated types (`api/gen/`, `frontend/src/gen/`, `db/query/`) or the command to generate them.

If any of these are missing, the handoff is incomplete.

## Output of an implementer

For each feature, the implementer produces:

- Production code matching the feature spec and style contract.
- Tests covering non-trivial logic and authorization boundary. Trivial code needs no test.
- Updated generated code if the contract changed.
- A short PR description that lists the files added, the authorization rule, and any `ponytail:` simplifications with their ceilings.

## How implementers read the docs

1. Start with the feature spec. Identify the actor, the entry points, the authorization rule, and the error cases.
2. Open the reference implementation. Copy its structure, not by rote but by matching the pattern.
3. Use generated types. Do not define parallel structs.
4. Write the smallest check that fails if the logic is wrong, then the code. For trivial code, no check is needed. For authorization boundaries, the check is mandatory.
5. Run the test suite and lint before marking ready.

## Escalation rule

> If a spec is ambiguous, stop and ask. Do not guess.

Examples of ambiguity:

- The feature spec mentions "admins can do X" but the role matrix does not list "admin."
- The API contract omits a status code for an error case.
- The DB query does not specify ordering or pagination.
- The WS event channel is unclear.

When in doubt, open a comment on the feature spec or ping the leader model. The leader fixes the spec; the implementer does not improvise.
