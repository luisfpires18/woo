# 035: Layered village art and admin integration gate
Task key: VILLAGE-VISUAL-GATE
Status: TODO
Updated: 2026-10-05
Milestone: Assets
Dependencies: VILLAGE-UI, HOTSPOTS, IMAGE-CACHE, CONFIG

## Goal
Complete steps 9–11 of the visual proof: demonstrate that the small Arkazia village can be maintained through admin, with repeatable art and correct published player rendering.

## Required reading
- [Ordered visual proof](../design/village-visual-prototype.md)
- [Asset specification](../design/village-asset-specification.md)
- [Art production](../design/village-art-production.md)
- [Renderer](../technical/village-scene.md)
- [Scene editor](../technical/village-scene-editor.md)
- [Asset lifecycle](../technical/asset-lifecycle.md)
- [Prompt protocol](../workflow/implementation-protocol.md)

## Acceptance criteria
Use the reviewed local slice and real protected admin/storage flow. Save/reload/publish terrain, building variants, wall connectors, bridge and tree. Replace forge image/level mapping without code changes and preserve anchors, hit shape, depth and selected inspector.
Test independent cache refresh, shared references, failed upload/save, abandoned draft cleanup and coherent published revision. A failed edit leaves the good published scene usable.
Verify resize/zoom/touch/list keyboard access and representative device performance against budgets agreed at dispatch. Record production effort and defects; owner reviews whether art workflow can scale.
A complete screenshot or shell with placeholders does not pass art/integration criteria. Do not expand art families or implement next gameplay tasks automatically.

## Implementation boundary
One pinned reviewed task. No new gameplay, full roster, free placement, animation system or production provisioning. Owner handles Git image uploads. Resolve missing approvals/budgets before dispatch.

## Execution record
Task-spec commit: Not dispatched
Expected dev base: Not dispatched
Implementation branch: Not created
Implementation commit: None
Review: Pending
Merged dev commit: None
Deployment: Not started; applicability defined in prompt
Issues: Approved modular assets, editor and integration evidence absent.

## TLDR
Next: dispatch after dependencies and art review.
Done: specification only.
Issues: all acceptance checks pending.
