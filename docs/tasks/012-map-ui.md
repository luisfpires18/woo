# 012: Map interaction prototype

Task key: MAP-UI
Status: TODO
Updated: 2026-10-05
Milestone: UI
Dependencies: GAME-SHELL

## Goal

Prototype regional map viewing, district selection, ownership overlays and details.

## Acceptance criteria

Follow steps 1–9 of the map visual guide using local fixtures. Use labelled mock geography, not invented canon coordinates. Align authored terrain with explicit world geometry; settlement sprites, marker labels and district ownership remain separate. Prove owner recolour without terrain regeneration, independently replaceable settlement art, zoom/pan/selection, drag-versus-click, search/filter/recentre and a React inspector. Check touch, keyboard/list access, resize, both themes, missing-image fallback and representative label density/performance. Record art alignment/repair effort and actual measurements. No saved admin, R2, movement/pathfinding or capture backend yet. Do not declare whole-world seamless art proven by one regional image.

## Required reading

- [docs/gameplay/map-and-activities.md](../../docs/gameplay/map-and-activities.md)
- [Map visual guide](../design/map-visual-prototype.md)
- [Map artwork runbook](../design/map-art-production.md)
- [Map scene contract](../technical/map-scene.md)
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
