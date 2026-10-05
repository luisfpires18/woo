# 057: Player and settlement profiles

Task key: PROFILES
Status: TODO
Updated: 2026-10-05
Milestone: Community
Dependencies: MEMBERSHIP, VILLAGE-CREATE, SCOUTING, CAPTURE, ASSETS

## Goal

Allow players to inspect other players and their settlements within the selected world. Use settlement/settlements in player-facing UI; existing VILLAGE task keys remain stable.

## Acceptance criteria

Provide a world-scoped player page with nickname, image/symbol fallback, kingdom, approved public statistics and settlement links. Each settlement page shows name, authorised owner/kingdom, visible location and approved public development summary. Profile links work from visible map markers and kingdom lists; leaderboard integration follows LEADERBOARDS.

Provide distinct self, same-kingdom and opponent views under an explicitly approved visibility policy. Server projections omit hidden troop counts, equipment, resources, queues, exact private activity and account/admin information. No scouting bypass via direct IDs or profile search. Define whether discovered settlements, all owned settlements or only public ones appear before dispatch.

Counts, world membership and ownership reflect capture correctly; stale or removed/captured settlements have deliberate navigation behaviour. Public profile text/image fields are limited to already supported editable fields unless separately approved; do not invent messaging, biographies or social features.

Support accessible responsive lists, search within permitted world data, image loading/fallback, both themes and version footer. Test direct links, unknown IDs, permissions, world isolation, capture transfer and missing images. Settlement quantities, population and scoring formulas are not invented. No new economic or combat rules.

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
