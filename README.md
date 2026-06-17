# agent_bay

`agent_bay` is a shared playbook for building software with a human plus two classes of AI models: a strong **leader** model that designs, and cheaper **implementer** models that build. It documents the exact workflow, stack-specific templates, reusable prompts, and real project examples so that every feature ships with explicit error handling, type safety, authorization, tests, and generated service boundaries.

## What lives where

- `workflow/` — The end-to-end manufacturer flow, from discovery and tech-stack selection through design, style contracts, feature specs, reference implementations, codegen, implementer handoff, and review.
- `templates/` — Stack-specific starter recipes. The default is `webdev-go-react-psql/` (Go + Fiber or standard HTTP, React/TypeScript, PostgreSQL, Centrifuge WebSockets). Add a new folder here when a project needs a different stack.
- `prompts/` — Copy-pasteable prompts for leader and implementer models. leader prompts produce `design.md`, `style.md`, feature specs, and scaffold setup. implementer prompts produce backend handlers, frontend components, and DB migrations.
- `examples/` — Real projects documented with this workflow. `examples/carryover/` is the first one.

## Using this repo as a submodule

`agent_bay` is intended to live as a git submodule inside a parent project. In the parent project root:

```bash
git submodule add <agent_bay-remote-url> agent_bay
git submodule update --init --recursive
```

Keep the submodule pointer at whatever version the parent project currently uses. Do not modify this submodule inside the parent project's commits; evolve `agent_bay` in its own repository and update the submodule pointer explicitly.

When you reference playbook files from the parent project, use relative paths, e.g.:

```text
agent_bay/templates/webdev-go-react-psql/design.template.md
agent_bay/prompts/leader_model/01-design-prompt.md
```

## Adding a new project stack template

1. Copy `templates/webdev-go-react-psql/` to `templates/<your-stack-name>/`.
2. Update the README with the stack's rationale and quickstart.
3. Rewrite the design, style, and feature templates for the new stack's rules, conventions, and generated-code locations.
4. Add or replace the backend/frontend/DB/WS/codegen reference sections.
5. Add stack-specific prompts in `prompts/` if they differ enough from the default set.
6. Update `templates/README.md` with a link to the new template.
