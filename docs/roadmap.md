# Development milestones

Updated: 2026-10-05. No application implementation has started. Current numbered queue is [tasks_roadmap.md](../tasks_roadmap.md).

## Current path

1. Claude setup, selected dependency versions and pragmatic architecture; create a skeleton and first running page.
2. Visible app version/build identity, real build/test CI and an early Azure dev deployment.
3. Shared themes and UI foundations; polish landing, game, village/map/combat/forge, admin and onboarding screens with explicit mock data.
4. SQLite local/dev persistence and backup/redeploy verification; real login, admin authorisation, Worlds and player kingdom membership.
5. Configurable medieval content and uploaded assets, independent image revisions and cleanup; saved scene editor and end-to-end village visual gate.
6. Persisted village economy, timed actions, construction, recruitment, ordinary forging and talents.
7. Authoritative map/routes/scouting/movement, scoped activities, numerical counters, battles and reports.
8. Kingdom campaigns, siege/capture/recovery and useful cooperation.
9. Offline/live state reconciliation, balance/mobile/reliability validation and a friends alpha with feedback-driven task additions.

Do not require a full-game medieval specification before the skeleton/first deployment. Decide each gameplay feature's rules before implementing that feature. Prototype UI does not claim working persistence/authentication.

## Continuous refinement

Each implementation cycle reads the current queue and pinned task detail. Split broad work, insert discoveries using whole numbers, record decisions and preserve immutable task keys. All current tasks are TODO; numbers may change. Dependencies govern readiness, not merely position.

## Future releases

Gradual rune emergence and advanced fantasy mechanics are separate updates. Naval/air/tunnel movement and other kingdoms await their own scope. Friends alpha does not imply production release. PostgreSQL and master production Actions are deferred until an explicit production decision.

## Village visual feasibility
During VILLAGE-UI, prove the small layered Arkazia slice locally before expanding art families. Do not wait for persistence/R2 to discover inconsistent perspective or wall seams. Later HOTSPOTS and VILLAGE-VISUAL-GATE verify admin-controlled configuration, uploads, revisions and cleanup. Read [the ordered proof](design/village-visual-prototype.md). All steps remain TODO; documentation and preview generation are not implementation.
