# Research and inspirations

Updated: 2026-10-05. Status: design documentation; no game implementation yet.

Qualitative Travian research suggests recurring concerns around spending advantage, attendance pressure, repetitive administration, onboarding, cheating/enforcement trust, mobile parity and player density. These are not representative survey results.

Recommendations: no purchased combat power; useful queues and shared truce; complete mobile browser flow; understandable battle reports; clear enforcement/appeals; recovery after defeat; support contributions. Daily play target proposed at 20-40 minutes across 2-3 visits, not a guarantee. Test with competitive veterans, lapsed players and newcomers.

Draft arithmetic defects: resource starting rate inconsistent (3/s versus formula 7/s); starting 1,600 storage fills from 500 in about 2.6 minutes at 7/s; 9,200 maximum storage cannot afford approximately 199,500 lumber for a level-20 field. Existing numbers are not adopted.

## BFME adaptation assessment

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

## Research links

- BFME II manual reproduction: https://manuals.plus/m/f958931bd298d60eacdaad759963ecae735f4f96b71f848404bc8b5468796c88
- BFME II developer Create-a-Hero Q&A: https://www.gamespot.com/articles/the-lord-of-the-rings-the-battle-for-middle-earth-ii-qanda-fan-questions-part-i/1100-6133390/
- BFME II contemporary guide (secondary): https://www.gamespot.com/articles/the-lord-of-the-rings-the-battle-for-middle-earth-ii-walkthrough/1100-6146219/
- Travian village management feedback: https://blog.travian.com/2024/06/travian-loop-village-management/
- Kingdoms community feedback: https://blog.kingdoms.com/feedback-loop-community-topics-3/
- Kingdoms night protection: https://support.kingdoms.com/en/articles/3-game-versions-speed-and-specials
- Kingdoms world concentration: https://blog.kingdoms.com/new-game-world-schedule/

## Technical and cost context

Earlier research found historical original Travian PHP/MySQL/HTML/CSS, modern Legends React/TypeScript-related hiring signals and historical Kingdoms AngularJS/Node/PHP/data stores. These do not establish today's full production stack. Actual Travian hosting costs are private. Small-world infrastructure estimates are planning hypotheses, not measured operator costs; price/provider/date verification required before budgeting.

Qualitative feedback is not a representative player survey. Do not present these recommendations as proven demand or copy Travian/BFME mechanics wholesale.
