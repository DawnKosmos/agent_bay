# carryover Reference Code Map

This directory does not contain copied code. It points to the actual reference code in the parent project, which demonstrates how the `agent_bay` patterns are applied.

```text
../../
├── apps/carryover/ws_api/ws_api.go
│   └── Defines the generic WS envelope `Msg[T]`, event type constants,
│       and concrete payloads `ChatMessage` and `Notification`.
│       This is the source of truth for `libs/gogen`.
│
├── apps/carryover/ws/ws.go
│   └── Contains the transport-level `Msg[T]` envelope.
│
├── apps/carryover/ws/publish.go
│   └── Shows how to publish to a user-scoped Centrifuge channel.
│       Demonstrates explicit error handling and JSON marshaling.
│
├── libs/gogen/gogen.go
│   └── The custom Go struct → TypeScript generator.
│       Review this to understand how `ws_api` types become frontend types.
│
├── libs/proto/levi/levi.proto
│   └── Protobuf contract between carryover and Levi.
│
└── apps/carryover/wal_handler/
    └── Contains handlers for WAL/CDC style events.
```

## What each piece demonstrates

- `ws_api/` — Define shared contracts in plain Go structs.
- `ws/publish.go` — Keep publishing logic isolated; services can trigger events without knowing Centrifuge details.
- `libs/gogen` — Generate TypeScript so the frontend never drifts from the backend model.
- `libs/proto/levi/levi.proto` — gRPC contract for service boundaries.

As carryover grows, add new entries here for canonical handlers, components, and migrations.
