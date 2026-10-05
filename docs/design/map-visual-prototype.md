# World map visual and editing prototype: step-by-step guide
Updated: 2026-10-05. Status: proposed implementation plan; no code or production assets.

## Purpose and boundaries
Prove a convincing interactive world map and a square-based authoring workflow in the dedicated MAP-UI (012) task. This guide does not implement village interiors or resource landscapes. Resource names, counts, production fields and screen layout are undecided.

Confirmed constraints: kingdom ownership, admin image slots/fallbacks, safe cleanup, dark/light UI, restrained textures and owner-managed Git image uploads. RPG-Maker-like square painting is a candidate editing approach, not a confirmed movement grid. The latest map preview awaits visual approval; its names, geography and numbers are illustrative.

## Read first
Read [world-map research](../technical/world-map-research.md), [map artwork production](map-art-production.md) and [map scene contract](../technical/map-scene.md). Keep paint cells, districts, settlements, render tiles and streaming chunks distinct.

## Small representative slice
Use synthetic geography containing a coast, river bend and explicit crossing, forest edge, foothills/cliff, a few districts and separately placed medieval settlement/camp icons. Crimson, green and blue demonstrate ownership without establishing canon borders. Add one adjacent chunk to expose seam and neighbour-rule problems. No entire Bellum atlas, village interior, resource fields, naval simulation or fantasy progression.

## Ordered work
| Step | Work | Exit evidence |
|---|---|---|
| 1 | Pin bounded coordinates, camera/projection, cell model and sample districts/crossings | Labelled geometry and schema; outstanding decisions visible |
| 2 | Make a schematic player view and local square-painting harness | Terrain paint, explicit crossing, marker placement, grid toggle and numeric controls |
| 3 | Configure a small connected ground/water transition family | Assembled corners/edges and defined neighbour rules, not isolated sprite samples |
| 4 | Compare tiled terrain with an illustrated baseline on the same geometry | Same camera/markers; alignment, style and asset-repair effort recorded |
| 5 | Add independent multi-cell props and settlement icons | Stable anchors, correct overlaps and independently replaceable art |
| 6 | Add district ownership, selection, pan/zoom and inspector/list | Recolour without terrain replacement; accessible navigation and correct picking |
| 7 | Prove undo/redo, deterministic export/reload and boundary edits | One stroke/undo unit; neighbour updates and unchanged IDs |
| 8 | Check real view sizes, touch, keyboard, themes, missing art and density | Readable labels; drag does not select; image failure preserves navigation |
| 9 | Measure loading, frame behaviour and texture memory; review approach | Named devices, agreed budgets, defects and owner visual verdict |
| 10 | Later integrate protected drafts, publication and world state | DISTRICTS owns saved map configuration; matching revisions and safe live-edit policy |

012 owns steps 1–9 with local fixtures only. It does not depend on Identity, SQLite or R2. Step 10 is separate from village HOTSPOTS. Refine/split its map task before implementation if needed; no automatic successor work.

## Success and fallback
The grid can be hidden without making the map confusing. Connected terrain looks coherent at actual scale; ownership and markers change independently. Picking survives camera transforms. Multi-cell props and two-chunk boundaries work at fractional zoom. An exported fixture reloads consistently.

The owner reviews side-by-side visuals and art effort before adopting a production method. Missing matching assets can prove controls, but leave aesthetic criteria pending. Tile seams require repair or a reviewed illustrated fallback; illustration misalignment requires repair or smaller scope. Do not claim whole-world expansion or flawless asset generation from one attractive region.

## Evidence and storage
Record reference IDs, attempted/accepted art, coordinates, asset dimensions/revisions, repairs, owner verdict, screenshots, interactions and device measurements in the dispatched task's evidence location. Owner uploads images to Git. Do not commit generated images on their behalf. R2 later holds referenced assets and necessary used derivatives under the lifecycle policy.
