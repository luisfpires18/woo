# Claude skills and plugins inventory

Updated: 2026-10-05. Requested inventory; installation status on the owner's computer is UNKNOWN. No installations were performed here.

## Health check

The first implementation task audits Claude Code version, installed skills/plugins, executable helpers and browser integrations on the actual machine. For each requested item, record exact name, source URL/repository, version/ref, install location, availability, invocation and a minimal smoke check. Distinguish a skill from a plugin, MCP server and CLI.

Install missing items only from verified owner-approved sources. Names alone are insufficient for an installer. Retain working installations; avoid duplicate/conflicting versions. Credentials and user-specific config remain local. Installation alone is not use.

## Requested inventory

| Category | Requested names | Initial state |
|---|---|---|
| Backend | aspnet-core-guidance; graphify | UNKNOWN |
| Frontend | web-design-guidelines; impeccable; taste-skill; image-to-code-skill; emilkowalski; huashu-design-guidelines; ui-ux-pro-max; awesome-design-md | UNKNOWN |
| PixiJS / 2D rendering | Entire official pixijs/pixijs-skills collection; see section below | UNKNOWN |
| Phase control | phase-workflow | UNKNOWN |
| Wording | humanizer | UNKNOWN |
| Browser | Claude Browser MCP; playwright-cli | UNKNOWN |
| Workflow helpers | Caveman; Ponytail; RTK; Understand Anything | UNKNOWN |

Graphify appears once despite being requested in multiple categories. Exact sources/aliases and platform support remain to verify. Others may be added later.

Use backend guidance for .NET, focused design/accessibility skills for UI, image-to-code only when relevant art exists, browser tooling for actual UI checks and graph tools when they provide useful repository understanding. Where skill advice conflicts, owner task constraints and repository instructions take precedence. Do not invoke every frontend skill ceremonially.

LP WORK belongs to ChatGPT planning/review. Do not tell Claude to install or execute the unavailable ChatGPT LP AI WORK skill. The repository workflow is expressed through CLAUDE.md and these protocol documents.

## PixiJS / 2D game rendering

Owner-approved addition, 2026-10-05: install the **entire official PixiJS skills collection**, not only the eight initially highlighted skills. Installation and verification on the owner's machine remain UNKNOWN.

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

During SKILLS, inspect the current source and installed collection, record revision/install scope, confirm discovery of every listed skill and run a small relevant rendering smoke check. Report missing or renamed skills explicitly; do not mark availability verified from this inventory alone. No installation was performed by this documentation update.

### Use in WOO

Use the entry skill to load appropriate guidance for village, map and combat tasks. All skills should be available, but only task-relevant ones need loading. Initial priorities are Application, Assets, Container, Sprite, Graphics, Events, Math, Performance and Accessibility. Advanced shaders, particles, GIFs and migration guidance do not become game requirements merely because they are installed.

Read the [village prototype](../design/village-visual-prototype.md) and [map prototype](../design/map-visual-prototype.md) for WOO-specific contracts. Skills help implement rendering; they do not establish lore, art quality, combat mechanics or balance. React remains responsible for the normal interface, and ASP.NET Core remains authoritative for game state.
