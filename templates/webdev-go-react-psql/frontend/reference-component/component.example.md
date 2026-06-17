# Reference Component Example

Below is a complete `CreateChatMessageForm` component plus a custom hook that consumes generated types and handles loading/error states.

```tsx
// src/components/features/CreateChatMessageForm.tsx
import { useState, FormEvent } from 'react';
import { useCreateChatMessage } from '../../hooks/useCreateChatMessage';

interface Props {
  projectId: string;
}

export function CreateChatMessageForm({ projectId }: Props) {
  const [text, setText] = useState('');
  const { mutate, isPending, error } = useCreateChatMessage(projectId);

  const onSubmit = (e: FormEvent) => {
    e.preventDefault();
    if (!text.trim()) return;
    mutate({ text: text.trim() }, {
      onSuccess: () => setText(''),
    });
  };

  const errorMessage = error instanceof Error ? error.message : null;

  return (
    <form onSubmit={onSubmit}>
      <textarea
        value={text}
        onChange={(e) => setText(e.target.value)}
        maxLength={4000}
        placeholder="Type a message..."
        aria-label="Message text"
      />
      <button type="submit" disabled={isPending || !text.trim()}>
        {isPending ? 'Sending...' : 'Send'}
      </button>
      {errorMessage && (
        <p role="alert" style={{ color: 'red' }}>
          {errorMessage}
        </p>
      )}
    </form>
  );
}
```

```ts
// src/hooks/useCreateChatMessage.ts
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { apiClient } from '../api/client';
import type { CreateChatMessageResponse } from '../gen/models';

interface Input {
  text: string;
}

export function useCreateChatMessage(projectId: string) {
  const qc = useQueryClient();

  return useMutation<CreateChatMessageResponse, Error, Input>({
    mutationFn: async (input) => {
      const { data } = await apiClient.post<CreateChatMessageResponse>(
        `/api/projects/${encodeURIComponent(projectId)}/chat/messages`,
        input,
      );
      return data;
    },
    onError: (err) => {
      // Centralized error logging can go here.
      // eslint-disable-next-line no-console
      console.error('create chat message failed', err);
    },
    onSuccess: () => {
      // Optimistic or invalidation strategy.
      void qc.invalidateQueries({ queryKey: ['projects', projectId, 'messages'] });
    },
  });
}
```

## Hook design notes

- The hook owns the API contract shape (`Input`) and imports the generated response type.
- Errors are always surfaced as `Error` instances to the component.
- The component never calls `fetch` directly.
- The submit button disables while pending; the form does not reset on error.

## Auth integration

Assume `apiClient` injects the JWT header from the auth provider. The route guard ensures this component is rendered only for authenticated users. If the API returns `401`, the global error boundary or query client redirects to login.

## Tests

```ts
// tests/CreateChatMessageForm.test.tsx
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { CreateChatMessageForm } from '../src/components/features/CreateChatMessageForm';

test('sends a message', async () => {
  // Mock API with MSW.
});

test('disables submit when text is empty', async () => {
  // Assert UX behavior.
});
```
