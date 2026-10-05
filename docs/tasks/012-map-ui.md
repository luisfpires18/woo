# 012: World map visual and editing proof

Task key: MAP-UI
Status: TODO
Updated: 2026-10-05
Milestone: UI
Dependencies: GAME-SHELL

## Goal

Independently prove a readable, attractive world map and a bounded RPG-Maker-style square-editing workflow using local fixtures. This task excludes village interiors, resource screens and resource catalogue decisions. Square authoring does not approve square-based movement.

## Acceptance criteria

Follow the dedicated map visual guide and research checklist using labelled synthetic geography. Compare tiled terrain and an illustrated-terrain baseline on the same geometry, camera and markers. Prove coast/river/cliff transitions, multi-cell props, cross-chunk neighbour handling and geometry alignment; report asset generation/repair effort, not only screenshots.

Local fixture editing supports terrain painting, explicit crossing placement, marker placement, grid toggle, undo/redo and versioned export/reload. Validate unsupported combinations rather than silently guessing. This is not a saved protected admin editor.

Prove independent ownership recolour, replaceable settlement markers, zoom/pan/selection, drag-versus-click, search/filter/recentre and React inspector/list access. Check touch, keyboard, resize, both themes, missing-image fallback, dense labels and measured loading/rendering/memory. Capture a shared-boundary case at fractional zoom and viewport edge. Agree device, density and performance budgets before dispatch.

Owner reviews side-by-side visuals and production effort before adopting a terrain approach. Missing matching art leaves the aesthetic gate pending. No village interior, resource-site fixture, R2/admin persistence, authoritative movement/pathfinding, capture backend or whole-world art production. A single beautiful region does not prove seamless world expansion.

## Required reading

- [docs/gameplay/map-and-activities.md](../../docs/gameplay/map-and-activities.md)
- [Map visual guide](../design/map-visual-prototype.md)
- [Map artwork runbook](../design/map-art-production.md)
- [World-map research](../technical/world-map-research.md)
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
