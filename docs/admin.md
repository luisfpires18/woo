# Administration

Updated: 2026-10-05. Status: design documentation; no game implementation yet.

## Confirmed

Login plus owner admin account and a dedicated workspace. Quickly modify supported content and balance without code. Mandatory configurable asset slots and emoji/symbol fallbacks for every visual entity. Replace/delete unused storage objects as documented.

## Proposed workspace areas

Worlds/seasons; faction availability; units/buildings/resources/recipes; talents and supported effects; map districts/routes; NPC placement; assets/hotspots; player support; reports; configuration history.

Draft/preview/publish/rollback and audit are recommended. New mechanics still require code. Deleting old image binaries means configuration rollback cannot promise restoration of retired artwork.

## Access requirements

Owner role assigned securely by server provisioning, never public signup. Server checks authorisation on every admin action. Account setup, bootstrap mechanism and authentication stack details remain open.

## Open configuration policy

Which fields editable during a season; validation of costs/prerequisites; immediate versus new-order versus next-season changes; treatment of active queues; destructive content deletion; audit retention; emergency support actions.

## Acceptance scenarios to define

Normal player denied admin access; owner edits supported content without deployment; invalid changes rejected; published revisions auditable; running-order effects predictable; replaced image visible and unused old object removed.
