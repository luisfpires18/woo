# 013: Combat formation and replay visual proof

Task key: COMBAT-UI
Status: TODO
Updated: 2026-10-05
Milestone: UI
Dependencies: GAME-SHELL

## Goal

Prove the preferred combat layout: army rosters beside a central illustrated battlefield, readable regiment formations and a short explanatory replay. Own this proof independently from map, village and resources.

## Acceptance criteria

Follow steps 1–9 of the combat visual guide with explicit local fixtures. Four synthetic medieval roles only: infantry, archers, cavalry and pikemen. These are test roles, not approved kingdom rosters. Pre-battle placement and stance controls demonstrate accepted direction; exact slots, effects, timing and casualty rules remain proposals.

Prove rosters, frontline/backline/flank role readability, regiment selection/inspector and an illustrated battlefield. Start with symbols, then independently replace representative art. Use a bounded number of visible figures per regiment with authoritative-looking counts labelled as mock state. No one-sprite-per-soldier requirement.

Play a schema-versioned recorded fixture through pause/resume, replay speed, restart, skip-to-report and a bounded phase seek. Test out-of-order/invalid fixture rejection, skipped time, repeated playback, reduced motion, missing art and renderer disposal. Outcome and recorded casualties cannot change with FPS, playback speed, seeks or asset changes. A slow or hidden tab must not trigger an unbounded catch-up loop.

Keep event-derived explanations and totals consistent between battlefield, event list and report. Counter messages appear only when the fixture records the relevant interaction; no fabricated causal claim or fictitious percentage attribution. Unknown opponent information is clearly marked in a restricted-view fixture.

Check mouse/touch/keyboard navigation, responsive roster/inspector, both themes, labels/icons beyond colour, no forced camera shake, accessible report and app-version footer. Record actual sprite camera/anchor consistency, attempts/repairs, density and named-device loading/frame/memory measurements. Set budgets and owner-approved visual reference before dispatch.

Provide owner-reviewed captures of pikes meeting a charge, ranged support and an exposed flank. These illustrate playback only, not a balanced combat engine. Static cutouts with movement/effects are a valid comparison; missing animation art leaves the animation-quality gate pending. No live micro, fake attacks, manual waves, saved orders, server-authoritative resolution, map conquest, siege, heroes, runes or all-kingdom art production.

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
