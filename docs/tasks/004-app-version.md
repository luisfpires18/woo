# 004: App version footer and build identity

Task key: VERSION
Status: TODO
Updated: 2026-10-05
Milestone: Foundation
Dependencies: SKELETON

## Goal

Implement the visible initial version, such as 0.0.1-dev, and a single build identity source.

## Acceptance criteria

Landing/game/admin layouts share a footer version. Backend and frontend identify the same build; deployed version can be compared with its commit. Define version bump policy and distinguish dependency versions from app release versions. Docs-only updates need not bump app version.

## Required reading

- [docs/technical/versioning-and-cache.md](../../docs/technical/versioning-and-cache.md)
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
