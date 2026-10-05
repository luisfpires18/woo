# Screens and themes

Updated: 2026-10-05. Status: design documentation; no game implementation yet.

## Screens

Resources, Village, Kingdom Map, Combat, Forge and Kingdom Duties have concept previews. Login and kingdom selection now have approved conversation previews; admin still needs designs. The owner will upload their images to Git. Preview numbers, ore choices and map layouts are illustrative.

React handles menus, forms, counters and accessible controls; PixiJS is proposed for interactive scenes. Village/map/combat interaction should remain usable with fallback symbols before art exists.

## Confirmed theme direction

Independent dark/light mode. Arkazia crimson/black; Sylvara green/gold; Veridor blue/silver. Exact shades and other secondary colours open. Kingdom accent colours remain separate from readable text/surfaces and semantic feedback.

## Proposed interaction requirements

No UI text/counters baked into art. Selected/hover/focus/disabled states, touch targets, companion accessible lists and readable contrast. Ownership uses labels/icons alongside colour. Same art across themes; do not globally darken it to imitate dark UI.

## Early UI delivery order

Build shared theme/layout foundations, polish landing, establish game navigation and village/map/combat/forge prototypes, then admin/onboarding prototypes. Use explicit mock data until database/account APIs exist. All shared layouts include the app-version footer. Actual login and World/Kingdom association come after SQLite, not as a pretend session in the UI prototype.

## Approved landing page, 2026-10-05

See [mockup/00-landing-page.md](../../mockup/00-landing-page.md) and its light/dark screenshots. Slim navbar with anvil branding and auth control only; supplied title/banner below navbar; general info, metrics and world cards underneath. Login uses a separate page/layout. Logged-in navbar uses avatar/nickname with Profile, Settings and Log out. Settings opens a separate page for appearance; no appearance controls in the dropdown and no navbar theme toggle. The written correction overrides the light screenshot's expanded appearance section. Theme and authentication state are independent. Metrics/world names are illustrative, not actual content.

## Latest visual direction, 2026-10-05
The owner prefers the early game-screen style: dark charcoal navigation, warm ivory content, crimson actions and artwork integrated into the village/inspector; minimal textures. Dark/light support remains required. The ivory reference is not the final dark-theme design.
Centered login over full-screen forge artwork is approved. Kingdom selector is approved: Arkazia/Veridor/Sylvara with image placeholders for later admin uploads; other future playable kingdoms disabled; NPC factions not selectable.
The latest remade village screen is awaiting approval. Its costs, timings, labels, castle, river and building roster are illustrative. Do not treat visual generation as gameplay approval or production art.
Read the [village visual proof](village-visual-prototype.md). UI labels/counters must remain live; scene art is not a flattened application.
The owner now handles all Git image uploads; the older committed landing references remain historical until owner replacement. Match newer confirmed written directions when these differ from older pixels.

## Combat preparation
V1 is a static battalion board plus final result, as clarified by the owner. Retain the preferred illustrated central layout and readable army rosters, with static battalion icons/cards and optional battlefield backdrop. No animation, replay, timeline or narrated log is required. Follow [combat board proof](combat-visual-prototype.md) and [static artwork production](combat-art-production.md). Exact formation slots, stance effects, counters and losses remain pending. React/CSS is the first choice for this simpler screen; PixiJS is not mandatory.
