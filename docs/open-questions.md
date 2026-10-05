# Open questions

Updated: 2026-10-05. Status: design documentation; no game implementation yet.

Use this list to resolve implementation blockers before writing final formulas/contracts. An Open item is not an adopted rule.

| ID | Question | Blocks |
|---|---|---|
| O-01 | How many villages; how is expansion earned? | Ownership, map capacity, economy |
| O-02 | Are upgrades individual, batch, regiment or unit-type; do recruits inherit them? | Forging, troop data model |
| O-03 | Which units can swap weapon categories? | Equipment compatibility |
| O-04 | Capture/defeat, main building, refuge, infrastructure and equipment survival? | Conquest and recovery |
| O-05 | Automatic battle rules, timing, offline defence and truce? | Combat engine and scheduling |
| O-06 | Who declares kingdom campaigns and assigns captured villages? | Cooperation and permissions |
| O-07 | Which starting resources, buildings, units and provisional numbers? | First content catalogue |
| O-08 | Talent point earning, branches, resets and dual specialisation? | Medieval progression |
| O-09 | What ends the medieval first season and what resets? | Alpha scope and retention |
| O-10 | Map topology, capacities, travel and visibility? | World generation and orders |
| O-11 | Exact first-version duties/camps, NPCs, hero and naval scope? | Scope freeze |
| O-12 | Select recommended stack, versions, hosting and account bootstrap? | Technical setup |
| O-13 | Which admin edits affect active orders; content rollback semantics? | Config/event consistency |
| O-14 | R2 confirmation, file limits, delivery and cleanup design? | Asset storage |
| O-15 | Rune emergence within a season or via later releases? | Future roadmap |
| O-16 | Rune availability/scarcity, bearer bonds, death and recoil/armour interaction? | Future advanced forging |
| O-17 | Resolve recorded canon conflicts, especially soul/blood and Order L3? | Future lore-dependent mechanics |
| O-18 | Moraphys victory, Chaos imprisonment/use and rival defeat meaning? | Future finale |
| O-19 | Sea/air/tunnel access and counterplay for every kingdom? | Future domains |

Recommended next order: O-01/O-02, O-04/O-05/O-06, O-07/O-08/O-09, then technical contracts. The owner will provide rosters later; original drafts do not resolve O-07.

Future questions need not block a purely medieval prototype unless its data model would make the future feature impractical.
