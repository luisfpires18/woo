# 049: Persisted automatic battle resolution

Task key: BATTLES
Status: TODO
Updated: 2026-10-05
Milestone: Combat
Dependencies: COUNTERS, MOVEMENT, ORDERS

## Goal

Resolve a bounded automatic encounter with authoritative inputs/outcomes.

## Acceptance criteria

Freeze eligible troop/equipment/stance inputs under an approved commitment policy; persist battle identity, rules/content versions and RNG details if used. Resolve independently of UI or replay duration. Worker recovery must handle due battles while the web host is idle; no always-on process assumption for Azure F1.

Use database-enforced uniqueness/idempotency and concurrency handling with atomic troop adjustments, result and audit persistence. Simultaneous battles cannot spend the same troops; define reservation/conflict ordering. Test duplicate jobs, competing workers, restart before/after commit, stale versions, zero surviving forces and bounded retry failures. Do not hold a database write transaction for an entire long simulation.

Damage/losses and any approved retreat/loot apply once; repeated processing cannot grant rewards or casualties twice. SQLite-safe transactions and application-managed concurrency tokens are required where appropriate; do not assume SQL Server rowversion. Persist enough immutable data to audit historical outcomes after balance changes. Notifications follow commit and can retry without replaying economic effects. Outcome independent of animation. Capture intentionally deferred.

## Required reading

- [Combat authority and persistence research](../technical/combat-research.md)

- [docs/gameplay/combat-and-conquest.md](../../docs/gameplay/combat-and-conquest.md)
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
