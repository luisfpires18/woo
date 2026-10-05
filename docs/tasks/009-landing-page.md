# 009: First landing page UI polish

Task key: LANDING
Status: TODO
Updated: 2026-10-05
Milestone: UI
Dependencies: DESIGN-SYSTEM

## Goal

Design and implement a clear WOO landing page using the approved visual direction.

## Acceptance criteria

Desktop/mobile page introduces the medieval game, offers clear entry actions and includes version footer. No claims of implemented future features; review actual screenshots/browser behavior. Placeholder art allowed. Follow the approved mockups and their mandatory written correction: no inline login, no Overview/Worlds/theme navbar links, profile menu Settings opens a separate page and never embeds appearance controls. All screenshot activity/world data remains illustrative.

## Required reading

- [docs/design/screens-and-themes.md](../../docs/design/screens-and-themes.md)
- [Approved landing mockup and corrections](../../mockup/00-landing-page.md)
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
