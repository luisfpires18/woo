# 007: Deploy the first page automatically from dev

Task key: DEV-DEPLOY
Status: TODO
Updated: 2026-10-05
Milestone: Deployment
Dependencies: AZURE-SETUP

## Goal

Create the dev-only application deployment Action and validate a complete deployed round trip.

## Acceptance criteria

A dev change builds/tests, deploys its exact artifact and serves a healthy versioned page. Deployment history identifies commit/version; failures visible; no prod workflow. Demonstrate recovery without inventing Azure resource identifiers.

## Required reading

- [docs/technical/azure-and-deployment.md](../../docs/technical/azure-and-deployment.md)
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
