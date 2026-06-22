# Implementer Prompt: DB Migration

Copy-paste the context into an implementer model chat. The output must be a reversible PostgreSQL migration and the sqlc queries for the feature.

---

You are implementing a database change. I will give you:

- The feature spec (`feature/<name>.md`).
- The project `style.md`.
- The existing schema patterns (`db/schema-patterns.md`).
- The current `design.md` data model overview.

## Ponytail pre-flight

Before writing the migration, ask:
- Does this table need to exist? Can an existing table hold this data with a new column?
- Fewest tables possible. Fewest columns possible.
- Only write the migration if the data doesn't fit anywhere existing.

## Your task

Write a migration pair (up + down) and any new sqlc queries.

## Required sections in your output

1. **Migration up file**
   - `CREATE TABLE` or `ALTER TABLE`.
   - UUID primary keys, `created_at`/`updated_at`, foreign keys with explicit actions.
   - Indexes on foreign keys and any lookup columns.
   - Named constraints.

2. **Migration down file**
   - Reverse the up migration.
   - Avoid data loss where possible.

3. **sqlc queries**
   - One `.sql` file with `-- name:` annotations.
   - Use `:one`, `:many`, `:exec` as appropriate.
   - Explicit column lists; no `SELECT *`.
   - Parameter and return types inferable by sqlc.

## Constraints

- All migrations reversible unless explicitly noted as irreversible.
- Foreign keys use appropriate `ON DELETE` actions.
- Use `uuid` for IDs and `timestamptz` for timestamps.
- No speculative columns, no nullable fields "for later." Add columns when they're needed.
- Mark simplifications with `ponytail:` comments.

## Forbidden shortcuts

- No unindexed foreign keys.
- No raw string concatenation in queries (use sqlc parameters).
- No migration without a down file.
- Do not guess column types or sizes if the spec is ambiguous; stop and ask.

## Output discipline

SQL first. At most 3 lines of explanation. No prose defending the schema.

Return only the file contents inside markdown code blocks with file paths as headers.
