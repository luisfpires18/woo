# 012: Game UI structure and navigation polish

Task key: GAME-SHELL
Status: TODO
Updated: 2026-10-16
Milestone: UI
Dependencies: LANDING

## Goal

Build the reusable game layout with navigation, resource-header placeholders and page structure.

## Acceptance criteria

Village/map/forge/army/duties views reachable using explicit mock data. Layout works at mobile widths and both themes; shared version footer. No fake saved state or world membership.

## Required reading

- [docs/design/screens-and-themes.md](../../docs/design/screens-and-themes.md)
- [Scope](../scope.md)
- [Prompt protocol](../workflow/implementation-protocol.md)

## Implementation boundary

A separately reviewable planning task. ChatGPT inspects current code and resolves this task's blockers before issuing one pinned Claude prompt. Preserve medieval v1. Split/refine based on actual development; do not implement successors automatically.

## Validation and administration

Run proportionate checks against acceptance criteria, including relevant failure/retry scenarios. Report actual local evidence. ChatGPT checks remote merge/deployment where applicable. Content/images follow admin configuration and cleanup rules; new mechanics remain code changes.

## Dispatch history

1. 2026-10-16: dispatched as task 011 against dev 006be69e689180703c9acd3d826aab4476801d22 on branch `feat/011-game-shell`; local implementation commit d4f4bbe7f827aae208f611c514bcf7d84d84fd8a (game shell under `/game/`, screenshots supplied). Not reviewed, not pushed.
2. 2026-10-16: paused by the owner while the legacy title-screen migration (issue [#65](https://github.com/luisfpires18/woo/issues/65)) was inserted as [011 LEGACY-TITLE](011-legacy-title-screen.md); GAME-SHELL became 012 with its key and scope unchanged. The branch `feat/011-game-shell` is preserved unchanged at d4f4bbe for later review and rebase; its name is historical and is not reused. On resumption it must be rebased onto the LEGACY-TITLE result (the landing page moves to `/worlds`, and its own task records and roadmap edits were written before this renumbering).

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
