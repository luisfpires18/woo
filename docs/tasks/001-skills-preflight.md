# 001: Skills and plugins health check

Task key: SKILLS
Status: TODO
Updated: 2026-10-05
Dependencies: None

## Goal

Audit the actual Claude environment, identify verified sources and install missing approved tools.

## Required reading

- [docs/workflow/skills-and-plugins.md](../../docs/workflow/skills-and-plugins.md)
- [Prompt protocol](../workflow/implementation-protocol.md)

## Acceptance criteria

Inventory records sources, versions and smoke checks; missing sources and access are explicit. Do not claim installation from repository names alone.

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
