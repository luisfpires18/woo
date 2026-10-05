# Working instructions

Updated: 2026-10-05. Status: design documentation; no game implementation yet.

Read README.md, docs/INDEX.md, docs/scope.md, docs/decisions/README.md and the relevant feature documents before changing this project.

- This is currently a documentation project. Do not scaffold or implement the game without a user request.
- Preserve the simple medieval first version. Future features are not launch requirements.
- Label statements as Confirmed, Accepted direction, Proposal, Open or Superseded when ambiguity matters. Accepted directions do not confirm numbers.
- The player is a blacksmith with control of their own villages and troops. They are not the kingdom ruler.
- No independent player alliances, fake attacks or manually synchronised attack waves.
- User-provided earlier game drafts are exploratory. Do not import their claimed stack, implemented status, stats or rosters as facts.
- Story canon and game adaptations are separate. Do not silently reconcile conflicting lore. User defines named characters and weapon assignments.
- Keep docs/open-questions.md and docs/decisions/README.md consistent with feature changes. Update docs/INDEX.md when adding files.
- GitHub is the primary specification. docs/reference/brainstorming-snapshot.md is historical and cannot override current documents or later user instructions.
- Before editing remotely, read the latest branch/files. Preserve other work; never force-push to resolve a conflict.
- Mandatory admin-controlled image slots and symbol/emoji fallbacks apply to all visual entities. Remove unreferenced storage objects after successful replacement; do not preserve unused image versions.
- Verify links and document consistency for documentation edits. When code is authorised, use appropriate tests and record actual validation. Never claim unrun tests or working features.
- Do not invent final formulas, costs, timings, APIs or database schemas where decisions remain open. Propose them explicitly.
- Do not commit credentials, real account passwords or personal player data. Assets will be configured later.

Feature documents should cover purpose, scope/status, player actions, requirements, costs/timing, outcomes, failure/recovery, admin configuration, dependencies, open questions and acceptance scenarios. Missing details are visible gaps, not permission to invent them.
