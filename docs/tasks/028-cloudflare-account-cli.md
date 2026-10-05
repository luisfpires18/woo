# 028: Cloudflare account, access and CLI configuration

Task key: CLOUDFLARE
Status: TODO
Updated: 2026-10-05
Milestone: Assets
Dependencies: ROLES

## Goal

Prepare the owner's Cloudflare account and verified CLI/API access with a step-by-step command handoff.

## Acceptance criteria

Check account/R2 availability and current CLI version; authenticate through owner handoff; verify the selected account. Give one dependent concrete command at a time. Do not ask the owner to paste secrets or assume billing activation.

## Required reading

- [docs/technical/cloudflare-r2.md](../../docs/technical/cloudflare-r2.md)
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
