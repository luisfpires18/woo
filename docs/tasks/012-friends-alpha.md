# 012: Friends alpha and production-release gate

Task key: ALPHA
Status: TODO
Updated: 2026-10-05
Dependencies: FRONTIER

## Goal

Run medieval alpha validation and decide readiness for production promotion.

## Required reading

- [docs/roadmap.md](../../docs/roadmap.md)
- [docs/gameplay/seasons-and-victory.md](../../docs/gameplay/seasons-and-victory.md)
- [Prompt protocol](../workflow/implementation-protocol.md)

## Acceptance criteria

Findings logged and blockers addressed; owner approves alpha release before master promotion/prod Action. No automatic public launch.

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
