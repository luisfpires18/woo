# 043: Authoritative map districts and ownership

Task key: DISTRICTS
Status: TODO
Updated: 2026-10-05
Milestone: Map
Dependencies: VILLAGE-CREATE, MAP-UI, ASSETS, CONFIG

## Goal

Persist approved world-map definitions, district topology and political ownership. Own protected map authoring/publication integration separately from the village editor.

## Acceptance criteria

Village-centred conquest geometry independent of illustration; three playable homelands visible. Canon geography respected; exact coordinates, cell/district relationship and topology supplied/approved.

Refine the dispatch into bounded implementation slices if saved map editing and world-state integration cannot be reviewed together; use whole-number tasks with stable keys. No village-interior or resource-site editing here.

Build on the accepted MAP-UI terrain approach. Protected map drafts can save/reload geometry and asset references, preview with the player renderer, validate IDs/edges/crossings/districts and publish one atomic revision. Handle concurrent edits and rejected imports; failed publication keeps the last good definition. Track references and cache revisions under ASSETS/CONFIG.

Live-season geometry changes require an explicit compatibility policy for villages, ownership and travel. Never silently remap or delete occupied cells/districts. Route/pathfinding rules remain ROUTES, not inferred from painted rivers or roads.

## Required reading

- [World-map research and saved-editing boundary](../technical/world-map-research.md)
- [Map scene contract](../technical/map-scene.md)

- [docs/world/kingdoms-and-geography.md](../../docs/world/kingdoms-and-geography.md)
- [Scope](../scope.md)
- [Prompt protocol](../workflow/implementation-protocol.md)

## Implementation boundary

A separately reviewable planning task. ChatGPT inspects current code and resolves this task's blockers before issuing one pinned Claude prompt. Preserve medieval v1. Split/refine based on actual development; do not implement successors automatically.

## Validation and administration

Run proportionate checks against acceptance criteria, including relevant failure/retry scenarios. Report actual local evidence. ChatGPT checks remote merge/deployment where applicable. Content/images follow admin configuration and cleanup rules; new mechanics remain code changes.

## Execution record

Task-spec commit: Not dispatched
Expected dev base: Not dispatched
Implementation branch: Not created
Implementation commit: None
Review: Pending
Merged dev commit: None
Deployment: Not started; applicability defined in prompt
Issues: Relevant decisions and implementation details must be resolved before dispatch.

## TLDR

Next: prepare when prerequisites are ready.
Done: task recorded, no implementation.
Issues: acceptance checks not yet run.
