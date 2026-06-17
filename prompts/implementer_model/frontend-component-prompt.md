# Implementer Prompt: Frontend Component

Copy-paste the context into an implementer model chat. The output must be a React/TypeScript component plus hook, matching the project's design, style, and feature spec.

---

You are implementing a frontend feature. I will give you:

- The feature spec (`feature/<name>.md`).
- The project `style.md`.
- A reference component/hook example.
- The generated TypeScript types in `frontend/src/gen/models.ts`.

## Your task

Write a React component file and a hook file. Follow the reference component example exactly.

## Required sections in your output

1. **Hook file**
   - Typed mutation or query using the generated response type.
   - Error typed as `Error`.
   - API endpoint built from props.
   - Invalidation or optimistic update strategy.

2. **Component file**
   - Props interface.
   - Local state for form inputs.
   - Loading state disables submit.
   - Error state rendered with `role="alert"`.
   - No direct `fetch` calls.

## Constraints

- Strict TypeScript; no implicit `any`.
- Import generated types, do not redefine them.
- UI text matches the feature spec.

## Forbidden shortcuts

- Do not put data fetching inside the component.
- Do not skip error display.
- Do not use `any` for API responses.
- Do not guess behavior if the spec is ambiguous; stop and ask.

Return only the file contents inside markdown code blocks with file paths as headers.
