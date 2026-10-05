# Data model planning

Updated: 2026-10-05. Status: design documentation; no game implementation yet.

## Proposal only

Candidate concepts: Account/Role; World/Season; Kingdom; PlayerProfile/Blacksmith/TalentSelection; Village/Building/ResourceState; UnitDefinition/TroopContingent; Equipment/Recipe/Upgrade; District/Route; Order/ScheduledEvent; Campaign/Battle/Report; AssetReference; ConfigurationRevision/AuditEntry.

An individual Soldier/BondedWeapon model may be needed later, but equipment scale must be decided before choosing per-soldier persistence. A unit-definition upgrade and an individually bonded Artifact are different concepts.

## Invariants to specify

Server validates ownership, resource sufficiency and prerequisites. Spending and order creation must be atomic. Completion must apply once. Capture must transfer a district and village coherently. Shared assets must not be deleted while referenced. Configuration activation must state its effect on running orders.

## Open

Troop granularity, village count, world/account ownership, equipment inheritance, battle resolution, resource precision and maximums, durable scheduling, FK/deletion rules and migration strategy.

Final schemas and endpoint payloads are intentionally not invented at this stage.

## Confirmed storage direction

SQLite for local and dev. Production PostgreSQL is a later consideration; no simultaneous provider implementation required now. Use account, World, Player membership and Kingdom as distinct concepts. A player's membership is scoped to a World; multiple-world and kingdom-switch policies must be explicitly chosen.

Introduce persistence after UI prototypes but before actual account/world state. Incremental schema/migrations follow authorised slices; do not require the full combat/forging model before first deployment. Choose ORM support, ownership keys, constraints, transactions, precision and concurrency behavior against SQLite itself.

Dev database survives redeploy/restart, resides outside app artifacts and has a verified consistent backup/restore path. Azure filesystem/journal compatibility is an early persistence gate. No blind WAL enablement or multi-instance SQLite scaling assumption. See [Azure deployment](azure-and-deployment.md).
