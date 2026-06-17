# Example: carryover

This directory documents `carryover`, the first real project built using the `agent_bay` workflow. `carryover` is a self-hosted Jira-style issue tracker with real-time features and an AI assistant named Levi.

## Parent project layout

The actual code lives up two directories from here, at the root of the `carryover` repository. Paths below are relative to this file (`agent_bay/examples/carryover/`):

```text
../../
├── README.md
├── Techstack.md
├── apps/
│   ├── carryover/      # Main monolith application (Go backend + frontend served from here)
│   └── levi/           # AI assistant service
├── libs/
│   ├── gogen/          # Go struct → TypeScript generator
│   ├── proto/          # Protobuf contracts (e.g., levi)
│   ├── cdc/            # Change-data-capture helpers
│   ├── uuid/           # UUID utilities
│   ├── ts_util/        # Typesense utilities
│   └── wrap/           # Framework wrappers
└── frontend/           # React/TypeScript frontend
```

## What this example covers

- `design.md` — draft design doc with carryover-specific placeholders.
- `style.md` — carryover-specific style rules referencing `libs/gogen`, sqlc, and the WS packages.
- `implementation/README.md` — planned features to spec out.
- `reference-code/README.md` — map to actual reference code already in the repo.

As `carryover` evolves, this example should stay in sync and grow into a full project case study.
