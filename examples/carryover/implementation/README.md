# carryover Implementation Plan

This file tracks known and planned features. Each feature should get its own `feature/<name>.md` document before implementation begins.

## Phase 1: Identity & membership

1. User registration/login via GitHub OAuth.
2. JWT issuance and refresh.
3. User profile retrieval.

## Phase 2: Projects

4. Create project.
5. Add/remove project members.
6. List projects for the current user.

## Phase 3: Issues & comments

7. Create issue.
8. Update issue status / priority / assignee.
9. Add comment with mentions.
10. List issues in a project.

## Phase 4: Real-time

11. Real-time notification on mention or assignment (via Centrifuge `personal.{userID}`).
12. Real-time chat or comment thread updates if required by UX.

## Phase 5: Search & AI

13. Typesense indexing for issues and comments.
14. Levi gRPC methods for querying project state.
15. Levi-assisted command routing.

## Notes

- Feature docs should reference the authorization matrix once it is defined.
- Each feature doc must include the API contract, DB queries, WS events (if any), authz checks, errors, and tests.
