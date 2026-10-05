# 011: Village and resources screen prototype

Task key: VILLAGE-UI
Status: TODO
Updated: 2026-10-05
Milestone: UI
Dependencies: GAME-SHELL

## Goal

Prove a readable village/resources interface and a small repeatable layered Arkazia art workflow before expanding the building family. This is local mock-state work, not persisted game/admin behaviour.

## Acceptance criteria

Follow steps 1–8 of the visual proof. Record approved reference/camera/plots before asset production. Assemble clean terrain, forge and upgrade, two buildings, tree, modular wall junction/gate/tower and river/bridge. Replace the forge at a stable anchor; inspect alpha, scale, perspective, seams and occlusion.
Selection, labels and companion list work across zoom/resize/touch/keyboard. Mock buildings have configurable image slots/fallbacks. React inspector uses illustrative costs/timings; no real orders/auth or pretend saved admin.
Report actual screenshots, tested viewports/devices, asset attempt/repair effort and measured rendering behaviour. Resolve provisional performance budgets at dispatch. Modular proof cannot pass on a flattened scene alone.
If approved art is missing, placeholders can demonstrate the shell but art-consistency criteria remain pending; obtain owner review rather than claiming success.

## Required reading

- [docs/design/artwork-pipeline.md](../../docs/design/artwork-pipeline.md)
- [Visual prototype](../design/village-visual-prototype.md)
- [Asset specification](../design/village-asset-specification.md)
- [Art runbook](../design/village-art-production.md)
- [Renderer design](../technical/village-scene.md)
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
