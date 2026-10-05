# Historical brainstorming snapshot

Updated: 2026-10-05. Status: design documentation; no game implementation yet.

Imported from WOO_TRAVIAN_ROADMAP.md version 12 on 2026-10-05. This preserves earlier discussion, including superseded statements. It is not the current specification. Read [scope](../scope.md) and [decisions](../decisions/README.md) first.

# Weapons of Order: Seasonal Browser Strategy Roadmap

Last updated: 2026-10-05 (Europe/Lisbon)
Status: brainstorming; no development authorised by this document.

## Read this first in any continuing chat

This is the working record for a browser strategy game based on Weapons of Order (WOO). Travian is a genre reference, not a design to copy. Battle for Middle-earth (BFME) is a requested inspiration. Preserve explicit user decisions, distinguish story canon from game adaptations, and label unapproved ideas as proposals. Do not treat uploaded development drafts as canon, even where they call themselves authoritative or implemented. Repository implementation claims have not been verified.

Historical snapshot only. Current repository topic documents and decision register supersede this record.

## 1. Explicit user direction

- Build an original game in the persistent browser strategy genre.
- Base kingdom identities and weapon rules on WOO, while treating earlier game drafts as exploratory.
- Current player identity: master blacksmith, NOT the ruler (user correction 2026-10-05). The kingdom ruler is most likely a separate unique hero; exact roster and inclusion remain open.
- Runeforging is central. Discover existing Chaos weapons on the map through seasonal gates; later discover/create Order weapons to counter Chaos. Exact gates and victory rules remain open.
- No Travian-style player alliances, fake attacks or manually synchronised attack waves.
- Players must defend their chosen kingdom.
- Losing a territory/village transfers territory to the conquering kingdom and expands its borders.
- Investigate map distribution, daily missions, animal camps and BFME mechanics.
- Maintain this roadmap for later development and access from other chats.

## 2. Canon references and limits

Original canon filenames are listed in ../world/lore-constraints.md. Private file identifiers omitted.

User identified these attachments as story canon, subject to their own WIP labels:

| File | Library identity |
|---|---|

The kingdoms, units, buildings, core-mechanics and scaling drafts are idea sources, not confirmed game specifications. Their Go/backend implementation claims do not establish the stack for this new game.

### Canon constraints relevant to gameplay

- Natural runes are found as Runestones. Only Chaos and Order are crafted principles; neither functions alone.
- One rune per weapon defines its identity. Elemental fusion occurs during runeforging (e.g. Fire + Water = Steam).
- L0 Dormant, L1 Enhanced/Conduit, L2 Artifact/Aspect. L2 has staged embodiment at 25/50/75/100%.
- L3 Dreadform/Ascendant requires Chaos/Order fusion; soul requirements conflict across files and need clarification.
- Mystic runes, mythical animal runes, and primal Light/Dark have uniqueness constraints.
- Rune families constrain weapon forms: Mystic to staff; Animal to claw gauntlets; Physical to pummelers/three-section staff; Material excludes fist/staff/flexible categories.
- Blood infusion binds a weapon to its wielder. Soul infusion is distinct, incompatible with blood infusion, and links weapon/wielder death. Soul entrapment is another technique.
- Chaos weapons contain trapped spirits, corrupt the bearer and can manifest brief giant Dreadforms with lasting terrain effects.
- Order list is explicitly WIP with some entries labelled confirmed. Purified runes are future Book 2 discoveries.
- Technic runes belong to the later Thalori return.
- Moraphys is a devastated peninsula with Null/Residual zones and unstable routes, not an established central NPC empire.

### Clarifications to preserve

1. Order described as containing no soul versus soul-entrapment/sacrificed-soul requirements for L3.
2. Victura: Necro versus Necrosis.
3. Vantashields: Meteorite/Metorite versus Tarnish.
4. Hailbane: Ice versus Frost.
5. Entropy and Protect appear in weapon descriptions but not the corresponding rune lists.
6. Light location unknown in rune list versus installed in Sol's heliostat in map reference.
7. Preserve geographical and temporal distinctions within canon; do not silently reconcile contradictions.

## 3. Recommended game structure

Acceptance update (2026-10-05): User accepted the preceding recommendations for village-centred map districts, homeland/frontier/expedition areas, useful Kingdom Duties, differentiated animal/monster/mythical camps, and BFME-inspired battalions, pre-battle stances, fortress modules, hero progression and capture objectives. These are accepted design directions, not final mechanics or numbers. Exact geography, population distribution, balance, rewards, victory and lore-dependent rules remain open.

### Kingdoms and player cooperation

Kingdom = competing faction, culture, roster and visual identity. All players in that kingdom cooperate through shared objectives, defence requests and contribution tools. No independent alliance layer. Avoid allowing one player to command or confiscate everyone else's troops.

Membership stays fixed during a season, apart from carefully constrained balancing options if needed. Population balance must consider active participation, not just registration totals. Use join incentives/capacity checks before granting emergency combat buffs.

### Map distribution

Use Bellum's canon geography as the macro layout, with procedurally varied sites inside designed regions. Homeland labels and political ownership are separate layers: Arkazian conquest does not physically turn a forest into mountains.

- Sylvara: forest heartland, rivers and meadows.
- Arkazia: mountain realm west of the Dark Reach; defensible passes and crossing routes.
- Veridor: northern inland roads/rivers and southern coastline/ports.
- Draxys: rocky desert, wadis, oases and settled cores.
- Nordalh: northern fjords, green belts, ice and hot springs.
- Zandres: subterranean districts with limited surface entrances; full layered travel can wait.
- Lumus: island with defined maritime connections; prevent isolation or dependence on Veridor for all access.
- Drakanith: volcanic region; playable/NPC status not decided.
- Moraphys: hazardous late expedition region; not automatically the map centre or a playable kingdom.

Starting recommendation: 55% secure homeland districts, 30% frontier districts, 15% neutral expedition/hazard districts, measured across playable land districts. These are initial test proportions, not canon. Capital districts may be protected initially; capital immunity and elimination remain open.

Use village-centred districts as conquest units. Borders expand on district capture. Connect conquest through adjacent districts or explicitly designed naval routes. Multiple frontier routes and some bypasses prevent one permanent bottleneck. Balance usable village capacity, early materials, objective access and travel times rather than identical land area. Do not choose exact coordinates without the canonical map geometry.

A small prototype should use 3 factions, several frontier villages and one expedition region; final faction count remains undecided. Real path costs should account for intervening terrain and route connectivity.

### Conquest and recovery

Proposed declared sieges with visible commitment and battle windows, replacing fake waves and precise manual timing. Players pledge contingents and tactical orders; battles resolve automatically. Limit simultaneous active campaigns so participation does not require constant attendance.

On capture, district/village ownership changes kingdom. Exact defender outcomes are open. Recommended refuge/reassignment preserves blacksmith progression and personal equipment while losing some local infrastructure/resources. Define infrastructure persistence, occupation recovery and the new governor; never assume automatic village destruction.

### Daily missions: kingdom duty board

Generate useful tasks from actual world conditions rather than arbitrary chores:

- Repair a threatened village gate or produce supplies for its garrison.
- Forge equipment requested by a frontier contingent.
- Scout an approaching siege route or unexplored rune site.
- Escort a convoy, clear a dangerous den or survey a ruin.
- Recover wounded troops after a kingdom battle.

Offer several options; completing any 2-3 is sufficient for the intended daily reward budget. Allow a small backlog (e.g. 3 days) and no punitive login streaks. Rewards: modest supplies, crafting experience, contribution reputation and cosmetic progress. No unique rune/Order progression locked behind daily attendance. Shared task reservations prevent wasted duplicate donations. Give support and defensive work credit, not only kills.

### Animal camps and expeditions

Distinguish ordinary dens, dangerous monster sites and unique mythical encounters. Regional ecology follows canon tendencies. Ordinary animals are not all magical, hostile or rune sources.

Missions can include scouting, driving predators away, rescue/escort, observing a habitat, or retrieving a Runestone discovered at a site. Animal-rune acquisition details require confirmation; do not assume killing an animal generates its rune. Mythical runes remain unique, not endlessly respawning loot.

Publish difficulty, estimated losses and mission duration after scouting. Use bounded expedition rewards per player/time period to prevent nonstop farming. Camps replenish slowly or relocate; returning players should have viable objectives. No kill stealing: ownership/reservations or individually credited expeditions where appropriate. Taming and mounts require separate lore and gameplay decisions.

## 4. BFME inspiration assessment

Research inspected BFME II manual reproduction, developer Q&A, and contemporary guide. Exact balance is patch dependent; adaptations below are proposals, not claims of identical BFME mechanics.

| Inspiration | WOO adaptation | Priority |
|---|---|---|
| Battalions and combat stances | Named regiments; aggressive/balanced/hold-ground orders selected before combat | High |
| Unit counters | Pike/cavalry/archer/frontline roles with terrain and rune interactions | High |
| Fortress expansion plots | Choose limited village defensive modules: towers, gateworks, recovery station, forge ward | High |
| Hero progression | Personal blacksmith abilities, leadership and weapon mastery | High |
| Capture-and-hold objectives | Watchtowers, crossings and supply depots embedded in village campaigns | High |
| Battalion veterancy | Surviving regiment identity, modest capped benefits and recruit replacement | Medium |
| Creep encounters | Animal/monster expeditions with clear risks and bounded rewards | Medium |
| Power trees | Small seasonal doctrine choices; avoid global destructive clicks and kill-only snowballing | Later |
| Ring hunt | Inspiration for transporting/contesting a unique relic; WOO Chaos rules remain authoritative | Later |
| Free-form live RTS building/micro | Would change scope and daily attention substantially | Defer |

Avoid universal leadership stacking, permanent runaway veterancy and wipe-everything powers. Defences and attack plans should remain legible without real-time attendance.

## 5. Forging and season progression proposals

Settlements support blacksmiths; blacksmiths shape military identity; expeditions and conquest supply knowledge/materials. Do not let settlement management drown out forging.

Suggested season 10-12 weeks: foundations, discovery, Chaos access, Order breakthrough, final conflict. Dates, duration and victory are unapproved. Order research must remain accessible to kingdoms without a Chaos weapon. Several relic access opportunities reduce first-arrival monopolies.

Chaos custody, equipping and Dreadform activation are separate. Proposed predictable corruption costs reflect possession/consumption, rather than arbitrary offline destruction. Personal runeforged equipment versus regiment-wide rune equipment needs clarification against rune scarcity.

Victory must suit kingdom territorial warfare and the Order/Chaos story. Options still open: final territory/objective score with forging contribution; stabilising multiple ancient sites; confronting a scenario-specific threat. Earlier player-alliance victory proposals are superseded. Moraphys holding every Chaos weapon is not adopted.

Season reset proposal: reset villages, troops, territory and unique relic possession; retain history, cosmetics and titles without combat bonuses.

## 6. Technology recommendation (not selected/implemented)

React + TypeScript + Vite UI; PixiJS village/map; ASP.NET Core API; PostgreSQL; durable scheduled events processed by .NET Worker; SignalR notifications; Docker; object storage/CDN; off-site backups. Begin with modular monolith and separate API/worker processes, no Kubernetes requirement.

Village rendering: illustrated background, predefined building slots, sprites with limited visual upgrade stages. World rendering: viewport tiles/chunks, terrain/location/ownership/activity layers, zoom-dependent detail. Use consistent perspective, anchors and transparent sprites. Server owns visibility and all game decisions.

## 7. Player research implications

Qualitative Travian research suggests recurring concerns around spending advantage, attendance pressure, repetitive administration, onboarding, cheating/enforcement trust, mobile parity and player density. These are not representative survey results.

Recommendations: no purchased combat power; useful queues and shared truce; complete mobile browser flow; understandable battle reports; clear enforcement/appeals; recovery after defeat; support contributions. Daily play target proposed at 20-40 minutes across 2-3 visits, not a guarantee. Test with competitive veterans, lapsed players and newcomers.

Draft arithmetic defects: resource starting rate inconsistent (3/s versus formula 7/s); starting 1,600 storage fills from 500 in about 2.6 minutes at 7/s; 9,200 maximum storage cannot afford approximately 199,500 lumber for a level-20 field. Existing numbers are not adopted.

## 8. Development path, pending design validation

1. Resolve foundational rules: playable factions/era, conquest/defeat, combat windows, weapon bearers, Order soul conflict, victory.
2. Simulate economy and automatic battles with a small unit/rune set.
3. Build vertical slice: village, forge, one expedition, one declared border siege and ownership transfer.
4. Closed short test with 3 kingdoms; measure active faction balance, recovery, attendance and forging usefulness.
5. Expand map, weapon repertoire, daily board and seasonal gates based on results.
6. Public season only after reliability, recovery, moderation and mobile checks.

## 9. Research links

- BFME II manual reproduction: https://manuals.plus/m/f958931bd298d60eacdaad759963ecae735f4f96b71f848404bc8b5468796c88
- BFME II developer Create-a-Hero Q&A: https://www.gamespot.com/articles/the-lord-of-the-rings-the-battle-for-middle-earth-ii-qanda-fan-questions-part-i/1100-6133390/
- BFME II contemporary guide (secondary): https://www.gamespot.com/articles/the-lord-of-the-rings-the-battle-for-middle-earth-ii-walkthrough/1100-6146219/
- Travian village management feedback: https://blog.travian.com/2024/06/travian-loop-village-management/
- Kingdoms community feedback: https://blog.kingdoms.com/feedback-loop-community-topics-3/
- Kingdoms night protection: https://support.kingdoms.com/en/articles/3-game-versions-speed-and-specials
- Kingdoms world concentration: https://blog.kingdoms.com/new-game-world-schedule/

## 10. Decision log

- 2026-10-05: Initial research, genre assessment and draft review.
- 2026-10-05: User accepted ruler-blacksmith identity.
- 2026-10-05: User rejected independent alliances/fakes/wave planning; selected kingdom defence and village-territory conquest.
- 2026-10-05: Canon attachments supplied; discrepancies retained for clarification.
- 2026-10-05: Roadmap created with map/duty-board/animal-camp/BFME recommendations marked as proposals.

## 11. Screen concepts

2026-10-05: User requested simple, clean and understandable concept screens. Initial visual exploration covers Resources, Village, Kingdom Map, Combat, Forge and Kingdom Duties in a consistent Arkazia-themed interface. Mockup numbers, layouts, timing and ore/material choices are illustrative rather than final balance. Map is a local frontier illustration, not a new canonical world geography. Forge shows one rune per weapon. Combat depicts advance commitments and stances, with no fake attacks or live-micro requirement. Generated visuals are concept previews, not implemented UI.

- 2026-10-05: User accepted map/duties/camps/BFME recommendations; requested six initial screen concepts.

## 12. Visual production roadmap without a designer

User direction (2026-10-05): Artwork is the principal project concern. Build an explicit quality-controlled workflow using ChatGPT-generated assets uploaded into the game, initially as clickable scenes. Support independent dark/light appearance and kingdom identity. User specifies Arkazia crimson/black, Sylvara green/gold, Veridor blue/silver. Exact shades and other kingdom secondary palettes remain open; this is game visual direction, not a rewrite of story canon. No code examples needed in brainstorming.

### Theme system

Keep neutral surface/text/border tokens separate from kingdom primary/secondary tokens and semantic warning/success/enemy colours. Dark mode uses charcoal/slate surfaces; light mode uses ivory/light stone surfaces. Kingdom identity accents navigation, buttons, selected borders, banners and small decorations. Black is not Arkazia's text colour everywhere; silver/gold are not low-contrast body text. Ownership uses kingdom colours plus labels/patterns/icons. Same layouts and artwork work in both themes; do not regenerate every image for light/dark. Do not darken scene art globally to simulate dark UI. Verify text contrast, focus, hover and selected states.

### Phase A: establish a visual standard

Approve one Arkazia master village scene plus a representative building, map marker and regiment portrait. Record camera angle, lighting, line/shading treatment, architecture, palette, scale and detail budget. Maintain reference images, reusable prompt briefs and versioned approved assets. Existing LF2 fighter style is not automatically the environment style.

### Phase B: asset import and interaction editor

Development proposal: upload PNG/WebP scenes; preserve originals, validate dimensions/alpha/file size, create display derivatives, and version replacements. Draw labelled polygon hotspots over buildings or districts; store coordinates normalised to image dimensions. Configure hotspot actions, tooltip, highlight and panel. Preview desktop/mobile and both themes. Hotspots follow image transforms on resize/zoom; they do not depend on guessed screen pixels. Validate overlapping shapes, touch targets and keyboard-accessible companion lists. Changes to art geometry require hotspot review; never assume regenerated scenes preserve positions.

### Phase C: clickable village first

Start with one complete village painting WITHOUT embedded labels/buttons/resource counters. Add UI labels, selection outlines, construction badges and details in the app. Buildings are polygon hotspots on the painting. Live levels/queues are data overlays, not painted values. This provides rich visual quality with fewer independently generated parts. It does not yet provide visually interchangeable individual buildings. Later create an empty foundation scene and separate matched transparent buildings, using fixed footprints/anchors and limited visual stages. Test one replacement building before committing to modular conversion.

### Phase D: illustrated map with authoritative geometry

First design district geometry and routes from approved geography; create a corresponding reference layout for the illustration. Start with one bounded regional map painting. Ownership colours/borders, village markers, sieges and expeditions remain dynamic overlays. AI art cannot determine adjacency, traversability or political boundaries. Annotate terrain masks/routes explicitly. At conquest change overlays; keep geography. If a full scene has baked flags/ownership, it is unsuitable for changing political control. Provide unflagged terrain art or separate settlement patches. Limit initial zoom; generate regional detail/chunks only after testing seams and coordinate alignment. Seamless procedural tiles are deferred until a small tileset proves workable.

### Phase E: combat visual prototype

One landscape plus regiment tokens/portraits, stance indicators and formation slots. React panels show commands and casualties; rendered scene previews deployment. Replay uses authoritative server events. Add a few software-created effects (movement, projectiles, flashes, rune halos); defer individually animated soldiers and giant transformation sequences. Portraits are achievable assets but don't imply coherent directional animation is solved.

### Phase F: visual quality gate and expansion

A scene is accepted only after it works at normal screen size, mobile scale and both themes, with crisp artwork, coherent style, correct interaction alignment and understandable states. Review asset alpha/edges, unwanted text/symbols, blur, geometry drift and incorrect canon traits. Reject/repair individual failures before generating the next family. Asset tracker stores ID, kingdom, role, prompt/reference, source version, status, dimensions, anchor/hotspots, usage and known issues. Expand to Sylvara and Veridor only after Arkazia's village/map/combat workflow is proven.

### Limits and next milestone

Generated art can support static clickable scenes; consistent modular art and animation require iteration and may need cleanup. Flawless output cannot be promised. Next proposed milestone: one real Arkazia village screen using an approved clean scene, six clickable locations, a functioning details panel and a dark/light toggle. No full game development authorised by this roadmap update.

## 13. Accounts and administration

User requirement: login and an administrator account for the owner, with a dedicated admin workspace to change game aspects quickly without code edits. Recommended stack, not yet formally selected: ASP.NET Core Identity, server-enforced policies, secure HttpOnly browser cookies, React/TypeScript workspace and versioned PostgreSQL configuration. Scope includes worlds/seasons, balance, content, assets/hotspots, map editing and player support. Draft/preview/publish/rollback and audit history; declare whether changes apply immediately, to new actions, or next season. New mechanics still require code; supported effects and content become configurable. Owner admin role is assigned securely server-side, never through public registration.

## 14. First alpha season: user-defined scope

2026-10-05 explicit direction supersedes earlier open faction-count/NPC-role proposals:

- Closed alpha with the owner and friends.
- Playable: Arkazia, Veridor, Sylvara.
- Other future playable kingdoms (Draxys, Nordalh, Zandres, Lumus) appear greyed out and cannot be selected. Their degree of map presence remains open.
- Drakanith is NPC-only, neutral to all for now, with some outposts near volcanoes. Its Drakani are human-dragon hybrids. Neutral does not yet define whether players may attack or trade with them.
- Moraphys is NPC-only and hostile to all, with some outposts throughout the world. User describes its visual identity as evil black grass. This is a game scenario direction; retain the canon homeland's devastated peninsula/Null/Residual geography rather than silently rewriting it.
- Some Chaos weapons may initially belong to Moraphys characters. Other weapons' locations, availability, custody and transfer rules remain undecided.
- User will define the characters entering the world. Do not invent named NPCs or assign canonical Chaos weapons to them.
- Clarified: user meant HEROES, not artists or Artifacts. They are considering named story characters, e.g. Nail Ark or commander Akron Wright, stationed at the main building to help fight. Inclusion is still under consideration, not a confirmed alpha feature. User will define roster, era, titles and assignments.

Recommendations only: distribute friends across the three playable factions; use frontier NPC defenders where player numbers are insufficient; make initial Chaos bearer access subject to agreed season gates rather than granting early capturable weapons accidentally. NPC behaviour complexity remains open.

### Named kingdom heroes: proposal pending decision

Keep named heroes distinct from each player's master blacksmith. Recommended unique kingdom-level characters, not a duplicated Akron in every village. Main building exposes a kingdom hero panel with current station, availability, role and recovery status. A hero can join one defence/campaign at a time; deployment and selection rules remain open. For the alpha, start with at most one hero per playable kingdom if adopted, using fixed simple roles, capped leadership and one signature ability. Avoid manual last-second activations, purchases of power, universal bonus stacking and autonomous deployment by an unrestricted first-click race. Nail Ark as a leadership/frontline example and Akron Wright as forge support/battlefield example are proposals, not fixed stats or confirmed roster. Preserve player forging relevance; NPC masters must not supply all progression automatically. Temporary recovery after defeat recommended instead of permanent canonical death in alternate seasonal scenarios.

## 15. Medieval start and Artifact bonding: updated user rules

2026-10-05: Begin with normal medieval forging and equipment, gradually introduce runes and advanced runeforging. Ordinary uninscribed weapons are distinct from L0 Dormant weapons with an inscribed rune. L1 Conduit and L2 Aspect produce advanced military expressions: Animal L2 shapeshifters of the relevant animal, Mystic L2 unique rune-bearing wizards (nine Mystic runes), and Nature L2 expressions matching the essence.

Explicit L2 rule: Artifact forging requires the intended wielder's blood in the process, creating a soldier-weapon bond. Blacksmith skill determines recoil/self-damage during use. Ordinary smiths may produce weapons that damage their bearer with each swing; the best master can perform the process without side effects. This clarification supersedes earlier generic assumptions about L2 bonding/recoil. Do not independently decide its relationship to the older Soul Infusion passage: preserve that as a lore clarification to resolve.

L3 remains exclusive to Chaos/Order weapons. User anticipates a similar approach but has not defined it. Do not automatically give Chaos and Order identical recoil/blood requirements or let smith skill remove Chaos's trapped-spirit corruption.

Design proposals, not accepted rules: recruit ordinary troops and advance compatible bearers via crafted weapons and training, rather than mass-producing all L2 unit types directly. Display forecast Artifact power, recoil, intended bearer and crafting requirements before committing. Translate per-swing recoil into actual combat resolution phases/attacks; do not base gameplay damage on cosmetic replay frames. Skill progression should reward meaningful forging/training rather than unlimited cheap-item spam. Reforging or improving an existing bonded Artifact could preserve invested progression, pending blood/ownership rules. Define whether ordinary smiths have L2 access, how forging difficulty affects skill, survivability and healing interactions, and what happens to a bonded weapon after bearer death.

## 16. Player identity correction

2026-10-05: User replaced the combined ruler-blacksmith concept. The player is a master blacksmith, not a ruler. The ruler will most likely be a separate unique kingdom hero. Earlier player-rulership proposals are superseded. Systems will be refined incrementally. Exact starting mastery remains undecided; master blacksmith is the player identity, not an automatic maximum skill value. Authority over villages, troops and kingdom decisions must be defined separately. Ruler roster and leadership rules remain open.

## 17. Smith specialisations and possible armour requirement for L2

2026-10-05 user direction: smithing has Weaponsmith and Armorsmith specialisation categories. Runes enter Weaponsmith through runeforging. Armorsmith's advanced rune system remains a standby concept. Whether a player chooses one, learns both, or splits mastery is undecided.

User proposal (not yet final): armour does not grant independent powers, but both a runeforged weapon and runeforged armour may be necessary to achieve the very powerful L2 Aspect state. Do not silently adopt this as canon or define rune counts/identities for armour. Existing L2 weapon blood-bond and craftsmanship-dependent recoil rules remain recorded; their interaction with armour requires design.

Recommended exploration: weapon determines the rune identity and abilities; compatible armour supports safe embodiment and aura containment rather than adding another ability set. Weaponsmith mastery governs weapon execution/recoil; armorsmith mastery may govern protection/compatibility/stability, with exact interaction open. Distinguish ordinary physical armour protection from proposed advanced Aspect-support equipment. Resolve how armour behaves during Animal shapeshift, Physical fusion and Material living-armour Aspect; Mystic robes may need compatible crafted attire rather than universal plate. These are proposals. Alpha can retain ordinary armour crafting while advanced Aspect armour is deferred; do not assume alpha exclusion unless user confirms.
## 18. Magical armour classes and Aspect transformations

2026-10-05 user clarification: armour is magical. Always retain Light, Medium and Heavy categories, with cloth, leather and plate equivalents respectively. The categories persist as fantasy forms evolve; exact mechanics and rune compatibility remain open.

User examples: Mystic armour uses cloth; leather transforms into fur for an Animal shapeshifter and into vines for the Life rune. These are examples, not a complete compatibility table. Define each scenario later; do not assign unconfirmed heavy-armour transformations.

This supersedes any suggestion that advanced armour is merely mundane. Magical transformation does not itself confirm independent armour powers. The possible weapon-plus-armour requirement for L2 remains exploratory; armour rune identity/count, bonding and interaction with weapon recoil remain unresolved.

## 19. Mandatory admin-managed visual assets

2026-10-05 explicit requirement: every visual entity (units, buildings, runes, weapons, armour, heroes and other relevant content) must have a configurable 2D sprite/image slot in the admin workspace. The slot and fallback support are mandatory; an uploaded image is not required. Use an emoji or symbol until artwork is supplied. Asset changes must work without code edits.

User intends Amazon S3-style storage through Cloudflare (interpreted as Cloudflare R2, exact service to confirm when the account is created). User will create the account later; no provisioning requested now. Store asset references in game data and bytes in the bucket.

Replacement policy: upload and validate the new image, save the entity reference successfully, then delete the previous image if no other entity uses it. Failed replacements preserve the working image and clean up the unsuccessful upload. Entity deletion or clearing its image also removes unreferenced objects. Shared assets require reference checks before deletion. No permanent historical image copies or unused uploads; this supersedes earlier asset-version retention proposals, while configuration/audit metadata may still be retained. Keep every retained object associated with content, with cleanup retries/reconciliation for interrupted operations. Brief overlap during a safe replacement is acceptable; steady-state orphan files are not.

Implementation guidance: use a new object key for each successful replacement so cached artwork updates reliably; expose replacement as one admin action. Store role-specific display geometry (dimensions, anchor and relevant hotspots) and review it when changing artwork. These are implementation recommendations serving the confirmed upload/cleanup requirements, not new gameplay rules.

## 20. Player control, unit equipment and outstanding season rules

2026-10-05 user clarification:

- Player is a blacksmith in identity but controls their settlement development, resources, recruitment and military actions, broadly like a Travian player. This supersedes uncertainty about whether blacksmith identity restricts settlement/troop control. It does not grant authority over other players or settle shared kingdom campaign governance. Number of villages per player remains pending.
- Units arrive with their default weapons. Players can upgrade them through smithing. Whether some unit types may change weapon categories remains pending; do not assume unrestricted swapping or derive a final roster from exploratory files.
- The main building is the settlement's central pillar. Defeat, capture, destruction and recovery rules are not decided. No automatic main-building destruction/elimination rule adopted.
- Moraphys is the intended final villain faction, analogous in endgame role to Travian's Natars. Tentative goal: obliterate Moraphys and defeat Chaos weapon wielders; recovered Chaos weapons could be imprisoned or used with severe recoil. Tentative kingdom victory involves defeating Moraphys and all rival kingdoms. User explicitly asks to improve this and not set it in stone. Define what defeat/obliteration means, how weapon containment works and how the finale avoids prolonged elimination or runaway dominance. This does not assign all Chaos weapons to Moraphys or resolve canon corruption versus gameplay recoil.
- Combat timing should be balanced; no schedule, truce duration or daily attendance target confirmed.
- Blacksmith progression should use a talent tree with choices purchased using available points. Point sources, reset rules, branches and specialisation limits remain open.
- User will provide unit roster and other alpha content later. Earlier supplied unit/kingdom files contain possible ideas but remain exploratory, not a confirmed roster.

Next design priorities: village count/expansion incentives; weapon-upgrade versus weapon-swap compatibility; defeat/recovery; practical endgame victory; combat windows; talent-point progression. Detailed armour transformations and final unit data can follow later.

## 21. Kingdom military identities, geography and counter design

2026-10-05 explicit user direction: kingdom rosters should be creative and distinct, with every unit having strengths and weaknesses. Aim for viable faction choices rather than one dominant faction/meta; this is a balancing goal, not a guarantee. Combine readable type/counter relationships inspired by Pokémon with numerical stats in the browser-strategy tradition. Horse-mounted troops can have an advantage over appropriate foot infantry but lose to suitable pikemen, also drawing on BFME. Do not interpret cavalry as beating all infantry universally.

User military examples: Veridor has a strong shipwright identity, sea fighting will exist, and a unit will have a water mount (creature and mechanics undefined). Sylvara will have Eagle Raiders; precise mount, abilities and availability require definition. Neither example automatically implies a rune-based transformation or unrestricted flight. Sea combat is intended, but scope for the first alpha remains open.

User geographic clarification for game design:

- Southern Veridor is strongly coastal; do not make all Veridor coastal.
- Lumus is an island.
- Zandres is underground, with Middle-earth dwarf-realm inspiration, but its people are human, not dwarves. Drakanith's previously established human-dragon hybrids remain distinct.
- Nordalh resembles Iceland: green areas, mixed green/snow regions, fully snowy areas and icy zones. It is not uniformly frozen.
- Arkazia is comparatively conventional, with some mountains; avoid making every district mountainous.
- Sylvara is jungle; Draxys is desert; Drakanith is volcanic.
- Dark Reach is a passage between Sylvara's jungle and Arkazia's border granting access to Moraphys's dark peninsula. Exact routes/geometry and whether other access exists remain undefined; do not invent a sole choke point.

Balance recommendations, not adopted numerical rules: separate movement domain (land/sea/air), combat role (infantry/pike/cavalry/ranged/siege etc.), armour class (light/medium/heavy), and rune identity. Use a small number of moderate counter modifiers alongside attack, defence, speed and range; avoid a giant multiplier chart or armour alone deciding outcomes. Counter effects should depend on engagement conditions, such as prepared pikes facing cavalry or vulnerable ranged troops caught by cavalry.

All playable kingdoms need practical responses to enemy domains, without identical rosters: fleets/coastal defences for naval threats; ranged or other defined counters for aerial threats; contestable entrances/routes for underground movement. Strong specialties should carry costs or limitations such as carrying capacity, supply, deployment or access. Ships should matter at sea and ports, flying units need defined counters and capture rules, and underground routes must not enable unanswerable conquest. Design these rules before declaring specific unit statistics or immunity.

Map balance should compare effective route times, resources, objectives and usable frontier access. Test mixed armies and siege participation across home terrain, neutral terrain and enemy terrain. Define alternative access or contestable connections where island/underground geography risks isolating a kingdom. Exact faction mechanics, travel penalties and counters remain open for roster design.

## 22. Evolving release scope: medieval first

2026-10-05 explicit clarification: the many ideas describe an evolving game, not features to deliver all at once. The first playable version should be simple and fully medieval. Runes begin to emerge later as the user evolves the game. This supersedes any assumption that L1/L2/L3, Chaos/Order finale, magical armour transformations, flying mounts or every unique kingdom system must ship in version one.

Keep ordinary forging, default medieval unit equipment and upgrades as the initial identity. First-alpha playable faction direction remains Arkazia, Veridor and Sylvara. Exact first-version buildings, units, combat and conquest scope still require definition; NPC outpost timing and naval features are not automatically first-release commitments. Preserve advanced concepts as future design material, not mandatory launch requirements.

Distinguish release progression (adding game systems over time) from seasonal progression (unlocking existing systems within a season). User has not yet decided whether rune emergence occurs during the first season, in a later season, or with a separate update. Do not convert this statement into an automatic first-season fantasy schedule or fixed release dates. The previously suggested 10-12-week Chaos/Order season arc is a future proposal, not the first-version specification.
