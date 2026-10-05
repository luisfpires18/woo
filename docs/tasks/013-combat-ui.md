# 013: Static battalion board and result proof

Task key: COMBAT-UI
Status: TODO
Updated: 2026-10-05
Milestone: UI
Dependencies: GAME-SHELL

## Goal

Prove a simple static battalion board and final result for v1. No animation, replay, timeline, narrated combat log or event-playback controls are required. Keep this separate from village, map and resources.

## Acceptance criteria

Follow the combat visual guide using local mock state. Define a readable board with frontline, backline and flank slots as a provisional layout. Place representative infantry, archer, cavalry and pikeman battalion cards/icons with name, role, count, side and relevant equipment. These generic fixtures do not define kingdom rosters or counter formulas.

Show pre-battle placement/stance controls only where agreed, with clearly different editable versus committed/read-only states. Display a static opposing arrangement; unknown enemy data stays unknown. Select battalions through the board or accessible list and inspect their details.

Use a separate final-result state with winner/draw as supported by the chosen rules, starting forces, survivors and approved loss categories. No requirement to watch a battle or read narrative text. Reconcile every displayed total; loss types and formulas remain pending. Mock results are explicitly labelled and do not claim live authority.

Use independently replaceable static images with symbol fallbacks. Test mouse/touch/keyboard, narrow screens, both themes, readable labels beyond colour, missing images and version footer. Record actual screenshots, layout/asset consistency and basic load/interaction measurements on named devices. A polished board does not need a PixiJS renderer: use React/CSS unless an agreed requirement justifies canvas.

No animation/spritesheet production, replay controls, event timeline, live micro, fake attacks, manual waves, persisted orders, server resolution, siege, conquest, heroes or runes. Obtain owner review of the board and result before expansion. Final balance belongs to COUNTERS/BALANCE.

## Required reading

- [Combat visual proof](../design/combat-visual-prototype.md)
- [Combat asset production](../design/combat-art-production.md)
- [Combat simulation and replay research](../technical/combat-research.md)

- [docs/gameplay/combat-and-conquest.md](../../docs/gameplay/combat-and-conquest.md)
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
