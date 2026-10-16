# Task roadmap

Updated: 2026-10-16. 001 SKILLS to 012 GAME-SHELL are DONE, 013 VILLAGE-UI is WIP and all other tasks remain TODO. MCP-SETUP was approved by ChatGPT, verified on remote dev at 99abf03 with documentation CI passing, and issue #64 is closed. STACK was first dispatched as 002 (issue #2) against e8703aa on feat/002-tech-stack, paused on 2026-10-15 for MCP-SETUP, resumed as 003 against base 99abf03 on feat/003-tech-stack, then approved and merged at 588c5a0 with documentation CI passing; see its [dispatch history](docs/tasks/003-tech-stack.md#dispatch-history). SKELETON (issue #3) was approved on 2026-10-16 and merged at 77e206c with documentation CI passing. VERSION (issue #4) was approved on 2026-10-16 and merged at dc407cf with documentation CI passing. BUILD-CI (issue #5) was approved on 2026-10-16 and merged at 0e69cad; its first Ubuntu application and documentation runs passed. AZURE-SETUP was approved on 2026-10-16 and merged at dbb07cf after provisioning the West Europe F1 dev app. DEV-DEPLOY was approved on 2026-10-16, merged at 44d63f4 and deployed successfully in run 38891810702. DESIGN-SYSTEM was approved on 2026-10-16, merged at ed1f3c2 and deployed in run 38900180639. LANDING was approved on 2026-10-16, merged at 006be69 and deployed in run 38906266025. GAME-SHELL was dispatched as 011 against 006be69 on feat/011-game-shell and delivered local commit d4f4bbe, then paused on 2026-10-16 for the owner-requested legacy title-screen migration (issue #65), inserted as 011 LEGACY-TITLE; GAME-SHELL became 012 with its key, scope and preserved branch unchanged (see its [dispatch history](docs/tasks/012-game-ui-shell.md#dispatch-history)). LEGACY-TITLE was approved on 2026-10-16, merged and deployed at 2e295cf, and the old weapons-of-order staging resource group was deleted. GAME-SHELL resumed on 2026-10-16 against 2e295cf on its preserved branch, was approved, merged at 8a743e4 and deployed. VILLAGE-UI is dispatched against 8a743e4 on feat/013-village-ui. This expanded sequence replaced the original 12 broad tasks before any dispatch, so no branch history was affected.

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
| 007 | AZURE-SETUP | Azure dev provisioning through step-by-step commands | DONE | [007-azure-dev-setup.md](docs/tasks/007-azure-dev-setup.md) |
| 008 | DEV-DEPLOY | Deploy the first page automatically from dev | DONE | [008-dev-deployment.md](docs/tasks/008-dev-deployment.md) |
| 009 | DESIGN-SYSTEM | Shared UI foundations, themes and kingdom colours | DONE | [009-design-system.md](docs/tasks/009-design-system.md) |
| 010 | LANDING | First landing page UI polish | DONE | [010-landing-page.md](docs/tasks/010-landing-page.md) |
| 011 | LEGACY-TITLE | Port legacy title screen, then retire old staging Azure | DONE | [011-legacy-title-screen.md](docs/tasks/011-legacy-title-screen.md) |
| 012 | GAME-SHELL | Game UI structure and navigation polish | DONE | [012-game-ui-shell.md](docs/tasks/012-game-ui-shell.md) |
| 013 | VILLAGE-UI | Village scene and interaction prototype | WIP | [013-village-ui.md](docs/tasks/013-village-ui.md) |
| 014 | MAP-UI | World map visual and editing proof | TODO | [014-map-ui.md](docs/tasks/014-map-ui.md) |
| 015 | COMBAT-UI | Static battalion board and result proof | TODO | [015-combat-ui.md](docs/tasks/015-combat-ui.md) |
| 016 | FORGE-UI | Forge, equipment and talent UI prototypes | TODO | [016-forge-ui.md](docs/tasks/016-forge-ui.md) |
| 017 | ADMIN-UI | Admin dashboard UI and content-editor structure | TODO | [017-admin-dashboard-ui.md](docs/tasks/017-admin-dashboard-ui.md) |
| 018 | AUTH-UI | Login, world selection and kingdom-choice UI | TODO | [018-login-onboarding-ui.md](docs/tasks/018-login-onboarding-ui.md) |
| 019 | SQLITE | SQLite persistence and migration foundation | TODO | [019-sqlite-foundation.md](docs/tasks/019-sqlite-foundation.md) |
| 020 | DB-RECOVERY | Dev database persistence, backup and restore | TODO | [020-sqlite-backup-recovery.md](docs/tasks/020-sqlite-backup-recovery.md) |
| 021 | IDENTITY | Real accounts, sessions and login | TODO | [021-accounts-login.md](docs/tasks/021-accounts-login.md) |
| 022 | ROLES | Owner bootstrap and admin authorisation | TODO | [022-owner-admin-access.md](docs/tasks/022-owner-admin-access.md) |
| 023 | WORLDS | World creation, availability and alpha factions | TODO | [023-worlds-admin.md](docs/tasks/023-worlds-admin.md) |
| 024 | MEMBERSHIP | Associate player with World and chosen kingdom | TODO | [024-player-world-kingdom.md](docs/tasks/024-player-world-kingdom.md) |
| 025 | RESOURCES-ADMIN | Admin resource definitions | TODO | [025-resource-definitions.md](docs/tasks/025-resource-definitions.md) |
| 026 | BUILDINGS-ADMIN | Admin building definitions and prerequisites | TODO | [026-building-definitions.md](docs/tasks/026-building-definitions.md) |
| 027 | UNITS-ADMIN | Admin unit definitions and default equipment | TODO | [027-unit-definitions.md](docs/tasks/027-unit-definitions.md) |
| 028 | EQUIPMENT-ADMIN | Admin equipment and ordinary forge recipes | TODO | [028-equipment-recipes.md](docs/tasks/028-equipment-recipes.md) |
| 029 | TALENTS-ADMIN | Admin talent tree definitions | TODO | [029-talent-definitions.md](docs/tasks/029-talent-definitions.md) |
| 030 | CLOUDFLARE | Cloudflare account, access and CLI configuration | TODO | [030-cloudflare-account-cli.md](docs/tasks/030-cloudflare-account-cli.md) |
| 031 | R2-SETUP | Create and verify the R2 dev asset bucket | TODO | [031-r2-bucket-creation.md](docs/tasks/031-r2-bucket-creation.md) |
| 032 | R2-CONFIG | R2 upload credentials, CORS and image delivery | TODO | [032-r2-upload-delivery.md](docs/tasks/032-r2-upload-delivery.md) |
| 033 | ASSETS | Entity image upload, replacement and cleanup | TODO | [033-asset-upload-cleanup.md](docs/tasks/033-asset-upload-cleanup.md) |
| 034 | IMAGE-CACHE | Image cache revisions and immediate replacements | TODO | [034-image-cache-revisions.md](docs/tasks/034-image-cache-revisions.md) |
| 035 | HOTSPOTS | Admin village anchors and clickable hotspots | TODO | [035-scene-hotspot-editor.md](docs/tasks/035-scene-hotspot-editor.md) |
| 036 | CONFIG | Configuration publish, audit and active-order policy | TODO | [036-content-publish-policy.md](docs/tasks/036-content-publish-policy.md) |
| 037 | VILLAGE-VISUAL-GATE | Layered village art and admin integration gate | TODO | [037-village-visual-gate.md](docs/tasks/037-village-visual-gate.md) |
| 038 | VILLAGE-CREATE | Initial village placement and player settlement | TODO | [038-initial-village.md](docs/tasks/038-initial-village.md) |
| 039 | RESOURCE-STATE | Resource production, storage and offline accrual | TODO | [039-resource-production.md](docs/tasks/039-resource-production.md) |
| 040 | ORDERS | Durable action queues and catch-up processing | TODO | [040-durable-action-queues.md](docs/tasks/040-durable-action-queues.md) |
| 041 | CONSTRUCTION | Building construction, upgrades and queues | TODO | [041-building-construction.md](docs/tasks/041-building-construction.md) |
| 042 | RECRUIT | Recruitment and troop inventory | TODO | [042-troop-recruitment.md](docs/tasks/042-troop-recruitment.md) |
| 043 | FORGING | Ordinary weaponsmith and armorsmith upgrades | TODO | [043-medieval-forging.md](docs/tasks/043-medieval-forging.md) |
| 044 | TALENTS | Blacksmith experience, points and talent choices | TODO | [044-blacksmith-progression.md](docs/tasks/044-blacksmith-progression.md) |
| 045 | DISTRICTS | Authoritative map districts and ownership | TODO | [045-world-districts.md](docs/tasks/045-world-districts.md) |
| 046 | ROUTES | Route connectivity and travel calculations | TODO | [046-routes-and-travel.md](docs/tasks/046-routes-and-travel.md) |
| 047 | SCOUTING | Scouting and map information visibility | TODO | [047-scouting-and-visibility.md](docs/tasks/047-scouting-and-visibility.md) |
| 048 | MOVEMENT | Troop movement, stationing and return | TODO | [048-troop-movement.md](docs/tasks/048-troop-movement.md) |
| 049 | DUTIES | Useful kingdom duties | TODO | [049-kingdom-duties.md](docs/tasks/049-kingdom-duties.md) |
| 050 | COUNTERS | Unit counters and battle-model validation | TODO | [050-battle-model.md](docs/tasks/050-battle-model.md) |
| 051 | BATTLES | Persisted automatic battle resolution | TODO | [051-battle-resolution.md](docs/tasks/051-battle-resolution.md) |
| 052 | REPORTS | Final battle results and army totals | TODO | [052-battle-reports.md](docs/tasks/052-battle-reports.md) |
| 053 | CAMP | Ordinary animal camps and expeditions | TODO | [053-ordinary-camps.md](docs/tasks/053-ordinary-camps.md) |
| 054 | CAMPAIGNS | Declared kingdom campaigns and defence commitments | TODO | [054-kingdom-campaigns.md](docs/tasks/054-kingdom-campaigns.md) |
| 055 | SIEGE | Village siege and defensive building effects | TODO | [055-village-siege.md](docs/tasks/055-village-siege.md) |
| 056 | CAPTURE | Capture, occupation and district transfer | TODO | [056-territory-capture.md](docs/tasks/056-territory-capture.md) |
| 057 | RECOVERY | Defeat, wounded troops and player recovery | TODO | [057-defeat-and-recovery.md](docs/tasks/057-defeat-and-recovery.md) |
| 058 | COOP | Kingdom overview, defence requests and contribution credit | TODO | [058-kingdom-cooperation.md](docs/tasks/058-kingdom-cooperation.md) |
| 059 | PROFILES | Player and settlement profiles | TODO | [059-player-settlement-profiles.md](docs/tasks/059-player-settlement-profiles.md) |
| 060 | LEADERBOARDS | World and season leaderboards | TODO | [060-world-leaderboards.md](docs/tasks/060-world-leaderboards.md) |
| 061 | OFFLINE | Offline reconciliation and live UI notifications | TODO | [061-offline-ui-updates.md](docs/tasks/061-offline-ui-updates.md) |
| 062 | BALANCE | Economy, counters and faction balance pass | TODO | [062-economy-and-faction-validation.md](docs/tasks/062-economy-and-faction-validation.md) |
| 063 | QA | Mobile, accessibility, persistence and permission validation | TODO | [063-mobile-and-reliability.md](docs/tasks/063-mobile-and-reliability.md) |
| 064 | FRIENDS | Invite-only friends alpha setup and support | TODO | [064-friends-alpha-setup.md](docs/tasks/064-friends-alpha-setup.md) |
| 065 | ALPHA-FEEDBACK | Run medieval alpha and refine next tasks | TODO | [065-alpha-feedback.md](docs/tasks/065-alpha-feedback.md) |

## Milestones

| Milestone | Tasks | Outcome |
|---|---|---|
| Foundation | 001, 002, 003, 004, 005, 006 | Running versioned skeleton |
| Deployment | 007, 008 | First page deployed to dev |
| UI | 009, 010, 011, 012, 013, 014, 015, 016, 017, 018 | Polished screen prototypes |
| Persistence | 019, 020 | Durable SQLite local/dev |
| Accounts | 021, 022, 023, 024 | Identity and world membership |
| Content | 025, 026, 027, 028, 029, 036 | Independently reviewed content slices |
| Assets | 030, 031, 032, 033, 034, 035, 037 | Cloudflare/R2, upload, cache and hotspots |
| Village | 038, 039, 040, 041 | Independently reviewed village slices |
| Military | 042, 043, 044 | Independently reviewed military slices |
| Map | 045, 046, 047, 048 | Independently reviewed map slices |
| Activities | 049, 053 | Independently reviewed activities slices |
| Combat | 050, 051, 052 | Independently reviewed combat slices |
| Conquest | 054, 055, 056, 057 | Independently reviewed conquest slices |
| Kingdom | 058 | Independently reviewed kingdom slices |
| Community | 059, 060 | Player/settlement inspection and world-season rankings |
| Reliability | 061 | Independently reviewed reliability slices |
| Validation | 062, 063 | Independently reviewed validation slices |
| Alpha | 064, 065 | Friends test and roadmap refinement |

Dependencies in detailed files govern dispatch. Task order is a current plan; new tasks are inserted using whole numbers.

## Insertion and status rules

Insert a task by renumbering later rows/files and updating links atomically. Preserve stable keys, status and any dispatched branch/prompt/commit history. Future branches use the current number. Small internal steps are checklists; separately prompted/reviewed/merged work gets its own whole number.

TODO -> WIP when one implementation prompt is dispatched. WIP remains through fixes/review/merge. A dispatched task paused for an inserted task stays WIP with a `Paused:` line in its detail file; only one unpaused task may be WIP (checked by scripts/check_docs.py). DONE requires APPROVED, verified dev merge and required deployment verification. Missing evidence does not count as completion.

## Deferred future systems

Rune discovery, Conduit, Aspect, magical armour, Chaos/Order, heroes, additional kingdoms and naval/air/tunnel systems remain future design material, not active medieval implementation tasks. Production/PostgreSQL/master deployment gets its own explicit tasks when the owner chooses that release, not merely because the friends alpha starts.

## Village visual proof
The queue contains 65 tasks: 001 SKILLS to 012 GAME-SHELL are DONE, 013 VILLAGE-UI is WIP, the rest TODO. VILLAGE-UI (013) is the early local layered-art proof, HOTSPOTS (035) is the saved admin editor, and VILLAGE-VISUAL-GATE (037) is the end-to-end acceptance gate after CONFIG (036). Later tasks renumbered with immutable keys preserved; none had been dispatched. VILLAGE-CREATE requires this gate. Internal proof steps are checklists, not letter-suffixed tasks.
Read [step-by-step preparation](docs/design/village-visual-prototype.md). Flat illustration hotspots alone do not prove interchangeable buildings or modular walls.

## Scene scope boundaries

013 (VILLAGE-UI), 035 (HOTSPOTS) and 037 (VILLAGE-VISUAL-GATE) cover village interiors only. 014 (MAP-UI) independently proves the world-map rendering and local square-editing workflow; 045 (DISTRICTS) later owns saved map definitions and authoritative ownership. Resource types, production-site visuals and resource-screen layout are not decided and must not be invented during either scene proof. 025 and 039 remain gated by an owner-approved economy direction. Shared rendering utilities may be reused, but acceptance evidence and content remain separate.

## Combat proof and integration

015 (COMBAT-UI) proves a static battalion board and final-result layout. 050 (COUNTERS) validates numerical rules; 051 (BATTLES) safely persists server outcomes; 052 (REPORTS) displays authorised final results and reconciled army totals. V1 requires no animation, replay, timeline or narrated combat log. Prefer React/CSS for the board; PixiJS is optional only if a justified requirement emerges. Later visual improvements get separately scoped tasks. Read [combat preparation](docs/design/combat-visual-prototype.md). All remain TODO.

## Settlement terminology and community views

Settlement/settlements is the preferred player-facing term. Existing village filenames, VILLAGE keys and historical references remain technical aliases until a deliberate code/data naming decision. 059 (PROFILES) adds player and settlement pages; 060 (LEADERBOARDS) adds world/season rankings. Later tasks renumbered atomically, stable keys preserved. Ranking categories/formulas and profile visibility need approval before dispatch. Dev is the desired default branch; actual repository setting remains pending verification.
