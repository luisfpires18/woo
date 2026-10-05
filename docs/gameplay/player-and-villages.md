# Player and settlements

Updated: 2026-10-05. Status: design documentation; no game implementation yet.

## Confirmed

The player is thematically a master blacksmith and controls their own settlement development, resources, recruitment and military actions. The desired experience is warfare, conquest and building powerful equipped armies; blacksmithing supplies equipment/progression tools rather than restricting play to smith activities. It does not grant control of other players. Main building is the settlement's central pillar. Settlement capture transfers its territory to the conquering kingdom.

Light/medium/heavy equipment equivalents remain cloth/leather/plate.

## Open

Settlement count and expansion; initial settlement placement; main-building prerequisites; building slots and queues; upgrade/demolition cancellation; defender losses and refuge; ownership of captured infrastructure; capital protection and elimination.

## Proposal

Start with one settlement and a small expansion cap for the friends alpha, then measure workload. Preserve meaningful smith progression after defeat. Neither cap nor recovery is confirmed.

## Specification needed

Define construction request validation, resource consumption, completion times, offline completion, capture during construction, and admin changes to ongoing orders. Acceptance scenarios should include insufficient resources, duplicate requests, completed offline queues and ownership transfer.

## Terminology and inspection
The owner prefers settlement/settlements in the game UI. Historical village references and VILLAGE task keys refer to the same concept; do not create two entity types. Player and settlement profiles and world/season leaderboards are requested. Read [profiles and leaderboards](profiles-and-leaderboards.md). Settlement capacity and exact public metrics remain open.
