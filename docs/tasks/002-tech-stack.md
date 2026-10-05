# 002: Tech stack, dependency versions and architecture

Task key: STACK
Status: TODO
Updated: 2026-10-05
Milestone: Foundation
Dependencies: SKILLS

## Goal

Choose supported runtime/dependency versions, architecture boundaries and solution layout. Assess pragmatic DDD and a modular monolith.

## Acceptance criteria

Owner confirmed C# / ASP.NET Core and React/TypeScript on 2026-10-05. Preserve the backend comparison as selection rationale; assess supported versions and deployment compatibility against workload, account/admin integration, SQLite and Azure constraints. Record exact supported selected-backend/React/PixiJS/data-access versions, lockfile policy, candidate Domain/Application/Infrastructure/API boundaries and dependency direction. Confirm SQLite local/dev, Azure constraints and an initial 0.0.1-dev version contract. Avoid ceremonial aggregates or distributed services.

## Required reading

- [docs/technical/architecture.md](../../docs/technical/architecture.md)
- [Backend comparison](../technical/backend-comparison.md)
- [Hosting comparison](../technical/hosting-comparison.md)
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
