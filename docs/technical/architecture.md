# Architecture recommendation

Updated: 2026-10-05. Status: design documentation; no game implementation yet.

## Status

Recommended stack, not yet selected or implemented. This records the conversation, not a newly verified dependency/version audit. Confirm supported versions before development.

| Component | Recommendation | Purpose |
|---|---|---|
| UI | React, TypeScript, Vite | Panels, forms, responsive controls |
| Scene rendering | PixiJS | Village, map and combat presentation |
| API | ASP.NET Core | Server-authoritative game actions |
| Database | PostgreSQL | Accounts, content, world state, durable events |
| Worker | Separate .NET Worker process | Timed construction, recruitment, travel and combat |
| Notifications | SignalR | Push updates after state changes |
| Accounts | ASP.NET Core Identity, secure cookies | Login and server-enforced permissions |
| Images | Cloudflare R2, S3-compatible API (service interpretation pending) | Admin-uploaded assets |
| Hosting | Docker on a VPS initially | Simple API/worker deployment |
| Recovery | Off-site DB backups | Restore state after failures |

## Proposed boundaries

Modular monolith with separate API and worker processes. Database is authoritative. Browser graphics never decide outcomes. Persist actions and deadlines; processing must tolerate restarts/retries without duplicate spending or battles. SignalR is notification, not the durable state source.

## Open contracts

Database schema, API routes, identifiers, balance activation, event ordering, concurrency policy, deployment provider, dependency versions, world isolation and recovery objectives. See [data model](data-model.md). No Kubernetes requirement established.
