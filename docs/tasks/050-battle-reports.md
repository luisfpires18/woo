# 050: Battle reports and simple event replay

Task key: REPORTS
Status: TODO
Updated: 2026-10-05
Milestone: Combat
Dependencies: BATTLES, COMBAT-UI

## Goal

Connect the accepted COMBAT-UI representation and accessible reports to persisted, server-authorised battle events. Tokens remain the missing-art fallback.

## Acceptance criteria

Display composition, recorded interactions/modifiers, reconciled losses and final outcome. Do not invent precise causal percentages or expose hidden opponent information. Server projects an authorised report/event view for the requester; animation and downloaded assets cannot contain concealed data.

Replay uses stored events and presentation mappings, not today's ruleset to recalculate yesterday's battle. Old report/event schema compatibility is explicit. Use pause/speed/skip/restart and bounded seeking consistent with the proof; ending totals agree at every supported speed and frame rate. Reports remain immediately accessible without waiting for playback or successful art loading.

Test authorised and denied access, redacted reports, missing/stale assets, old versions, failed/retried fetches and switching battles during loading. Replaying/refreshing sends no resolution command and applies no economic effect. Optional notification integration must not depend on the later OFFLINE task to complete this slice. UI handles symbol fallback, reduced motion, accessible event list, both themes and version footer.

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
