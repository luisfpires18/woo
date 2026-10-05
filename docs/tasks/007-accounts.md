# 007: Login, kingdom selection and owner admin access

Task key: ACCOUNTS
Status: TODO
Updated: 2026-10-05
Dependencies: AZURE-DEV

## Goal

Implement account/auth roles and alpha faction availability.

## Required reading

- [docs/admin.md](../../docs/admin.md)
- [docs/world/kingdoms-and-geography.md](../../docs/world/kingdoms-and-geography.md)
- [Prompt protocol](../workflow/implementation-protocol.md)

## Acceptance criteria

Owner can access admin; normal players denied; Arkazia/Veridor/Sylvara selectable and later factions disabled.

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
