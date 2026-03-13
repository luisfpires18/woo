---
description: "Use when writing React components, hooks, pages, or feature modules. Covers Zustand vs React Query, component patterns, TypeScript strictness, kingdom theming, and config usage."
applyTo: "client/src/**/*.tsx"
---

# Frontend React Conventions

Full reference: `docs/04-frontend/frontend-guide.md`

## State Management (Two Layers)

- **Zustand** (`stores/`) — global game state: current village, player info, real-time data.
- **React Query** — server data: API responses, caching, refetching.
- Ask: "Is this game state (Zustand) or server-fetched data (React Query)?" Never mix.

## Component Structure

- **Reusable components** → `client/src/components/{Name}/{Name}.tsx` + `{Name}.module.css`
- **Feature components** → `client/src/features/{feature}/components/`
- **Pages** → `client/src/features/{feature}/pages/`
- If a UI element is used in 2+ places, extract it to `components/`.

## TypeScript Strictness

- Strict mode mandatory. No `any` types.
- All API responses typed in `types/api.ts`.
- Union types preferred over enums.
- Avoid `as unknown as` casts — use proper typing.

## Config Usage

- Building/troop/resource values come from `client/src/config/` (imports generated JSON).
- Use `useBuildingDisplayNames()` hook for building names — never hardcode.
- Military buildings (barracks, stable, etc.) → `BuildingTrainingModal`. Non-military → `BuildingDetailModal`.

## Kingdom Theming

- Themes driven by `data-kingdom` attribute on `<html>` — NOT a light/dark toggle.
- Use CSS custom properties (`--bg-primary`, `--accent-primary`, etc.) that change per kingdom.

## Accessibility

- Clickable non-button elements need `role="button"`, `tabIndex={0}`, and Enter/Space key handlers.
- Follow the pattern in `components/Card/Card.tsx`.

## WebSocket

- Reconnection: exponential backoff (1s → 2s → 4s → 8s, max 30s).
- On reconnect, re-sync state via REST fallback.

## Lazy Loading

- Route-level lazy loading for map, forge, and lore modules.
