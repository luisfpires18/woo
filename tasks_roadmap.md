# Task roadmap

Updated: 2026-10-05. All implementation tasks remain TODO. This expanded sequence replaces the original 12 broad tasks; none was dispatched, so no active branch history is affected.

Working branch: dev. Use sequential whole numbers, not 012a/012b. Numbers are current order; immutable keys preserve identity during insertion. Detailed files hold scope, dependencies and acceptance criteria. This is an evolving queue, not authorisation to implement everything now.

## Delivery approach

Get setup, a pragmatic architecture/skeleton, a versioned first page and a real Azure dev deployment working early. Polish the landing/game/admin/onboarding UI with explicit mock data, then implement SQLite and real account-to-World-to-Kingdom membership. Database comes before working login/world persistence even though its screens can be designed first. Resolve gameplay rules per feature, immediately before that feature's implementation; do not block the first page on the entire game's design.

Architecture must be defined in task 002. DDD is a candidate approach, not a demand for complex infrastructure. SQLite is confirmed for local/dev; production PostgreSQL is deferred. App version starts with an illustrative 0.0.1-dev footer; image revisions refresh independently.

| Order | Key | Task | Status | Detail |
|---|---|---|---|---|
| 001 | SKILLS | Claude setup and skills health check | TODO | [001-skills-preflight.md](docs/tasks/001-skills-preflight.md) |
| 002 | STACK | Tech stack, dependency versions and architecture | TODO | [002-tech-stack.md](docs/tasks/002-tech-stack.md) |
| 003 | SKELETON | Create skeleton and run the first page | TODO | [003-application-skeleton.md](docs/tasks/003-application-skeleton.md) |
| 004 | VERSION | App version footer and build identity | TODO | [004-app-version.md](docs/tasks/004-app-version.md) |
| 005 | BUILD-CI | GitHub build and focused test Actions | TODO | [005-build-ci.md](docs/tasks/005-build-ci.md) |
| 006 | AZURE-SETUP | Azure dev provisioning through step-by-step commands | TODO | [006-azure-dev-setup.md](docs/tasks/006-azure-dev-setup.md) |
| 007 | DEV-DEPLOY | Deploy the first page automatically from dev | TODO | [007-dev-deployment.md](docs/tasks/007-dev-deployment.md) |
| 008 | DESIGN-SYSTEM | Shared UI foundations, themes and kingdom colours | TODO | [008-design-system.md](docs/tasks/008-design-system.md) |
| 009 | LANDING | First landing page UI polish | TODO | [009-landing-page.md](docs/tasks/009-landing-page.md) |
| 010 | GAME-SHELL | Game UI structure and navigation polish | TODO | [010-game-ui-shell.md](docs/tasks/010-game-ui-shell.md) |
| 011 | VILLAGE-UI | Village scene and interaction prototype | TODO | [011-village-ui.md](docs/tasks/011-village-ui.md) |
| 012 | MAP-UI | World map visual and editing proof | TODO | [012-map-ui.md](docs/tasks/012-map-ui.md) |
| 013 | COMBAT-UI | Static battalion board and result proof | TODO | [013-combat-ui.md](docs/tasks/013-combat-ui.md) |
| 014 | FORGE-UI | Forge, equipment and talent UI prototypes | TODO | [014-forge-ui.md](docs/tasks/014-forge-ui.md) |
| 015 | ADMIN-UI | Admin dashboard UI and content-editor structure | TODO | [015-admin-dashboard-ui.md](docs/tasks/015-admin-dashboard-ui.md) |
| 016 | AUTH-UI | Login, world selection and kingdom-choice UI | TODO | [016-login-onboarding-ui.md](docs/tasks/016-login-onboarding-ui.md) |
| 017 | SQLITE | SQLite persistence and migration foundation | TODO | [017-sqlite-foundation.md](docs/tasks/017-sqlite-foundation.md) |
| 018 | DB-RECOVERY | Dev database persistence, backup and restore | TODO | [018-sqlite-backup-recovery.md](docs/tasks/018-sqlite-backup-recovery.md) |
| 019 | IDENTITY | Real accounts, sessions and login | TODO | [019-accounts-login.md](docs/tasks/019-accounts-login.md) |
| 020 | ROLES | Owner bootstrap and admin authorisation | TODO | [020-owner-admin-access.md](docs/tasks/020-owner-admin-access.md) |
| 021 | WORLDS | World creation, availability and alpha factions | TODO | [021-worlds-admin.md](docs/tasks/021-worlds-admin.md) |
| 022 | MEMBERSHIP | Associate player with World and chosen kingdom | TODO | [022-player-world-kingdom.md](docs/tasks/022-player-world-kingdom.md) |
| 023 | RESOURCES-ADMIN | Admin resource definitions | TODO | [023-resource-definitions.md](docs/tasks/023-resource-definitions.md) |
| 024 | BUILDINGS-ADMIN | Admin building definitions and prerequisites | TODO | [024-building-definitions.md](docs/tasks/024-building-definitions.md) |
| 025 | UNITS-ADMIN | Admin unit definitions and default equipment | TODO | [025-unit-definitions.md](docs/tasks/025-unit-definitions.md) |
| 026 | EQUIPMENT-ADMIN | Admin equipment and ordinary forge recipes | TODO | [026-equipment-recipes.md](docs/tasks/026-equipment-recipes.md) |
| 027 | TALENTS-ADMIN | Admin talent tree definitions | TODO | [027-talent-definitions.md](docs/tasks/027-talent-definitions.md) |
| 028 | CLOUDFLARE | Cloudflare account, access and CLI configuration | TODO | [028-cloudflare-account-cli.md](docs/tasks/028-cloudflare-account-cli.md) |
| 029 | R2-SETUP | Create and verify the R2 dev asset bucket | TODO | [029-r2-bucket-creation.md](docs/tasks/029-r2-bucket-creation.md) |
| 030 | R2-CONFIG | R2 upload credentials, CORS and image delivery | TODO | [030-r2-upload-delivery.md](docs/tasks/030-r2-upload-delivery.md) |
| 031 | ASSETS | Entity image upload, replacement and cleanup | TODO | [031-asset-upload-cleanup.md](docs/tasks/031-asset-upload-cleanup.md) |
| 032 | IMAGE-CACHE | Image cache revisions and immediate replacements | TODO | [032-image-cache-revisions.md](docs/tasks/032-image-cache-revisions.md) |
| 033 | HOTSPOTS | Admin village anchors and clickable hotspots | TODO | [033-scene-hotspot-editor.md](docs/tasks/033-scene-hotspot-editor.md) |
| 034 | CONFIG | Configuration publish, audit and active-order policy | TODO | [034-content-publish-policy.md](docs/tasks/034-content-publish-policy.md) |
| 035 | VILLAGE-VISUAL-GATE | Layered village art and admin integration gate | TODO | [035-village-visual-gate.md](docs/tasks/035-village-visual-gate.md) |
| 036 | VILLAGE-CREATE | Initial village placement and player settlement | TODO | [036-initial-village.md](docs/tasks/036-initial-village.md) |
| 037 | RESOURCE-STATE | Resource production, storage and offline accrual | TODO | [037-resource-production.md](docs/tasks/037-resource-production.md) |
| 038 | ORDERS | Durable action queues and catch-up processing | TODO | [038-durable-action-queues.md](docs/tasks/038-durable-action-queues.md) |
| 039 | CONSTRUCTION | Building construction, upgrades and queues | TODO | [039-building-construction.md](docs/tasks/039-building-construction.md) |
| 040 | RECRUIT | Recruitment and troop inventory | TODO | [040-troop-recruitment.md](docs/tasks/040-troop-recruitment.md) |
| 041 | FORGING | Ordinary weaponsmith and armorsmith upgrades | TODO | [041-medieval-forging.md](docs/tasks/041-medieval-forging.md) |
| 042 | TALENTS | Blacksmith experience, points and talent choices | TODO | [042-blacksmith-progression.md](docs/tasks/042-blacksmith-progression.md) |
| 043 | DISTRICTS | Authoritative map districts and ownership | TODO | [043-world-districts.md](docs/tasks/043-world-districts.md) |
| 044 | ROUTES | Route connectivity and travel calculations | TODO | [044-routes-and-travel.md](docs/tasks/044-routes-and-travel.md) |
| 045 | SCOUTING | Scouting and map information visibility | TODO | [045-scouting-and-visibility.md](docs/tasks/045-scouting-and-visibility.md) |
| 046 | MOVEMENT | Troop movement, stationing and return | TODO | [046-troop-movement.md](docs/tasks/046-troop-movement.md) |
| 047 | DUTIES | Useful kingdom duties | TODO | [047-kingdom-duties.md](docs/tasks/047-kingdom-duties.md) |
| 048 | COUNTERS | Unit counters and battle-model validation | TODO | [048-battle-model.md](docs/tasks/048-battle-model.md) |
| 049 | BATTLES | Persisted automatic battle resolution | TODO | [049-battle-resolution.md](docs/tasks/049-battle-resolution.md) |
| 050 | REPORTS | Final battle results and army totals | TODO | [050-battle-reports.md](docs/tasks/050-battle-reports.md) |
| 051 | CAMP | Ordinary animal camps and expeditions | TODO | [051-ordinary-camps.md](docs/tasks/051-ordinary-camps.md) |
| 052 | CAMPAIGNS | Declared kingdom campaigns and defence commitments | TODO | [052-kingdom-campaigns.md](docs/tasks/052-kingdom-campaigns.md) |
| 053 | SIEGE | Village siege and defensive building effects | TODO | [053-village-siege.md](docs/tasks/053-village-siege.md) |
| 054 | CAPTURE | Capture, occupation and district transfer | TODO | [054-territory-capture.md](docs/tasks/054-territory-capture.md) |
| 055 | RECOVERY | Defeat, wounded troops and player recovery | TODO | [055-defeat-and-recovery.md](docs/tasks/055-defeat-and-recovery.md) |
| 056 | COOP | Kingdom overview, defence requests and contribution credit | TODO | [056-kingdom-cooperation.md](docs/tasks/056-kingdom-cooperation.md) |
| 057 | OFFLINE | Offline reconciliation and live UI notifications | TODO | [057-offline-ui-updates.md](docs/tasks/057-offline-ui-updates.md) |
| 058 | BALANCE | Economy, counters and faction balance pass | TODO | [058-economy-and-faction-validation.md](docs/tasks/058-economy-and-faction-validation.md) |
| 059 | QA | Mobile, accessibility, persistence and permission validation | TODO | [059-mobile-and-reliability.md](docs/tasks/059-mobile-and-reliability.md) |
| 060 | FRIENDS | Invite-only friends alpha setup and support | TODO | [060-friends-alpha-setup.md](docs/tasks/060-friends-alpha-setup.md) |
| 061 | ALPHA-FEEDBACK | Run medieval alpha and refine next tasks | TODO | [061-alpha-feedback.md](docs/tasks/061-alpha-feedback.md) |

## Milestones

| Milestone | Tasks | Outcome |
|---|---|---|
| Foundation | 001, 002, 003, 004, 005 | Running versioned skeleton |
| Deployment | 006, 007 | First page deployed to dev |
| UI | 008, 009, 010, 011, 012, 013, 014, 015, 016 | Polished screen prototypes |
| Persistence | 017, 018 | Durable SQLite local/dev |
| Accounts | 019, 020, 021, 022 | Identity and world membership |
| Content | 023, 024, 025, 026, 027, 034 | Independently reviewed content slices |
| Assets | 028, 029, 030, 031, 032, 033, 035 | Cloudflare/R2, upload, cache and hotspots |
| Village | 036, 037, 038, 039 | Independently reviewed village slices |
| Military | 040, 041, 042 | Independently reviewed military slices |
| Map | 043, 044, 045, 046 | Independently reviewed map slices |
| Activities | 047, 051 | Independently reviewed activities slices |
| Combat | 048, 049, 050 | Independently reviewed combat slices |
| Conquest | 052, 053, 054, 055 | Independently reviewed conquest slices |
| Kingdom | 056 | Independently reviewed kingdom slices |
| Reliability | 057 | Independently reviewed reliability slices |
| Validation | 058, 059 | Independently reviewed validation slices |
| Alpha | 060, 061 | Friends test and roadmap refinement |

Dependencies in detailed files govern dispatch. Task order is a current plan; new tasks are inserted using whole numbers.

## Insertion and status rules

Insert a task by renumbering later rows/files and updating links atomically. Preserve stable keys, status and any dispatched branch/prompt/commit history. Future branches use the current number. Small internal steps are checklists; separately prompted/reviewed/merged work gets its own whole number.

TODO -> WIP when one implementation prompt is dispatched. WIP remains through fixes/review/merge. DONE requires APPROVED, verified dev merge and required deployment verification. Missing evidence does not count as completion.

## Deferred future systems

Rune discovery, Conduit, Aspect, magical armour, Chaos/Order, heroes, additional kingdoms and naval/air/tunnel systems remain future design material, not active medieval implementation tasks. Production/PostgreSQL/master deployment gets its own explicit tasks when the owner chooses that release, not merely because the friends alpha starts.

## Village visual proof
The queue now contains 61 TODO tasks. VILLAGE-UI (011) is the early local layered-art proof, HOTSPOTS (033) is the saved admin editor, and VILLAGE-VISUAL-GATE (035) is the end-to-end acceptance gate after CONFIG (034). Later tasks renumbered with immutable keys preserved; none dispatched. VILLAGE-CREATE requires this gate. Internal proof steps are checklists, not letter-suffixed tasks.
Read [step-by-step preparation](docs/design/village-visual-prototype.md). Flat illustration hotspots alone do not prove interchangeable buildings or modular walls.

## Scene scope boundaries

011 (VILLAGE-UI), 033 (HOTSPOTS) and 035 (VILLAGE-VISUAL-GATE) cover village interiors only. 012 (MAP-UI) independently proves the world-map rendering and local square-editing workflow; 043 (DISTRICTS) later owns saved map definitions and authoritative ownership. Resource types, production-site visuals and resource-screen layout are not decided and must not be invented during either scene proof. 023 and 037 remain gated by an owner-approved economy direction. Shared rendering utilities may be reused, but acceptance evidence and content remain separate.

## Combat proof and integration

013 (COMBAT-UI) proves a static battalion board and final-result layout. 048 (COUNTERS) validates numerical rules; 049 (BATTLES) safely persists server outcomes; 050 (REPORTS) displays authorised final results and reconciled army totals. V1 requires no animation, replay, timeline or narrated combat log. Prefer React/CSS for the board; PixiJS is optional only if a justified requirement emerges. Later visual improvements get separately scoped tasks. Read [combat preparation](docs/design/combat-visual-prototype.md). All remain TODO.
