# Task roadmap

Updated: 2026-10-05. Implementation has not begun. All tasks are TODO; initial sequencing below is proposed and may be refined.

Current working branch: dev. Production branch: master. Read [workflow](docs/workflow/development-workflow.md) and [prompt protocol](docs/workflow/implementation-protocol.md).

| Order | Key | Task | Status | Detail |
|---|---|---|---|---|
| 001 | SKILLS | Skills and plugins health check | TODO | [001-skills-preflight.md](docs/tasks/001-skills-preflight.md) |
| 002 | STACK | Select tech stack and Azure-compatible layout | TODO | [002-tech-stack.md](docs/tasks/002-tech-stack.md) |
| 003 | RULES | Specify medieval alpha rules | TODO | [003-medieval-rules.md](docs/tasks/003-medieval-rules.md) |
| 004 | DATABASE | Database model and persistence design | TODO | [004-database.md](docs/tasks/004-database.md) |
| 005 | FOUNDATION | Application scaffold and local validation | TODO | [005-application-foundation.md](docs/tasks/005-application-foundation.md) |
| 006 | AZURE-DEV | Azure dev setup and deployment Action | TODO | [006-azure-dev.md](docs/tasks/006-azure-dev.md) |
| 007 | ACCOUNTS | Login, kingdom selection and owner admin access | TODO | [007-accounts.md](docs/tasks/007-accounts.md) |
| 008 | ADMIN-ASSETS | Content administration and image lifecycle | TODO | [008-admin-assets.md](docs/tasks/008-admin-assets.md) |
| 009 | VILLAGE | Playable village and economy | TODO | [009-village-economy.md](docs/tasks/009-village-economy.md) |
| 010 | FORGING | Default troops, medieval forging and talents | TODO | [010-troops-and-forging.md](docs/tasks/010-troops-and-forging.md) |
| 011 | FRONTIER | Map, kingdom battles and conquest | TODO | [011-map-and-combat.md](docs/tasks/011-map-and-combat.md) |
| 012 | ALPHA | Friends alpha and production-release gate | TODO | [012-friends-alpha.md](docs/tasks/012-friends-alpha.md) |

## Insertion and renumbering

If STACK is 002 and DATABASE is 003, inserting a stack refinement at 003 moves DATABASE to 004. Rename its detail file and every affected link/number in the same commit. Keep its immutable DATABASE key and any existing execution branch/commit history. Status travels with the task. No duplicate order numbers or competing files.

## Status policy

TODO -> WIP when the implementation prompt is dispatched. WIP remains during review/fixes/merge. DONE requires approval, verified dev merge and deployment verification where required. Review outcomes are separate: APPROVED / NEEDS CHANGES / REFUSAL. Do not mark work DONE from Claude's summary alone.

One implementation task at a time. Future fantasy systems are tracked in design docs until their release scope is chosen, rather than preloading this executable queue.
