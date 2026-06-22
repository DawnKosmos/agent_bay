# 05 — Reference Implementations

The leader writes 1–2 reference implementations before handing off to implementer models. These are not throwaway prototypes; they are the canonical example of how the style contract applies to real code in this project.

## Why reference code matters

- It proves the `design.md` and `style.md` are actually achievable.
- It gives implementer models a concrete local example to imitate rather than abstract rules.
- It catches mismatches between the feature spec and the project layout early.

## What a good reference includes

1. **Happy path** — the normal flow from request to response.
2. **Error path** — at least one validation failure, one authorization failure, and one downstream error.
3. **Authorization check** — the exact call and what happens on denial.
4. **Validation** — input validated with generated types or an explicit validator, not ad-hoc string checks.
5. **Generated types usage** — structs from sqlc, gogen, or Protobuf, not hand-rolled duplicates.
6. **Tests** — one happy-path test and one error-path test.
7. **Docstring/comment** — what the file does, the authorization rule, and the expected caller.

## Where to save them

In the target project, create a `reference/` or `docs/reference/` directory:

```text
docs/reference/
├── backend/
│   └── create_chat_message_handler.go
├── frontend/
│   └── CreateChatMessageForm.tsx
└── db/
    └── create_chat_message.sql
```

Reference code can be removed or relocated after the first real implementation lands, but keep it intact until review passes.

Skip reference implementations for simple features. A one-endpoint CRUD handler doesn't need a reference; the style contract is enough. Write references for features that introduce a new pattern, a new service boundary, or non-trivial authorization logic. Reference implementations are themselves subject to the ladder — don't write a reference that's more complex than the features it demonstrates.

## Review checkpoint

Before handing off, ask:

- Does the reference pass the style contract from `style.md`?
- Does it import generated types from the correct directories?
- Does every test fail if authorization is removed?
