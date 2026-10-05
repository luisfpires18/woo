# 033: Admin scene anchors and clickable hotspots

Task key: HOTSPOTS
Status: TODO
Updated: 2026-10-05
Milestone: Assets
Dependencies: ASSETS, VILLAGE-UI, MAP-UI

## Goal

Implement the protected saved scene editor for village templates, layered assets, anchors and hotspots; retain map hotspot support without inferring gameplay topology from artwork.

## Acceptance criteria

Owner can choose/upload terrain and building variants, edit plot position/scale/ground anchor/footprint/hit polygon/label/depth, save and reload. Test wall connectors and bridge/foreground layout. Numeric/object-list controls complement dragging.
Hotspots align across resize/zoom/touch; replacement requires geometry review. Keyboard-accessible player building navigation exists. Server authorisation and asset reference/cleanup rules apply to drafts and shared assets.
Do not rotate perspective sprites arbitrarily or implement free-placement gameplay/pathfinding. Declare limited activation model until CONFIG; later end-to-end gate requires CONFIG. No terrain adjacency inferred from pixels.

## Required reading

- [docs/design/artwork-pipeline.md](../../docs/design/artwork-pipeline.md)
- [Admin scene editor](../technical/village-scene-editor.md)
- [Visual prototype](../design/village-visual-prototype.md)
- [Asset specification](../design/village-asset-specification.md)
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
