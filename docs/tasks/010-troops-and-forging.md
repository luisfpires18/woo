# 010: Default troops, medieval forging and talents

Task key: FORGING
Status: TODO
Updated: 2026-10-05
Dependencies: VILLAGE

## Goal

Recruit approved ordinary units and implement equipment upgrades/talent choices.

## Required reading

- [docs/gameplay/forging-and-progression.md](../../docs/gameplay/forging-and-progression.md)
- [docs/content/catalogues.md](../../docs/content/catalogues.md)
- [Prompt protocol](../workflow/implementation-protocol.md)

## Acceptance criteria

Default weapons and chosen upgrade granularity work; resources/points spent correctly; no runeforging introduced.

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
