# 052: Village siege and defensive building effects

Task key: SIEGE
Status: TODO
Updated: 2026-10-05
Milestone: Conquest
Dependencies: CAMPAIGNS, CONSTRUCTION

## Goal

Resolve contested villages with agreed fortress/main-building rules.

## Acceptance criteria

Choose win conditions, module effects and infrastructure damage; defence works offline within timing policy. Main building's pillar role explicit, not automatic total deletion.

## Required reading

- [docs/gameplay/player-and-villages.md](../../docs/gameplay/player-and-villages.md)
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
