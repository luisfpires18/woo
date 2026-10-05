# 008: Content administration and image lifecycle

Task key: ADMIN-ASSETS
Status: TODO
Updated: 2026-10-05
Dependencies: ACCOUNTS

## Goal

Implement agreed configurable definitions and mandatory visual slots/fallbacks.

## Required reading

- [docs/admin.md](../../docs/admin.md)
- [docs/technical/asset-lifecycle.md](../../docs/technical/asset-lifecycle.md)
- [Prompt protocol](../workflow/implementation-protocol.md)

## Acceptance criteria

Replace/delete cleans unreferenced objects; failed replacement preserves old art; shared refs protected; invalid edits rejected.

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
