# Map scene and interaction
Updated: 2026-10-05. Status: proposed technical contract; exact APIs/versions resolved at implementation.

## Responsibilities
React renders navigation, generic status placeholders, search/filter controls, accessible location list, inspector and actions. PixiJS renders terrain, geography-aligned political polygons, settlement/site sprites, markers and selection. ASP.NET Core owns authoritative world state, permissions, visibility, movement and capture. The client never determines ownership by reading image colours.

## Data model to prove with fixtures
- Map definition: stable ID, revision, world-coordinate bounds, terrain asset references, authored geometry and coordinate convention.
- District: stable ID, valid land polygon(s), village slot, adjacency/connection references and terrain classification. Exact district-to-village cardinality remains open.
- Location: stable ID, kind, world position, optional district ID, sprite asset/revision and inspector metadata.
- Route/connection: explicit endpoints, allowed movement type and restrictions. Distance/cost rules belong to gameplay decisions.
- World state: world ID, district owner, location status, authorised visibility and revision. Separate static definition from per-season state.

Use a stable bounded world coordinate plane, not viewport percentages. Normalised image coordinates may be imported into that plane once. All overlays and terrain share its transform. Ownership changes state only; geometry editing is a map-definition publication, not a capture operation.

## Render order
Terrain first; optional roads/static geometry; low-opacity district fill; settlement/site sprites; thin political borders; selected district/location highlight; markers and labels. Screen-space controls remain in React. Adjust ordering if fills obscure artwork. District land polygons exclude rivers/sea where appropriate.

A world container supplies pan/zoom to geometry and sprites. Labels/markers may use a separate projected overlay or compensate scale to retain readable screen size. Choose one tested strategy; do not apply the camera twice. Event coordinates must transform from canvas/screen coordinates into world coordinates before polygon hit-testing. Account for device pixel ratio and canvas offset.

District outlines come from authored polygons. Shared edges should use consistent shared vertices to avoid gaps; avoid drawing shared borders repeatedly if it produces dark seams. Update polygon styles on owner change, not every animation frame. Arbitrary concave shapes/holes need verified triangulation/hit-testing support, not rectangular assumptions.

## Interaction contract
Click/tap a marker selects a location and opens React inspector by stable ID. Empty land can select a district. Marker hits take precedence over district hits; decorative terrain does not consume interactions. Drag moves the camera, with a tested threshold preventing accidental clicks on release. Wheel zoom is bounded and centred around the pointer; touch supports pan and pinch with defined behaviour.

Search and filters use authorised data. Clicking a search result recentres and selects it. Recentring a hidden/filtered location requires an explicit consistent policy. Provide an equivalent keyboard-accessible list, clear selected state and focus handling. On phones the inspector can become a bottom sheet; it must not obscure all navigation.

Hover feedback is optional; essential information must not require hover. Labels reduce detail at far zoom, prioritising selection and important locations; reveal other names on focus/selection. Icons retain meaningful touch targets independently of the small painted village footprint.

## Graph versus art
A river painted on the image does not automatically block movement. A bridge painted on the image does not automatically allow it. Connections and movement rules govern that. District polygons support selection and ownership; the graph governs legal travel. The early prototype displays sample connections but implements no combat, authoritative pathfinding or naval simulation.

Capture later updates district ownership and its village marker. The terrain remains unchanged. A render preview may include dark hostile outpost art, but future fantasy objects require normal scope/lore decisions.

## Large map strategy
Do not begin with one enormous texture for all Bellum. Prove one bounded region. For expansion, use world-aligned terrain chunks with manifests, viewport-aware loading/culling and a tested zoom-detail policy. Load visible chunks plus a small margin; release unused textures carefully if shared. Bound caches and cancel stale requests when panning or switching worlds.

Maximum useful zoom is limited by actual art resolution. GPU memory depends on decoded pixels and mipmaps, not just compressed download size. Measure device texture limits and realistic frame/memory costs before selecting chunk dimensions. Defer atlases, filters and animations until they solve measured problems.

## Administration and publishing
Later protected admin tools can upload terrain/sprite assets, place markers, trace/import district polygons, set connections and validate a draft. Separate map-definition editing from player actions and season ownership.

Validate unique IDs, bounds, polygon validity, overlap/gaps by intended topology, shared borders, land/water consistency, connection endpoints, required assets and location-to-district links. Preview before publishing an atomic definition revision. Publishing geometry into an active world needs explicit compatibility handling for existing locations, travel and districts; never silently delete referenced IDs.

Maintain last good published configuration on failed upload/save. Track references from active drafts and published worlds. Clean only unreferenced image objects after successful metadata changes; concurrent editing and shared assets follow the existing lifecycle contract. Pin revisions/cache URLs so replacing terrain does not misalign old geometry with a new image.

## Server integration and safety
MAP-UI uses local fixtures only. Later district/capture tasks connect real state. Validate server permissions, kingdom/world membership and action eligibility. Hidden scouting information must be omitted by the server, not merely hidden with CSS or opaque overlays. Responses should carry revisions; ignore stale results after selection/world changes and re-fetch after missed notifications.

Assets loading slowly must not block basic list access. Show skeleton/error/retry and symbol fallbacks. Canvas unavailable: React list and inspector remain useful. Battle timing, ownership transfers and equipment are server operations, not canvas events.

## Proof evidence
Show terrain and separately replaced settlement, district recolour without terrain reload, marker and polygon hits under zoom/resize, drag versus selection, touch/keyboard alternative, missing-image fallback and representative label density. Record actual viewport/device, content counts, texture sizes, load time and rendering measurements. No universal FPS claim without measurements.

## Primary implementation references
Checked 2026-10-05:
- [PixiJS scene objects and transforms](https://pixijs.com/8.x/guides/components/scene-objects)
- [PixiJS events and hit areas](https://pixijs.com/8.x/guides/components/events)
- [PixiJS asset loading](https://pixijs.com/8.x/guides/components/assets)

These support the rendering primitives. The map data model, editor and gameplay contracts above are WOO proposals, not built-in PixiJS features. Read [prototype sequence](../design/map-visual-prototype.md) and [art runbook](../design/map-art-production.md).

## Square authoring candidate and separate scope
Read [world-map research](world-map-research.md) before implementing MAP-UI. Compare connected tiles against an illustrated baseline on identical geometry. Square paint cells are not automatically districts, village slots, streaming chunks or movement nodes. Projection and cell/district relationships remain open.

Store semantic cells, canonical shared edges/crossings, stable district/location IDs and independent decoration. Choose a primary geometry source; derive secondary boundaries rather than maintaining conflicting cell memberships and polygons. Tiled imports require a documented supported subset and explicit conversions. Editor terrain rules do not run automatically in PixiJS.

The local map proof includes undo/redo and versioned export/reload, not saved admin. Map draft/publication integration belongs to DISTRICTS, while HOTSPOTS covers village interiors only. Resources and production-site imagery remain undecided and excluded. Neighbour rules, cross-chunk props, fractional-zoom seams and bounded texture loading require evidence before expansion.
