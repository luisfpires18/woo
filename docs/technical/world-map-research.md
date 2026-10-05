# World map: square editing, terrain production and proof
Updated: 2026-10-05. Status: researched proposal, not implemented or benchmarked.

## Owner direction and hard boundaries
The owner wants a convincing cel-shaded world map with an RPG-Maker-like square authoring workflow and independently uploaded transparent assets. The map has its own task: MAP-UI (012). Village interiors belong to VILLAGE-UI/HOTSPOTS; resource types and resource presentation remain undefined. Do not invent a resource catalogue, production fields or a Travian resource ring.

Square painting is a candidate editor interaction, not an accepted gameplay topology. A paint cell, render tile, streaming chunk, political district and settlement location are different entities. No one-village-per-square or diagonal-travel rule is implied.

Research can reduce risks; flawless production cannot be certified before representative art and browser interaction have passed a practical proof. No packages were installed and no images were added during this research.

## Recommendation to test
Use React for controls, inspector and accessible lists; PixiJS for world composition. Prove a square-addressed authoring model with the grid hidden in the normal player view. Compare two ways to draw the same synthetic region:
1. Connected terrain tiles with independent multi-cell trees, mountains and settlement icons.
2. An authored terrain illustration aligned to the same cells/geometry, with identical independent markers and ownership overlays.

The first can support repeatable painting but requires a reusable transition art family. The second can look closer to the illustrative mockup with fewer initial pieces, but arbitrary terrain edits may require image repair. A hybrid of connected ground tiles and larger authored landmarks is a likely production candidate, not a decided outcome. Avoid building two production renderers; the comparison is a bounded proof. Keep logical geometry independent of either art source.

Top-down versus oblique/isometric projection needs one owner-approved camera reference. Square data can be projected into diamond shapes, but projection changes placement, picking and art requirements. Do not change projection midway or infer exact isometry from a generated illustration.

## Tools and what each actually provides
| Tool | Useful contribution | Boundary |
|---|---|---|
| PixiJS v8 | Sprites, geometry, transforms, hit areas and texture loading | No terrain editor, topology model or autotile engine |
| @pixi/react | React integration; documented React 19 / PixiJS 8 support | Pin and prove lifecycle with the selected versions |
| pixi-viewport | Pan, pinch, zoom and camera constraints | Camera plugin, not world geometry or chunk manager |
| @pixi/tilemap | Batched rectangular tile rendering; v5 targets PixiJS v8 | Low-level renderer, not an RPG Maker replacement or Tiled importer |
| Tiled Terrain Sets | Neighbour-aware corner/edge terrain painting | Requires matching transition tiles and configured terrain metadata |
| Tiled Automapping | Rules for secondary layers and repeated decoration | Authoring rules are not automatically executed by PixiJS |
| Tiled JSON export | Structured tile/object layers and chunk data | Convert a documented supported subset into WOO data |
| AssetPack | Optional atlas/manifests and build-time derivatives | Not automatically applied to later admin uploads |
| PixiJS Devtools / official skills | Debugging and implementation guidance | Cannot approve composition, seams or device performance |
| Playwright | Repeatable picking/editor checks and screenshot comparisons | Real touch/GPU and aesthetic owner review still required |

For the first authoring comparison, use Tiled as a reference/offline authoring option instead of rebuilding its entire feature set. Prove a small local WOO editing harness for painting, crossings and export. Decide one primary production authoring route after evidence. No mandatory extra game engine or unverified all-purpose editor.

## Why transparent sprites are only part of the artwork
Transparent trees, beasts, mountains and village icons are suitable independent props. Ground tiles need continuous surface coverage; terrain transition overlays may use alpha. Independently generating random grass, water and cliff squares does not establish matching edges.

Create a controlled terrain family from an approved palette/camera/light sheet:
- Ground interiors with a few subtle variations, avoiding obvious repeated motifs.
- Grass/water edges, outside and inside corners, narrow inlets and islands.
- River straights, bends, junctions if needed, source/mouth and banks.
- Road straights, bends and junctions; explicit bridge alignment.
- Cliff faces/corners/endcaps with consistent height and lighting.
- Forest edge/interior composition using multi-cell foliage rather than identical trees centred in every square.

This is an art inventory proposal, not a confirmed biome list. Start with one ground/water pair and a forest edge. Only add the pieces needed by the fixture.

Tiled documents 16 combinations for two-terrain corner sets and 16 for edge sets; a mixed set has 256 theoretical combinations, often reduced to a 47-tile blob subset. This does not mean 47 assets per biome will automatically suffice: multiple terrain families, cliffs, roads and special intersections introduce further cases. Choose a restricted transition vocabulary first. Do not enable arbitrary sprite rotations to save art when lighting or cliff direction would become wrong.

Owner-generated art should follow an explicit edge template and be assembled at actual scale. An attractive contact sheet is insufficient: place the tiles together, check edges/corners on light/dark terrain, and repair seams before expanding. Record attempts and repair time. If this workload is excessive, compare a bounded illustrated region rather than silently weakening the quality gate.

## Data layers proposed for the proof
| Layer | Information | Do not infer |
|---|---|---|
| Cells | Stable coordinates and semantic terrain/elevation tags | Production yields, movement cost or kingdom ownership from texture name |
| Edges/connections | Explicit barrier/crossing/road connectivity | A painted bridge automatically allowing armies through |
| Districts | Stable IDs and cell membership or approved polygons | One cell equals one district |
| Locations | Stable IDs, positions, marker kind and asset reference | Pixel bounds equal territory or village interior |
| Decoration | Sprite footprint, anchor, layer, scale and allowed variants | Decorative mountain equals impassable terrain |
| Season state | Owners, visibility and statuses supplied by server | Admin paint changing live conquest state |

Decide which geometry is primary. If districts are cell unions, derive shared boundaries from membership; do not maintain a conflicting second polygon source. If authored polygons are primary, validate them explicitly and do not pretend they follow cell edges. Decorative borders can be softened visually without changing selection or legal connections.

Edges need a canonical representation so the east side of one cell is the same edge as the west side of its neighbour. This helps crossings and barriers remain consistent. Rivers may occupy cells, edges or an authored corridor; choose one representation before implementation. Wide seas and narrow border rivers need not share one rule. Rendering may hide the grid while debug overlays expose cells, edges, IDs and crossings.

A district travel graph may coexist with square authoring; ROUTES decides authoritative travel later. Local proof connections illustrate connectivity only.

## Local editor sequence for MAP-UI
1. Pin bounded fixture size, camera/projection, cell coordinates and export schema version. Use synthetic names and geometry.
2. Paint semantic ground/water cells with a small brush; select explicit crossing edges and place location markers.
3. Apply terrain rules in authoring/export or implement the chosen small rule subset locally. Keep the source semantics plus ruleset revision; do not treat arbitrary rendered tile indexes as permanent game IDs.
4. Add deterministic decoration variation using stored choices or a stable seed. Reopening the map must not reshuffle its appearance.
5. Expose layers, grid toggle, selected IDs and numeric controls; group one brush stroke into one undo operation. Test undo/redo and interrupted drags.
6. Export/reload the fixture and verify equivalent geometry and visuals. Reject unsupported layers/properties rather than silently losing them.
7. Render the same fixture in the player view, grid hidden, with accessible location selection and ownership overlays.
8. Compare the tile and illustration approaches with the same geometry, zoom and marker density; record defects and art effort.

No authenticated draft storage or R2 uploads in 012. Saved map editing/publication belongs to DISTRICTS (043), independently of village HOTSPOTS (033). Refine or split 043 into whole-number tasks before dispatch if its implementation becomes too large.

## Import contract and traps
If Tiled is used, define a deliberately small supported subset: one chosen orientation, selected layer types, documented custom properties and known tilesets. Handle or reject external tilesets, compressed/encoded layer data, layer offsets, tile transforms and object alignment. Infinite maps can include negative coordinates and chunks. Preserve stable WOO IDs separately from Tiled tile IDs and object-export details.

Terrain Sets and Automapping help prepare the exported appearance. Loading JSON in PixiJS does not reproduce the authoring tools. Either import resolved tiles or implement an explicit compatible subset; do not accidentally claim general Tiled support. Keep raw authoring files out of authoritative server gameplay until converted and validated.

## Chunk boundaries, zoom and memory
Editor chunks, runtime geometry chunks and texture pages need not have identical dimensions. Begin with a bounded region and a second adjacent chunk to test a real boundary. Neighbour-dependent rules need surrounding cells beyond the edited chunk; use an overlap/halo at least as wide as the chosen ruleset's dependency radius. Rebuild affected neighbours after boundary edits. Store/render each object's identity once even if its visual bounds cross chunks.

Test seams at fractional camera positions and zoom, device pixel ratios, viewport edges and min/max approved zoom. Padding/gutters and appropriate sampling must be validated; atlas padding alone does not guarantee seamless geography. Multi-cell trees/cliffs must remain visible when their anchor is offscreen but their image enters the viewport.

Far zoom should simplify decoration and labels, keeping ownership and important sites readable. Detailed view should not expose excessive blur or huge sprites hiding nearby markers. Set a bounded maximum zoom. Dark/light applies mainly to UI chrome and readable overlays; do not require two entirely repainted maps before proof acceptance.

PixiJS culling is opt-in and skips offscreen rendering. It is not a complete texture-streaming/unloading policy. Load visible chunks plus a bounded margin, cancel or disregard stale requests and retain shared atlases while still referenced. Measure before adding culling to a small scene. A decoded 4096 x 4096 RGBA texture is about 64 MiB before mipmaps/copies, regardless of PNG download size. Report cold-load bytes, decoded texture estimates, actual frame timings and navigation memory behaviour on named devices.

## Saved map publication, later
The protected map editor needs server-checked permissions, saved drafts, revision conflicts, validation, preview and atomic publication. The player renderer should also render the preview. Publish matching geometry, tileset/ruleset, manifest and asset revisions together.

Validate unique IDs, bounded coordinates, complete tile references, matching shared edges, legal bridge endpoints, district coverage/overlap policy and location membership. Errors identify a specific cell/edge/object. A bad upload/import/publish leaves the last good map available. Draft, live and referenced assets follow the existing cleanup contract.

Editing a template for a future season is different from changing an occupied live world. Define live edit compatibility for villages, travel and ownership first; do not shift a coastline through active settlements or silently delete referenced districts. Geometry editing never performs capture. Keep ownership state and hidden scouting data server-owned.

## Acceptance matrix and stop conditions
| Concern | Required evidence |
|---|---|
| Aesthetic quality | Owner reviews identical-region tile/illustration comparison at real play sizes; connected coasts, cliffs and props |
| Editing | Paint, crossing edit, undo/redo, reload and deterministic appearance; invalid case shows actionable error |
| Picking | Selected ID matches under pan/zoom/resize; marker priority and drag-versus-click; accessible list alternative |
| Topology | Shared edge consistency and crossings survive export; recolour does not alter geometry/art |
| Chunk seams | Two adjacent chunks; edge/corner edits update neighbours; fractional zoom and overlapping props |
| Density | Agreed realistic location count; far-zoom label reduction; selected location remains discoverable |
| Failures | Missing art, stale requests, renderer disposal, retry and canvas fallback |
| Performance | Named desktop and actual mobile browser; cold/warm loads, texture estimates, frame measurements and repeated navigation |

Set numeric device/density/load/frame budgets at dispatch, rather than claiming universal FPS. A screenshot cannot pass editor persistence or interaction criteria. No matching art means aesthetic acceptance stays pending. If transition production fails, review a bounded illustrated-map fallback. If illustration fails geometry alignment or expansion seams, repair or reduce region scope. If neither is practical, resolve the production method before whole-world asset generation.

## Primary references
Checked 2026-10-05; these establish capabilities, not WOO feasibility:
- [Tiled Terrain Sets and transition counts](https://doc.mapeditor.org/en/stable/manual/terrain/)
- [Tiled Automapping and rule radius](https://doc.mapeditor.org/en/latest/manual/automapping/)
- [Tiled infinite-map authoring](https://doc.mapeditor.org/en/stable/manual/using-infinite-maps/)
- [Tiled JSON map format](https://doc.mapeditor.org/en/stable/reference/json-map-format/)
- [Red Blob Games: grid edges, barriers and crossings](https://www.redblobgames.com/grids/edges/)
- [PixiJS tilemap package and compatibility](https://github.com/pixijs-userland/tilemap)
- [PixiJS React integration](https://github.com/pixijs/pixi-react)
- [pixi-viewport camera](https://github.com/pixijs-userland/pixi-viewport)
- [PixiJS culling plugin](https://pixijs.com/8.x/guides/components/application/culler-plugin)
- [PixiJS Assets loading, caching and unloading](https://pixijs.com/8.x/guides/components/assets)
- [AssetPack atlas padding and trimming](https://pixijs.io/assetpack/docs/guide/pipes/texture-packer/)
- [Playwright visual comparisons](https://playwright.dev/docs/test-snapshots)

## TLDR
Next: dispatch the dedicated local map proof only when its references and budgets are pinned.
Done: researched tools, art inventory, data boundaries, editor sequence, chunk risks and later publishing requirements.
Issues: projection, topology, final art and measured device performance remain unproven; resources remain undefined.
