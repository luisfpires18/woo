# Village visual prototype
Updated: 2026-10-05. Status: planned proof; no code or production assets created.

## Purpose and decision status
The owner wants a detailed, step-by-step proof that a beautiful interactive village can be produced without a designer before expanding the game. Documentation is authorised; implementation still follows the normal task workflow.
Confirmed: admin image slots/fallbacks, independent image revisions, safe cleanup, dark/light support, restrained textures and owner-managed Git image uploads.
Preferred visual reference: early dark navigation / warm ivory content / crimson actions, illustrated village and selected-building inspector. The latest remade village screenshot is awaiting approval; style preference does not approve its building roster, prices or layout.
Proposal to validate: fixed village plots, reusable layered artwork, React controls and a PixiJS scene. Exact framework versions belong to STACK. Fixed plots are not a confirmed gameplay restriction; free placement is excluded from this proof.

## Read in order
1. [Asset and composition specification](village-asset-specification.md).
2. [Art production runbook](village-art-production.md).
3. [Scene and interaction design](../technical/village-scene.md).
4. [Admin scene editor](../technical/village-scene-editor.md).
5. [Storage lifecycle](../technical/asset-lifecycle.md).
6. [Researched village tooling and production plan](../technical/village-resource-scene-research.md).

## Small representative slice
Use Arkazia only. Produce clean terrain; a forge at two visual levels; two other ordinary buildings; a tree; straight wall, corner, gate and tower pieces; a river and bridge. Names and levels are visual samples, not final gameplay content. Test at least one overlapping foreground object. No runes, armies walking through the village, free building placement, combat simulation or nine kingdom art families.
A flat whole-village illustration with hotspots remains a fallback/comparison mode. It does not prove independent upgrades or modular walls. A composed image cannot automatically be separated into clean reusable assets.

## Ordered delivery and evidence
| Step | Work | Exit evidence |
|---|---|---|
| 1 | Select the reference, camera, working canvas, plots and palette | Owner-approved composition/style sheet; unresolved values visible |
| 2 | Produce empty terrain with roads, banks and plot footprints | No baked buildings/UI; reserved space survives at play size |
| 3 | Generate forge, transparency and shadow treatment | Ground anchor, clean edges and correctly matching perspective |
| 4 | Assemble terrain + forge in a local scene | Screenshot at actual desktop/mobile size; selection opens React inspector |
| 5 | Replace forge with its next visual level | Same plot/anchor, no jump, new hit shape and no clipping |
| 6 | Add two buildings and a tree | Consistency and foreground occlusion proven |
| 7 | Assemble wall/corner/gate/tower and bridge | No seams; correct depth and plausible river crossing |
| 8 | Check zoom, resize, touch and keyboard alternative | Pointer coordinates and labels align; pan does not trigger selection |
| 9 | Persist editor configuration and integrate R2 | Admin places/replaces art without game-code edits; reload restores layout |
| 10 | Replace assets, interrupt failures and re-test | Good published scene remains usable; shared assets survive; unused objects cleaned |
| 11 | Review production effort and rendering measurements | Continue / needs changes / stop verdict with actual evidence |

VILLAGE-UI covers steps 1–8 using local fixture assets and mock state. HOTSPOTS covers the persistent authorised editor. VILLAGE-VISUAL-GATE covers steps 9–11 end to end. Early proof must not wait for R2, Identity or the full database. An internal local placement fixture is not a pretend saved admin feature.

## Success gate
- Independently replace a building and its visual level without changing game code.
- Maintain consistent camera, scale, lighting, ground contact and readable silhouettes.
- Keep hit areas, selection and labels aligned under zoom, resizing and touch.
- Assemble representative wall pieces and a bridge; render overlaps deliberately.
- Provide an accessible building list and inspector even if canvas/images fail.
- Demonstrate actual admin saving, permission checks, cache refresh and cleanup later.
- Report generation attempts, repairs, accepted assets, time spent and measured performance. No claim that prompts alone guarantee production quality.

An attractive screenshot is insufficient. If building consistency fails, repair the camera/asset contract before making more content. If wall seams fail, test larger precomposed wall groups as a fallback. If layered art is impractical, explicitly review reduced visual upgrade scope; do not silently declare the modular proof passed.

## Review evidence record
Keep an evidence folder/document in the implementation task's agreed location: approved reference identifiers, asset dimensions/revisions, placement metadata, screenshots or short captures, devices/viewports, timings, defects, repair attempts and verdict. Owner uploads approved mockups to Git; do not upload images on their behalf.
Retain review/source files only in agreed source storage. The live R2 delivery bucket keeps referenced assets, not unused trial generations or binary history.

## Outstanding decisions
Working canvas/camera numbers; final alpha roster/plots; reference and artwork rights; asset limits/formats; shadow convention; depth groups; performance budgets; publishing safety. See [open questions](../open-questions.md). Do not block the skeleton on these; resolve before the affected prototype step.

## Related map guide
The [map visual prototype](map-visual-prototype.md) applies the same separation of artwork, geometry, interaction and game state at regional scale. Its authored terrain plus dynamic districts/markers is distinct from the village's independently upgraded building composition.

## Cel-shaded sprite preparation and scope
Owner plans transparent 2D cel-shaded buildings, units, beasts, trees and props. No final village art family should expand before assembled perspective, lighting and anchor checks pass. See the [tool comparison and proof sequence](../technical/village-resource-scene-research.md). Its filename is historical; the current guide covers the village only.

Resource types and resource-scene layout remain undecided and are excluded from this proof. World-map rendering and square editing have their own task and guide. Sharing technical utilities does not combine these deliverables.
