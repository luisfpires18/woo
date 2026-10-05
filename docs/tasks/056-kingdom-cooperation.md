# 056: Kingdom overview, defence requests and contribution credit

Task key: COOP
Status: TODO
Updated: 2026-10-05
Milestone: Kingdom
Dependencies: CAMPAIGNS, DUTIES, RECOVERY

## Goal

Make faction cooperation understandable without an alliance hierarchy.

## Acceptance criteria

Own/allied visibility respected; actionable requests, contributions and heroes-unavailable states clear. No one can seize others' forces. Extra social/chat features need explicit scope.

## Required reading

- [docs/vision.md](../../docs/vision.md)
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
