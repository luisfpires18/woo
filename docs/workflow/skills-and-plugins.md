# Claude skills and plugins inventory

Updated: 2026-10-13. Audited on the owner's Windows machine by task 001 SKILLS (see [audit results](#audit-results-2026-10-13)). Only the PixiJS collection was installed; everything else was recorded as found.

## Health check

The first implementation task audits Claude Code version, installed skills/plugins, executable helpers and browser integrations on the actual machine. For each requested item, record exact name, source URL/repository, version/ref, install location, availability, invocation and a minimal smoke check. Distinguish a skill from a plugin, MCP server and CLI.

Install missing items only from verified owner-approved sources. Names alone are insufficient for an installer. Retain working installations; avoid duplicate/conflicting versions. Credentials and user-specific config remain local. Installation alone is not use.

## Requested inventory

| Category | Requested names | State after 2026-10-13 audit |
|---|---|---|
| Backend | aspnet-core-guidance; graphify | aspnet-core-guidance UNRESOLVED; graphify CLI verified, Claude skill missing |
| Frontend | web-design-guidelines; impeccable; taste-skill; image-to-code-skill; emilkowalski; huashu-design-guidelines; ui-ux-pro-max; awesome-design-md | Five installed (discovery verified); emilkowalski, huashu-design-guidelines, ui-ux-pro-max UNRESOLVED |
| PixiJS / 2D rendering | Entire official pixijs/pixijs-skills collection; see section below | VERIFIED (26 skills) |
| Phase control | phase-workflow | UNRESOLVED |
| Wording | humanizer | Installed (discovery verified) |
| Browser | Claude Browser MCP; playwright-cli | Browser MCP VERIFIED; playwright-cli skill installed, executable missing |
| Workflow helpers | Caveman; Ponytail; RTK; Understand Anything | Caveman, Ponytail VERIFIED; RTK CLI verified without Claude hook; Understand Anything UNRESOLVED |

Graphify appears once despite being requested in multiple categories. Exact sources/aliases and platform support remain to verify. Others may be added later.

Use backend guidance for .NET, focused design/accessibility skills for UI, image-to-code only when relevant art exists, browser tooling for actual UI checks and graph tools when they provide useful repository understanding. Where skill advice conflicts, owner task constraints and repository instructions take precedence. Do not invoke every frontend skill ceremonially.

LP WORK belongs to ChatGPT planning/review. Do not tell Claude to install or execute the unavailable ChatGPT LP AI WORK skill. The repository workflow is expressed through CLAUDE.md and these protocol documents.

## PixiJS / 2D game rendering

Owner-approved addition, 2026-10-05: install the **entire official PixiJS skills collection**, not only the eight initially highlighted skills. Installed and verified on the owner's machine on 2026-10-13; see [PixiJS install record](#pixijs-install-record).

Source: [pixijs/pixijs-skills](https://github.com/pixijs/pixijs-skills).
Official overview: [PixiJS Skills](https://pixijs.com/llms).
The collection currently targets PixiJS v8. Verify compatibility with the version selected in STACK before implementation.

| Skill | Coverage |
|---|---|
| pixijs | Entry point; routes to relevant specialist skills |
| pixijs-accessibility | Keyboard and screen-reader accessibility |
| pixijs-application | Application setup, resize and lifecycle |
| pixijs-assets | Assets, bundles, manifests and spritesheets |
| pixijs-blend-modes | Layer compositing and blend modes |
| pixijs-color | Colour handling and conversions |
| pixijs-core-concepts | Renderer architecture and rendering lifecycle |
| pixijs-create | Scaffolding or adding PixiJS to a project |
| pixijs-custom-rendering | Custom shaders, filters and rendering |
| pixijs-environments | Workers, offscreen and nonstandard environments |
| pixijs-events | Pointer/touch input, hit areas and dragging |
| pixijs-filters | Built-in and custom visual effects |
| pixijs-html-source | Experimental HTML-in-Canvas DOM textures; found upstream 2026-10-13, absent from the original list |
| pixijs-math | Coordinates, transforms, shapes and hit testing |
| pixijs-migration-v8 | Migration from v7 to v8 |
| pixijs-performance | Profiling, culling, batching and cleanup |
| pixijs-scene-container | Grouping, transforms and ordering |
| pixijs-scene-core-concepts | Scene graph, layers, masking and render groups |
| pixijs-scene-dom-container | HTML elements aligned with canvas content |
| pixijs-scene-gif | Animated GIF rendering and lifecycle |
| pixijs-scene-graphics | Polygons, paths, fills and borders |
| pixijs-scene-mesh | Custom geometry and meshes |
| pixijs-scene-particle-container | Batched particles |
| pixijs-scene-sprite | Sprites and sprite animations |
| pixijs-scene-text | Canvas text and bitmap text |
| pixijs-ticker | Frame updates and render-loop timing |

### Setup and verification

The official universal installer command, run on the owner's development machine, is:

```bash
npx skills add https://github.com/pixijs/pixijs-skills
```

Select Claude Code and the full collection in the installer. The official alternative is adding its Claude Code marketplace via `/plugin marketplace add pixijs/pixijs-skills`; adding a marketplace is not itself evidence that its plugin is installed/enabled. Choose one installation route and avoid duplicate installations.

During SKILLS, inspect the current source and installed collection, record revision/install scope, confirm discovery of every listed skill and run a small relevant rendering smoke check. Report missing or renamed skills explicitly; do not mark availability verified from this inventory alone. SKILLS used the `npx skills` route; the marketplace route was not added, so do not add it as well.

### Use in WOO

Use the entry skill to load appropriate guidance for village, map and combat tasks. All skills should be available, but only task-relevant ones need loading. Initial priorities are Application, Assets, Container, Sprite, Graphics, Events, Math, Performance and Accessibility. Advanced shaders, particles, GIFs and migration guidance do not become game requirements merely because they are installed.

Read the [village prototype](../design/village-visual-prototype.md) and [map prototype](../design/map-visual-prototype.md) for WOO-specific contracts. Skills help implement rendering; they do not establish lore, art quality, combat mechanics or balance. React remains responsible for the normal interface, and ASP.NET Core remains authoritative for game state.

## Audit results (2026-10-13)

Task 001 SKILLS, branch feat/001-skills-preflight. Machine: owner's Windows 11 Pro 10.0.26200. Claude Code 2.1.296 (npm global `@anthropic-ai/claude-code`, used from the desktop app Code tab). Node v22.23.2, npm 10.9.8, Python 3.13.14, git and gh present. `~` means the owner's user profile; no credentials or personal settings are recorded.

Result labels: **VERIFIED** = exercised in a real check; **INSTALLED** = present and listed in a Claude session's skill list, invocation not exercised; **PARTIAL** = one part works, another is missing; **UNRESOLVED** = not found locally and no owner-approved source, so nothing was installed.

| Item | Kind | Source | Version/ref | Location/scope | Invocation | Smoke check | Result |
|---|---|---|---|---|---|---|---|
| aspnet-core-guidance | ? | none known | - | not found | - | local search found nothing | UNRESOLVED |
| graphify | CLI (+ optional skill) | PyPI `graphifyy`, homepage github.com/Graphify-Labs/graphify | 0.9.56 | `~/.local/graphify-venv`, on PATH | `graphify <command>` | `graphify --version` → 0.9.56 | PARTIAL: CLI verified; Claude skill not installed (`graphify install --platform claude` not run, source not owner-approved) |
| web-design-guidelines | skill | github.com/vercel-labs/agent-skills `skills/web-design-guidelines` | metadata 1.0.0; folder hash 3116f3e6 | `~/.claude/skills`, user | `/web-design-guidelines` | listed in session skill list | INSTALLED |
| impeccable | plugin | `impeccable@local-desktop-app-uploads` (desktop upload; homepage impeccable.style; upstream repo not recorded) | 4.3.1 | user, enabled | `/impeccable:impeccable` | listed in session skill list | INSTALLED |
| taste-skill | plugin (15 skills) | `taste-skill@local-desktop-app-uploads`; manifest repo github.com/leonxlnx/taste-skill | 1.0.0 | user, enabled | `taste-skill:<skill>` | listed in session skill list | INSTALLED |
| image-to-code-skill | skill inside taste-skill | as taste-skill | 1.0.0 | user | `taste-skill:image-to-code-skill` | listed in session skill list | INSTALLED |
| emilkowalski | ? | none known | - | not found (only a taste-skill sponsor image matches "emil") | - | local search | UNRESOLVED |
| huashu-design-guidelines | ? | none known | - | not found | - | local search | UNRESOLVED |
| ui-ux-pro-max | ? | none known | - | not found | - | local search | UNRESOLVED |
| awesome-design-md | skill | github.com/wearescientist/awesome-design-md | folder hash 6d732af3 | `~/.claude/skills`, user | `/awesome-design-md` | listed in session skill list | INSTALLED |
| PixiJS collection | 26 skills | github.com/pixijs/pixijs-skills | main @ 83760c6f53462ca9cecd68055041f5a8c94758ce | `~/.claude/skills/pixijs*`, user | `/pixijs` (router) | see below | VERIFIED |
| phase-workflow | ? | none known | - | not found | - | local search | UNRESOLVED |
| humanizer | plugin | `humanizer@humanizer`, github.com/blader/humanizer | 2.8.2 (1b48564898e9); marketplace clone at 3.0.0 not applied to user scope | user, enabled | `/humanizer:humanizer` | listed in session skill list | INSTALLED |
| Claude Browser MCP | built-in desktop app browser pane (`mcp__Claude_Browser__*`) | Claude desktop app | ships with app | app built-in | browser tools | opened localhost fixture, read page text, captured screenshot (PixiJS check) | VERIFIED |
| playwright-cli | skill + CLI | github.com/microsoft/playwright-cli `skills/playwright-cli` | folder hash 6b746b43 | skill in `~/.claude/skills`, user | `/playwright-cli` | skill listed; `playwright-cli` not on PATH, not in global npm | PARTIAL: skill installed, executable missing (not installed; source not owner-approved) |
| Caveman | plugin + hooks | `caveman@caveman`, github.com/JuliusBrussee/caveman | 655b7d9c5431 | user, enabled; also standalone `~/.claude/hooks/caveman-*.js` and global npm `@juliusbrussee/caveman-code` 0.65.2 | `/caveman` | SessionStart/UserPromptSubmit context injected in audit session | VERIFIED; activation text appeared twice at SessionStart (likely plugin + standalone hook overlap; left unchanged) |
| Ponytail | plugin | `ponytail@ponytail`, github.com/DietrichGebert/ponytail | 4.10.0 (e3ba2aa6f1e6) | user, enabled | `/ponytail` | SessionStart context injected in audit session | VERIFIED |
| RTK | CLI | not recorded locally | 0.48.0 | `~/.local/bin/rtk.exe` | `rtk <command>` | `rtk git status` ran; reported no Claude hook installed | PARTIAL: CLI verified; Claude hook not installed (`rtk init -g` changes global config; source unverified) |
| Understand Anything | ? | none known | - | not found | - | local search | UNRESOLVED |

### PixiJS install record

- Route: universal installer, matching how the existing user skills were installed (`~/.agents/.skill-lock.json`). Command, run from a temporary directory outside the repository with `DISABLE_TELEMETRY=1 DO_NOT_TRACK=1`: `npx -y skills@1.7.2 add https://github.com/pixijs/pixijs-skills -g -a claude-code --skill '*' --copy -y`.
- Source revision: GitHub API reported main at 83760c6f53462ca9cecd68055041f5a8c94758ce before and after installation. The skills CLI records per-skill folder hashes (entry skill: 8c1854dd), not a commit SHA.
- Installed 2026-10-13T21:22:30Z; 26 skills: the 25 originally requested plus `pixijs-html-source`. Each folder name matches its SKILL.md `name`. No missing or renamed skills.
- Discovery: a fresh non-interactive Claude session (`claude -p`) listed all 26 (`COUNT=26`); the audit session also listed all 26 and invoked `pixijs`, which loaded its router.
- Render check: temporary HTML outside the repository, served on 127.0.0.1, loaded pixi.js 8.21.0 from jsDelivr, ran `Application.init` and drew a red `Graphics` rect. Output: `OK version=8.21.0 renderer=webgl pixel=255,0,0,255`. Fixture, server and browser tab removed afterwards.
- Compatibility: collection targets PixiJS v8; npm latest observed was 8.22.0. STACK still chooses the exact version.

### Other integrations found (not requested)

- Plugins: frontend-design and superdesign (claude-plugins-official), user scope. Some plugins also have project-scope entries for other local repositories; none are scoped to WOO.
- MCP: claude.ai connectors (Claude Docs, Figma, Gmail, Google Drive, Google Calendar connected; Notion needs authentication); user-scope `unity-mcp` (connected, tool fetch timed out); Godot MCP tools exposed in the session. WOO has no project `.mcp.json` or `.claude/` directory.
- Claude in Chrome: tools available, but no connected Chrome browser (empty browser list). Not exercised.
- An `lp-ai-work` folder exists in `~/.agents/skills`. Claude does not load it, and per this document it stays ChatGPT-side.
