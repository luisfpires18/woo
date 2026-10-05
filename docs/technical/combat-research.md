# Combat: simulation, readable replay and validation research
Updated: 2026-10-05. Status: researched proposals; no combat engine, production art or benchmarks.

## Owner direction and scope
Warfare, conquest and cool armies with powerful weapons are central. The owner prefers the initial central illustrated battlefield and army rosters. Regiment identity, counters and pre-battle stances are accepted direction; exact formation slots, bonuses, targeting, timing and casualty rules remain open. No live micro, fake attacks or manually synchronised waves.

COMBAT-UI (013) owns a local visual/replay proof. COUNTERS (048) chooses and validates numerical rules. BATTLES (049) persists authoritative results safely. REPORTS (050) connects authorised stored events to the accepted presentation. Village, map and resource proofs are separate. No resource catalogue or complete faction roster is inferred.

Research below supports component capabilities and principles. WOO-specific mechanics are recommendations to test, not features supplied by the sources. No code/packages/assets were created or installed.

## Recommended first model
Use a small regiment-level resolver, not an individual-soldier physics simulation. Compare a phase-based model first: deployment, ranged exchange where applicable, engagement, further exchanges and resolution. Exact phases and simultaneous-versus-sequential damage need owner review. A discrete phase model can compute when a battle is due; it does not require a continuously running 60 Hz server.

Separate three things:
- World commitment time: deadlines, troop reservations, travel and actual encounter time.
- Resolution order: chosen discrete phases/steps and numerical rules.
- Presentation time: the length/speed of a replay.

The browser presents committed events. An arrow colliding with a sprite must not decide casualties; pausing the video must not pause the war. No realtime lockstep networking or physical collision engine is justified by this proof.

Fiedler's timestep work explains why simulation must not inherit varying display frame time, while his determinism work warns that fixed steps and a seed alone are insufficient for identical results across environments. For WOO, retain authoritative results/events; the browser does not independently reproduce the server's arithmetic.

## Numerical model questions, before COUNTERS
Pin unit roles separately from movement domain, armour and equipment. Use a few contextual interactions; do not adopt a vast type matrix or universal cavalry bonus. Prepared pikes stopping an eligible frontal charge is a fixture hypothesis, not an unconditional law that pikes always win.

Decide targeting/retargeting, flank access, ranged protection, eligible counter conditions, stacking limits, overkill, rounding, count aggregation and termination. A stance needs an observable tradeoff; do not silently give one stance all benefits. Introduce morale, fatigue, ammunition, dodge chance and terrain bonuses only when they answer a demonstrated gameplay need.

Wesnoth separates terrain defence and damage resistance. That is a useful reference for keeping mechanics distinct, not a reason to copy its dodge probabilities or turn-based formulas. WOO still needs its own equations.

Ordinary equipment upgrades should make a perceptible difference without erasing counters. Role remains recognisable when a weapon changes. The first proof does not approve arbitrary weapon swaps, rune tiers, magical armour or heroes.

If randomness is selected, pin a named/versioned generator, seed, call ordering and stable target tie-breaks. Freeze input/content/rules versions; do not query today's unit stats while replaying an old result. Prefer explicit integer/fixed-point conventions where useful, with overflow and rounding tests. This is a proposal to evaluate, not a guarantee of determinism.

## Explanation and balance
Record observed interactions and applied modifiers with relevant actor/target IDs. Explain facts such as a prepared-pike modifier activating during a recorded charge. Do not turn that into an unsupported statement that the modifier alone caused victory, or invent a percentage contribution without a defined attribution method.

Reports reconcile starting forces, surviving forces and each approved loss category. Wounded/captured/retreated categories are not adopted until their rules exist. Equal headcount and equal budget answer different questions; budget comparisons wait for approved costs.

COUNTERS needs a scenario harness:
| Scenario | Question / invariant |
|---|---|
| Matched forces, sides swapped | Any accidental attacker/iteration-order advantage? |
| Prepared pikes versus cavalry | Does the intended contextual counter activate and stay bounded? |
| Mixed army versus single-role spam | Are several compositions useful under intended conditions? |
| Protected versus exposed archers | Does protection matter without invented invulnerability? |
| Ordinary equipment upgrade | Is its effect meaningful without removing counterplay? |
| Tiny, large, zero and uneven forces | No negative counts, overflow or nontermination |
| Same frozen input/rules/seed | Identical canonical results/events under the supported resolver |
| Simultaneous engagements | No double commitment or casualties applied twice |

Use curated expected cases and invariants first. If RNG exists, evaluate multiple seeds and record distributions, not one favourable example. Save failing seeds. Full faction/cost balance is later BALANCE, not proven by four generic roles. No kingdom roster deletions/merges are authorised here.

## Event record and browser playback
Proposed contract:
| Record | Required meaning |
|---|---|
| Battle identity/version | Stable battle ID, report/event schema and relevant rules/content revisions |
| Initial state | Stable regiment IDs, side, role, approved counts/equipment and formation |
| Event | Ordered sequence, logical phase/time, actor/target IDs, event type and recorded state changes |
| Final state | Authoritative outcome and reconciled force totals |
| Presentation mapping | Versioned phase duration, lanes/anchors, asset references and effect vocabulary |

Keep authoritative audit data separate from the view returned to a player. The server must redact hidden units, upgrades, seeds or event details where disclosure would reveal scouting information. An opaque canvas overlay cannot secure a full unredacted JSON response.

Use one replay clock for regiment movement, frame animation, event text and count changes. Drive visual state from playback position and ordered events. Do not let each AnimatedSprite autonomously run a different time stream if pause/seek needs coordination. Interpolate display poses between events without calculating damage.

For seeking, rebuild from initial state or approved checkpoints and deterministically apply recorded events. Avoid using attack/death callbacks as the only source of count changes: skipped frames and seeking would lose effects. Speed changes, FPS and skip-to-end must give the same report totals. Define equal-timestamp ordering and invalid/unknown-event handling.

A hidden tab can pause presentation or resume from an explicit time, but must not execute thousands of catch-up updates on return. Restart clears transient effects; repeat playback applies no game command. Preserve report access if animation or art fails.

## Readability and controls
Riot's clarity guidance emphasises recognisable silhouettes, visual importance and restrained noise. Apply those principles to WOO: pike, bow and mounted silhouettes need to differ at actual screen size; decisive charge/counter events should draw more attention than ordinary hits.

Keep rosters/counts outside the busy battlefield. Display side, regiment role and selection using labels/icons as well as kingdom colour. Group minor repeated hits; show decisive events and a readable transcript. Decorative projectiles are representative, not a claim that every visual arrow maps to one recorded casualty.

Recommend play/pause, slower/faster playback, restart, bounded phase seek and skip-to-report. Reduced motion can use static formations and event highlights. Game Accessibility Guidelines supports adjustable speed; Xbox guidance addresses camera motion/distractions. Default to no camera shake, avoid forced flashing, and provide a text report independently of animation. Replay controls affect presentation only.

Keyboard/list navigation selects a regiment without canvas precision. On a narrow screen use tabs/sheets without covering the whole battle; preserve the outcome and controls. Do not require hover or sound to understand a counter. Optional sound has independent mute and visible equivalents.

## Authority, persistence and hosting
BATTLES freezes eligible committed inputs under an approved reservation policy. A pure calculation produces result/events; an atomic database commit checks battle/troop versions and records effects. A conflicting troop change must be rejected/retried under a defined policy, not overwritten. Avoid holding a write lock across long calculations.

Microsoft's EF Core transaction and concurrency docs support atomic persistence and conflict detection, but do not by themselves provide battle idempotency. WOO needs database-enforced unique battle/effect identity and tested retry semantics. SQLite lacks database-generated concurrency tokens; use a suitable application-managed token rather than SQL Server rowversion.

Test duplicate due jobs, concurrent workers and crashes before/after commit. Result/event persistence and troop losses must agree. Notifications occur after commit and retries cannot repeat losses/rewards. Do not add a full event-sourcing framework solely to save replays.

Azure App Service can unload an idle app with Always On disabled. The F1 dev plan cannot be treated as a reliable continuous combat clock. Persist due times and define recovery/catch-up in ORDERS; choose a supported worker/scheduling arrangement before promising exact wall-clock execution. UI countdowns display estimates/status and do not resolve battles.

## Rendering and performance
React owns rosters, formation controls, inspector and report. PixiJS owns the composed battlefield and presentation. AnimatedSprite/spritesheets can provide frame animation; they do not supply gameplay events. Start with a capped number of representative figures per regiment and reuse textures. Counts are explicit overlays, not deduced from visible figures.

Use stable ground anchors and test opposing facing; mirror only when shield side/lighting remains acceptable. Batched sprites and restrained effects are preferable to large per-unit filters. Avoid React state updates for every moving figure on each frame. Atlas trimming must preserve original offsets/anchors; independent uploaded assets must remain supported.

Name the target device, viewport, visible regiment/figure count, concurrent effects and art sizes at dispatch. Measure load bytes, decoded texture estimates, frame time, memory across repeat navigation and report/control latency. Do not import benchmark claims as WOO FPS guarantees. A large army can have large numerical counts with a bounded visual representation.

## Decisions and stop conditions
Before 013: owner-reviewed camera/layout reference, fixture events, representation density, art subset and numeric browser budgets. Before 048: approved rules, targeting, casualty categories, rounding and RNG choice. Before 049: commitment/withdrawal, concurrency, recovery and worker policy. Before 050: report visibility, event compatibility and retention policy.

Stop expanding art if facing/animation drifts or roles become indistinguishable. Fix event/report reconciliation before adding effects. If animated art is impractical, review static cutouts plus motion/effects; never mark full animation quality passed on placeholders.

## Primary sources
Checked 2026-10-05:
- [Riot: Clarity in League](https://www.leagueoflegends.com/en-us/news/dev/clarity-in-league/)
- [Wesnoth: defence and resistance](https://wiki.wesnoth.org/Defense_and_resistance)
- [Fiedler: Fix Your Timestep](https://gafferongames.com/post/fix_your_timestep/)
- [Fiedler: Deterministic Lockstep](https://gafferongames.com/post/deterministic_lockstep/)
- [PixiJS AnimatedSprite API](https://pixijs.download/release/docs/scene.AnimatedSprite.html)
- [PixiJS performance guidance](https://pixijs.com/8.x/guides/concepts/performance-tips)
- [AssetPack atlas padding and trim](https://pixijs.io/assetpack/docs/guide/pipes/texture-packer/)
- [Game Accessibility Guidelines: adjustable speed](https://gameaccessibilityguidelines.com/include-an-option-to-adjust-the-game-speed/)
- [Xbox accessibility: motion/distractions](https://learn.microsoft.com/en-us/xbox/accessibility/xbox-accessibility-guidelines/117)
- [EF Core transactions](https://learn.microsoft.com/en-us/ef/core/saving/transactions)
- [EF Core concurrency](https://learn.microsoft.com/en-us/ef/core/saving/concurrency)
- [EF Core SQLite limitations](https://learn.microsoft.com/en-us/ef/core/providers/sqlite/limitations)
- [Azure App Service settings and Always On](https://learn.microsoft.com/en-us/azure/app-service/configure-common)

## TLDR
Next: a small independently reviewed local combat proof; numerical rules follow separately.
Done: researched replay, readability, sprite production, balance and persistence requirements.
Issues: art quality, rules, timing and actual device performance still need practical evidence.
