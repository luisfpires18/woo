# 002: Select tech stack and Azure-compatible layout

Task key: STACK
Status: TODO
Updated: 2026-10-05
Dependencies: SKILLS

## Goal

Choose supported frontend/backend/runtime/deployment versions and document architecture tradeoffs.

## Required reading

- [docs/technical/architecture.md](../../docs/technical/architecture.md)
- [docs/technical/azure-and-deployment.md](../../docs/technical/azure-and-deployment.md)
- [Prompt protocol](../workflow/implementation-protocol.md)

## Acceptance criteria

Decisions recorded; F1 runtime, database hosting and sleeping-worker implications addressed. No service provisioned as part of stack selection.

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
