# Admin village scene editor
Updated: 2026-10-05. Status: proposed editor scope serving confirmed configurable images.

## Purpose
Let the owner assemble and maintain approved artwork without editing game code for each building. Early VILLAGE-UI uses local fixture metadata; HOTSPOTS adds real protected saving/upload integration. Framework/schema/routes remain task-level choices.

## Editor workflow
1. Open an authorised village-layout template and select its kingdom/visual variant.
2. Set the terrain reference and canonical dimensions; preview empty layout.
3. Add a building plot or decoration with stable ID and logical definition reference.
4. Choose/upload approved art; set position, scale, ground anchor, label anchor and depth group.
5. Draw/edit its footprint and clickable polygon separately. Offer numeric inputs and list selection alongside dragging.
6. Configure visual-level variants; compare/swap levels at the same plot.
7. Position wall pieces using visible connector guides; verify gate/bridge/foreground intersections.
8. Preview desktop/mobile, zoom and both UI themes with player controls.
9. Validate all references and geometry, save draft, then publish deliberately.
10. Reload player preview and replace one image to verify persistence, cache and cleanup.

## Controls and data
Selection outline, zoom/pan, move/scale, snapping aids, layer/object list, hide/lock for editing only, reorder/manual depth, variant switch, reset and unsaved-change warning. Undo/redo for edits is proposed; do not imply retained binary asset histories.
Do not implement arbitrary rotation for fixed-perspective sprites; it breaks camera/lighting. Horizontal flipping requires approved art support. No collision/pathfinding engine or player free placement in this task.
Separate shared template geometry from village gameplay state; moving a visual plot must not move its world district or change production. Record layout schema/revision, references, publish timestamp and actor. Final DB contracts are resolved before dispatch.

## Validation and publication
Validate allowed formats, limits, referenced definitions, finite bounded geometry, nonzero dimensions, polygons and connector metadata. Reject broken references or incompatible geometry rather than publish a corrupted layout. Decide whether intentional off-canvas objects are permitted.
Review geometry on image replacement, especially changed aspect ratio/anchor. Validate both building variants. Published scene changes should activate as one coherent revision; a player must not see new terrain with old hotspots.
Draft/published edits need concurrency checking to avoid overwriting another admin's work. Server authorisation protects reading/editing/publishing; hiding UI is insufficient. No credentials in client scene definitions.
Publishing/audit integration belongs to CONFIG. Until then, the editor task must state its limited activation model; VILLAGE-VISUAL-GATE waits for CONFIG. Rollback cannot promise deleted historical images: metadata rollback requires still-referenced assets or an explicit re-upload. Follow the owner cleanup requirement.

## Storage references and failure
An assigned draft asset is a real reference; an abandoned draft must release it. Shared sprites survive until their last reference is removed. Preview/failed uploads must be temporary and cleaned, not archived indefinitely.
Upload/validate new bytes, stage geometry review, persist reference atomically as designed, then delete unreferenced prior bytes; if save fails retain the working published image and clean failed upload. Handle cleanup retry and reconciliation.
Owner Git mockup upload policy is separate from authorised admin R2 uploads in the future game.

## Acceptance
Owner uploads terrain/building, edits placement/hit shape, saves, reloads and selects the correct building in player preview. Replaces forge variant without code edit, sees correct cache refresh and validates old object cleanup. Unauthorised user cannot call save/upload/publish APIs. Failure preserves published layout. Shared references and abandoned drafts pass cleanup checks. Provide list/numeric controls as a usable alternative to precise dragging.

See [asset lifecycle](asset-lifecycle.md), [scene design](village-scene.md) and [prototype](../design/village-visual-prototype.md).
