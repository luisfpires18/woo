# Frontend Guide

> All conventions for the React + TypeScript frontend. Read before writing any client code.

---

## Technology Stack

| Tool | Version | Purpose |
|------|---------|---------|
| React | 19 | UI framework |
| TypeScript | 5+ | Type safety (strict mode enabled) |
| Vite | 5+ | Build tool and dev server |
| Zustand | 4+ | Global game state management |
| TanStack Query (React Query) | 5+ | Server state fetching/caching |
| React Router | 6+ | Client-side routing |
| CSS Modules | built-in | Scoped component styling |
| Vitest | latest | Unit testing framework |
| React Testing Library | latest | Component testing |

---

## Folder Structure

```
client/src/
├── components/             # Reusable UI components (used in 2+ places)
│   ├── Button/
│   │   ├── Button.tsx
│   │   └── Button.module.css
│   ├── Modal/
│   ├── Card/
│   ├── ResourceBar/
│   ├── Tooltip/
│   ├── LoadingSpinner/
│   ├── ErrorBoundary/
│   └── Layout/
│       ├── Header/
│       ├── Sidebar/
│       └── Footer/
├── config/                 # Game config imports
│   ├── generated/          # Auto-generated JSON from Go config (DO NOT EDIT)
│   │   ├── buildings.json
│   │   ├── troops.json
│   │   └── resources.json
│   ├── buildings.ts        # TS wrapper for buildings.json
│   ├── troops.ts           # TS wrapper for troops.json
│   └── resources.ts        # TS wrapper for resources.json
├── features/               # Feature modules (domain-specific)
│   ├── auth/
│   │   ├── components/     # Auth-specific components (LoginForm, RegisterForm)
│   │   ├── hooks/          # Auth-specific hooks (useLogin, useRegister)
│   │   └── pages/          # Auth pages (LoginPage, RegisterPage)
│   ├── village/
│   │   ├── components/     # VillageView, BuildingPanel, ResourcePanel, etc.
│   │   ├── hooks/          # useVillage, useBuilding, useResources
│   │   └── pages/          # VillagePage
│   ├── map/
│   │   ├── components/     # MapRenderer, MapControls
│   │   ├── hooks/          # useMapData
│   │   └── pages/          # MapPage
│   ├── kingdom/            # Kingdom selection
│   ├── admin/              # Admin panel, map editor, template system
│   ├── profile/            # Player profile
│   ├── season/             # Season/world management
│   └── landing/            # Landing page
├── hooks/                  # Shared custom hooks
│   ├── useBuildingDisplayNames.ts  # Kingdom-specific building display names
│   ├── useResourceTicker.ts        # Client-side resource tick interpolation
│   └── ...
├── services/               # API and WebSocket layer
│   ├── api.ts              # REST API client (typed fetch wrapper)
│   ├── auth.ts             # Auth-specific API calls
│   ├── village.ts          # Village API calls
│   ├── training.ts         # Troop training API calls
│   ├── map.ts              # Map API calls
│   ├── admin.ts            # Admin API calls
│   ├── season.ts           # Season API calls
│   ├── template.ts         # Map template API calls
│   └── player.ts           # Player API calls
├── stores/                 # Zustand stores
│   ├── authStore.ts        # Auth state (user, token, login status)
│   ├── gameStore.ts        # Real-time game state (current resources tick)
│   ├── mapStore.ts         # Map viewport and visible tiles
│   ├── assetStore.ts       # Sprite/asset management
│   └── themeStore.ts       # Kingdom theme management
├── styles/                 # Global styles
│   ├── globals.css         # CSS reset, base styles, CSS variables
│   ├── themes.css          # Kingdom theme variables
│   └── typography.css      # Font imports and text styles
├── types/                  # TypeScript interfaces (shared)
│   ├── api.ts              # API request/response types + union types (BuildingType, TroopType, Kingdom)
│   ├── game.ts             # Game entity types (Village, Building, Troop, etc.)
│   ├── websocket.ts        # WebSocket message types
│   └── map.ts              # Map-related types
├── utils/                  # Pure utility functions
│   ├── format.ts           # Number formatting, date formatting, time-ago
│   ├── calculations.ts     # Display-only calculations (resource ETA, countdown timers)
│   └── constants.ts        # Client-side constants (building names, troop names, etc.)
├── App.tsx                 # Root component with React Router
├── main.tsx                # Entry point (renders App)
└── vite-env.d.ts           # Vite type declarations
```

---

## Component Rules

### 1. Reuse is Mandatory

If a UI element appears in **2 or more places**, it **must** be extracted to `client/src/components/`. Feature-specific components live inside their feature folder (`features/village/components/`). Shared components live in the top-level `components/` folder.

### 2. Component Structure

Every component follows this pattern:

```
ComponentName/
├── ComponentName.tsx           # Component logic and JSX
├── ComponentName.module.css    # Desktop styles + mobile overrides (@media ≤ 768px)
└── ComponentName.test.tsx      # Unit tests (optional for simple presentational components)
```

### 3. Component Guidelines

```tsx
// Good: Typed props, destructured, no inline styles
interface ButtonProps {
  label: string;
  variant?: 'primary' | 'secondary' | 'danger';
  disabled?: boolean;
  onClick: () => void;
}

export function Button({ label, variant = 'primary', disabled = false, onClick }: ButtonProps) {
  return (
    <button
      className={`${styles.button} ${styles[variant]}`}
      disabled={disabled}
      onClick={onClick}
    >
      {label}
    </button>
  );
}
```

### Anti-Patterns (DO NOT)

- **No inline styles**: Always use CSS Modules.
- **No `any` types**: All props, state, and API responses must be typed.
- **No direct DOM manipulation**: Use React refs only when absolutely necessary.
- **No business/game logic in components**: Components render; hooks and services calculate.
- **No hardcoded strings**: Use constants for building names, troop names, resource names, etc.

---

## State Management

### When to Use What

| Data Type | Tool | Why |
|----------|------|-----|
| Server data (fetched once, cached) | **React Query** | Automatic caching, refetching, loading/error states |
| Real-time game state (ticking resources) | **Zustand** | Synchronous updates from WebSocket, no re-fetch overhead |
| Auth state (user, token) | **Zustand** | Needs to persist across components, updated on login/logout |
| UI state (modal open, tab selection) | **React useState** | Component-local, no need for global state |
| Form state (inputs, validation) | **React useState** or React Hook Form | Component-local |

### Zustand Store Example

```tsx
// stores/villageStore.ts
import { create } from 'zustand';
import type { Village } from '../types/game';

interface VillageStore {
  activeVillage: Village | null;
  setActiveVillage: (village: Village) => void;
  updateResources: (resources: Partial<Village['resources']>) => void;
}

export const useVillageStore = create<VillageStore>((set) => ({
  activeVillage: null,
  setActiveVillage: (village) => set({ activeVillage: village }),
  updateResources: (resources) =>
    set((state) => ({
      activeVillage: state.activeVillage
        ? { ...state.activeVillage, resources: { ...state.activeVillage.resources, ...resources } }
        : null,
    })),
}));
```

### React Query Example

```tsx
// features/village/hooks/useVillage.ts
import { useQuery } from '@tanstack/react-query';
import { api } from '../../../services/api';
import type { Village } from '../../../types/game';

export function useVillage(villageId: number) {
  return useQuery<Village>({
    queryKey: ['village', villageId],
    queryFn: () => api.get(`/villages/${villageId}`),
    staleTime: 30_000, // 30 seconds (server state doesn't change that fast)
  });
}
```

---

## Routing Structure

```tsx
// App.tsx (simplified)
<Routes>
  {/* Public routes */}
  <Route path="/login" element={<LoginPage />} />
  <Route path="/register" element={<RegisterPage />} />

  {/* Protected routes (require auth) */}
  <Route element={<ProtectedLayout />}>
    <Route path="/" element={<DashboardPage />} />
    <Route path="/village/:id" element={<VillagePage />} />
    <Route path="/map" element={<MapPage />} />
    <Route path="/forge" element={<ForgePage />} />
    <Route path="/alliance" element={<AlliancePage />} />
    <Route path="/profile" element={<ProfilePage />} />
  </Route>

  {/* Single-player lore mode */}
  <Route path="/lore" element={<LoreExplorerPage />} />
</Routes>
```

---

## TypeScript Rules

1. **Strict mode**: `"strict": true` in `tsconfig.json`. No exceptions.
2. **No `any`**: Use `unknown` if type is truly unknown, then narrow. Use generics when needed.
3. **API types**: Every API response must have a corresponding TypeScript interface in `types/`.
4. **Enum alternatives**: Prefer union types (`type Kingdom = 'veridor' | 'sylvara' | 'arkazia'`) over TypeScript enums.
5. **Null safety**: Always handle `null` and `undefined` cases. Use optional chaining (`?.`) and nullish coalescing (`??`).

### Type Example

```tsx
// types/game.ts
export type Kingdom = 'veridor' | 'sylvara' | 'arkazia' | 'draxys' | 'nordalh' | 'zandres' | 'lumus';
export type ResourceType = 'food' | 'water' | 'lumber' | 'stone';

export interface Resources {
  food: number;
  water: number;
  lumber: number;
  stone: number;
}

export interface Village {
  id: number;
  name: string;
  x: number;
  y: number;
  isCapital: boolean;
  resources: Resources;
  buildings: Building[];
}

export interface Building {
  id: number;
  type: string;
  level: number;
}
```

---

## WebSocket Integration

### Connection Management

```tsx
// services/websocket.ts
class WebSocketService {
  private ws: WebSocket | null = null;
  private listeners: Map<string, Set<(data: unknown) => void>> = new Map();

  connect(token: string) {
    this.ws = new WebSocket(`${WS_URL}?token=${token}`);
    this.ws.onmessage = (event) => {
      const message = JSON.parse(event.data);
      this.emit(message.type, message.data);
    };
  }

  subscribe(type: string, callback: (data: unknown) => void) {
    if (!this.listeners.has(type)) {
      this.listeners.set(type, new Set());
    }
    this.listeners.get(type)!.add(callback);
    return () => this.listeners.get(type)?.delete(callback); // cleanup
  }

  send(type: string, data: unknown) {
    this.ws?.send(JSON.stringify({ type, data }));
  }

  private emit(type: string, data: unknown) {
    this.listeners.get(type)?.forEach((cb) => cb(data));
  }
}

export const wsService = new WebSocketService();
```

### Reconnection Strategy

The server does **not** replay missed events. When a WebSocket disconnects, the client must:

1. **Detect disconnection** via the `onclose` event.
2. **Retry with exponential backoff**: 1s → 2s → 4s → 8s → max 30s.
3. **Re-authenticate**: Reconnect with current JWT (refresh if expired).
4. **Re-subscribe** to all topics.
5. **Fetch current state** via REST API to sync any missed updates.

```ts
// Simplified reconnection logic inside WebSocketService
private reconnect(token: string) {
  const maxDelay = 30_000;
  let delay = 1_000;

  const attempt = () => {
    this.connect(token);
    this.ws!.onerror = () => {
      setTimeout(attempt, delay);
      delay = Math.min(delay * 2, maxDelay);
    };
    this.ws!.onopen = () => {
      delay = 1_000; // reset on success
      this.resubscribeAll();
    };
  };
  attempt();
}
```

### Using WebSocket in Components

```tsx
// hooks/useWebSocket.ts
import { useEffect } from 'react';
import { wsService } from '../services/websocket';

export function useWebSocketEvent<T>(type: string, handler: (data: T) => void) {
  useEffect(() => {
    const unsubscribe = wsService.subscribe(type, handler as (data: unknown) => void);
    return unsubscribe;
  }, [type, handler]);
}
```

---

## Performance Guidelines

1. **Code splitting**: Use `React.lazy()` + `Suspense` for route-level components. The map, forge, and lore modules should be lazy-loaded.
2. **Memoization**: Use `React.memo()` for map tiles and list items that re-render frequently. Use `useMemo` and `useCallback` judiciously (don't over-optimize).
3. **Virtual scrolling**: If troop lists or building lists grow large, use virtual scrolling (e.g., `@tanstack/react-virtual`).
4. **Image optimization**: Use WebP format for game assets. Lazy-load images below the fold.
5. **Bundle analysis**: Regularly check bundle size with `vite-plugin-visualizer`.

---

## Accessibility (a11y) Basics

- All interactive elements must be keyboard-accessible.
- Use semantic HTML: `<button>`, `<nav>`, `<main>`, `<section>`, `<header>`.
- Add `aria-label` to icon-only buttons.
- Ensure sufficient color contrast (especially across kingdom themes).
- Use `role` attributes where semantic HTML doesn't suffice.

---

## Changelog

| Date | Change |
|------|--------|
| 2026-03-03 | Initial creation of frontend guide |
| 2026-03-03 | Removed .mobile.css convention; mobile overrides now live inside .module.css via @media queries |
| 2026-03-07 | Docs sync: noted that villageStore.ts doesn't exist (replaced by gameStore.ts), hooks/useWebSocket.ts and services/websocket.ts are planned but not yet implemented, ResourceBar/Tooltip/ErrorBoundary components are planned but not yet created. Vite config uses strictPort: true. Added assetStore.ts to stores list. Admin panel has map editor with template system. Kingdom selection page shows admin link for admin users. |
| 2026-03-08 | Troops training UI refactored: removed sidebar `TrainingPanel` component. Training now happens inside `BuildingTrainingModal` — clicking a military building (barracks, stable, colosseum) opens a modal filtered to that building's trainable troops. Non-military buildings open the standard `BuildingDetailModal` for upgrades. `BuildingCard` gained `isMilitary` + `onUpgradeClick` props to show a small ⬆ upgrade icon on military cards. `BuildingGrid` passes `onUpgradeClick` for military buildings. Config helpers added: `getTroopsForBuilding()` and `isMilitaryBuilding()` in `config/troops.ts`. |
| 2026-03-10 | API types strengthened: `BuildingType`, `TroopType`, `Kingdom` union types replace bare strings throughout `types/api.ts`, eliminating ~13 redundant `as` casts in components. App.tsx inline styles extracted to `App.module.css`. Mobile responsive `@media (max-width: 768px)` overrides added to 6 CSS module files (LoadingSpinner, LoginPage, RegisterPage, LandingLayout, BuildingCard, BuildingDetailModal). Shared map utilities extracted to `features/map/mapUtils.ts` (TILE_SIZE, hexColor, hexColorAlpha, screenToTile, tileHash, extractBaseName) — used by MapRenderer.tsx and AdminMapEditorPage.tsx. |
| 2026-03-10 | Full docs sync: React 18→19. Folder structure updated: added `config/` + `config/generated/` for codegen pipeline, removed planned-but-unbuilt `combat/`, `forge/`, `lore/` features, added actual features (`kingdom/`, `admin/`, `profile/`, `season/`, `landing/`). Services updated to match real files (village.ts, training.ts, map.ts, admin.ts, season.ts, template.ts, player.ts). Stores: replaced villageStore with assetStore + themeStore. Hooks: replaced planned hooks with actual (useBuildingDisplayNames, useResourceTicker). Types: kingdom.ts → map.ts. Kingdom type expanded to all 7 playable kingdoms. ResourceType updated from iron/wood/stone/food to food/water/lumber/stone. |
