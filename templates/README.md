# Templates

Templates are stack-specific starter recipes. They translate the generic `workflow/` into concrete conventions, file paths, and tooling commands for one technology stack.

## Available templates

- `webdev-go-react-psql/` — Default web stack. Go backend + React/TypeScript frontend + PostgreSQL + Centrifuge WebSockets.

## How to copy and rename a template

1. Copy the template folder into the target project under `docs/templates/` or alongside `docs/`:

   ```bash
   cp -r agent_bay/templates/webdev-go-react-psql myproject/docs/stack-template
   ```

2. Rename the files if needed, but keep the section structure.
3. Fill in `design.template.md` with project-specific actors, modules, and security model.
4. Fill in `style.template.md` with project-specific package layout and generated-code paths.
5. Use `feature.template.md` for every feature.

## How to add a new template

1. Create `templates/<stack-name>/`.
2. Add a `README.md` explaining when to use the stack and the quickstart.
3. Provide `design.template.md`, `style.template.md`, and `feature.template.md` tailored to the stack.
4. Add backend/frontend/DB/WS/codegen reference sections that match the stack.
5. Update the list in this file.
