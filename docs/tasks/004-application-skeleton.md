# 004: Create skeleton and run the first page

Task key: SKELETON
Status: DONE
Updated: 2026-10-15
Milestone: Foundation
Dependencies: STACK

## Goal

Create the agreed solution and frontend/backend startup with one minimal page.

## Acceptance criteria

A clean checkout runs with documented commands; browser displays the page and reaches a basic backend health endpoint. No account, DB, game content or polished marketing page required.

## Required reading

- [docs/technical/architecture.md](../../docs/technical/architecture.md)
- [Scope](../scope.md)
- [Prompt protocol](../workflow/implementation-protocol.md)

## Implementation boundary

A separately reviewable planning task. ChatGPT inspects current code and resolves this task's blockers before issuing one pinned Claude prompt. Preserve medieval v1. Split/refine based on actual development; do not implement successors automatically.

## Validation and administration

Run proportionate checks against acceptance criteria, including relevant failure/retry scenarios. Report actual local evidence. ChatGPT checks remote merge/deployment where applicable. Content/images follow admin configuration and cleanup rules; new mechanics remain code changes.

## Result

Minimal runnable slice, following the approved [stack](../technical/stack.md):

- `Woo.slnx`, `global.json` (SDK 10.0.402, roll-forward disabled, Microsoft Testing Platform runner), `Directory.Build.props` (net10.0, nullable, NuGet lock files), `Directory.Packages.props` (central exact versions).
- `src/Woo.Server`: `GET /api/health` returns `{"status":"ok"}`; any other `/api/...` path returns a 404 problem response, never the SPA page; published builds serve the React build from `wwwroot` with an `index.html` fallback for client routes. Publishing fails with a clear message if `web/dist` has not been built.
- `tests/Woo.Server.Tests`: host tests for health, unknown API routes and client-route fallback.
- `web/`: React 19.3.1 + Vite 8.3.4 + TypeScript 6.0.3 page that fetches `/api/health` and shows loading, reachable and unavailable states (with retry and a 5 second timeout). Vite proxies `/api` to `http://localhost:5080` in development.
- No VERSION file, build metadata, footer, database, accounts, game content, PixiJS, CI or deployment (later tasks).

Commands are in the [README](../../README.md#run-locally).

## Toolchain

On the owner's second computer (2026-10-15), the pinned toolchain was installed side-by-side, user-local, without administrator rights or replacing existing versions:

- .NET SDK 10.0.402 (runtime 10.0.13) from the official zip `https://builds.dotnet.microsoft.com/dotnet/Sdk/10.0.402/dotnet-sdk-10.0.402-win-x64.zip`, SHA-512 matching the 10.0 `releases.json`, extracted to `F:\Workspace\Tools\dotnet-sdk-10.0.402`.
- Node.js 24.22.0 (npm 11.19.1) from `https://nodejs.org/dist/v24.22.0/node-v24.22.0-win-x64.zip`, SHA-256 matching `SHASUMS256.txt`, extracted to `F:\Workspace\Tools\node-v24.22.0-win-x64`.
- They are used by putting those folders first on `PATH` (and `DOTNET_ROOT`) per shell; no system `PATH` or registry change was made. The existing SDKs (8.0.413, 9.0.307, 10.0.100) and Node 22.18.0 remain installed. Drive C: had about 3 GB free, so F: was used.

## Execution record

Task-spec commit: 588c5a09ab7ecce6aa596d79e63579345897c3ea (issue [#3](https://github.com/luisfpires18/woo/issues/3) dispatch)
Expected dev base: 588c5a09ab7ecce6aa596d79e63579345897c3ea (verified equal to local dev and origin/dev, 2026-10-15)
Implementation branch: feat/004-application-skeleton
Implementation commit: 77e206cd47f5a444485ac0dfb3aba194ad059d7f
Review: APPROVED by ChatGPT, 2026-10-15
Merged dev commit: 77e206cd47f5a444485ac0dfb3aba194ad059d7f (remote dev; documentation CI succeeded)
Deployment: Not applicable to this task (deployment is DEV-DEPLOY)
Issues: browser checks covered the Vite proxy (HTTP 502) failure path; the timeout and network-error messages were not exercised in a browser.

## TLDR

Next: none; VERSION (005) builds on this skeleton.
Done: approved and merged to dev at 77e206c with documentation CI passing.
Issues: see execution record.
