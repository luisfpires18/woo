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

## Superseded narrow first-art proposal

One real Arkazia village scene with six configurable hotspots, working detail panels and dark/light mode. This earlier proposal is superseded as the complete feasibility gate by the layered visual proof below. No game code or production-ready modular assets exist. Preview mockups have been generated.

## Village feasibility work, 2026-10-05
Read the [step-by-step visual prototype](village-visual-prototype.md) before producing the village family. Prove a clean terrain, separately replaceable forge/upgrade, two buildings, tree, modular walls and river/bridge integration. VILLAGE-UI tests local art/interaction; HOTSPOTS implements saved administration; VILLAGE-VISUAL-GATE validates end-to-end operation before expanding village gameplay. Static whole-scene hotspots are a fallback, not proof of modular upgrades.
The owner now uploads images to Git personally. ChatGPT maintains documentation and generates previews; it must not upload approved or unapproved images on the owner's behalf. Admin R2 uploads are a separate future game feature.

## Map feasibility work, 2026-10-05
Read the [map prototype sequence](map-visual-prototype.md), [art runbook](map-art-production.md) and [scene contract](../technical/map-scene.md). Begin with geometry, then clean authored terrain, separate settlements and dynamic political overlays. MAP-UI proves one local region; full-world chunks and saved admin publishing are later integrations. The generated map preview is awaiting approval and cannot be used as final terrain because labels, ownership and villages are baked in.
