# 004: Database model and persistence design

Task key: DATABASE
Status: TODO
Updated: 2026-10-05
Dependencies: RULES

## Goal

Define schema, ownership, scheduled orders, concurrency and migration/backup strategy.

## Required reading

- [docs/technical/data-model.md](../../docs/technical/data-model.md)
- [Prompt protocol](../workflow/implementation-protocol.md)

## Acceptance criteria

Equipment granularity follows agreed rules; atomic spending/idempotent completion specified; storage provider chosen. Actual migrations belong to an authorised implementation prompt.

## Scope boundaries

Only this task's goal. Preserve the medieval first release. Missing decisions require refinement before implementation. This document is a planning specification, not an instruction to start coding now.

## Validation

ChatGPT defines focused checks in the dispatch prompt after inspecting current code/state. Claude reports actual local checks; ChatGPT reviews evidence and checks remote merge/pipeline where available.

## Execution record

Task-spec commit: Not dispatched
Expected dev base: Not dispatched
Implementation branch: Not created
Implementation commit: None
Review: Pending
Merged dev commit: None
Deployment: Not started; applicability decided in prompt
Issues: Dependencies and detailed implementation requirements pending.

## TLDR

Next: refine prerequisites, then issue one prompt when authorised.
Done: task specification skeleton only.
Issues: implementation and validation have not begun.
