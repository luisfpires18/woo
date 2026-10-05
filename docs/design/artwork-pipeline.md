# Artwork pipeline

Updated: 2026-10-05. Status: design documentation; no game implementation yet.

## Confirmed priorities

Designer-free workflow using generated assets uploaded through admin. Artwork quality and clickable village/map/combat are major concerns. Static art is the first approach; consistency and animation require iteration.

## Proposed production stages

Approve an Arkazia visual reference; make a clean village scene with normalised polygon hotspots; add live UI overlays; later test separate building sprites with fixed footprints/anchors. Design map geometry/routes first, then illustration; use dynamic political overlays. Combat begins with background and tokens/portraits with server-event replay, not individually animated armies.

## Quality gate

Test at actual desktop/mobile sizes and both themes. Check style, perspective, transparency, unwanted symbols/text, interaction alignment and lore traits. Reject or repair failures before expanding families. Existing LF2 character sprite style is not automatically an environment standard.

## Current storage policy

Keep only referenced images and required used derivatives in the game bucket. Earlier image-version retention proposals are superseded. Prompt/style/approval metadata can remain in documentation, but unused binary histories must not accumulate. See [asset lifecycle](../technical/asset-lifecycle.md).

## Proposed first art milestone

One real Arkazia village scene with six configurable hotspots, working detail panels and dark/light mode. This is a proposal; no artwork production or code has been started by this documentation migration.
