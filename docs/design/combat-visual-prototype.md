# Static battalion board and final result: visual proof
Updated: 2026-10-05. Status: owner-confirmed v1 presentation; layout details are proposals.

## V1 boundary
The owner chose a board showing battalions and a final result. No animation, replay, event timeline or narrated combat log. These can be considered later in separate tasks. COMBAT-UI (013) covers local mock presentation only; COUNTERS/BATTLES/REPORTS later supply validated rules and persisted outcomes. Map, village and resources remain separate.

Use React/CSS for the static board, cards, controls and result. A static battlefield image may sit beneath it. PixiJS is optional if a concrete requirement justifies the added machinery.

## Proposed board
A central board with two opposing sides and clearly labelled frontline, backline and flank slots. Army rosters sit beside it on desktop and in accessible tabs/sections on mobile. Exact slot count and assignment restrictions require review.

A battalion is one logical card/icon with stable ID, role/name, troop count, side and relevant equipment. Do not draw one sprite for every soldier. Infantry, archers, cavalry and pikemen are synthetic proof roles, not faction-roster decisions. Labels/icons supplement kingdom colour.

Selection opens details. Pre-battle placement and stance controls demonstrate agreed direction only; committed armies are read-only unless rules explicitly allow changes. Opponent information can be unknown. Do not expose hidden counts just to fill the board.

## Ordered proof
| Step | Work | Evidence |
|---|---|---|
| 1 | Pin board reference, provisional slots and fixture IDs | Layout sketch and unresolved mechanics visible |
| 2 | Build static cards/icons and rosters using symbols | Role, side, counts and selection readable |
| 3 | Add accessible placement/stance controls where agreed | Editable and committed states distinguishable |
| 4 | Add optional battlefield backdrop and representative static art | Consistent scale, true alpha, stable placement |
| 5 | Build final-result state | Outcome and starting/surviving/lost totals reconcile |
| 6 | Check mobile, themes, keyboard/touch and unknown data | No hover dependence or colour-only identification |
| 7 | Check loading failures and repeated navigation | Symbol/list fallback, correct selection and no stale result |
| 8 | Owner reviews board/result captures and measured behaviour | Continue or needs changes; limitations explicit |

No event fixture, animation clock, seeking or spritesheet is needed. A fixture result is clearly labelled mock; it does not validate balance.

## Result
A concise outcome panel and per-battalion totals. Show only approved categories; wounded/captured/retreated are not invented. No text narration is necessary. Initial and final totals must agree with the fixture. Server integration later authorises/redacts results and supplies frozen historical outcomes.

## Artwork and acceptance
Read [static combat art production](combat-art-production.md) and [combat research](../technical/combat-research.md). Static unit image slots are independently replaceable. Missing art preserves role/count access. Test role silhouettes at actual card size and on light/dark surfaces.

Record viewports/devices, screenshot identifiers, fixture values, art references/repairs and basic load/interaction measurements. Owner handles Git image uploads. Do not expand faction art until the board is reviewed.

## TLDR
Next: a simple static board and result proof.
Done: v1 simplified; replay and animation deferred.
Issues: exact formations and numerical rules remain open.
