# Screens and themes

Updated: 2026-10-05. Status: design documentation; no game implementation yet.

## Screens

Resources, Village, Kingdom Map, Combat, Forge and Kingdom Duties have concept previews. Login, kingdom selection and the admin workspace also need designs. Preview numbers, ore choices and map layouts are illustrative.

React handles menus, forms, counters and accessible controls; PixiJS is proposed for interactive scenes. Village/map/combat interaction should remain usable with fallback symbols before art exists.

## Confirmed theme direction

Independent dark/light mode. Arkazia crimson/black; Sylvara green/gold; Veridor blue/silver. Exact shades and other secondary colours open. Kingdom accent colours remain separate from readable text/surfaces and semantic feedback.

## Proposed interaction requirements

No UI text/counters baked into art. Selected/hover/focus/disabled states, touch targets, companion accessible lists and readable contrast. Ownership uses labels/icons alongside colour. Same art across themes; do not globally darken it to imitate dark UI.

## Early UI delivery order

Build shared theme/layout foundations, polish landing, establish game navigation and village/map/combat/forge prototypes, then admin/onboarding prototypes. Use explicit mock data until database/account APIs exist. All shared layouts include the app-version footer. Actual login and World/Kingdom association come after SQLite, not as a pretend session in the UI prototype.
