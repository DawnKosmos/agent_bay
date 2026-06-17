# Implementer Prompt: Backend Handler

Copy-paste the context into an implementer model chat. The output must be a Go HTTP handler file plus tests, matching the project's design, style, and feature spec.

---

You are implementing a backend HTTP handler. I will give you:

- The feature spec (`feature/<name>.md`).
- The project `style.md`.
- A reference handler example showing the pattern.
- The generated types available (`db/query/`, `ws_api/`).

## Your task

Write the Go handler file and its test file. Follow the reference handler pattern exactly.

## Required sections in your output

1. **Handler file**
   - Request struct with JSON tags.
   - Response struct or use generated type.
   - Handler method.
   - Actor extraction from context.
   - Input validation (explicit, no panics).
   - Authorization check using the existing `service.Authorizer`.
   - Service call; explicit error mapping.
   - JSON response.

2. **Test file**
   - Table-driven happy path test.
   - Table-driven authorization failure test.
   - Table-driven validation error test.
   - Any required mocks/fakes.

## Constraints

- No `panic` in request paths.
- No direct DB access from the handler.
- Use generated types where possible.
- Every error path has a test.

## Forbidden shortcuts

- Do not skip the authorization call.
- Do not use `interface{}` where a generated type exists.
- Do not swallow errors.
- Do not guess behavior if the spec is ambiguous; stop and ask.

Return only the file contents inside markdown code blocks with file paths as headers.
