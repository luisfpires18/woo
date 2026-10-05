# 030: R2 upload credentials, CORS and image delivery

Task key: R2-CONFIG
Status: TODO
Updated: 2026-10-05
Milestone: Assets
Dependencies: R2-SETUP

## Goal

Configure server upload/delete access, selected delivery URLs and applicable browser/CDN settings.

## Acceptance criteria

Scoped credentials stay server-side in local/Azure configuration; CORS uses verified local/dev origins where needed. Distinguish S3 API endpoint from public/custom-domain delivery and test put/read/delete while removing test objects. Provide concrete commands based on previous outputs.

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
