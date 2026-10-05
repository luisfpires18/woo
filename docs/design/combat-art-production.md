# Combat battlefield and sprite production
Updated: 2026-10-05. Status: proposed asset workflow; owner-generated cel-shaded art expected.

## Required composition contract
Approve one camera, scale, outline treatment, light/shadow direction and colour reference. Combat sprites need their own role/facing readability; a village building camera or character portrait does not automatically fit. The owner supplies/approves images and handles Git uploads.

Begin with one infantry master on one ground patch. Inspect facing, feet, weapon clearance, contrast and size before expanding to archers, pikemen and cavalry. Mounted sprites need rider/horse contact and a consistent ground footprint. Character identity must survive each animation pose.

## Asset separation
| Asset | Contract |
|---|---|
| Battlefield ground | Clean terrain, no troops/UI/counts/baked faction ownership |
| Regiment/unit cutout | True alpha, stable ground anchor, readable weapon silhouette |
| Optional frame sequence | Same character/camera/scale; preserved padding and named action |
| Effects | Separate arrow/dust/impact; limited visual importance |
| Live overlays | Names, counts, role/status, selection and events rendered in UI |
| Fallback symbol | Available even if image request or animation fails |

Decoration must not imply a real terrain bonus unless the fixture/model records one. A forest painted behind the battlefield is not automatically cover.

## Minimal progression
1. Static infantry cutout, symbolic other roles and clean ground.
2. Four approved static roles; compare roster readability at minimum play size.
3. Move representative formations with limited impact/projectile effects.
4. Add a small action set for representative roles only: idle, advance and attack/impact where useful. Frame count and duration are selected after testing, not mandated now.
5. Add removal/fade or an approved casualty pose if it improves explanation.
6. Expand animation coverage only after owner review of assembled footage and production effort.

Do not make full walk/run/attack/death sheets for every kingdom before this proof. Existing sprite skills may help generate consistent source frames, but their Unity export and fighting-game motion rules are not automatically the WOO browser specification. No sprite-generation skill is invoked by this documentation update.

## Per-asset metadata
Stable logical ID, role/variant, revision, dimensions, original canvas dimensions, ground anchor, footprint, allowed facing/mirroring, scale, frame names/durations, animation loops and layer. Record source/usage rights, approved reference and repair effort.

If trimming/packing, preserve original frame offsets so feet do not jump. Atlas rotation is storage packing, not permission to rotate the displayed fighter. Opposite facing may need separate art: mirroring swaps shield/weapon hands and light direction. Approve that explicitly rather than silently flipping all units.

Use transparent PNG sources first. Check alpha on charcoal, ivory and real terrain. Remove painted checkerboards, unwanted glow halos and background remnants. Reserve space for a spear, bow draw and horse without clipping. Do not resize each frame independently to its visible bounding box.

## Regiment representation
Use a capped display count per regiment and reuse one approved texture/frame set. Visible figures are representative; overlays state fixture numbers. Disappearing figures must not falsely claim one-to-one casualties. Use consistent display-count mapping and exact report totals.

Offset repeated figures within a fixed formation footprint; avoid uncontrolled random jitter and overlaps. A single representative figure per group is an acceptable first comparison. Select the regiment as a group with an accessible roster entry, not tiny soldier targets.

## Quality gate
Role/side/facing readable; no floating feet or doubled shadow; no frame-to-frame size/identity drift; no spear clipping; cavalry remains distinguishable; important hits visible without effect overload. Inspect movement at normal and slow replay speed on desktop and mobile.

A static image is not proof of animation consistency. A nice GIF is not proof of replay seeking, casualty totals or browser performance. Record both art review and interaction evidence in COMBAT-UI.

## Replacement and runtime
Independent admin image/animation references follow ASSETS and IMAGE-CACHE later. Source frames, atlases and necessary derivatives need reference-aware cleanup. Do not introduce a mandatory giant atlas for all runtime uploads. Local replacement in 013 proves visual flexibility only; saved administration remains later work.
