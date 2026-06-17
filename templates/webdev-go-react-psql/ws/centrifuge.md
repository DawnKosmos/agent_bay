# Centrifuge WebSockets

This stack uses [Centrifuge](https://centrifugal.dev/) for live updates. The goal is a type-safe pipeline from Go events to TypeScript handlers.

## Recommended architecture

### 1. Define a shared envelope in Go

Keep the generic envelope in `ws_api/` so `libs/gogen` can mirror it if needed:

```go
package ws_api

// Msg is the generic WebSocket envelope.
type Msg[T any] struct {
	Type string `json:"type"`
	Data T      `json:"data"`
}
```

### 2. Define concrete payloads in `ws_api/`

```go
package ws_api

const (
	ChatMessageReceived = "message.received"
	NotificationCreated = "notification.created"
)

type ChatMessage struct {
	MessageID string `json:"messageId"`
	FromID    string `json:"fromId"`
	Text      string `json:"text"`
	CreatedAt string `json:"createdAt"`
}

type Notification struct {
	NotificationID string `json:"notificationId"`
	Title          string `json:"title"`
	Body           string `json:"body"`
	CreatedAt      string `json:"createdAt"`
}
```

### 3. Backend publishing

The `ws/` package owns channel publishing:

```go
package ws

import (
	"encoding/json"
	"fmt"

	"github.com/DawnKosmos/carryover/apps/carryover/ws_api"
	"github.com/centrifugal/centrifuge"
)

func PersonalChannel(userID string) string { return ws_api.PersonalPrefix + userID }

func PublishPersonalChat(node *centrifuge.Node, userID string, msg ws_api.ChatMessage) error {
	payload, err := json.Marshal(ws_api.Msg[ws_api.ChatMessage]{
		Type: ws_api.ChatMessageReceived,
		Data: msg,
	})
	if err != nil {
		return fmt.Errorf("marshal chat message: %w", err)
	}
	_, err = node.Publish(PersonalChannel(userID), payload)
	return err
}
```

- `ws_api/` owns the contract.
- `ws/` owns the transport and publishers.
- Services call `ws.Publish*`; handlers/services never import Centrifuge directly except via the `ws` package.

### 4. Generate TypeScript interfaces

Run `libs/gogen` against `ws_api/` to produce `frontend/src/gen/models.ts`. Then the frontend imports `ChatMessage`, `Notification`, and the generated `Msg<T>` envelope.

### 5. Frontend dispatch

Use a type-safe dispatch map keyed by event type:

```ts
// src/ws/dispatch.ts
import type { ChatMessage, Notification } from '../gen/models';

type WsEventType =
  | 'message.received'
  | 'notification.created';

interface WsEventData {
  'message.received': ChatMessage;
  'notification.created': Notification;
}

type WsDispatcher = {
  [K in WsEventType]: (data: WsEventData[K]) => void;
};

export function createDispatcher(handlers: Partial<WsDispatcher>): (type: string, data: unknown) => void {
  return (type, data) => {
    const handler = handlers[type as WsEventType];
    if (handler) {
      // At runtime this should already be validated; in dev, assert shape.
      handler(data as never);
    }
  };
}
```

Subscription lifecycle:

- Subscribe on component mount.
- Unsubscribe on unmount.
- Re-subscribe on token refresh.

### 6. Authorization

- Channels are scoped to the authenticated user: `personal.{userID}`.
- The Centrifuge token endpoint verifies the caller's JWT and signs a subscription token only for the requested user's channel.
- Never allow a user to subscribe to another user's `personal` channel.
- For project channels, validate membership before returning a token.

## Anti-patterns

- Do not send WS events to clients that have not subscribed to the channel.
- Do not put business logic inside `ws/` publishers.
- Do not hand-maintain TypeScript copies of `ws_api` structs.
