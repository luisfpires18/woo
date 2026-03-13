---
description: "Use when editing game config files (buildings, troops, resources), running codegen, or working with generated JSON. Covers the Go → JSON → TypeScript pipeline."
applyTo:
  - "server/internal/config/*.go"
  - "client/src/config/**"
---

# Config Codegen Pipeline

## Source of Truth

Go config files in `server/internal/config/` are the **single source of truth**:
- `buildings.go` — building configs (costs, scaling, prerequisites)
- `troops.go` — troop configs (stats, costs, kingdoms)
- `resources.go` — resource economy constants (starting values, rates, storage)

## Pipeline

```
Go config (edit here) → npm run gen-config → generated JSON → TS wrappers (import JSON)
```

1. Edit Go config in `server/internal/config/*.go`
2. Run codegen: `npm run gen-config` (repo root) or `cd server && go run ./cmd/genconfig`
3. Generated JSON lands in `client/src/config/generated/` (buildings.json, troops.json, resources.json)
4. TS wrappers in `client/src/config/` import from generated JSON

## Rules

- **NEVER edit files in `client/src/config/generated/`** — always edit Go source and re-run codegen.
- **Never duplicate values** between Go and TS. TS imports from generated JSON.
- **After any Go config change**: run `npm run gen-config` and commit updated JSON.
- **Parity test**: `go test ./internal/config/` verifies committed JSON matches Go config.
- **DTO types** in `server/internal/config/generated_types.go` define the shared JSON schema.

## Verification Checklist

After changing any config:
1. `cd server && go run ./cmd/genconfig` — regenerate JSON
2. `cd server && go test ./internal/config/` — verify parity
3. Check the diff — only expected fields should change
4. Commit both Go source and generated JSON together
