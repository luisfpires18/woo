# 050: Final battle results and army totals

Task key: REPORTS
Status: TODO
Updated: 2026-10-05
Milestone: Combat
Dependencies: BATTLES, COMBAT-UI

## Goal

Connect the accepted static battalion board and final result to persisted, server-authorised battle outcomes. Animation and replay are deferred beyond v1.

## Acceptance criteria

Display authorised composition, outcome, starting forces, survivors and approved loss categories with reconciled totals. No narrated log, event timeline or replay required. Counter/modifier details may be compact optional data only where recorded and useful; do not invent causal percentages.

The server redacts hidden opponent data before delivery. Persist sufficient frozen input/rules references and final outcomes for audit; render historical results without recalculating using today's stats. Define result-schema compatibility.

Test allowed/denied access, restricted results, missing/stale images, failed/retried requests and switching battles during loading. Viewing/refreshing a result sends no resolution command and applies no economic effect. Static image/symbol fallbacks, accessible lists, both themes and version footer work. Notification integration must not depend on the later OFFLINE task.

## Required reading

- [Combat visual proof](../design/combat-visual-prototype.md)
- [Replay contract and research](../technical/combat-research.md)

- [docs/design/artwork-pipeline.md](../../docs/design/artwork-pipeline.md)
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
