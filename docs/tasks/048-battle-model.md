# 048: Unit counters and battle-model validation

Task key: COUNTERS
Status: TODO
Updated: 2026-10-05
Milestone: Combat
Dependencies: RECRUIT, FORGING, SCOUTING, COMBAT-UI

## Goal

Specify/test numerical combat with readable role, armour and terrain interactions.

## Acceptance criteria

Choose and record targeting, phase/order rules, simultaneous-versus-sequential damage, counters, armour/terrain interactions, equipment effects, rounding and casualties. Only approved mechanics are numerical rules; morale, ammunition, retreat and RNG are unresolved until selected.

Use a pure bounded resolver with immutable inputs, ruleset/content revisions and stable ordering. Prefer a small regiment/phase model before individual-soldier physics. If RNG is approved, pin its algorithm/seed and test distributions; seed alone does not guarantee reproducibility. Frame time and sprite collisions cannot affect results.

Validate cavalry/pikes/mixed compositions, exposed ranged units, equipment upgrades, side-swapped equal forces, force-size extremes, zero forces, overkill and configured termination. Define equal-headcount versus equal-cost comparisons after costs are approved. Prevent negative counts, duplicate troop commitments and accidental first-side advantage. Test deterministic reruns, rounding thresholds and bounded runtime; preserve failing seeds where relevant.

Outputs include observed interaction/modifier events and reconciled totals so reports can explain facts. Do not claim a stance caused a win without supporting evidence. No assumed universal immunity or one unstoppable army; final roster/faction balance belongs to BALANCE. No full kingdom battle engine here.

## Required reading

- [Combat model and replay research](../technical/combat-research.md)

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
