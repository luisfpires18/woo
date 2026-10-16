# 005: App version footer and build identity

Task key: VERSION
Status: WIP
Updated: 2026-10-16
Milestone: Foundation
Dependencies: SKELETON

## Goal

Implement the visible initial version, such as 0.0.1-dev, and a single build identity source.

## Acceptance criteria

Landing/game/admin layouts share a footer version. Backend and frontend identify the same build; deployed version can be compared with its commit. Define version bump policy and distinguish dependency versions from app release versions. Docs-only updates need not bump app version.

## Required reading

- [docs/technical/versioning-and-cache.md](../../docs/technical/versioning-and-cache.md)
- [Scope](../scope.md)
- [Prompt protocol](../workflow/implementation-protocol.md)

## Implementation boundary

A separately reviewable planning task. ChatGPT inspects current code and resolves this task's blockers before issuing one pinned Claude prompt. Preserve medieval v1. Split/refine based on actual development; do not implement successors automatically.

## Validation and administration

Run proportionate checks against acceptance criteria, including relevant failure/retry scenarios. Report actual local evidence. ChatGPT checks remote merge/deployment where applicable. Content/images follow admin configuration and cleanup rules; new mechanics remain code changes.

## Result

Implemented as recorded in [versioning and cache](../technical/versioning-and-cache.md#implementation-task-005), following the [stack contract](../technical/stack.md#build-identity-contract-001-dev): root `VERSION` (`0.0.1-dev`), coordinated `node scripts/build.mjs`, `GET /api/build`, the reusable `AppFooter` with the update notice, explicit cache headers and the [bump policy](../technical/versioning-and-cache.md#version-bump-policy). Commands are in the [README](../../README.md#run-locally). No CI, deployment, accounts, database, gameplay, PixiJS or image revision work.

## Execution record

Task-spec commit: 77e206cd47f5a444485ac0dfb3aba194ad059d7f (issue [#4](https://github.com/luisfpires18/woo/issues/4) dispatch)
Expected dev base: 77e206cd47f5a444485ac0dfb3aba194ad059d7f (verified equal to local dev and origin/dev, 2026-10-16)
Implementation branch: feat/005-app-version (local only)
Implementation commit: local commit on that branch; SHA given in the implementation report
Review: Pending ChatGPT review
Merged dev commit: None
Deployment: Not applicable to this task (deployment is DEV-DEPLOY)
Issues: no frontend unit-test runner exists, so footer/notice logic was verified in the browser, not by automated tests. Tab-visibility re-checks were triggered by dispatching the browser event, not by switching real tabs.

## TLDR

Next: ChatGPT review of the local commit.
Done: VERSION file, coordinated build, build-identity API, footer, update notice, cache headers and bump policy, with host tests and browser checks.
Issues: see execution record.
