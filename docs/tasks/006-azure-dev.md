# 006: Azure dev setup and deployment Action

Task key: AZURE-DEV
Status: TODO
Updated: 2026-10-05
Dependencies: FOUNDATION

## Goal

Set up verified dev resources and a real dev-branch build/deploy pipeline.

## Required reading

- [docs/technical/azure-and-deployment.md](../../docs/technical/azure-and-deployment.md)
- [Prompt protocol](../workflow/implementation-protocol.md)

## Acceptance criteria

Dev artifact deploys and health check passes; quota handling reported; F1/B1 SKU verified; no prod workflow. Commands handed off step by step.

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
