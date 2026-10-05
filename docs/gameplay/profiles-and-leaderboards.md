# Player profiles, settlement profiles and leaderboards
Updated: 2026-10-05. Status: owner-requested features; details remain proposals.

## Confirmed
Players should be able to inspect player profiles and settlements. Settlement is the preferred player-facing term instead of village. Include browser-strategy-style leaderboards. This does not adopt Travian's scoring formulas, population model, alliances or resource catalogue.

The same entity is meant by older village references and settlement. Preserve stable task keys and filenames for now; agree code/entity names before implementation.

## Player and settlement profiles: PROFILES (057)
Proposed player summary: nickname, image/symbol fallback, kingdom, approved public achievements/statistics and settlement links scoped to the current world. Account identity and per-world membership are distinct; do not mix several seasons' ownership.

Proposed settlement summary: name, current owner/kingdom, authorised location and approved development summary. Link from map markers to settlement and from settlement to owner. Own-settlement management remains an existing game screen, not an unrestricted action on another player's profile.

Define self/friendly/opponent/non-member visibility first. Public fields are server-projected. Never expose resources, troop/equipment inventory, queues or concealed locations solely because a user visits a profile URL. Unknown data stays unknown. Search/direct IDs obey the same policy as map discovery.

After capture, show current ownership consistently and keep historical outcomes separate. Decide old-link/deleted-settlement behaviour. No messaging, alliance groups or additional biography editor is implied.

## Rankings: LEADERBOARDS (058)
Provide world/season player and kingdom rankings. Proposed navigation: category tabs, rank/name/kingdom/score, search/filter, pagination, current-player position and clear last-update time. Link players to profiles. No alliance ranking because independent alliances are excluded.

Candidate categories to review, not adopted metrics:
| Category | Question before adoption |
|---|---|
| Settlement development | Which approved development measure exists? No invented population system |
| Attack contribution | What action earns credit, and how is repeated opponent farming handled? |
| Defence contribution | How do multiple defenders receive meaningful, fair credit? |
| Kingdom territory | Current held districts, area or another approved measure? Avoid repeated-capture inflation |

Prioritise categories supporting warfare/conquest and kingdom defence. A smith/progression category can be considered later if useful. Do not equate leaderboard rank with season victory unless victory rules explicitly say so. No mandatory leaderboard rewards or all-time global score.

Use authoritative recorded facts. Score rules have a version and defined change timing; historical completed seasons need stable interpretation. Ties, participant eligibility, NPC exclusion, resets, retention and snapshot cadence are dispatch decisions.

Do not reveal hidden army strength indirectly through a military-power ranking. Count a battle/capture once despite retries. Distinguish current territory from historical capture contribution. Shared battles require explicit attribution; leaderboards should not reward merely claiming a killing blow by accident.

## Technical boundaries and validation
Use authorised world-scoped queries and bounded pagination suitable for SQLite. Start with straightforward aggregate/read models; no separate analytics infrastructure is justified for friends alpha. Decide freshness/caching from measured queries. Labels expose freshness; private data never reaches the client.

Check ownership transfer, cross-world isolation, unknown/deleted targets, restricted search, tied ranks, empty/new seasons, retry idempotency, finished-season snapshots, accessibility and narrow screens. Public image slots/fallbacks follow ASSETS. Ranking formulas wait for approved gameplay and economy definitions.

## Tasks and next decisions
057 owns profile and settlement inspection. 058 owns rankings, and later BALANCE/QA checks them. They remain TODO. Resolve profile visibility and ranking metrics before their implementation prompts; foundation work need not wait for these decisions.

## TLDR
Next: choose visible fields and useful score definitions before dispatch.
Done: terminology and requested profile/ranking scope recorded.
Issues: formulas, attribution and privacy policies remain open.
