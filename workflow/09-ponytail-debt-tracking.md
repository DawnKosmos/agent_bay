# 09 — Ponytail Debt Tracking

Every deliberate shortcut the Ponytail mindset leaves behind is marked with a `ponytail:` comment. This collects them into a ledger so a deferral can't quietly become permanent.

## The `ponytail:` comment convention

When you intentionally simplify something — skip an abstraction, use a naive algorithm, defer a feature — mark it:

```go
// ponytail: global lock, per-account locks if throughput matters
```

```typescript
// ponytail: O(n²) scan, switch to indexed lookup when list > 1000
```

The format is: `ponytail: <ceiling>, <upgrade path>`

- **Ceiling** — the known limit of the shortcut (global lock, O(n²), naive heuristic, no pagination).
- **Upgrade path** — the trigger or condition that should prompt revisiting (when throughput matters, when list > 1000, when a second implementation appears).

A `ponytail:` comment without an upgrade path is a rot risk — it names a shortcut but gives no trigger to revisit it.

## Harvesting the ledger

To collect all `ponytail:` markers into a debt report:

```bash
grep -rnE '(#|//) ?ponytail:' .
```

Each hit is one ledger row. Group by file:

```
<file>:<line>, <what was simplified>. ceiling: <the limit named>. upgrade: <the trigger to revisit>.
```

Flag any marker with no upgrade path:

```
<file>:<line>, <what was simplified>. ceiling: <the limit named>. upgrade: NONE — rot risk.
```

End with: `<N> markers, <M> with no trigger.`

If nothing is found: `No ponytail: debt. Clean ledger.`

## When to revisit

The upgrade path in each comment is the trigger. When that condition is met — throughput matters, the list grows past the threshold, a second implementation appears — the shortcut is no longer acceptable and the upgrade should be implemented.

## Scope

This is a read-only audit. It changes nothing. To persist the ledger, write it to a file (e.g., `PONYTAIL-DEBT.md`) and commit it.
