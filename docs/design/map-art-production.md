# Map artwork production
Updated: 2026-10-05. Status: proposed workflow, no production art.

Read [map prototype](map-visual-prototype.md) first. Generation assists art production; it does not guarantee exact geometry, seamless tiles or consistent perspective.

## Asset separation
| Asset | Contents | Must exclude |
|---|---|---|
| Regional terrain | Land, water, static roads, mountains, woodland and reserved settlement spaces | UI, text labels, faction colours, pins, mutable settlements |
| Settlement sprites | Small medieval village, fort or outpost on transparent background | Giant ground rectangles, labels, ownership flags baked permanently |
| Site sprites | Ordinary camp/den or other approved site | Unapproved rune/monster canon |
| UI markers | Crisp shield/pin, selection ring, status icon | Detailed painted text or illegible ornament |
| Inspector thumbnail | Approved close-up location illustration | Any assumption that this is the actual zoomed map terrain |

One consistent regional terrain is acceptable for the first proof. Do not generate every tree or river segment separately. Reuse small consistent settlement sprites; use different markers to express ownership. UI markers can be code-native SVG, while terrain and settlement paintings are raster assets.

The art is neutral geographically; political borders are data-driven overlays. Canonically distinctive vegetation/architecture can remain, but capture must not require regenerating landscape.

## Production sequence
1. Approve a simple geometry sketch with land/water, river crossing, district boundaries and reserved village slots. Lock a top-down or shallow pictorial overhead convention; do not mix village-scale isometric buildings with a different map angle.
2. Produce the terrain with the sketch as a structural reference. Prompt for composition preservation, restrained painterly detail, clean settlement clearings, no UI/text/borders/pins. Inspect actual alignment rather than relying on the prompt.
3. Repair coastline, crossing, roads and clearings where necessary. Preserve known coordinates. Large layout changes require geometry review; do not move gameplay silently to accommodate artwork.
4. Produce one settlement sprite and one camp sprite matching perspective, lighting and ground scale. Test genuine alpha, edge halos and ground anchor in the actual scene.
5. Assemble the map with temporary markers. Review at normal play size before making many variants. Individual house-level detail should not overwhelm the regional view.
6. Create a few variants only after the base sprite passes. An Arkazia sprite must not inherit random magical symbols or unapproved heraldry.
7. Measure asset dimensions, transfer size and decoded GPU memory. Select formats and size limits in the implementation proof based on devices and texture capabilities.

## Generation prompt scaffolds
Terrain: "Create a clean medieval regional terrain painting following this geometry sketch. Preserve river banks, crossing position, coastline and empty settlement clearings. Consistent shallow overhead view and lighting. Static geography only. No text, UI, flags, faction colours, pins or settlements."

Settlement: "Create one compact medieval settlement sprite in the approved regional-map camera and light direction. Entire silhouette visible, matching ground anchor and scale, transparent background. No text, markers or painted ownership tint."

These are starting briefs, not guaranteed assets. Include the approved style reference and measured coordinate sketch in the actual production prompt. Record any hand repair required.

## Why screenshot cropping is insufficient
The preview contains labels, political washes, settlements and terrain already flattened together. Cropping it leaves baked ownership and villages. It also provides no hidden background beneath buildings and no extra detail beyond its original resolution. A clean terrain asset must be generated/edited specifically for reuse. Do not automatically extract modular terrain from the preview.

## Scaling to multiple regions
Start with one image. Later, artwork chunks need the same world scale, palette, light and shared river/road/coast edge specifications. Cropping an approved large master into delivery chunks guarantees continuity within that master, but independent image generation does not guarantee matching neighbours.

For new adjacent regions, use overlap references and deliberate edge repair. Define actual chunk bounds and overlap cropping in a manifest; avoid visible double-painted seams. Generate overview/downsample derivatives from approved geometry-aligned art. All published variants must be referenced and counted in asset cleanup.

## Review checklist
Inspect normal zoom, maximum allowed zoom, touch size and dark/light controls. Check river crossings, marker ground contact, village silhouettes, faction colour contrast, unlabelled landmarks and ownership wash opacity. Do not recolour an entire forest crimson to hide border ambiguity. Pair colours with names/icons. Exact art dimensions and performance targets remain pending the proof.

See [asset lifecycle](../technical/asset-lifecycle.md) and [village asset specification](village-asset-specification.md) for shared anchor, revision and cleanup discipline.
