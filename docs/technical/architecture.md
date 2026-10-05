# Architecture recommendation

Updated: 2026-10-05. Status: design documentation; no game implementation yet.

## Status

Recommended stack, not yet selected or implemented. This records the conversation, not a newly verified dependency/version audit. Confirm supported versions before development.

| Component | Recommendation | Purpose |
|---|---|---|
| UI | React, TypeScript, Vite | Panels, forms, responsive controls |
| Scene rendering | PixiJS | Village, map and combat presentation |
| API | ASP.NET Core | Server-authoritative game actions |
| Database | SQLite for local/dev; PostgreSQL considered for future production | Current persistence direction |
| Job processing | Persisted due actions; in-process initially or separate worker when hosting justifies it | Timed construction, recruitment, travel and combat |
| Notifications | SignalR | Push updates after state changes |
| Accounts | ASP.NET Core Identity, secure cookies | Login and server-enforced permissions |
| Images | Cloudflare R2, S3-compatible API (service interpretation pending) | Admin-uploaded assets |
| Hosting | Azure App Service; F1 preferred for dev | Owner-selected direction; exact runtime/layout pending |
| Recovery | Off-site DB backups | Restore state after failures |

## Proposed boundaries

Proposed modular monolith, initially one deployable application; keep job-processing logic separable. A separate worker is not an automatic early requirement. Database is authoritative. Browser graphics never decide outcomes. Persist actions and deadlines; processing must tolerate restarts/retries without duplicate spending or battles. SignalR is notification, not the durable state source.

## Open contracts

Database schema, API routes, identifiers, balance activation, event ordering, concurrency policy, Azure layout, dependency versions, world isolation and recovery objectives. See [data model](data-model.md). No Kubernetes requirement established.

## Hosting correction

Azure hosting is confirmed by the owner, superseding the VPS proposal. A separate always-running worker is still a design proposal and is not guaranteed on sleeping F1. See [Azure constraints](azure-and-deployment.md) before selecting deployment and database architecture.

## Revised early path

Define architecture and exact supported dependency versions first, then create a minimal skeleton and running page. Assess pragmatic DDD with explicit boundaries and dependency direction (Domain/Application/Infrastructure/API as a candidate layout). DDD was offered as an example, not a confirmed demand for every tactical pattern. Choose the smallest useful structure; no speculative microservices.

First landing/game/admin UI can be polished with mock data before SQLite. Real login, owner roles, World membership and kingdom choice require persisted state. See [version/cache contract](versioning-and-cache.md). PostgreSQL provider/server is not implemented now; record migration concerns rather than developing two providers prematurely.

## Backend reassessment, 2026-10-05
Read [backend comparison](backend-comparison.md). ASP.NET Core remains the recommendation, not a confirmed backend choice. Java/Spring, Go and TypeScript/Node are credible alternatives. React/TypeScript is the owner's preferred frontend; rendering remains independent of backend language. Final selection occurs in STACK. Early REST/polling and in-process persisted job processing are proposals; SignalR/another worker require justification, not automatic provisioning.
