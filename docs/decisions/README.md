# Decision register

Updated: 2026-10-05. Status: design documentation; no game implementation yet.

This register records current user decisions. Feature documents explain details. Future choices should include date, status, reason and superseded rules.

| ID | Status | Decision |
|---|---|---|
| D-01 | Confirmed | Original WOO kingdom strategy; no independent alliances, fake attacks or manual waves |
| D-02 | Confirmed | Blacksmith player controls own settlement/troops; separate kingdom ruler |
| D-03 | Confirmed | Village capture transfers district to conquering kingdom |
| D-04 | Confirmed | Friends alpha: Arkazia, Veridor, Sylvara; four future kingdoms greyed out |
| D-05 | Confirmed direction | Drakanith neutral NPC volcanic outposts; Moraphys hostile NPC outposts; timing open |
| D-06 | Confirmed | Recruits have default weapons; smithing upgrades; swapping undecided |
| D-07 | Confirmed | Weaponsmith/Armorsmith categories and points-based blacksmith talent tree |
| D-08 | Confirmed future rule | L2 bearer blood bond; smith quality reduces crafting recoil to zero at best mastery |
| D-09 | Confirmed future direction | Magical light/medium/heavy equivalents; scenario-specific transformations |
| D-10 | Confirmed | Dark/light UI independent of kingdom palette |
| D-11 | Confirmed | Owner admin workspace and mandatory image slots with fallbacks |
| D-12 | Confirmed | Successful replacement/deletion cleans unreferenced images; no unused histories |
| D-13 | Confirmed | Simple fully medieval first version; advanced ideas evolve later |
| D-14 | Confirmed | This repository is the main ongoing specification |
| D-15 | Accepted direction | District map, useful duties, camp classes, BFME-inspired counters/stances/modules |
| D-16 | Open | Village limits, defeat, combat timings, equipment granularity and final victory |
| D-17 | Proposal | Weapon plus runeforged armour required for L2 |
| D-18 | Proposal | 10–12 week fantasy season arc and particular reset rules |
| D-19 | Proposal | React/PixiJS/.NET/PostgreSQL architecture |

All initial records dated 2026-10-05. Superseded: player as ruler-blacksmith; treating all future fantasy features as first-release requirements; retaining unused image versions. Neither a fixed victory condition nor an approved numeric balance exists.

## Workflow decisions, 2026-10-05

D-20 Confirmed: dev integration; feat/NNN branches from dev; master production at alpha release.
D-21 Confirmed: numbered TODO/WIP/DONE queue with detailed task files, insertion renumbering and preserved task identity.
D-22 Confirmed: ChatGPT prompt -> Claude local implementation/report -> ChatGPT APPROVED/NEEDS CHANGES/REFUSAL -> owner merge/push -> remote deployment verification.
D-23 Confirmed: Azure host, F1 preference and conditional temporary B1 path; subscription/quota eligibility must be checked. This supersedes VPS hosting.
D-24 Confirmed: one dependent command at a time; no guessed placeholders; short TLDR at task-response end.
D-25 Confirmed: skills/plugin health check before implementation; requested inventory not presumed installed.
D-26 Confirmed: shared prompt protocol can be refined in another chat; active tasks pin a specification revision.

Exact referenced LP AI WORK skill remains unavailable; an older reconstructed LP WORK was read as a reference, without claiming equivalence. Workflow is based on the user's current instructions.

## Roadmap refinement, 2026-10-05

D-27 Confirmed: SQLite local/dev; PostgreSQL only a future production consideration. This supersedes the current PostgreSQL dev recommendation.
D-28 Confirmed: broad evolving roadmap with small whole-number tasks; no 012a/012b. Insert and renumber while preserving stable task identity/history.
D-29 Confirmed: prioritise Claude setup -> defined tech/architecture -> skeleton/first page -> CI/Azure dev -> UI polish -> persistence-backed accounts/world/kingdom -> game slices. Do not make all medieval design decisions a prerequisite to first page.
D-30 Confirmed: visible app-version footer, exemplified by 0.0.1-dev; dependency versions defined early; image changes must account for caches.
D-31 Proposal: pragmatic DDD modular boundaries; exact structure chosen in STACK. DDD was a user example, not a complete adopted architecture.

Original 12 TODO task files were replaced by smaller specifications before any dispatch. No active task branch, completed implementation or release version exists.

D-32 Confirmed: include Cloudflare configuration and command-guided R2 bucket creation for admin-uploaded images, before upload integration. No provisioning requested in the roadmap update itself.

## Village artwork preparation, 2026-10-05
D-33 Confirmed: owner personally uploads images to Git from now on. Assistant maintains text documentation and previews; do not commit image files for the owner.
D-34 Confirmed preference: restrained textures; early dark-navigation / ivory-content / crimson-action game UI. Dark/light remains required. Latest village remake still awaits approval.
D-35 Accepted planning direction: make village visual feasibility explicit before broad asset production; document an early local proof plus later admin/storage integration gate.
D-36 Proposal to validate: fixed plots, layered terrain/buildings/walls/bridge/props, stable anchors and geometry; React/PixiJS responsibilities. This does not settle free-placement gameplay, exact camera/schema, budgets or final artwork.
Approved in conversation: centered full-background login and placeholder-based kingdom selector. Preview sample data is not balance canon.
