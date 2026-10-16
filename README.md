# Weapons of Order

Updated: 2026-10-05. Status: design documentation; no game implementation yet.

A seasonal browser strategy game set in Bellum. The player is a blacksmith who controls village development and troops, serving one kingdom. Forging should shape military decisions.

The first playable version is simple and medieval. Runes and advanced fantasy systems belong to the evolving game, with release timing still open. The closed alpha is for the owner and friends, with Arkazia, Veridor and Sylvara playable.

Start with [the documentation index](docs/INDEX.md), [first-version scope](docs/scope.md), [open questions](docs/open-questions.md) and [development roadmap](docs/roadmap.md). Coding agents must read [AGENTS.md](AGENTS.md).

This repository is now the primary game-design record. Update the relevant documents and decision log as discussions evolve. Documentation authorisation does not authorise starting game implementation.

## Development orchestration

The working integration branch is **dev**. Read [tasks_roadmap.md](tasks_roadmap.md), [workflow](docs/workflow/development-workflow.md) and [prompt protocol](docs/workflow/implementation-protocol.md). Claude Code implements one reviewed task at a time; ChatGPT prepares prompts and reviews; the owner runs merge/Azure commands. Azure is the selected host direction. Production deployment from master is deferred until alpha release.

GitHub's repository default branch may still be master; explicitly select dev to read the latest workflow. Changing the default branch is an outstanding repository setting, not claimed complete.

## Run locally

Status: tasks 004 SKELETON and 005 VERSION: a minimal page, `/api/health`, `/api/build` and a version footer; no game features. Versions and layout follow [the stack record](docs/technical/stack.md).

Prerequisites: .NET SDK 10.0.402 and Node.js 24.22.0 (npm 11.19.1) on `PATH`. `global.json` refuses any other SDK and `web/.npmrc` enforces the Node version. A side-by-side install can be put first on `PATH` for one Git Bash session, for example:

```bash
export PATH="/f/Workspace/Tools/dotnet-sdk-10.0.402:/f/Workspace/Tools/node-v24.22.0-win-x64:$PATH" DOTNET_ROOT='F:\Workspace\Tools\dotnet-sdk-10.0.402'
```

Commands below run from the repository root in Git Bash on Windows.

Restore exactly what the lock files record (clean checkout):

```bash
npm --prefix web ci
```

```bash
dotnet restore Woo.slnx --locked-mode
```

Build and test:

```bash
npm --prefix web run build
```

```bash
dotnet build Woo.slnx --no-restore
```

```bash
dotnet test --solution Woo.slnx --no-build
```

Develop with hot reload: start the API (http://localhost:5080) in one terminal and the Vite dev server (http://localhost:5173, proxies `/api` to the API) in another, then open http://localhost:5173.

```bash
dotnet run --project src/Woo.Server --launch-profile http
```

```bash
npm --prefix web run dev
```

Standalone builds and hot reload report the local fallback: `/api/build` returns version `0.0.1-dev`, environment `local` and empty commit/build time, and the footer reads `0.0.1-dev (local)`.

Coordinated build and publish (the single-origin app, React build served by ASP.NET Core). It requires committed tracked files, reads `VERSION`, captures the full commit and UTC build time once and passes the same values to the web build (`artifacts/web`) and the server publish (`artifacts/publish`). The footer then shows the version and short commit. `--environment` accepts `local` (default), `dev` or `production`.

```bash
node scripts/build.mjs
```

A standalone publish (local fallback metadata) needs `npm --prefix web run build` first and fails if `web/dist` is missing:

```bash
dotnet publish src/Woo.Server -c Release -o artifacts/publish
```

Run the published output from its own folder (it is the content root that holds `wwwroot`) and open http://localhost:5090; `curl http://localhost:5090/api/build` shows the build identity:

```bash
cd artifacts/publish && ASPNETCORE_URLS=http://localhost:5090 dotnet Woo.Server.dll
```

Generated output (`bin/`, `obj/`, `web/node_modules/`, `web/dist/`, `artifacts/`) and local runtime files are ignored by Git.
