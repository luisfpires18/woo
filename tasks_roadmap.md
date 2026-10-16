# Task roadmap

Updated: 2026-10-16. 001 SKILLS to 006 BUILD-CI are DONE, 007 AZURE-SETUP is WIP and all other tasks remain TODO. MCP-SETUP was approved by ChatGPT, verified on remote dev at 99abf03 with documentation CI passing, and issue #64 is closed. STACK was first dispatched as 002 (issue #2) against e8703aa on feat/002-tech-stack, paused on 2026-10-15 for MCP-SETUP, resumed as 003 against base 99abf03 on feat/003-tech-stack, then approved and merged at 588c5a0 with documentation CI passing; see its [dispatch history](docs/tasks/003-tech-stack.md#dispatch-history). SKELETON (issue #3) was approved on 2026-10-16 and merged at 77e206c with documentation CI passing. VERSION (issue #4) was approved on 2026-10-16 and merged at dc407cf with documentation CI passing. BUILD-CI (issue #5) was approved on 2026-10-16 and merged at 0e69cad; its first Ubuntu application and documentation runs passed. AZURE-SETUP is dispatched against 0e69cad and has provisioned the West Europe F1 dev app (see its task record). This expanded sequence replaced the original 12 broad tasks before any dispatch, so no branch history was affected.

Working branch: dev. Use sequential whole numbers, not 012a/012b. Numbers are current order; immutable keys preserve identity during insertion. Detailed files hold scope, dependencies and acceptance criteria. This is an evolving queue, not authorisation to implement everything now.

## Delivery approach

Get setup, a pragmatic architecture/skeleton, a versioned first page and a real Azure dev deployment working early. Polish the landing/game/admin/onboarding UI with explicit mock data, then implement SQLite and real account-to-World-to-Kingdom membership. Database comes before working login/world persistence even though its screens can be designed first. Resolve gameplay rules per feature, immediately before that feature's implementation; do not block the first page on the entire game's design.

Architecture is defined in task 003 ([stack decision record](docs/technical/stack.md), approved). DDD is a candidate approach, not a demand for complex infrastructure. SQLite is confirmed for local/dev; production PostgreSQL is deferred. App version starts with an illustrative 0.0.1-dev footer; image revisions refresh independently.

| Order | Key | Task | Status | Detail |
|---|---|---|---|---|
| 001 | SKILLS | Claude setup and skills health check | DONE | [001-skills-preflight.md](docs/tasks/001-skills-preflight.md) |
| 002 | MCP-SETUP | Microsoft Learn, Context7 and read-only GitHub integration | DONE | [002-mcp-setup.md](docs/tasks/002-mcp-setup.md) |
| 003 | STACK | Tech stack, dependency versions and architecture | DONE | [003-tech-stack.md](docs/tasks/003-tech-stack.md) |
| 004 | SKELETON | Create skeleton and run the first page | DONE | [004-application-skeleton.md](docs/tasks/004-application-skeleton.md) |
| 005 | VERSION | App version footer and build identity | DONE | [005-app-version.md](docs/tasks/005-app-version.md) |
| 006 | BUILD-CI | GitHub build and focused test Actions | DONE | [006-build-ci.md](docs/tasks/006-build-ci.md) |
| 007 | AZURE-SETUP | Azure dev provisioning through step-by-step commands | WIP | [007-azure-dev-setup.md](docs/tasks/007-azure-dev-setup.md) |
| 008 | DEV-DEPLOY | Deploy the first page automatically from dev | TODO | [008-dev-deployment.md](docs/tasks/008-dev-deployment.md) |
| 009 | DESIGN-SYSTEM | Shared UI foundations, themes and kingdom colours | TODO | [009-design-system.md](docs/tasks/009-design-system.md) |
| 010 | LANDING | First landing page UI polish | TODO | [010-landing-page.md](docs/tasks/010-landing-page.md) |
| 011 | GAME-SHELL | Game UI structure and navigation polish | TODO | [011-game-ui-shell.md](docs/tasks/011-game-ui-shell.md) |
| 012 | VILLAGE-UI | Village scene and interaction prototype | TODO | [012-village-ui.md](docs/tasks/012-village-ui.md) |
| 013 | MAP-UI | World map visual and editing proof | TODO | [013-map-ui.md](docs/tasks/013-map-ui.md) |
| 014 | COMBAT-UI | Static battalion board and result proof | TODO | [014-combat-ui.md](docs/tasks/014-combat-ui.md) |
| 015 | FORGE-UI | Forge, equipment and talent UI prototypes | TODO | [015-forge-ui.md](docs/tasks/015-forge-ui.md) |
| 016 | ADMIN-UI | Admin dashboard UI and content-editor structure | TODO | [016-admin-dashboard-ui.md](docs/tasks/016-admin-dashboard-ui.md) |
| 017 | AUTH-UI | Login, world selection and kingdom-choice UI | TODO | [017-login-onboarding-ui.md](docs/tasks/017-login-onboarding-ui.md) |
| 018 | SQLITE | SQLite persistence and migration foundation | TODO | [018-sqlite-foundation.md](docs/tasks/018-sqlite-foundation.md) |
| 019 | DB-RECOVERY | Dev database persistence, backup and restore | TODO | [019-sqlite-backup-recovery.md](docs/tasks/019-sqlite-backup-recovery.md) |
| 020 | IDENTITY | Real accounts, sessions and login | TODO | [020-accounts-login.md](docs/tasks/020-accounts-login.md) |
| 021 | ROLES | Owner bootstrap and admin authorisation | TODO | [021-owner-admin-access.md](docs/tasks/021-owner-admin-access.md) |
| 022 | WORLDS | World creation, availability and alpha factions | TODO | [022-worlds-admin.md](docs/tasks/022-worlds-admin.md) |
| 023 | MEMBERSHIP | Associate player with World and chosen kingdom | TODO | [023-player-world-kingdom.md](docs/tasks/023-player-world-kingdom.md) |
| 024 | RESOURCES-ADMIN | Admin resource definitions | TODO | [024-resource-definitions.md](docs/tasks/024-resource-definitions.md) |
| 025 | BUILDINGS-ADMIN | Admin building definitions and prerequisites | TODO | [025-building-definitions.md](docs/tasks/025-building-definitions.md) |
| 026 | UNITS-ADMIN | Admin unit definitions and default equipment | TODO | [026-unit-definitions.md](docs/tasks/026-unit-definitions.md) |
| 027 | EQUIPMENT-ADMIN | Admin equipment and ordinary forge recipes | TODO | [027-equipment-recipes.md](docs/tasks/027-equipment-recipes.md) |
| 028 | TALENTS-ADMIN | Admin talent tree definitions | TODO | [028-talent-definitions.md](docs/tasks/028-talent-definitions.md) |
| 029 | CLOUDFLARE | Cloudflare account, access and CLI configuration | TODO | [029-cloudflare-account-cli.md](docs/tasks/029-cloudflare-account-cli.md) |
| 030 | R2-SETUP | Create and verify the R2 dev asset bucket | TODO | [030-r2-bucket-creation.md](docs/tasks/030-r2-bucket-creation.md) |
| 031 | R2-CONFIG | R2 upload credentials, CORS and image delivery | TODO | [031-r2-upload-delivery.md](docs/tasks/031-r2-upload-delivery.md) |
| 032 | ASSETS | Entity image upload, replacement and cleanup | TODO | [032-asset-upload-cleanup.md](docs/tasks/032-asset-upload-cleanup.md) |
| 033 | IMAGE-CACHE | Image cache revisions and immediate replacements | TODO | [033-image-cache-revisions.md](docs/tasks/033-image-cache-revisions.md) |
| 034 | HOTSPOTS | Admin village anchors and clickable hotspots | TODO | [034-scene-hotspot-editor.md](docs/tasks/034-scene-hotspot-editor.md) |
| 035 | CONFIG | Configuration publish, audit and active-order policy | TODO | [035-content-publish-policy.md](docs/tasks/035-content-publish-policy.md) |
| 036 | VILLAGE-VISUAL-GATE | Layered village art and admin integration gate | TODO | [036-village-visual-gate.md](docs/tasks/036-village-visual-gate.md) |
| 037 | VILLAGE-CREATE | Initial village placement and player settlement | TODO | [037-initial-village.md](docs/tasks/037-initial-village.md) |
| 038 | RESOURCE-STATE | Resource production, storage and offline accrual | TODO | [038-resource-production.md](docs/tasks/038-resource-production.md) |
| 039 | ORDERS | Durable action queues and catch-up processing | TODO | [039-durable-action-queues.md](docs/tasks/039-durable-action-queues.md) |
| 040 | CONSTRUCTION | Building construction, upgrades and queues | TODO | [040-building-construction.md](docs/tasks/040-building-construction.md) |
| 041 | RECRUIT | Recruitment and troop inventory | TODO | [041-troop-recruitment.md](docs/tasks/041-troop-recruitment.md) |
| 042 | FORGING | Ordinary weaponsmith and armorsmith upgrades | TODO | [042-medieval-forging.md](docs/tasks/042-medieval-forging.md) |
| 043 | TALENTS | Blacksmith experience, points and talent choices | TODO | [043-blacksmith-progression.md](docs/tasks/043-blacksmith-progression.md) |
| 044 | DISTRICTS | Authoritative map districts and ownership | TODO | [044-world-districts.md](docs/tasks/044-world-districts.md) |
| 045 | ROUTES | Route connectivity and travel calculations | TODO | [045-routes-and-travel.md](docs/tasks/045-routes-and-travel.md) |
| 046 | SCOUTING | Scouting and map information visibility | TODO | [046-scouting-and-visibility.md](docs/tasks/046-scouting-and-visibility.md) |
| 047 | MOVEMENT | Troop movement, stationing and return | TODO | [047-troop-movement.md](docs/tasks/047-troop-movement.md) |
| 048 | DUTIES | Useful kingdom duties | TODO | [048-kingdom-duties.md](docs/tasks/048-kingdom-duties.md) |
| 049 | COUNTERS | Unit counters and battle-model validation | TODO | [049-battle-model.md](docs/tasks/049-battle-model.md) |
| 050 | BATTLES | Persisted automatic battle resolution | TODO | [050-battle-resolution.md](docs/tasks/050-battle-resolution.md) |
| 051 | REPORTS | Final battle results and army totals | TODO | [051-battle-reports.md](docs/tasks/051-battle-reports.md) |
| 052 | CAMP | Ordinary animal camps and expeditions | TODO | [052-ordinary-camps.md](docs/tasks/052-ordinary-camps.md) |
| 053 | CAMPAIGNS | Declared kingdom campaigns and defence commitments | TODO | [053-kingdom-campaigns.md](docs/tasks/053-kingdom-campaigns.md) |
| 054 | SIEGE | Village siege and defensive building effects | TODO | [054-village-siege.md](docs/tasks/054-village-siege.md) |
| 055 | CAPTURE | Capture, occupation and district transfer | TODO | [055-territory-capture.md](docs/tasks/055-territory-capture.md) |
| 056 | RECOVERY | Defeat, wounded troops and player recovery | TODO | [056-defeat-and-recovery.md](docs/tasks/056-defeat-and-recovery.md) |
| 057 | COOP | Kingdom overview, defence requests and contribution credit | TODO | [057-kingdom-cooperation.md](docs/tasks/057-kingdom-cooperation.md) |
| 058 | PROFILES | Player and settlement profiles | TODO | [058-player-settlement-profiles.md](docs/tasks/058-player-settlement-profiles.md) |
| 059 | LEADERBOARDS | World and season leaderboards | TODO | [059-world-leaderboards.md](docs/tasks/059-world-leaderboards.md) |
| 060 | OFFLINE | Offline reconciliation and live UI notifications | TODO | [060-offline-ui-updates.md](docs/tasks/060-offline-ui-updates.md) |
| 061 | BALANCE | Economy, counters and faction balance pass | TODO | [061-economy-and-faction-validation.md](docs/tasks/061-economy-and-faction-validation.md) |
| 062 | QA | Mobile, accessibility, persistence and permission validation | TODO | [062-mobile-and-reliability.md](docs/tasks/062-mobile-and-reliability.md) |
| 063 | FRIENDS | Invite-only friends alpha setup and support | TODO | [063-friends-alpha-setup.md](docs/tasks/063-friends-alpha-setup.md) |
| 064 | ALPHA-FEEDBACK | Run medieval alpha and refine next tasks | TODO | [064-alpha-feedback.md](docs/tasks/064-alpha-feedback.md) |

## Milestones

| Milestone | Tasks | Outcome |
|---|---|---|
| Foundation | 001, 002, 003, 004, 005, 006 | Running versioned skeleton |
| Deployment | 007, 008 | First page deployed to dev |
| UI | 009, 010, 011, 012, 013, 014, 015, 016, 017 | Polished screen prototypes |
| Persistence | 018, 019 | Durable SQLite local/dev |
| Accounts | 020, 021, 022, 023 | Identity and world membership |
| Content | 024, 025, 026, 027, 028, 035 | Independently reviewed content slices |
| Assets | 029, 030, 031, 032, 033, 034, 036 | Cloudflare/R2, upload, cache and hotspots |
| Village | 037, 038, 039, 040 | Independently reviewed village slices |
| Military | 041, 042, 043 | Independently reviewed military slices |
| Map | 044, 045, 046, 047 | Independently reviewed map slices |
| Activities | 048, 052 | Independently reviewed activities slices |
| Combat | 049, 050, 051 | Independently reviewed combat slices |
| Conquest | 053, 054, 055, 056 | Independently reviewed conquest slices |
| Kingdom | 057 | Independently reviewed kingdom slices |
| Community | 058, 059 | Player/settlement inspection and world-season rankings |
| Reliability | 060 | Independently reviewed reliability slices |
| Validation | 061, 062 | Independently reviewed validation slices |
| Alpha | 063, 064 | Friends test and roadmap refinement |

Dependencies in detailed files govern dispatch. Task order is a current plan; new tasks are inserted using whole numbers.

## Insertion and status rules

Insert a task by renumbering later rows/files and updating links atomically. Preserve stable keys, status and any dispatched branch/prompt/commit history. Future branches use the current number. Small internal steps are checklists; separately prompted/reviewed/merged work gets its own whole number.

TODO -> WIP when one implementation prompt is dispatched. WIP remains through fixes/review/merge. DONE requires APPROVED, verified dev merge and required deployment verification. Missing evidence does not count as completion.

## Deferred future systems

Rune discovery, Conduit, Aspect, magical armour, Chaos/Order, heroes, additional kingdoms and naval/air/tunnel systems remain future design material, not active medieval implementation tasks. Production/PostgreSQL/master deployment gets its own explicit tasks when the owner chooses that release, not merely because the friends alpha starts.

## Village visual proof
The queue contains 64 tasks: 001 SKILLS to 006 BUILD-CI are DONE, 007 AZURE-SETUP is WIP, the rest TODO. VILLAGE-UI (012) is the early local layered-art proof, HOTSPOTS (034) is the saved admin editor, and VILLAGE-VISUAL-GATE (036) is the end-to-end acceptance gate after CONFIG (035). Later tasks renumbered with immutable keys preserved; none had been dispatched. VILLAGE-CREATE requires this gate. Internal proof steps are checklists, not letter-suffixed tasks.
Read [step-by-step preparation](docs/design/village-visual-prototype.md). Flat illustration hotspots alone do not prove interchangeable buildings or modular walls.

## Scene scope boundaries

012 (VILLAGE-UI), 034 (HOTSPOTS) and 036 (VILLAGE-VISUAL-GATE) cover village interiors only. 013 (MAP-UI) independently proves the world-map rendering and local square-editing workflow; 044 (DISTRICTS) later owns saved map definitions and authoritative ownership. Resource types, production-site visuals and resource-screen layout are not decided and must not be invented during either scene proof. 024 and 038 remain gated by an owner-approved economy direction. Shared rendering utilities may be reused, but acceptance evidence and content remain separate.

## Combat proof and integration

014 (COMBAT-UI) proves a static battalion board and final-result layout. 049 (COUNTERS) validates numerical rules; 050 (BATTLES) safely persists server outcomes; 051 (REPORTS) displays authorised final results and reconciled army totals. V1 requires no animation, replay, timeline or narrated combat log. Prefer React/CSS for the board; PixiJS is optional only if a justified requirement emerges. Later visual improvements get separately scoped tasks. Read [combat preparation](docs/design/combat-visual-prototype.md). All remain TODO.

## Settlement terminology and community views

Settlement/settlements is the preferred player-facing term. Existing village filenames, VILLAGE keys and historical references remain technical aliases until a deliberate code/data naming decision. 058 (PROFILES) adds player and settlement pages; 059 (LEADERBOARDS) adds world/season rankings. Later tasks renumbered atomically, stable keys preserved. Ranking categories/formulas and profile visibility need approval before dispatch. Dev is the desired default branch; actual repository setting remains pending verification.
