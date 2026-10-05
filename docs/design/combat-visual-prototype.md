# Combat formation and replay: step-by-step visual proof
Updated: 2026-10-05. Status: proposed preparation; no implementation or production assets.

## Direction and boundaries
The owner prefers the original combat screen: army rosters around an illustrated central battlefield with formation roles and pre-battle stances. Keep clean dark navigation, readable content, kingdom accents and restrained texture. Dark/light themes are independent of kingdom identity.

COMBAT-UI (013) is its own local proof. It does not implement village interiors, world-map editing, resource design, battle authority or conquest. Screen unit names, schedules, numbers and bonuses are illustrative. Four generic medieval fixture roles are not a final roster.

Read [combat research](../technical/combat-research.md) and [art production](combat-art-production.md) before dispatch.

## Small slice
One neutral battlefield, two opposing sides, infantry, archers, cavalry and pikemen. A few bounded representative figures stand for each regiment; labelled counts communicate actual fixture forces. No heroes, rune effects, ships, flying troops, siege destruction or nine kingdom art families.

Pre-battle controls can demonstrate role placement and stance choices, but their numerical effects are not invented. Playback is a recorded fixture, explicitly labelled. It is not a battle engine or live troop order.

## Ordered work
| Step | Work | Exit evidence |
|---|---|---|
| 1 | Pin reference, camera, lanes/formation proposal, fixture IDs/events and budgets | Written source of truth; unapproved rules visible |
| 2 | Build rosters, formation controls and battlefield with symbols | Selection/inspector and keyboard list work without art |
| 3 | Define initial/events/final fixture with schema and stable ordering | Totals reconcile; invalid event references are rejected |
| 4 | Implement one clock and recorded playback | Pause/speed/restart/skip and bounded phase seek agree |
| 5 | Place static cel-shaded role art on clean battlefield | Roles distinguishable; facing, ground contact and overlap pass |
| 6 | Add a minimal approved animation/effect subset | Charge/counter/ranged support follows recorded events |
| 7 | Connect event highlights, roster counts and accessible report | Explanations reflect recorded facts; final totals match |
| 8 | Test narrow screens, themes, touch, keyboard and failure cases | Reduced-motion/static view, missing-art fallback, hidden-tab policy |
| 9 | Review captures, art effort and named-device measurements | Owner verdict on readability/feel; unresolved defects documented |
| 10 | Later connect approved numerical model and stored server events | COUNTERS/BATTLES/REPORTS separately own real game integration |

No R2 or saved admin is needed for steps 1–9. Use local fixtures, approved source images and independent asset slots. Preview replacement locally without claiming it is saved. Owner handles Git image uploads.

## Scenarios to show
- A cavalry approach meeting a prepared pike formation, with a recorded counter highlight.
- Ranged support behind a frontline, followed by a fixture illustrating exposed ranged troops.
- A mixed-army exchange with clear targets, losses and final outcome.
- Restricted opponent information, showing unknown data without revealing it in labels.
- Missing sprites and reduced motion, where the report remains understandable.

These scripted scenarios prove presentation and report consistency. They do not prove pikes always win, validate balance or adopt ammunition/morale/retreat rules.

## Acceptance gate
Recognise role, side and facing at actual display size. Select a regiment using mouse, touch or keyboard/list. Stance/placement controls have clear pre-battle versus committed/read-only state. Formation overlays do not imply playable live movement.

Replaying the same fixture at different speeds/FPS and seeking/skipping yields identical totals and event order. Count updates do not depend on animation callbacks. Repeated playback sends no order and grants no reward.

Important events attract attention without hiding neighbouring formations. Report lists what happened without invented attribution. Missing art or canvas failure preserves text access. Reduced motion removes optional motion; no forced camera shake.

Record source dimensions, anchors, art attempts/repairs, supported facing, viewports, devices, figure/effect counts, load/frame/memory measurements and owner verdict. Resolve exact budgets before dispatch. Placeholders can pass controls but not sprite/animation consistency.

## Fallback and evidence
Compare static cutouts with movement/highlights against minimal frame animation before commissioning an entire animated roster. Keep the representation that achieves approved readability and feel at sustainable art effort. A fallback needs owner review; it cannot silently pass animation criteria.

Store task evidence in its agreed location: fixture/schema versions, reference IDs, screenshot/video identifiers, repaired assets, interaction assertions, timings and limitations. No automatic image commits.

## TLDR
Next: a four-role, one-battlefield local proof.
Done: ordered production and acceptance plan.
Issues: fixture rules are illustrative; no simulation, art or device gate has passed.
