# 003: Tech stack, dependency versions and architecture

Task key: STACK
Status: WIP
Updated: 2026-10-15
Milestone: Foundation
Dependencies: SKILLS, MCP-SETUP

## Goal

Choose supported runtime/dependency versions, architecture boundaries and solution layout. Assess pragmatic DDD and a modular monolith.

## Acceptance criteria

Owner confirmed C# / ASP.NET Core and React/TypeScript on 2026-10-05. Preserve the backend comparison as selection rationale; assess supported versions and deployment compatibility against workload, account/admin integration, SQLite and Azure constraints. Record exact supported selected-backend/React/PixiJS/data-access versions, lockfile policy, candidate Domain/Application/Infrastructure/API boundaries and dependency direction. Confirm SQLite local/dev, Azure constraints and an initial 0.0.1-dev version contract. Avoid ceremonial aggregates or distributed services.

## Result

Decision record: [Tech stack, versions and architecture](../technical/stack.md). Documentation only; no scaffold, dependency installation or provisioning.

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

Task-spec commit: 99abf03422ad06f8ed64da75606c2bf0151e5cbf (this file at the resumed base; last edited in 854489d)
Expected dev base: 99abf03422ad06f8ed64da75606c2bf0151e5cbf (verified equal to local dev and origin/dev, 2026-10-15)
Implementation branch: feat/003-tech-stack (local only)
Implementation commit: local commit on that branch; SHA given in the implementation report
Review: Pending ChatGPT review
Merged dev commit: None
Deployment: Not applicable (documentation only)
Issues: subscription eligibility and Azure SQLite storage compatibility remain verification gates; see [stack](../technical/stack.md#azure-application-and-local-development). GitHub issue [#2](https://github.com/luisfpires18/woo/issues/2) still links the old 002 filename and needs updating by ChatGPT.

## Dispatch history

1. 2026-10-14: dispatched as task 002 (issue [#2](https://github.com/luisfpires18/woo/issues/2), `[002] STACK`) against dev e8703aa0c882c642dac8df0ad33cd6921f102ce7, branch `feat/002-tech-stack`. No implementation report was received and no such branch exists locally or on origin (checked 2026-10-15).
2. 2026-10-15: paused by the owner while [MCP-SETUP](002-mcp-setup.md) was inserted as 002; STACK became 003 with its key and scope unchanged.
3. 2026-10-15: resumed as task 003 against the new pinned base 99abf03 on branch `feat/003-tech-stack`, after the MCP-SETUP commits (854489d, 99abf03) were on origin/dev. The historical name `feat/002-tech-stack` stays reserved and was not reused.

## TLDR

Next: ChatGPT review of the local commit.
Done: versions, pinning, layout, boundaries, Azure route, gates and build identity contract documented.
Issues: two Azure verification gates open; local toolchain upgrades required before SKELETON.
