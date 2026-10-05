# Village asset and composition specification
Updated: 2026-10-05. Status: proposed production contract; values require validation in the visual proof.

## One composition contract
Use one approved camera, projection, horizon treatment, light direction and scale. A stylised isometric or oblique illustration is a candidate; the screenshot is not proof of mathematically exact isometry. Do not mix perspective styles. Record roof angle, verticals and visible building faces in the reference sheet.
Working proposal: 2048 × 1536 scene coordinates, rendered responsively. This is a coordinate/art target, not a mandated viewport or texture budget. Validate download size and GPU constraints before production.
Choose a fixed plot layout with reserved footprints, road access, riverbanks and bridge crossing. Leave room for the tallest tested upgrade. Architecture/style differences are permitted; ground scale and camera are shared.

## Layer separation
| Asset family | Separate when | Contract |
|---|---|---|
| Ground / banks | Stable village composition | No buildings, labels, smoke or duplicated shadows |
| River surface | Animation is justified | Static water is sufficient initially; fitted mask if separately animated |
| Roads / foundations | Reused or edited independently | Otherwise bake into terrain; plot entrances still documented |
| Buildings | Levels or occupancy change | Transparent cutout, own footprint/anchor/hit polygon |
| Walls | Upgrade/visibility require modular pieces | Straight orientations, corners, gate and tower with matching endpoints |
| Bridge | Crossing/occlusion differs from water | Separate sprite; deck and supporting geometry align to banks |
| Trees / rocks | Foreground depth or layout changes | Occluders separate; incidental background foliage can be baked |
| Shadows | Needed independently | One convention; do not bake into terrain and sprite simultaneously |
| Effects | Smoke/fire/construction presentation | Separate and optional; no game outcome comes from the effect |
| UI | Labels, counters, selection | Live controls/overlays; never baked into artwork |

Splitting every pixel increases production work. Choose separation based on needed state changes. River position is fixed in this prototype; village scenery does not establish world-map routes or naval rules.

## Per-asset handoff
Record a stable logical asset ID, building-definition ID where relevant, visual variant/level, dimensions, format, approved reference, prompt/repair notes, revision, intended render scale, crop/padding and usage.
Geometry: ground anchor in the image, footprint on ground, clickable polygon, label anchor, shadow policy, depth group, any manual ordering overrides. These are different concepts. A roof is not a footprint; transparent image bounds are not a hit shape.
Layout instance: stable plot/object ID, asset reference, scene position, scale, optional horizontal flip if approved, layer/depth and visibility. Store geometry relative to image or scene consistently, with an explicit schema version. Final schema belongs to implementation design.

## Transparency and export
Proposed first format: transparent PNG for buildings/props and opaque raster for terrain; validate formats/limits in ASSETS. A painted checkerboard is not transparency. No baked text or unexplained logos. Inspect alpha edges on white, charcoal and the real terrain. Avoid excessive empty padding; preserve ground anchor when trimming. Export consistent colour handling and review halo/fringe artefacts.
Do not arbitrarily mirror a sprite with directional lighting or asymmetrical connectors. Scaling repairs minor size differences; it cannot fix a wrong camera. Do not assume AI can cut buildings out of a flattened scene with clean hidden edges.

## Walls, bridge and depth
Wall assets need documented connector endpoints/orientations and matching wall thickness, height and light. Place neighbouring pieces by endpoints, not guessed image rectangles. Include gate opening/clearance, corner, tower junction and near/far wall cases in the first proof.
The bridge must meet both banks, with a deck above water. Rear/foreground rails may require split pieces if objects cross behind them later. No unit navigation is implemented for this proof.
Use feet/ground anchors as a starting point for building depth. Large walls, bridges and trees may span multiple depths: explicit groups, split sprites or occlusion masks may be required. Simple Y sorting does not solve every overlap.

## Completion checklist
- Correct perspective/lighting and matching scale.
- Ground contact, no floating foundations or doubled shadows.
- Readable building identity at minimum planned zoom.
- Space for the tested upgraded silhouette.
- Transparent edges pass light/dark and terrain checks.
- Footprint, hit shape and label anchor supplied.
- Connector seams and foreground overlaps pass.
- Source/usage permission and approval recorded; trial art not mislabelled final.

## Cel-shaded input clarification
The owner can generate transparent 2D cel-shaded assets. Match terrain, buildings and props to one approved camera/light/outline/shading sheet; cel shading alone does not ensure consistency. Validate a forge/tree/ground sample before producing the family. A few building visual stages may cover multiple numeric levels; exact mapping remains a proposal. Resource-site compositions are deferred until the resource model is defined. Read [researched production details](../technical/village-resource-scene-research.md).
