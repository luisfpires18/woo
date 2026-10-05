# Village and resource scene: research and production plan
Updated: 2026-10-05. Status: researched proposals, no implementation or production asset validation.

## Goal and owner's input
Make a coherent clickable village and resource landscape using owner-generated 2D cel-shaded transparent building, unit, beast, tree and prop sprites. The owner identifies this as a major feasibility concern. Production should be prepared thoroughly, but no framework, skill or prompt guarantees a flawless first pass.

This supplements the [ordered village proof](../design/village-visual-prototype.md), [asset contract](../design/village-asset-specification.md), [art runbook](../design/village-art-production.md) and [renderer design](village-scene.md). Owner uploads images to Git; no images are committed by this documentation update.

## Recommended implementation choices
| Tool | Proposed use | Limits / verification |
|---|---|---|
| React + TypeScript | Panels, resource numbers, inspector, queues, editor forms and accessible list | Ordinary UI stays HTML; avoid duplicating all UI in canvas |
| PixiJS v8 | Sprite composition, polygons, masks and overlays | Rendering library, not game state, editor or pathfinding |
| @pixi/react | Declarative scene components integrated with React | Official project documents React 19/PixiJS 8 support; pin exact compatible versions in STACK |
| pixi-viewport | Drag, pinch, wheel zoom, bounds and recenter | v6 targets PixiJS 8; verify selected release, lifecycle and event integration |
| PixiJS Assets | Revision-aware texture loading and reuse | Handle failures and shared texture ownership explicitly |
| PixiJS Devtools | Inspect scene objects, transforms and alignment during development | Debugger, not persistent owner-facing editor |
| Tiled | Optional offline object placement, polygons and metadata before the admin editor | Use a bounded import adapter; does not automatically render its map in PixiJS |
| AssetPack | Optional build-time optimisation, manifests and sprite atlases | Not required for early individual uploads; changes do not magically repack published assets |
| Playwright | Real interaction checks and controlled screenshot comparisons | Screenshots cannot alone prove selection, mobile feel or art quality |
| Official pixijs-skills | Focused Claude rendering guidance, already added to inventory | Available guidance is not proof of installed or working local dependencies |

Recommendation for the proof: React, PixiJS, @pixi/react, pixi-viewport and individual PNG assets. Add tools only where useful. Retain one independent scene-definition format so using React bindings or a direct PixiJS wrapper does not change the content contract. A direct wrapper is a fallback if integration problems are demonstrated, not a second renderer to build simultaneously.

Do not copy old tutorial APIs or broad installation commands from READMEs without checking selected versions. Current pages contain some legacy examples. Test compatible React, PixiJS, binding, viewport and devtool releases together. Do not replace Vite with another scaffold merely because a library README suggests one.

## What we should not require
A complete Phaser/GDevelop/Unity-like editor or physics engine is unnecessary for a static clickable settlement. Procedural terrain, skeletal animation, shaders, walking crowds and pathfinding are deferred. @pixi/tilemap is relevant for actual tile-heavy terrain; it adds little to a first scene composed from one terrain background and a small sprite set. PixiJS Layout can arrange canvas UI, but ordinary React panels already handle that need.

Third-party scene editors may be worth future experiments; do not make an unverified editor the critical dependency for alpha. None has been integrated or benchmarked here. Tiled's own maintained object format provides a more explicit optional authoring route.

## Resource view proposal
Two related views can share the renderer and selection model:
- Village: civic/military buildings, walls, roads, forge and main-building plots.
- Resource outskirts: farms/food sites, lumber sites, quarry and ore workings on fixed authored plots.

These are presentation proposals. Resource types, count of fields, plots, upgrade rules and whether the outskirts are a separate screen/tab or one zoomed landscape remain unapproved. The first proof can switch between two small fixture definitions. Do not assume Travian's resource-ring layout or exactly eighteen fields.

Each resource plot references a logical resource-site ID and optional building definition, rather than identifying itself by a sprite filename. Selection opens the same React inspector pattern: production, numeric level, upgrade preview, queue and relevant actions. Early fixtures use clearly illustrative values; real production and spending are server-owned later.

Use a compact site composition: a farm structure plus field patch; lumber structure plus reusable trees; quarry structure plus rock cut; mine entrance plus ore/rock props. These are illustrative asset families, not a confirmed resource catalogue. Repeated scenery can be decorative, while the production site has one meaningful selectable ID. Do not make every tree an independent production entity.

Level numbers and rates are real UI overlays. A small number of approved visual stages can cover many levels; exact mapping belongs to content configuration. Upgrading should retain position and ground contact, swapping the relevant sprite/cluster only. Bigger art must stay within reserved plot envelopes and avoid adjacent sites. Colour wash or tiny badges can communicate selection/status without regenerating terrain.

## Composition contract: the largest risk
Cel shading is a medium, not a shared perspective. First approve one representative forge and a style/camera sheet. All building sprites, trees, rocks and terrain must use that sheet:
- Same shallow overhead camera, visible face orientation and vertical-line treatment.
- Same light direction, shadow convention, outline thickness at play scale and shade-step language.
- Shared ground scale and footprint reference, not equal image widths.
- Defined maximum roof/canopy silhouette for plots and visual upgrades.
- Enough detail to recognise a building at normal play size without requiring extreme zoom.

Generate terrain to match the buildings, or buildings to match the approved terrain. Do not mix a painterly realistic ground with unrelated thick-outline cel-shaded cutouts unless a small assembled sample demonstrates coherence. UI themes change panels/contrast independently; terrain does not need separate light/dark image generations.

A human reference figure can define architectural scale on the private reference sheet. It need not be placed in every final image or turn into ambient unit animation.

## Per-sprite delivery
Required fields: logical ID, source dimensions, alpha bounds, ground anchor, footprint, hit polygon, label point, intended display scale, depth group, revision and approved visual stage. Layout stores position/scale and asset reference separately from gameplay state.

A canvas filled with transparent padding makes bottom-centre anchoring unreliable. Measure the real base contact, trim deterministically and preserve anchor offsets. Alpha silhouette, ground footprint and clickable geometry are distinct. Tree canopy should not invisibly capture clicks on every nearby building. For ambiguous overlapping sprites use explicit hit priority and list-based selection.

First use true-alpha PNG cutouts. Inspect on white, charcoal and final terrain; detect painted checkerboards, halos, clipped roofs and detached shadows. Separate ground shadows are the recommended first convention; if an accepted sprite includes a shadow, disable the additional shadow for that sprite. Ground blobs can help contact but cannot repair wrong perspective.

No arbitrary rotation of a painted 3D-looking building. Request approved orientations where needed. Horizontal mirroring can reverse lighting or asymmetric connectors. Avoid stretch-scaling to disguise camera errors.

## Terrain, walls, water and occlusion
Bake stable grass, river banks, distant vegetation and roads into one cel-shaded terrain image. Keep upgradeable buildings separate. Use static water initially; a fitted water mask/animation is optional later.

Prefer larger precomposed wall runs where that reduces seams. Retain representative straight/corner/gate/tower assets to validate modular placement. Document connector endpoints and gate clearance. Test near wall hiding a building base, rear wall behind roofs, tree canopy overlap and bridge rails. Sorting by ground-anchor Y works for ordinary objects but is not a universal solution; split foreground pieces or explicit ordering handles spanning geometry.

The river crossing is scenery for this proof. If units walk over a bridge later, split rear/deck/front portions and define navigation separately. Do not pay this production cost before moving units are scoped.

## Placement tooling
The early internal authoring harness should expose selected ID, sprite bounds, ground anchor, footprint, hit polygon, label point, layer and numeric position/scale; allow drag placement, camera reset and fixture export. This is a local tool, not a saved protected admin feature.

Tiled can substitute for some early placement work: image/object layers, free object placement, tile objects from individual images, custom properties and polygon objects. Choose a simple supported subset. Convert IDs, coordinates and properties into the WOO scene format and reject unsupported rotations/projections. Tiled's isometric object coordinates and image alignments require explicit conversion; importing JSON alone is not an implementation. Never rely on a random loader claiming full support without a proof.

Later HOTSPOTS adds authenticated saved administration, draft preview and geometry controls, with CONFIG publishing and R2 lifecycle. Share the player renderer for preview so published layouts do not differ from the editor. Add only WOO-required authoring functions rather than creating a general-purpose game engine editor.

## Asset processing and lifecycle
Local asset inspection can calculate dimensions, alpha bounds and suspicious empty padding, and preview anchors. It cannot automatically approve aesthetic consistency. Runtime uploads need normal file validation and usage/reference checks.

Start with independently served assets so changing one building does not require rebuilding everything. AssetPack can later create build-time atlases and manifests for stable shared props. Preserve original dimensions/trim offsets/anchors and verify padding against texture bleeding. Atlas support is optional; logical asset references must not expose atlas indexes as gameplay IDs.

Do not put all user-admin images in a giant mandatory atlas. Runtime replacements and generated derivatives require a publication/refcount policy; old source or atlas objects are removed only once no draft/published scene references them. Stale requests must not apply previous art after a newer selection/revision. Existing [asset lifecycle](asset-lifecycle.md) governs cleanup.

## Performance and interaction safeguards
Keep React updates for meaningful UI/state changes; camera motion and optional effects do not call React state setters each frame. Avoid recreating the Application or reloading textures whenever the resource counter changes. Stable object IDs preserve selection across updates. Rendering/ticker ownership must be explicit, especially under development mount/unmount cycles.

Cap device pixel ratio based on measurements, fit without stretching, and stop unnecessary animation/render work when hidden. Do not enable culling and filters blindly; profiling decides whether they help this small scene.

Measure decoded texture memory, not only PNG size. A single RGBA 2048 by 1536 texture is approximately 12 MiB before mipmaps/extra copies; a 4096 by 4096 texture is about 64 MiB. Repeated trees reuse one texture. Viewport texture limits and total memory are separate issues. Large generated images should be resized for their actual use and maximum approved zoom.

Provisional targets to agree before dispatch: responsive selection within 100 ms after assets are ready; normal interactions aiming at a 16.7 ms frame budget on the specified desktop, acceptable 33.3 ms mobile fallback; first proof scene delivery payload around 10 MB or less; explicitly measured cold-load and navigation memory. These are proposed targets, not measured claims or guaranteed budgets. Select actual desktop and iPhone Safari test conditions, cache/network state and maximum zoom; revise based on evidence.

## One small end-to-end proof
1. Approve camera/style sheet using one forge, tree and clean ground patch.
2. Implement temporary shapes and selection/camera/list before adding expensive artwork.
3. Assemble terrain, forge L1/L2, two buildings and a tree.
4. Add representative wall/gate/tower and bridge overlaps.
5. Add a small resource fixture with representative food/lumber/stone/ore sites using the same scene system.
6. Select sites on mouse, touch and keyboard; swap visual stage while retaining anchor and selection.
7. Check desktop, narrow portrait mobile, pan/pinch boundaries and both UI themes.
8. Simulate image failures, slow loads, rapid revision changes, repeated navigation and clean renderer disposal.
9. Produce controlled screenshots and interaction evidence. Review art repairs and generation effort.
10. Only after proof acceptance expand the building/site family and connect saved admin/R2 in the later gate.

Use fixed fixture state, settled assets/fonts, explicit scene-ready signal and disabled optional animation for screenshot baselines. Playwright browser/device emulation is useful but cannot replace an actual iPhone Safari check for touch and GPU behaviour. Verify inspector ID and world-coordinate interactions in addition to pictures. Snapshot baselines require review; accepting every changed screenshot defeats the check.

If resources composition, walls or alpha/perspective fails, repair that part and retest. Smaller fallback compositions can be reviewed explicitly. Placeholders can pass controls/shell checks but cannot pass final art consistency. Documentation itself is not a passed feasibility test.

## Claude handoff
Use official pixijs entry and relevant Application, Assets, Events, Math, Container, Sprite, Graphics, Performance and Accessibility guidance. Check actual installed sources and matching versions. Use relevant React design skills for HTML controls and browser tooling for validation. No additional unverified skill or plugin is made mandatory here.

The implementation prompt should pin fixture geometry, camera/style reference, asset manifest/revisions, scope and budgets. Ask Claude to build the proof and report limitations, not to improvise the entire production village from a screenshot. WOO remains in preparation until owner dispatches implementation.

## Primary sources and examples
Checked 2026-10-05:
- [Official PixiJS React integration and custom-component examples](https://github.com/pixijs/pixi-react)
- [React integration documentation](https://react.pixijs.io/getting-started)
- [pixi-viewport maintained project and live camera examples](https://github.com/pixijs-userland/pixi-viewport)
- [Tiled object placement and geometry](https://doc.mapeditor.org/en/stable/manual/objects/)
- [Tiled individual-image tilesets](https://doc.mapeditor.org/en/stable/manual/editing-tilesets/)
- [Tiled JSON format](https://docs.mapeditor.org/en/stable/reference/json-map-format/)
- [PixiJS Devtools scene inspection](https://pixijs.io/devtools/docs/guide/features/scene/)
- [AssetPack Pixi workflow](https://pixijs.io/assetpack/docs/guide/getting-started/pixi/)
- [AssetPack atlas padding and trim controls](https://pixijs.io/assetpack/docs/guide/pipes/texture-packer/)
- [PixiJS performance guidance](https://pixijs.com/8.x/guides/concepts/performance-tips)
- [Playwright visual comparisons](https://playwright.dev/docs/test-snapshots)
- [Official PixiJS skills](https://pixijs.com/llms)

Source examples establish component capabilities, not WOO art quality or measured performance. This research did not install, benchmark or visually validate these packages.
