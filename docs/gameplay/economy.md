# Economy

Updated: 2026-10-05. Status: design documentation; no game implementation yet.

## Confirmed context

The player controls resources and settlement development. Ordinary forging and equipment upgrades must be useful before runes.

## Open

Resource types; production/storage rates; costs; upkeep; trading/convoys; loot; crafting materials; offline production; queue limits; cancellation/refunds; growth curve and first-version building list. Food, lumber, stone and ore in concept screens or earlier fixtures are illustrative, not a confirmed catalogue. The owner explicitly wants an original resource model; no Travian field count, resource ring or copied economy is adopted. Resource-screen and production-site artwork are deferred until this direction is defined.

## Validation approach (proposal)

Simulate starting resources, storage overflow, time-to-upgrade, troop upkeep and forging affordability over representative daily visits. Check every upgrade can fit inside attainable storage. Use admin-managed, world-specific balance values with clearly defined change timing.

## Draft defects already found

resource starting rate inconsistent (3/s versus formula 7/s); starting 1,600 storage fills from 500 in about 2.6 minutes at 7/s; 9,200 maximum storage cannot afford approximately 199,500 lumber for a level-20 field. Existing numbers are not adopted.

Do not reuse those draft values. A battle/forging economy spreadsheet or simulation should precede final tuning.
