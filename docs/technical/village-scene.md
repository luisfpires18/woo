# Village scene rendering and interaction
Updated: 2026-10-05. Status: proposed implementation design, no renderer implemented.

## Responsibilities
React renders navigation, resources, building inspector, queues, forms and accessible lists. PixiJS is the proposed scene renderer; choose supported versions in STACK and validate APIs during implementation. This document describes behaviour rather than asserting a current API contract.
Server state owns buildings, levels, costs, jobs and permissions. Scene layout owns their presentation only. Selecting art never creates, upgrades or captures a building. Mock UI clearly labels illustrative state and never pretends jobs are saved.

## Rendering modes
Support the proof with a scene definition independent of React components.
Flat mode: one composed background plus associated hotspots; quick fallback.
Layered mode: empty terrain plus separately referenced sprites. The modular proof must use this mode for the forge upgrade, walls and bridge.
Both use stable plot/building IDs and the same selection interface; do not embed building IDs in hardcoded screen pixel positions.

## Coordinates and camera
Keep one canonical scene coordinate system independent of viewport pixels. Asset-local coordinates describe anchors and hit polygons; scene coordinates describe object placement. Normalised coordinates are an option, provided reference dimensions and conversions are explicit.
Use one camera transform for pan/zoom and map pointer positions through its inverse. Account for device pixel ratio, resize, canvas position and sprite transforms. HTML labels must derive screen positions from the same camera; alternatively render visual labels in scene with an accessible React counterpart.
Zoom limits and initial camera framing are task decisions. Bound panning so the village cannot disappear. Fit desktop/mobile viewports without stretching art. Gesture arbitration distinguishes drag from click; mobile tap opens inspector, not hover-only content.

## Depth and overlays
Suggested groups: terrain, water, roads/foundations, ground shadows, buildings/props, designated foreground occluders, effects, visual labels/selection. Validate actual overlap cases rather than assigning one universal order.
Within ordinary objects, sorting by ground anchor can work; fixed groups/manual overrides or split sprites handle walls, bridge rails and large trees. Prevent selection outline/label from vanishing behind an occluder unintentionally. Hit testing chooses the visible/topmost eligible object deliberately, with an accessible list for obscured buildings.
Keep labels legible and avoid overlapping labels: selected label prioritised, unselected labels can reduce with zoom. Do not bake UI into sprites.

## State and upgrades
Resolve asset variant from building definition and visual level mapping. Numeric level need not have unique art. Upgrade presentation replaces sprite and geometry atomically at stable plot anchor; refresh selection/label references without jumping the camera.
Missing art uses symbol and a readable list item, not a broken or invisible click target. Construction state may use a static scaffold overlay; damage variants are deferred unless scoped. Effects optional and reduced-motion aware.

## Loading and performance
Proposed requirements: load required scene assets once per revision, reuse textures for repeated pieces, avoid re-creating renderer/textures on every resource tick, dispose scene-owned resources on navigation, distinguish shared textures during cleanup. Validate ownership against selected Pixi version.
Show loading/failure state; never leave a blocking empty canvas. Large downloads need measured optimisation. Atlas packing is a later option, not a prerequisite for admin uploads or a reason to rebuild all art on replacement.
Before dispatch agree representative desktop/mobile/browser and budgets for first usable scene, payload, interaction frame time and memory. Record device/network/cache state. Do not claim a fixed FPS from screenshots. Test repeated navigation and replacement for leaks.

## Behavioural checks
- Resize and zoom leave selection/hotspots/labels aligned.
- Forge level swap preserves ground contact and updates hit shape.
- Near wall/tree overlaps work; bridge sits above water.
- Pan/pinch does not accidentally start upgrades or select buildings.
- Building list supports keyboard selection and inspector focus.
- Failed/missing image offers fallback without losing controls.
- Theme changes affect UI contrast, not the scene's lighting/art identity.
- Mock preview never changes authoritative gameplay; real upgrade requests enforce server rules later.

Read [prototype gates](../design/village-visual-prototype.md) and [editor](village-scene-editor.md).
