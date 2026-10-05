# Development roadmap

Updated: 2026-10-05. Status: design documentation; no game implementation yet.

## Current state

Documentation foundation created. No game code or infrastructure provisioned. Milestone contents below are proposals for execution after user authorisation.

## M0: Specify the medieval alpha

Resolve village count, equipment scale, basic economy/content, talent progression, combat/capture/recovery and first-season ending. Define map size/routes and config-change policy. Exit: each included feature has actions, costs/timing, outcomes, failure cases and acceptance scenarios; all remaining assumptions explicitly provisional.

## M1: Validate economy and battle model

Simulate resource/storage/forging/army growth and mixed unit matchups. Exit: no impossible upgrades, dominant trivial composition or unacceptable daily workload in tested scenarios; provisional values documented with results.

## M2: Accounts, admin and asset foundations

Select stack; implement owner permissions, content definitions, configurable visual slots/fallbacks and safe storage lifecycle. Exit: permissions, edits, image replacement and cleanup demonstrated. R2 account supplied by owner when needed.

## M3: One playable village

Resources, ordinary construction/recruitment, forging upgrades and talent choices; clickable scene and accessible fallback controls. Exit: complete actions persist through refresh/offline completion, display accurately, and cannot spend resources twice.

## M4: Small kingdom frontier

Three playable kingdoms, map travel, automatic battle/report and district capture with agreed recovery. Exit: end-to-end border campaign works, ownership updates consistently and defender can continue according to decided rules.

## M5: Friends alpha

Run a short agreed medieval test. Measure forging decisions, active faction balance, counters, recovery, attendance and admin workload. Exit: critical reliability/fairness issues addressed; findings determine next scope.

## Later releases

Add duties/camps or other missing medieval features based on test results; introduce runes gradually; later L1/L2, armour transformation, unique heroes, naval/air/tunnel systems, new kingdoms and Chaos/Order/Moraphys finale as individually scoped updates. No dates or order beyond medieval-first are committed.
