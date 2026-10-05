# Combat: static board, numerical model and safe resolution
Updated: 2026-10-05. Status: researched guidance; owner-confirmed static v1 presentation.

## Scope
V1 shows a battalion board and final result. No animation, replay, event timeline or narrated combat log is required. Earlier animated-replay recommendations are deferred. Warfare, conquest and powerful armies remain the central experience.

013 proves local board/result presentation. 048 defines and validates numbers. 049 commits authoritative outcomes. 050 displays authorised final results. Map, village and resources are independent. Exact slots, stances, targeting, timing, casualty categories and faction rosters remain open.

## Static board implementation
React/CSS is the first choice: a responsive grid of formation slots, battalion cards and a result panel. Use stable IDs, accessible selection/buttons and optional placement controls. A backdrop and transparent static unit icons can preserve the illustrated style. There is no need for a canvas, physics engine or spritesheet pipeline solely to show this board.

A battalion is a logical group with explicit count and equipment metadata. Render one representative icon, not one object per soldier. Label side and role beyond colour. Unknown opponent information remains unknown. Editable versus committed/read-only state must be obvious.

Riot's clarity guidance supports recognisable silhouettes, restrained noise and visual importance. Apply that to role icons and selected battalions. Test at real card sizes and both themes. Text narration is unnecessary; the result can consist of outcome and reconciled army totals.

## Numerical model recommendation
Prefer a small pure regiment/phase resolver over individual-soldier physics. Resolution need not run at 60 Hz or depend on browser frame time. Freeze inputs/content/rules revisions and produce an authoritative final result plus sufficient internal audit facts.

Decide targeting, phase/order semantics, simultaneous versus sequential damage, flank access, ranged protection, counters, stacking limits, overkill, rounding and termination. Accepted counter direction does not mean pikes always win or cavalry always beats infantry.

Keep role, movement domain, armour and equipment distinct. Wesnoth's separation of terrain defence and damage resistance is a useful reference, not a formula to copy. Morale, fatigue, ammunition, dodge and additional layers need separate decisions, not automatic inclusion.

If RNG is selected, pin generator/seed, call order and stable tie-breaks. Seed alone is not a determinism guarantee. Integer/fixed-point conventions may help but require overflow/rounding checks. The server is authoritative; the browser does not recalculate results.

## Model validation for 048
| Scenario | Check |
|---|---|
| Equal forces with sides swapped | No accidental iteration/first-side advantage unless intended |
| Prepared pikes versus cavalry | Contextual counter activates as defined, without universal immunity |
| Mixed army versus single-role spam | Useful alternatives under approved conditions |
| Protected/exposed ranged roles | Protection matters without invented invulnerability |
| Ordinary equipment upgrade | Meaningful improvement while preserving counters |
| Zero/tiny/large/uneven forces | No negative counts, overflow or nontermination |
| Frozen inputs/rules/seed | Same canonical outcome in supported resolver |
| Concurrent encounters | No double commitment of troops |

Use curated cases and invariants. If RNG exists, test distributions and preserve failing seeds. Equal-headcount and equal-cost tests answer different questions; costs await the owner-defined economy. Full faction balance remains BALANCE. No roster merges/deletions are authorised.

## Persistence and concurrency for 049
Define commitment/withdrawal and troop reservation rules before real orders. Persist battle ID, due time, frozen input/rules references and final result. Calculate without holding a long database write transaction, then atomically check versions and apply accepted losses with result/audit persistence.

EF Core transactions and concurrency mechanisms support this, but WOO still needs database-enforced idempotency and uniqueness. SQLite lacks database-generated concurrency tokens; use a suitable application-managed token. Test duplicate jobs, concurrent workers, stale input versions and crashes before/after commit. Losses and approved rewards apply once. Viewing a report applies nothing.

Keep sufficient audit data to investigate outcomes; a player-facing event stream and an event-sourcing framework are not required. Notifications follow commit and may retry without repeating effects.

Azure App Service may unload an idle application with Always On disabled. Persist due times and define recovery/catch-up under ORDERS; do not depend on an in-memory timer or UI countdown. Exact-time worker arrangements remain a hosting decision before live timed combat.

## Final results for 050
Display outcome, starting forces, survivors and approved loss categories. Reconcile totals. Wounded/captured/retreated remain unapproved until their rules exist. Compact recorded modifiers may be optional detail, but no narrative log or exact causal percentages are required.

Authorise/redact on the server before returning opponent data. Store historical results rather than recomputing using current balance. Define result schema/version compatibility and retention. Handle missing images, failed fetches, switching battles and unknown data with usable static fallback.

## Deferred visual research
AnimatedSprite, spritesheets, unified replay clocks, pause/speed/seek and event-driven playback were researched earlier. They are optional future enhancements, not v1 acceptance criteria. If reconsidered, scope a separate task after the board and numerical combat prove enjoyable. No animation asset generation now.

## Primary references
Checked 2026-10-05; capabilities and principles, not measured WOO feasibility:
- [Riot: Clarity in League](https://www.leagueoflegends.com/en-us/news/dev/clarity-in-league/)
- [Wesnoth: defence and resistance](https://wiki.wesnoth.org/Defense_and_resistance)
- [Fiedler: simulation versus display timing](https://gafferongames.com/post/fix_your_timestep/)
- [Fiedler: deterministic lockstep caveats](https://gafferongames.com/post/deterministic_lockstep/)
- [EF Core transactions](https://learn.microsoft.com/en-us/ef/core/saving/transactions)
- [EF Core concurrency](https://learn.microsoft.com/en-us/ef/core/saving/concurrency)
- [SQLite limitations](https://learn.microsoft.com/en-us/ef/core/providers/sqlite/limitations)
- [Azure App Service settings](https://learn.microsoft.com/en-us/azure/app-service/configure-common)
- [Optional future AnimatedSprite](https://pixijs.download/release/docs/scene.AnimatedSprite.html)
- [Optional future playback accessibility](https://gameaccessibilityguidelines.com/include-an-option-to-adjust-the-game-speed/)

## TLDR
Next: static board/result proof, then bounded numerical rules.
Done: animation and replay removed from v1 scope.
Issues: combat formulas, formations and timing remain to define.
