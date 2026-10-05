# Map visual prototype: step-by-step guide
Updated: 2026-10-05. Status: proposed implementation plan; no code or production assets.

## Purpose and status
Prepare an attractive interactive kingdom map without requiring a dedicated designer. This complements the [village visual guide](village-visual-prototype.md), not replaces it. The owner requested implementation guidance for the generated map preview. The preview remains awaiting explicit visual approval. Its village names, geography, numbers and buttons are illustrative, not lore or balance canon.

Confirmed constraints: independent kingdom ownership, admin image uploads/fallbacks, safe asset cleanup, dark/light UI, restrained textures, owner-managed Git image uploads. Proposal: an authored regional map with separate overlays and markers, rendered by PixiJS within React UI. Exact map topology, world capacity, travel rules and final artwork remain open.

## Recommended first approach
Start with one clean regional terrain illustration and separately rendered villages, location pins, labels and district polygons. Mountains, rivers and vegetation can be baked into static terrain; ownership, village markers and changing locations cannot. A screenshot of the complete UI is a visual reference, not a shippable terrain asset.

This is easier than procedural generation or building the entire world from individual tiles. It supports real pan/zoom and changing ownership while limiting the early art workload. It also has a deliberate limit: arbitrary new rivers, large terrain destruction or unrestricted settlement placement would need additional terrain/layout systems.

Use [map art production](map-art-production.md) for the asset workflow and [map scene implementation](../technical/map-scene.md) for coordinates, layering, interaction and game-state boundaries.

## Small proof region
Use synthetic, clearly labelled geography: one river crossing, foothills, forest edge and coastline. Include a few districts; several separately placed medieval settlement icons; one ordinary neutral camp; one optional hostile outpost. Crimson, green and blue demonstrate Arkazia/Sylvara/Veridor ownership without establishing canonical borders. No whole Bellum atlas, underground layers, naval simulation, rune discovery or army animation.

Use an explicit conceptual two-dimensional map plane, not true 3D terrain. Pictorial mountains are decoration: apparent ridges are not automatically impassable. Explicit topology defines crossings and movement later.

## Ordered work
| Step | Work | Exit evidence |
|---|---|---|
| 1 | Specify bounded region, camera convention, coordinates and sample district topology | Geometry sketch with IDs, polygons, anchors and explicit crossings |
| 2 | Draw a neutral schematic terrain in those coordinates | Roads/rivers and settlement slots readable without illustration |
| 3 | Generate clean terrain from the schematic | No UI, labels, flags, ownership paint or dynamic settlements baked into art |
| 4 | Align and repair artwork against geometry | Landmarks and bridge endpoints match; no inferred connections from accidental pixels |
| 5 | Place separate village/camp sprites and readable markers | Replace a village icon without replacing terrain; matching camera and ground contact |
| 6 | Draw district ownership and selection overlays | Change one district owner without touching image assets |
| 7 | Add pan, zoom, search, filters, selection and React inspector | Mouse/touch alignment; drag does not select; keyboard/list alternative works |
| 8 | Check label density, resize, both themes and loading failure | Legible at actual view sizes; bounded zoom; fallback list remains usable |
| 9 | Measure representative density and decide art approach | Report loaded textures, decoded memory, render timings, device and defects |
| 10 | Later integrate authorised saved configuration and world state | Admin publish/reload, cache refresh, asset cleanup, permissions and capture update proven |

MAP-UI (012) owns the early local/mock proof, steps 1–9. It must not wait for Identity, SQLite or R2. Step 10 belongs to later administration and authoritative district/capture slices; refine their task specs before dispatch. A local fixture editor is not a completed saved admin feature.

## Success gate
Selection remains aligned under pan/zoom/resizing. Ownership changes independently of terrain. Images can fail without losing location access. Representative map density stays readable. Asset creation effort and consistency are reported honestly. The art matches approved topology rather than forcing unapproved gameplay around a generated picture.

If an illustration is beautiful but cannot align with geometry, repair it or simplify the region. If zoom reveals unacceptable blur, cap zoom or produce approved higher-detail art; do not promise unlimited detail. If regional chunk seams cannot be repaired efficiently, keep a bounded single-region alpha rather than claiming seamless world expansion.

## Future expansion
Multiple authored regions can later use a world-aligned chunk manifest and lower-resolution overview. Fixed geographical districts may each receive a village at an authored slot. This population model is a proposal: player/village capacity and expansion rules must be decided before the live map schema. World-scale art and reusable tile production are separate future investments.

## Evidence and storage
Owner uploads approved reference images to Git. Do not commit generated images on their behalf. Record reference IDs, prompt attempts, accepted art, dimensions, coordinate alignment, device captures, timings and verdict in the dispatched task's evidence location. R2 holds referenced published/draft assets and necessary used derivatives, with cleanup governed by the existing lifecycle rules.
