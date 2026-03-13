---
description: "Run the config codegen pipeline: regenerate JSON from Go config, verify parity, and check the diff."
agent: "agent"
tools: [execute, read, search]
---

# Run Config Codegen

Execute the full config codegen pipeline and verify everything is in sync.

## Steps

1. **Kill running server** (if any):
   ```
   Get-Process -Name "server" -ErrorAction SilentlyContinue | Stop-Process -Force
   ```

2. **Run codegen**:
   ```
   cd d:\Workspace\WOO\server; go run ./cmd/genconfig
   ```

3. **Run parity tests**:
   ```
   cd d:\Workspace\WOO\server; go test ./internal/config/ -run Parity -v
   ```

4. **Check what changed** in the generated files:
   - `client/src/config/generated/buildings.json`
   - `client/src/config/generated/troops.json`
   - `client/src/config/generated/resources.json`

5. **Report** which fields changed and whether parity tests pass.

## When to Use

- After editing any file in `server/internal/config/*.go`
- Before committing config changes
- When frontend shows stale/wrong game values
