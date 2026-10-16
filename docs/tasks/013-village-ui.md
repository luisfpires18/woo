# 013: Village scene and interaction prototype

Task key: VILLAGE-UI
Status: WIP
Updated: 2026-10-16
Milestone: UI
Dependencies: GAME-SHELL

## Goal

Prove a readable village interface and a small repeatable layered Arkazia art workflow before expanding the building family. This is local mock-state work, not persisted game/admin behaviour.

## Acceptance criteria

Follow steps 1–8 of the visual proof. Record approved reference/camera/plots before asset production. Assemble clean terrain, forge and upgrade, two buildings, tree, modular wall junction/gate/tower and river/bridge. Replace the forge at a stable anchor; inspect alpha, scale, perspective, seams and occlusion.
Selection, labels and companion list work across zoom/resize/touch/keyboard. Mock buildings have configurable image slots/fallbacks. React inspector uses illustrative costs/timings; no real orders/auth or pretend saved admin.
Compare React/Pixi bindings and camera compatibility using the researched guide. Resource types, resource layouts and world-map editing are excluded. Optional placement tools are internal village fixtures, not saved admin.
Report actual screenshots, tested viewports/devices, asset attempt/repair effort and measured rendering behaviour. Resolve provisional performance budgets at dispatch. Modular proof cannot pass on a flattened scene alone.
If approved art is missing, placeholders can demonstrate the shell but art-consistency criteria remain pending; obtain owner review rather than claiming success.

## Required reading

- [docs/design/artwork-pipeline.md](../../docs/design/artwork-pipeline.md)
- [Visual prototype](../design/village-visual-prototype.md)
- [Asset specification](../design/village-asset-specification.md)
- [Art runbook](../design/village-art-production.md)
- [Renderer design](../technical/village-scene.md)
- [Village tooling research](../technical/village-resource-scene-research.md)
- [Scope](../scope.md)
- [Prompt protocol](../workflow/implementation-protocol.md)

## Implementation boundary

A separately reviewable planning task. ChatGPT inspects current code and resolves this task's blockers before issuing one pinned Claude prompt. Preserve medieval v1. Split/refine based on actual development; do not implement successors automatically.

## Result (first pass, 2026-10-16)

The Settlement screen (`/game/settlement`) is an interactive prototype: a PixiJS scene inside the GAME-SHELL stage, with a React building list and inspector beside it (below it on phones). **All scene visuals are temporary authored geometry, labelled "Prototype geometry · not final art"; no village artwork was generated, imported or committed.** The art-consistency criteria therefore remain open (see the art gate below), and the task stays WIP.

### Dependencies (exact pins)

- `pixi.js` 8.22.0: the STACK selection, still the current v8 release.
- `pixi-viewport` 6.0.3: current release, peers `pixi.js` >=8. It provides bounded drag, pinch, wheel zoom, deceleration and a drag threshold, and emits `clicked` only when the pointer did not pass that threshold, with the point already mapped to world coordinates. A hand-written camera would re-implement all of that.
- `@pixi/react`: **not used**. Its latest release (8.0.5) pins `react-reconciler` 0.31.0, built for React 19.0; WOO runs React 19.3.1. The scene is plain data rendered by a small imperative wrapper, so the binding would add a reconciler-version risk without simplifying anything. Recorded in [stack](../technical/stack.md).
- The renderer module is loaded with a dynamic `import()`, so PixiJS is downloaded only when the Settlement screen opens; title, worlds and other game screens are unaffected.

### Architecture

| Part | File | Role |
|---|---|---|
| Scene definition | `web/src/scenes/settlement/sceneDefinition.ts` | Plain data, no React or PixiJS: canonical size, camera convention, layers, objects, terrain, framing |
| Renderer | `web/src/scenes/settlement/SettlementRenderer.ts` | Imperative PixiJS: Application, viewport camera, layers, prototype shapes, image slots, selection overlay, labels, hit testing, resize, cleanup. Handle API: `select`, `setStage`, `zoomBy`, `resetCamera`, `destroy` |
| React view | `web/src/game/SettlementView.tsx` | Stage host with loading/failure states and camera buttons, building list, inspector, local mock state |
| Fixtures | `web/src/game/previewData.ts` (`previewBuildings`) | Illustrative mock levels and example upgrade figures, all labelled |
| Styles | `web/src/settlement.css` | Stage, list and inspector |

React state changes only on meaningful events (selection, visual stage, art status). Camera motion never updates React state. The Application is created once per mount and destroyed on unmount (canvas removed, ResizeObserver disconnected, scene-loaded textures unloaded). Creation can be cancelled before PixiJS finishes initialising, which StrictMode's double mount in development exercises.

### Scene and camera values

- Canonical scene: **2048 × 1536** units. Objects are placed in scene coordinates; the camera never changes them.
- One shallow-overhead oblique convention for every object: verticals stay vertical and ground depth runs **(0.35, −0.6)** per unit. Light comes from the upper left: front faces mid tone, right faces dark, front roof slopes lit, shadows offset down-right. The scene's colours do not change with the UI theme.
- Camera: `pixi-viewport` with clamp to the scene edges (underflow centred), zoom between **contain-fit of the whole scene** (minimum, recomputed on resize) and **2×** (maximum), drag threshold **8 px**, wheel and pinch zoom, deceleration (disabled under `prefers-reduced-motion`, like wheel smoothing).
- Initial framing: region x 230–1790, y 420–1420 on wide stages (scale 0.476 at 1280×800); region x 470–1450, y 430–1390 on stages narrower than 600 px (scale 0.348 at 375×812). The stage is 4:3 on desktop/tablet and 4:5 on phones.
- Resize: a ResizeObserver resizes the renderer (no stretching, device pixel ratio up to 2), recomputes the minimum zoom and keeps the camera centre.
- Draw layers, back to front: terrain, water, roads (with plots), bridge, shadows, objects (sorted by ground-anchor y), foreground (authored order), overlay (selection and labels).

### Asset contract (per object)

`id` (stable), `name`, `layer`, `anchor` (ground anchor in scene coordinates, fixed across visual stages), `scale`, `selectable`, and `stages` keyed by visual stage. Each visual has: `shape` (prototype geometry used when there is no image), optional `image` slot `{ src, anchor (image pixel on the ground anchor), scale, revision }`, `footprint` (ground polygon), `hit` (clickable polygon, not the image bounds) and `label` point. An image slot is loaded with `?v=<revision>`; when it loads, the sprite replaces the prototype shape at the same anchor. When it is missing or fails, the shape stays and the inspector says so. Nothing here is persisted; saved placement is HOTSPOTS (035) and uploads are ASSETS (033).

Texture ownership (review fix): image textures are reference-counted per URL. A load takes a reference when it starts. The reference is given back when the load fails, or when it resolves after its view was replaced or the scene was destroyed (it is then never attached). A live view gives it back when the view is replaced and when the scene is destroyed. `Assets.unload` runs only when the last reference goes, so a texture shared by two visual stages survives a swap (the new stage's load takes its reference before the old view gives its back). Only URLs loaded by the scene are ever unloaded. No callback reaches React after destruction or for a replaced view.

### Fixture composition (prototype content)

| ID | Name | Layer | Anchor | Notes |
|---|---|---|---|---|
| `forge` | Forge | objects | (780, 900) | Stage 1: 200 × 150 footprint, 105 high walls. Stage 2: 240 × 170, 150 high, chimney with forge glow. Different hit polygons and label points, same anchor |
| `cottage` | Cottage | objects | (1080, 760) | Image slot `/scene-assets/settlement/cottage-stage-1.png` (rev 1), intentionally absent: proves the missing-art fallback |
| `hall` | Hall | objects | (1060, 1258) | Base sits behind the south wall: foreground overlap |
| `oak` | Tree | foreground | (905, 1000) | Canopy overlaps the forge's front-right corner; not selectable (its canopy does not capture clicks) |
| `wall-west` | Wall | objects | (300, 1300) → (650, 700) | Straight run into the view |
| `wall-corner` | Wall corner | foreground | (300, 1300) | Junction of the west and south runs |
| `wall-south-west`, `wall-south-east` | Wall | foreground | y 1300 | Straight runs either side of the gate |
| `gate` | Gate | foreground | (840, 1300) | Arched opening, road passes through |
| `tower` | Tower | foreground | (1320, 1308) | End tower with roof |

Terrain: ground with tonal patches, a river (165 wide, banks either side) from the top edge to the bottom edge, roads from the gate to the square, the forge, the cottage, an empty plot and the bridge, four plot foundations (one empty), and a bridge (deck, back rail, front rail, posts) drawn above the water at y ≈ 935. Names and levels are visual samples only. Inspector fixtures: mock levels (Forge 1 at stage 1 and 5 at stage 2, Cottage 1, Hall 2) and, for the Forge only, an "Illustrative upgrade (mock)": Resource 1: 120, Resource 2: 80, Time: 12 min, marked "Example figures only; resources, costs and timings are not decided."

### Interaction and accessibility

- Pointer/tap: `clicked` → hit test of the current visual's hit polygon, topmost selectable object first (layer, then ground y). Non-selectable scenery (tree, walls) does not block selection behind it; the list remains the way to reach any building.
- Drag vs click: a drag that starts on the Cottage pans the camera and leaves the selection unchanged.
- Selection: outline of the hit polygon plus a tinted footprint, drawn in the overlay above every layer (visible through the tree), constant on-screen stroke width; the selected label turns gold. Labels keep a constant on-screen size at any zoom.
- Building list (React): one button per building, `aria-pressed` and a check mark when selected, 44 px targets, keyboard operable; a polite live region announces "<name> selected". The inspector follows the list in reading order.
- Inspector: name, "Mock values · nothing is saved", summary, mock level, visual stage, art status, the Forge's illustrative upgrade, and a "Prototype visual stage (local only, no upgrade is ordered)" control with Stage 1 / Stage 2 buttons (`aria-pressed`).
- Camera buttons: Zoom in, Zoom out, Reset view (44 px, keyboard operable). The canvas is `aria-hidden`; its wrapper is `role="img"` with a name and a description that points to the list.
- Failure states: "Preparing the scene…" while loading. If the renderer cannot start, an alert says so and the list and inspector keep working (checked with the dev-only `?scene=fail` fixture). Missing image art keeps the prototype shape and is reported in the inspector.

### Checks and measurements (2026-10-16, this Windows computer, headless Microsoft Edge with WebGL; Vite dev server unless stated)

- Pointer click on the Forge selected it; pointerup → list `aria-pressed` updated in **4–4.6 ms**; next painted frame after the click **9–27 ms** (four runs). Target < 100 ms met.
- Drag starting on the Cottage moved the camera centre from (1010, 920) to about (1268, 745) at the same scale with the selection unchanged.
- Forge Stage 1 → 2: selection kept, camera identical before and after, anchor unchanged; the world point (790, 640) hits the Forge at stage 2 and nothing at stage 1 (hit polygon follows the stage).
- Zoom in ×3 (scale 0.930) then click on the Forge: selected. Resize to 1000×720: canvas 807×602 CSS px = 807×602 device px (no stretch), camera kept, click on the Forge still selected it.
- Keyboard: focus on the Hall list button, Enter → Hall selected, live region "Hall selected", focus stays on the button with a visible ring.
- Touch: emulated tap on the Hall at 375×812 selected it.
- Frame timing: programmatic pan plus zoom for 240 frames, measured with `requestAnimationFrame`: average **7.7–17.9 ms**, p95 **10.3–25.9 ms**, max **14.6–40.2 ms** across four headless runs at 1280×800 (WebGL, resolution 1). Headless timing on this machine varies a lot between runs; this is evidence that a pan/zoom frame fits the 16.7 ms budget most of the time, not a measured FPS for real devices.
- Lifecycle: 10 cycles Settlement → Forge → Settlement through the rail: 0 canvases on Forge every time, exactly 1 on Settlement afterwards; JS heap 33.0 → 34.6 MB after forced GC. Browser Back from Settlement to Forge and back restores the right screen and canvas count. Reload of `/game/settlement` is ready in about 0.3 s.
- No horizontal overflow at 375×812, 768×1024 and 1280×800.
- Texture lifecycle (dev server, headless Edge; fixtures are 32×32 solid-colour PNG data URLs generated in the page, no artwork; a dev-only hook can hold a load after `Assets.load` resolves): all 13 checks passed. (1) A loaded image is adopted (one owner, in `Assets.cache`, one sprite, prototype shape hidden) and released on unmount (not cached, no owner, no canvas). (2) A load that resolves after the scene was unmounted is logged as abandoned and unloaded, never adopted. (3) Stage 1 loading, switch to Stage 2: the stage-2 image is adopted, and the stage-1 image arriving later is abandoned and unloaded and never attaches. Stage 2 is released on unmount. (4) The same image on both stages: no `Assets.unload` during the swap, still one owner and one sprite; released on unmount. (5) 10 mounts/unmounts with an adopted image: one owner while mounted, none afterwards, no stray canvas. (6) The missing cottage image still reports the fallback and owns nothing. No unexpected console errors or warnings.
- Production build (published app on 127.0.0.1, headless Edge, cache disabled, no HTTP compression): cold load of `/game/settlement` transferred **915 KiB (0.89 MiB)** in total: HTML 0.8 KB, main JS 251 KB, CSS 35 KB, PixiJS renderer chunks 650 KB; no scene images (none exist). Target ≤ 10 MB met. Scene ready **210–445 ms** after navigation start (three cold loads); reload 212 ms. The dev-only `window.__wooSettlement` handle is absent. The missing-art fallback reports the same status in production.
- Production lifecycle: 30 cycles Settlement → Army → Settlement: 0 canvases on Army every time, 1 on Settlement; JS heap after forced GC 5.5 → 7.3 → 7.7 → 8.0 MiB at 0/10/20/30 cycles. Growth is flattening (+1.8, +0.4, +0.3 MiB per 10 cycles), consistent with PixiJS's global caches warming up rather than a per-mount leak, but a longer soak and a heap snapshot are still to do. Console errors: only the existing `/favicon.ico` 404.
- Not tested: physical iPhone Safari, real pinch gestures (only emulated touch taps), GPU memory and low-end Android.

### Art gate (open, needs the owner)

Steps 1–3 and the art half of steps 4–7 of the [visual proof](../design/village-visual-prototype.md) need owner-approved assets, which do not exist yet: the reference/camera/style sheet, terrain, the forge at two visual levels, two buildings, a tree, wall/corner/gate/tower and bridge artwork. This pass proves the renderer, camera, selection, layering, upgrade swap and image-slot contract that those assets will plug into. Owner decisions still needed: the camera convention (this pass uses the oblique (0.35, −0.6) vector as a proposal), the plot layout and the building roster, the light direction, the shadow treatment, wall-piece connector rules, label style in the scene, and the 4:3 / 4:5 stage proportions.

### Screenshots (untracked, `screenshots/013/`)

`desktop-light`, `desktop-dark`, `mobile-light`, `mobile-dark`, `mobile-dark-inspector`, `tablet-light`, `forge-selected`, `forge-upgrade`, `walls-bridge`, `keyboard-selection`, `zoomed-scene`, `missing-art-fallback`, `renderer-failure-fallback`.

## Validation and administration

Run proportionate checks against acceptance criteria, including relevant failure/retry scenarios. Report actual local evidence. ChatGPT checks remote merge/deployment where applicable. Content/images follow admin configuration and cleanup rules; new mechanics remain code changes.

## Execution record

Task-spec commit: 8a743e4dbe00ab8568c51b84d5f1ce27105c4c39 (owner dispatch, 2026-10-16)
Expected dev base: 8a743e4dbe00ab8568c51b84d5f1ce27105c4c39 (verified equal to local dev and origin/dev, 2026-10-16)
Implementation branch: feat/013-village-ui (local only)
Implementation commit: local commit on that branch; SHA given in the implementation report
Review: Pending ChatGPT review
Merged dev commit: None
Deployment: Deploys through DEV-DEPLOY after integration; not yet run
Issues: first pass only; the art gate is open (no owner-approved village assets). The frontend/design/accessibility skills in the inventory are not installed on this computer; the PixiJS skills are installed but the implementation followed the installed package types and docs directly. No automated frontend tests.

## TLDR

Next: ChatGPT review of the first pass; owner approval of a camera/style reference and first assets for the art gate.
Done: PixiJS settlement scene with prototype geometry, camera, selection, Forge stage swap, image-slot contract, building list and inspector; measured checks.
Issues: art-consistency criteria cannot pass without approved assets; task stays WIP.
