# Documentation index

Updated: 2026-10-05. Status: design documentation; no game implementation yet.

GitHub is the primary working specification. Read current topic documents before the historical snapshot. Status labels: Confirmed = explicit user decision; Accepted direction = approved concept with details open; Proposal = unapproved recommendation; Open = unanswered; Superseded = replaced.

## Start here

- [Vision](vision.md)
- [Release scope](scope.md)
- [Decision register](decisions/README.md)
- [Open questions](open-questions.md)
- [Development roadmap](roadmap.md)

## Gameplay

- [Player and villages](gameplay/player-and-villages.md)
- [Economy](gameplay/economy.md)
- [Forging and progression](gameplay/forging-and-progression.md)
- [Combat and conquest](gameplay/combat-and-conquest.md)
- [Map and activities](gameplay/map-and-activities.md)
- [Seasons and victory](gameplay/seasons-and-victory.md)

## World and content

- [Kingdoms and geography](world/kingdoms-and-geography.md)
- [Lore constraints](world/lore-constraints.md)
- [Content catalogues](content/catalogues.md)

## Design and operations

- [Screens and themes](design/screens-and-themes.md)
- [Artwork pipeline](design/artwork-pipeline.md)
- [Village visual prototype: step-by-step entry point](design/village-visual-prototype.md)
- [Village asset specification](design/village-asset-specification.md)
- [Village art production runbook](design/village-art-production.md)
- [Village renderer and interaction](technical/village-scene.md)
- [Village scene: researched tools and production plan](technical/village-resource-scene-research.md)
- [Admin village scene editor](technical/village-scene-editor.md)
- [Map visual prototype: step-by-step guide](design/map-visual-prototype.md)
- [Map art production](design/map-art-production.md)
- [Map scene and interaction](technical/map-scene.md)
- [World-map square editing, terrain and validation research](technical/world-map-research.md)
- [Combat visual and replay proof](design/combat-visual-prototype.md)
- [Combat sprite and battlefield art production](design/combat-art-production.md)
- [Combat simulation, replay and validation research](technical/combat-research.md)
- [Administration](admin.md)
- [Architecture](technical/architecture.md)
- [Backend technology comparison](technical/backend-comparison.md)
- [Azure versus Cloudflare hosting](technical/hosting-comparison.md)
- [Data model planning](technical/data-model.md)
- [Asset lifecycle](technical/asset-lifecycle.md)

## Reference

- [Research](reference/research.md)
- [Historical brainstorming snapshot](reference/brainstorming-snapshot.md)

## Update workflow

Read latest repository documents, update the relevant topic, record decision changes and unresolved blockers, validate relative links, and commit a focused description. Do not maintain independent competing design copies. The earlier standalone roadmap is now a navigation pointer; its historical content is preserved here.

## Implementation orchestration

- [Task roadmap](../tasks_roadmap.md)
- [Claude instructions](../CLAUDE.md)
- [Development workflow](workflow/development-workflow.md)
- [Implementation prompt protocol](workflow/implementation-protocol.md)
- [Skills/plugin inventory](workflow/skills-and-plugins.md)
- [Azure and deployment](technical/azure-and-deployment.md)

Current working specification is on dev. Explicitly read dev rather than assuming the repository default. Individual task specifications are linked from tasks_roadmap.md.

- [App versions and image cache](technical/versioning-and-cache.md)

The numbered queue now contains smaller reviewable tasks. Read task dependencies and current milestone outline instead of treating the previous twelve groups as the execution plan.

- [Cloudflare/R2 setup and commands runbook](technical/cloudflare-r2.md)

## Approved visual references

- [Landing-page light/dark mockups and implementation corrections](../mockup/00-landing-page.md)

Read each mockup's matching notes before image-to-code. Owner approval is required before adding images to mockup/.
