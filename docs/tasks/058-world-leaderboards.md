# 058: World and season leaderboards

Task key: LEADERBOARDS
Status: TODO
Updated: 2026-10-05
Milestone: Community
Dependencies: PROFILES, BATTLES, CAPTURE, COOP, CONFIG

## Goal

Provide world/season rankings for players and kingdoms, inspired by familiar browser-strategy navigation but using WOO-specific approved metrics.

## Acceptance criteria

Approve available categories, score definitions, attribution, eligibility, reset/retention, update cadence and tie handling before dispatch. Candidate categories include development, attack contribution, defence contribution and kingdom-held territory; none is a confirmed formula. Do not import Travian population/raid metrics or create new resources solely to fill a ranking.

Show category tabs, rank, player/kingdom, score and snapshot/update time. Provide pagination, search, kingdom filters where meaningful and the current player's position; rows link to player profiles. Kingdom rankings are collective, with no alliance model introduced. Decide whether settlement rankings add useful information rather than adding them automatically.

Compute from server-owned approved facts with world/season isolation. Use stable ordering and documented equal-score ranks. Repeated battle processing, capture retries and profile refresh cannot award points twice. Define how shared defence is attributed and transferred territory counted; distinguish current territory from historical conquest.

Avoid exposing private military strength, resource stocks or hidden scout data through scores or filters. Evaluate repeated mutual attacks and recapture farming before adopting kill/conquest scores. No automated anti-cheat system or penalties invented in this task.

Test empty/new worlds, tied scores, multiple pages, score updates, retries, season reset/finished-world snapshots and access permissions. Bound queries for SQLite, record freshness and avoid per-row query patterns. Responsive accessible tables, both themes, loading/error states and version footer. No global all-time ranking or reward mechanic required for v1.

## Required reading

- [Profiles and settlements](../gameplay/profiles-and-leaderboards.md)
- [Player and settlements](../gameplay/player-and-villages.md)
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
