# 040: Ordinary weaponsmith and armorsmith upgrades

Task key: FORGING
Status: TODO
Updated: 2026-10-05
Milestone: Military
Dependencies: RECRUIT, EQUIPMENT-ADMIN

## Goal

Connect forge UI to actual medieval upgrades.

## Acceptance criteria

Agreed per-type/batch/individual rule implemented; costs/availability and compatibility enforced; define new-recruit inheritance and use during crafting. No generic unrestricted swaps.

## Required reading

- [docs/gameplay/forging-and-progression.md](../../docs/gameplay/forging-and-progression.md)
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
