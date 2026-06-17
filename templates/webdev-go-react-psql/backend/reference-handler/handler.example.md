# Reference Handler Example

Below is a complete, production-style Go HTTP handler for a hypothetical `POST /api/projects/{project_id}/chat/messages` endpoint. The comments explain the purpose of each block.

```go
package handlerhttp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"          // or your preferred router
	"github.com/yourorg/yourapp/db/query"
	"github.com/yourorg/yourapp/service"
)

// handler holds service dependencies. It is constructed once in main.
type handler struct {
	chat   service.Chat
	authz  service.Authorizer
}

// CreateChatMessageRequest is the decoded request. Untagged fields are ignored by json.
// Validation lives in the service so it can be reused by other callers.
type CreateChatMessageRequest struct {
	Text string `json:"text"`
}

// CreateChatMessageResponse mirrors the generated db/query.ChatMessage shape,
// but we define a handler-specific DTO if we need to hide or rename fields.
type CreateChatMessageResponse struct {
	MessageID string `json:"messageId"`
	FromID    string `json:"fromId"`
	Text      string `json:"text"`
	CreatedAt string `json:"createdAt"`
}

// CreateChatMessage handles POST /api/projects/{project_id}/chat/messages.
// Authorization: the caller must be a member of the project.
func (h *handler) CreateChatMessage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	projectID := chi.URLParam(r, "project_id")

	// 1. Identify actor. Extracted by auth middleware and placed in context.
	actor, ok := service.ActorFromContext(ctx)
	if !ok {
		respondError(w, http.StatusUnauthorized, "missing actor")
		return
	}

	// 2. Decode and validate request.
	var req CreateChatMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, fmt.Sprintf("invalid request: %v", err))
		return
	}

	if req.Text == "" || len(req.Text) > 4000 {
		respondError(w, http.StatusBadRequest, "text must be 1-4000 characters")
		return
	}

	// 3. Authorize. This is the only allowed path to membership checking.
	if err := h.authz.CanPostInProject(ctx, actor.UserID, projectID); err != nil {
		if errors.Is(err, service.ErrForbidden) {
			// Return 404 instead of 403 when the actor is not a member,
			// to avoid leaking project existence if policy requires.
			respondError(w, http.StatusNotFound, "project not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "authorization failed")
		return
	}

	// 4. Call the service. The service owns business logic and emits WS events.
	msg, err := h.chat.CreateMessage(ctx, service.CreateMessageInput{
		ProjectID: projectID,
		FromID:    actor.UserID,
		Text:      req.Text,
	})
	if err != nil {
		var derr *service.DomainError
		if errors.As(err, &derr) {
			respondError(w, derr.HTTPStatus, derr.Message)
			return
		}
		// Log the raw error once, return a generic message.
		respondError(w, http.StatusInternalServerError, "failed to create message")
		return
	}

	// 5. Encode the response using the generated/model shape.
	resp := CreateChatMessageResponse{
		MessageID: msg.MessageID,
		FromID:    msg.FromID,
		Text:      msg.Text,
		CreatedAt: msg.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(resp)
}

func respondError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
```

## Test expectations

```go
func TestCreateChatMessage(t *testing.T) {
    // Happy path: member posts a valid message -> 201.
    // Error path 1: non-member posts -> 404/403.
    // Error path 2: empty text -> 400.
    // Error path 3: service failure -> 500 generic response.
}
```

## What this example enforces

- Actor comes from context, never from the request body.
- Validation runs before authorization when safe, or authorization runs first when resource existence is sensitive — document the order in the feature spec.
- All errors are handled explicitly; nothing is swallowed.
- The handler does not call the DB directly.
