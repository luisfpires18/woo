# Tech stack, versions and architecture

Updated: 2026-10-15. Status: task 003 STACK decision record, APPROVED and merged to dev (588c5a0) on 2026-10-15. Task 004 SKELETON applies it; see [skeleton additions](#selected-by-skeleton-task-004).

C# / ASP.NET Core and React/TypeScript are Confirmed (D-37). The versions, layout and boundaries below are the approved STACK selection (D-39). Rationale for the backend and host stays in [backend comparison](backend-comparison.md) and [hosting comparison](hosting-comparison.md).

## Selected versions

Checked 2026-10-15 against the sources listed at the end. "Selected" means the exact version SKELETON and later tasks pin. Patch refreshes follow the [update rule](#pinning-and-lockfiles).

| Component | Selected exact version | Support status on 2026-10-15 | Why |
|---|---|---|---|
| .NET SDK | 10.0.402 | .NET 10 LTS, active, end of support 2028-11-14 | Latest 10.0 SDK; released with runtime 10.0.13 on 2026-10-13 |
| .NET / ASP.NET Core runtime | 10.0.13, target framework `net10.0` | Same LTS line | LTS outlives the friends alpha; .NET 8 and 9 end support 2026-11-10 |
| EF Core SQLite provider | `Microsoft.EntityFrameworkCore.Sqlite` 10.0.13 (`.Design` and `dotnet-ef` tool 10.0.13 when migrations start) | EF Core 10 LTS, supported until 2028-11-10 | EF 10 requires .NET 10; provider supports SQLite 3.46.1+ through its bundled native library |
| Node.js (build/dev only) | 24.22.0 "Krypton" | LTS; enters maintenance 2026-10-20, end of life 2028-04-30 | Only LTS line with long remaining life today; satisfies Vite 8 (`^20.19.0 \|\| >=22.12.0`) |
| Package manager | npm 11.19.1 (bundled with Node 24.22.0) | Ships with the selected Node | No extra tool; one lockfile format |
| React / React DOM | 19.3.1 | Current stable (`latest` tag, published 2026-10-14) | Owner-confirmed UI library |
| React type definitions | `@types/react` 19.3.0, `@types/react-dom` 19.3.0 | Current; `@types/react-dom` peers `@types/react` ^19.3.0 | Match React 19.3 |
| TypeScript | 6.0.3 | Latest 6.x; 6.0 is the maintained JS-API line beside 7.0 | See [TypeScript 6 versus 7](#typescript-6-versus-7) |
| Vite | 8.3.4 | Current stable (`latest`, 2026-10-08) | Dev server with `/api` proxy and hashed production build |
| `@vitejs/plugin-react` | 6.1.2 | Current; peers `vite` ^8.0.0 | Official React plugin for Vite 8 |
| PixiJS | `pixi.js` 8.22.0 | Current v8 stable (2026-10-01) | Matches the installed v8 PixiJS skills; added only by a scene task that needs it |

Not selected here: test frameworks, linters, the OpenAPI-to-TypeScript generator, Identity packages and the R2/S3 client. Each is chosen by the task that first needs it (SKELETON, BUILD-CI, IDENTITY, R2-CONFIG) under the same pinning rules, and must target the lines above.

### Deliberately not chosen

- .NET 11 is a release candidate (11.0.0-rc.2, go-live) shipping in November 2026 as an STS release. Stay on .NET 10 LTS; revisit only with a reason.
- Node 26 becomes LTS on 2026-10-28. Node 24 stays selected for this cycle; moving to 26 later is a deliberate, recorded bump, not drift.
- PostgreSQL / Npgsql: deferred (D-27). No second provider, package or migration set.

### TypeScript 6 versus 7

TypeScript 7.0.2 (the native compiler, `latest` on npm since 2026-07-08) ships without a programmatic API; the 7.0 announcement says API-based tools should keep using 6.0 and offers a side-by-side setup. typescript-eslint, the usual lint stack for a Vite React app, refuses to load under TypeScript 7. One package on 6.0.3 keeps `tsc`, editors and lint tooling on the same compiler. Revisit when TypeScript 7.1 ships its API and typescript-eslint supports it. Vite transpiles TypeScript itself and does not depend on this choice.

## This computer versus the selection

Recorded on the owner's second Windows 11 computer (profile `luisf`) on 2026-10-15. Nothing was installed or upgraded by this task.

| Tool | Installed here | Selected | Required action before SKELETON |
|---|---|---|---|
| .NET SDK | 8.0.413, 9.0.307, 10.0.100 | 10.0.402 | Install SDK 10.0.402; `global.json` will refuse 10.0.100 |
| .NET runtime / ASP.NET Core | 10.0.0 (plus 7.0.20, 8.0.x, 9.0.11) | 10.0.13 | Comes with SDK 10.0.402; 10.0.0 lacks a year of servicing fixes |
| Node.js / npm | 22.18.0 / 10.9.3 | 24.22.0 / 11.19.1 | Install Node 24.22.0 (Node 22 meets Vite's minimum but is not the pinned line) |
| pnpm | not installed | not used | None |
| Corepack | 0.33.0 | not used | None |
| `dotnet-ef` | not installed | 10.0.13 as a repository local tool | Restored from the tool manifest in task SQLITE, not installed globally |
| Azure CLI | not installed | version chosen in AZURE-SETUP | Needed by AZURE-SETUP, not by SKELETON |
| Python | 3.12.2 | not part of the app | Runs `scripts/check_docs.py` only |

The 2026-10-13 SKILLS audit computer reported Node 22.23.2; it was not re-inspected here.

Update 2026-10-15 (task 004): SDK 10.0.402 (runtime 10.0.13) and Node 24.22.0 (npm 11.19.1) were installed side-by-side on this computer from official archives, without replacing the existing versions; see the [SKELETON record](../tasks/004-application-skeleton.md#toolchain).

### Selected by SKELETON (task 004)

Checked 2026-10-15 on nuget.org and the xUnit.net documentation (via Context7, `/websites/xunit_net`).

| Component | Exact version | Notes |
|---|---|---|
| `xunit.v3.mtp-v2` | 4.0.2 (published 2026-10-09) | xUnit.net v3 on Microsoft Testing Platform v2; test projects are executables |
| `Microsoft.AspNetCore.Mvc.Testing` | 10.0.13 | In-memory host tests; matches the runtime line |
| Test runner | Microsoft Testing Platform, set in `global.json` `test.runner` | Required for `dotnet test` with MTP on the .NET 10 SDK; no VSTest packages (`Microsoft.NET.Test.Sdk`, `xunit.runner.visualstudio`) |

Not added: ESLint/typescript-eslint, a frontend test runner, an OpenAPI generator and PixiJS. Each arrives with the task that needs it.

Sources: <https://api.nuget.org/v3-flatcontainer/xunit.v3.mtp-v2/index.json>, <https://api.nuget.org/v3-flatcontainer/microsoft.aspnetcore.mvc.testing/index.json>, <https://xunit.net/docs/getting-started/v3/microsoft-testing-platform>.

## Pinning and lockfiles

- `global.json` at the root pins SDK 10.0.402 with `rollForward: disable`. Microsoft recommends `disable` together with NuGet lock files because a different SDK can change implicit package versions.
- NuGet: Central Package Management (`Directory.Packages.props`) with exact versions; `RestorePackagesWithLockFile` on every project; committed `packages.lock.json`; CI restores in locked mode.
- .NET tools: repository manifest `.config/dotnet-tools.json` (for example `dotnet-ef` 10.0.13), restored with `dotnet tool restore`.
- npm: exact versions in `package.json` (no `^`/`~`), `save-exact=true` in `.npmrc`, committed `package-lock.json`, CI installs with `npm ci`. `engines.node` (with `engine-strict=true`) and a root `.node-version` file name 24.22.0; `packageManager` names `npm@11.19.1`.
- Runtime patches: framework-dependent deployment rolls forward to the latest installed 10.0 patch on the host, as App Service patches its built-in stack. The SDK pin governs builds, not the host runtime.
- Update rule: a patch refresh within the selected lines (for example 10.0.13 to 10.0.14, Vite 8.3.4 to 8.3.5) is an explicit commit with release-note evidence. A minor or major change (Vite 9, TypeScript 7, Node 26, .NET 11) needs a task or a recorded owner decision. No automatic dependency bot until a task adds one.
- SKELETON may take a newer patch released after 2026-10-15 if it records the exact version and source; it may not change a line.

## Solution layout

Proposed layout. Only folders with real content are created. `*` = created by SKELETON; `v` = created by VERSION; unmarked = created by the task that first needs it.

```text
woo/
  global.json *                 SDK pin
  Directory.Build.props *       net10.0, nullable, implicit usings, lock files
  Directory.Packages.props *    central exact versions
  .config/dotnet-tools.json     when a .NET tool is first needed
  Woo.slnx *                    .NET 10 SDK default solution format
  VERSION v                     release version source (see below)
  scripts/build.mjs v           coordinated build and publish entry point
  src/
    Woo.Server/ *               ASP.NET Core host: endpoints, modules, EF Core, SPA serving
    Woo.Domain/                 created with the first real game rules (not in SKELETON)
  web/ *                        React + Vite + TypeScript app (package.json, package-lock.json)
  tests/
    Woo.Server.Tests/ *         API and host tests
    Woo.Domain.Tests/           created with Woo.Domain
```

## Module boundaries and dependency direction

Assessment: pragmatic DDD inside a modular monolith fits WOO; a full four-project Domain/Application/Infrastructure/API split does not, yet. Most early work is CRUD-like content administration and persistence plumbing. The rule-heavy parts (production accrual, costs, forging, counters, battle resolution, travel) benefit from pure, deterministic, testable code. So:

- `Woo.Domain`: pure C# game rules and value types. References nothing (no EF Core, ASP.NET Core or Azure packages). Deterministic: time comes in as arguments, not from a clock. Created by the first task with real rules; no empty project before that.
- `Woo.Server`: everything else, organised by feature module folders, each owning its endpoints, use-case code, EF Core mappings and tables. References `Woo.Domain` once it exists.
- Direction: `Woo.Server` -> `Woo.Domain`; tests -> the project they test. `Woo.Domain` never references `Woo.Server`. `web/` talks to `Woo.Server` only over HTTP.

Initial modules, created only when their task arrives: Platform (health, build identity), Identity, Admin content definitions, Worlds and membership, Settlements (resources, buildings), Military (recruitment, forging, talents), Map (districts, routes, movement, scouting), Orders (durable due actions), Combat, Assets (R2), Kingdom, Community (profiles, leaderboards).

Module rules:

- A module writes only its own tables; other modules call its public C# services in process. Cross-module consistency uses one database transaction, not messages.
- One `DbContext` on one SQLite database until a measured reason appears. Per-module contexts are speculative now.
- Interfaces only at real external or swappable boundaries: asset storage (R2 in Azure, local folder in development). Use the built-in `TimeProvider` for testable time. No interface with one implementation, generic repository over EF Core, mediator library, command bus, event bus, CQRS split or microservice.
- Extract a `Woo.Application` or `Woo.Infrastructure` project only when a second host (for example an always-running worker) must share use-case code, recorded as a task decision.

Frontend: `web/src/app` (shell, routing), `web/src/features/<feature>`, `web/src/shared` (UI primitives, API client). A `scenes/` folder appears with the first PixiJS task. API types are generated from the server's OpenAPI document once real endpoints exist; the generator is chosen then.

## Azure application and local development

Proposal for SKELETON, AZURE-SETUP and DEV-DEPLOY:

- One App Service app runs `Woo.Server`. The Vite production build is copied into the server's `wwwroot` at publish. ASP.NET Core serves the hashed assets, `/api/*` endpoints and an `index.html` fallback for client routes. One origin means same-site secure cookies, no CORS and one deployment.
- SKELETON provides only the minimal runnable page and a `/api/health` seam that reports the app is up; it carries no version or commit fields.
- `index.html` and the build-identity response (added by VERSION) are not cached long-term; hashed JS/CSS can be cached immutably (details in VERSION).
- Local development: Kestrel serves the API; the Vite dev server proxies `/api` to it (Vite `server.proxy`), so the browser sees one origin. Running `npm run build` then `dotnet run` reproduces the deployed single-app shape. No SPA proxy package is needed.
- Deployment route: GitHub Actions on pushes to `dev` builds the web app, publishes `Woo.Server` framework-dependent, deploys the zip to the dev app with OIDC sign-in (preferred, to be verified), then checks the health endpoint and that the build-identity response (from VERSION) reports the deployed commit. Only `dev` deploys; `master` production waits for alpha.
- OS: Linux with the built-in `DOTNETCORE|10.0` stack is the candidate (Linux runtime patch channels are documented); Windows also offers .NET 10. AZURE-SETUP decides with the SQLite gate below.
- F1 limits stay as recorded in [hosting comparison](hosting-comparison.md): shared compute, 60 CPU minutes per day, 1 GB memory, 1 GB storage, no Always On, idle unloading with cold starts, quota exhaustion stopping the app until reset, single instance.

Verification gates that stay open:

1. Subscription eligibility: whether the owner's Azure for Students subscription can create an F1 plan in the chosen region and OS (AZURE-SETUP, read-only checks first).
2. SQLite storage compatibility: a writable persistent path outside the deployment artifact, journal mode and locking on the actual App Service storage, single instance, survival across restart and redeploy, and backup/restore (SQLITE, DB-RECOVERY). WAL is generally unsafe on network file systems. A failure is reported to the owner; it does not silently switch the project to PostgreSQL.

## Persistence, assets and scheduling (unchanged directions)

- SQLite local/dev through EF Core 10; PostgreSQL deferred to a future production task (D-27).
- Admin-uploaded images live in Cloudflare R2 behind the Assets module, with immutable revisioned URLs ([asset lifecycle](asset-lifecycle.md), [versioning and cache](versioning-and-cache.md)).
- Timed actions are persisted with due timestamps and processed transactionally and idempotently, with catch-up on request and on startup. An in-process `BackgroundService` may speed processing while the app is awake; it is not the durability mechanism.
- A sleeping F1 host cannot run anything at an exact deadline while unloaded. Punctual shared-world timing needs a separately approved always-running tier or independent scheduler.

## Build identity contract (0.0.1-dev)

Implemented by task 005 VERSION (see [versioning and cache](versioning-and-cache.md#implementation-task-005) for routes, headers and the bump policy). SKELETON delivered the minimal page and `/api/health`; VERSION added the `VERSION` file, `scripts/build.mjs`, `GET /api/build` and the shared footer.

- Release version source: a root `VERSION` file whose only content is the release version, initially `0.0.1-dev` (SemVer 2.0 pre-release). No other copy is typed by hand.
- Fields: `version` (from `VERSION`), `commit` (full 40-character Git SHA), `builtAt` (UTC ISO 8601), `environment` (`dev`; `production` later; `local` for development).
- Coordinated build: one build entry point (the CI job, or a single local build command added by VERSION) reads `VERSION`, captures `commit` and `builtAt` exactly once, and passes those same values to both the web build and the server publish. Neither build reads Git or the clock on its own when values are supplied, so a deployed page and API carry identical metadata by construction.
- Separate builds never claim shared identity: a web or server build run on its own does not generate its own timestamp or commit and present it as the coordinated build's metadata.
- Local development fallback (`npm run dev`, `dotnet run`, IDE builds): `version` still comes from `VERSION`, `environment` is `local`, and `commit` and `builtAt` are empty rather than invented. The footer then reads `0.0.1-dev (local)`. Commit comparison only runs when both sides carry coordinated metadata.
- The server returns the fields from a build-identity response (`GET /api/build`). The web build embeds the same values. The footer shows the version plus the short commit; a commit mismatch between page and API means a stale client.
- The app version is independent of dependency versions and of image revisions. Documentation-only commits do not bump it. Increment policy: [bump policy](versioning-and-cache.md#version-bump-policy).

## Rendering

- Normal interface: React and CSS.
- Combat v1: a static React/CSS battalion board and a final result screen. No replay, animation, timeline or narrated log ([combat proof](../design/combat-visual-prototype.md)).
- PixiJS 8.22.0 is reserved for scenes with a justified need (layered village art, map pan/zoom and hit testing). It is not installed by SKELETON. The installed PixiJS skills target v8 and match this selection.
- Browser rendering never decides game outcomes; the server is authoritative.

## Sources (checked 2026-10-15)

- .NET release index (channels 8.0 to 11.0, latest SDK/runtime, support phase, end-of-life dates): <https://dotnetcli.blob.core.windows.net/dotnet/release-metadata/releases-index.json>
- EF Core releases and support, via Microsoft Learn MCP: <https://learn.microsoft.com/ef/core/what-is-new/>; EF Core 10 requires .NET 10: <https://learn.microsoft.com/ef/core/what-is-new/ef-core-10.0/whatsnew>; SQLite provider engines: <https://learn.microsoft.com/ef/core/providers/sqlite/>
- NuGet versions for `Microsoft.EntityFrameworkCore.Sqlite` (10.0.13 newest 10.0.x): <https://api.nuget.org/v3-flatcontainer/microsoft.entityframeworkcore.sqlite/index.json>
- SDK pinning and lock files, via Learn MCP: <https://learn.microsoft.com/dotnet/core/tools/global-json>, <https://learn.microsoft.com/dotnet/core/install/upgrade#version-pinning>, <https://learn.microsoft.com/nuget/consume-packages/package-references-in-project-files#locking-dependencies>
- SLNX default in .NET 10: <https://learn.microsoft.com/dotnet/core/compatibility/sdk/10.0/dotnet-new-sln-slnx-default>
- Static files and `MapStaticAssets` in ASP.NET Core 10: <https://learn.microsoft.com/aspnet/core/fundamentals/static-files?view=aspnetcore-10.0>
- App Service .NET 10 runtime stack and F1 in the quickstart: <https://learn.microsoft.com/azure/app-service/quickstart-dotnetcore>; Linux runtime patch channels: <https://learn.microsoft.com/azure/app-service/overview-patch-os-runtime>
- Node.js release schedule: <https://raw.githubusercontent.com/nodejs/Release/main/schedule.json>; Node 24.22.0 and bundled npm: <https://nodejs.org/dist/index.json>
- npm registry metadata (dist-tags, publish times, engines, peers) for `react`, `react-dom`, `typescript`, `vite`, `@vitejs/plugin-react`, `pixi.js`, `@types/react`, `@types/react-dom`: <https://registry.npmjs.org/>
- Vite 8 Node requirement and `server.proxy`, via Context7 (`/vitejs/vite`): <https://github.com/vitejs/vite/blob/main/docs/blog/announcing-vite8.md>, <https://github.com/vitejs/vite/blob/main/docs/config/server-options.md>
- TypeScript 7.0 announcement (no API in 7.0, 6.0 side-by-side): <https://devblogs.microsoft.com/typescript/announcing-typescript-7-0/>; typescript-eslint TS 7 guard, via Context7 (`/typescript-eslint/typescript-eslint`)
- Local inventory: `dotnet --list-sdks`, `dotnet --list-runtimes`, `node -v`, `npm -v`, `pnpm -v`, `corepack --version`, `dotnet ef --version`, `az version`, `python --version` on 2026-10-15.
