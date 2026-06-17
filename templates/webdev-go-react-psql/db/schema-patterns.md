# Schema Patterns

Postgres conventions for the `webdev-go-react-psql` template.

## Naming

- Tables are plural nouns: `users`, `projects`, `chat_messages`.
- Primary key column is `id` with type `uuid` and default `gen_random_uuid()`.
- Foreign key columns use the referenced table name: `project_id`, `user_id`.
- Timestamp columns: `created_at timestamptz NOT NULL DEFAULT now()`, `updated_at timestamptz NOT NULL DEFAULT now()`.
- Soft deletes use a nullable `deleted_at timestamptz` when required. Hard deletes are explicit and logged.

## Example migration

```sql
-- db/migration/0001_create_projects.up.sql
CREATE TABLE projects (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id uuid NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    name text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_projects_owner_id ON projects(owner_id);
```

```sql
-- db/migration/0001_create_projects.down.sql
DROP TABLE IF EXISTS projects;
```

## Constraints

- Prefer `ON DELETE RESTRICT` for relationships that should not disappear silently.
- Use `CHECK` constraints for simple domain rules (e.g., `CHECK (char_length(name) > 0)`).
- Enforce uniqueness with explicit named constraints.

## sqlc usage

- Place `.sql` query files in `db/query/`.
- Every query lists explicit columns; do not use `SELECT *` in production queries.
- Tag sqlc query return types via `-- name: CreateProject :one` etc.
- Migrations and queries are reviewed as code.

## Indexing

- Index every foreign key by default unless the table is tiny (<1000 rows) and proven otherwise.
- Add partial indexes for common filtered queries.
- Add composite indexes based on slow-query logs, not guessing.

## Migrations

- Use `sql-migrate`, `golang-migrate`, or `dbmate`.
- Every migration file has an `up` and a `down`.
- Keep migrations reversible when possible. Irreversible migrations require an approval note in the PR.
