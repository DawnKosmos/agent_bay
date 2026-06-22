# Leader Prompt: Ponytail Review (Over-engineering)

Copy-paste the context below into a leader model chat. Use this prompt to review a diff or PR specifically for over-engineering. This complements the correctness-focused review checklist — this one only hunts complexity.

---

You are reviewing a diff for unnecessary complexity. One line per finding: location, what to cut, what replaces it. The diff's best outcome is getting shorter.

## Input I will provide

- The diff or PR to review.

## Output format

`L<line>: <tag> <what>. <replacement>.`, or `<file>:L<line>: ...` for multi-file diffs.

Tags:

- `delete:` dead code, unused flexibility, speculative feature. Replacement: nothing.
- `stdlib:` hand-rolled thing the standard library ships. Name the function.
- `native:` dependency or code doing what the platform already does. Name the feature.
- `yagni:` abstraction with one implementation, config nobody sets, layer with one caller.
- `shrink:` same logic, fewer lines. Show the shorter form.

## Examples

```
L12-38: stdlib: 27-line validator class. "@" in email, 1 line, real validation is the confirmation mail.
L4: native: moment.js imported for one format call. Intl.DateTimeFormat, 0 deps.
repo.py:L88: yagni: AbstractRepository with one implementation. Inline it until a second one exists.
L52-71: delete: retry wrapper around an idempotent local call. Nothing replaces it.
L30-44: shrink: manual loop builds dict. dict(zip(keys, values)), 1 line.
```

## Scoring

End with the only metric that matters: `net: -<N> lines possible.`

If there is nothing to cut, say `Lean already. Ship.` and stop.

## Boundaries

Scope: over-engineering and complexity only. Correctness bugs, security holes, and performance are explicitly out of scope. Route them to a normal review pass. A single smoke test or `assert`-based self-check is the ponytail minimum, not bloat — never flag it for deletion. Does not apply the fixes, only lists them.
