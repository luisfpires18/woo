# 012: Game UI structure and navigation polish

Task key: GAME-SHELL
Status: WIP
Updated: 2026-10-16
Milestone: UI
Dependencies: LANDING

## Goal

Build the reusable game layout with navigation, resource-header placeholders and page structure.

## Acceptance criteria

Village/map/forge/army/duties views reachable using explicit mock data. Layout works at mobile widths and both themes; shared version footer. No fake saved state or world membership.

## Required reading

- [docs/design/screens-and-themes.md](../../docs/design/screens-and-themes.md)
- [Scope](../scope.md)
- [Prompt protocol](../workflow/implementation-protocol.md)

## Implementation boundary

A separately reviewable planning task. ChatGPT inspects current code and resolves this task's blockers before issuing one pinned Claude prompt. Preserve medieval v1. Split/refine based on actual development; do not implement successors automatically.

## Validation and administration

Run proportionate checks against acceptance criteria, including relevant failure/retry scenarios. Report actual local evidence. ChatGPT checks remote merge/deployment where applicable. Content/images follow admin configuration and cleanup rules; new mechanics remain code changes.

## Result

The game shell lives in `web/src/game/` with styles in `web/src/game.css`. It reuses the task 009 tokens/primitives and the task 010 navbar and footer. No dependency was added.

### Shell and navigation architecture

- **Command bar** (top, full width, dark charcoal in both appearances): the landing navbar in a full-width variant (its brand links to the main page `/worlds`), with the resource header between the brand and the account control. A 3px kingdom-accent stripe runs along its lower edge.
- **Navigation rail** (48rem and wider): a 7rem charcoal column on the left with icon-over-label items. A rail keeps the full content height for future Settlement/Map scenes and leaves room for later destinations; a top tab row would cost scene height and compete with the resource header. All five labels stay visible, so nothing depends on icons alone.
- **Bottom tab bar** (below 48rem): the same list becomes a fixed bar within thumb reach, with short labels (Settlement, Map, Forge, Army, Duties). The full label stays the accessible name (`aria-label`), and each short label is contained in it. The page reserves the bar height plus the safe-area inset, so the footer is never hidden. The resource header gets its own row (a 2-column grid) under the brand/account row. This is a recomposition, not a squeezed rail.
- **Current item**: `aria-current="page"`, a raised plate with a visible edge (at least 3:1), a bold label, a forge-lit icon and a marker bar (left on the rail, top on the bottom bar). Colour is never the only signal.
- **Page structure** (`GameView`): a heading row (h1 with a short accent rule, lead, tags), a framed **stage** (double rule with iron corner fittings drawn in CSS, no texture image) that later scenes fill, and an optional side column for inspectors or lists that moves under the stage below 64rem.
- **Kingdom accent**: `data-kingdom` on the shell root drives only the command-bar stripe, the current-item marker (`--kingdom-on-nav`, readable on charcoal in both appearances), the title rule and the stage's top edge. Surfaces, text and the crimson Log in action stay neutral. Without a kingdom everything is neutral grey.
- **Routing**: no router package. `web/src/routing.ts` (about 25 lines) exposes the pathname through `useSyncExternalStore` and navigates in-app with the History API, so the shell and footer stay mounted. Modified clicks (such as opening a new tab) keep normal browser behaviour. After in-app navigation, focus moves to the new h1 and `document.title` becomes `<View> · Weapons of Order`. Direct loads and reloads use the existing server fallback; a host test now covers `/game/settlement`.

### Route table

| Path | View | Notes |
|---|---|---|
| `/` | Legacy title screen (task 011) | unchanged; no navbar or footer |
| `/worlds` | Landing / main page | unchanged |
| `/game`, `/game/` | redirects (history replace) to `/game/settlement` | keeps the query string |
| `/game/settlement` | Settlement | player-facing term; legacy VILLAGE keys stay internal |
| `/game/map` | Kingdom Map | |
| `/game/forge` | Forge | |
| `/game/army` | Army | |
| `/game/duties` | Kingdom Duties | |
| `/game/<anything else>` | "Not available yet" inside the shell, no current item | link back to Settlement |
| other paths | existing "Not available yet" page | unchanged; links back to `/worlds` |

The view model is `gameViews` in `web/src/game/views.ts`: `id`, `path`, `label`, `shortLabel`, `Icon`, `lead` and `planned` (the later roadmap tasks that fill the view). Navigation, page titles and placeholders all read it.

### Resource-header contract

`ResourceBar` renders any list of `ResourceSlot` (`web/src/game/previewData.ts`): `{ id, label, amount: number | null, capacity?: number | null, iconUrl?: string }`. A `null` amount shows a dash and reads "Not available yet"; `capacity` adds "/ limit"; without `iconUrl` a neutral dashed fallback symbol is shown. Slots wrap on desktop and become a 2-column grid on phones, so the slot count is not a layout decision. Real definitions (025 RESOURCES-ADMIN) and amounts (039 RESOURCE-STATE) replace the preview list without changing the shell.

### Mock and fixture data (all labelled)

- `previewResourceSlots`: four numbered slots, "Resource 1" to "Resource 4", with no amounts, shown with the tag "Preview · resource types not decided". The count of four is arbitrary and is not a roster decision.
- View leads and "Planned for this screen" items only say which roadmap task fills each view; they state no mechanics, numbers or names.
- Every view shows "Layout prototype · nothing is saved". There is no session, player, world membership or saved settlement.
- Dev-server-only fixtures in `web/src/devFixtures.ts`, removed from production builds by `import.meta.env.DEV`: `?fixture=signed-in` (the existing generic "Preview player" navbar) and the new `?kingdom=arkazia|veridor|sylvara`, which shows the tag "Accent preview: <Kingdom>" so it cannot be read as membership. Production is signed out and neutral.

### Accessibility and responsive checks (2026-10-16, Vite dev server, built-in browser and headless Edge)

- Keyboard: Tab order is skip link → brand → account/Log in → the five navigation links → content. Pressing Enter on Kingdom Map changed the URL to `/game/map`, set `aria-current` on it, moved focus to the "Kingdom Map" h1 and changed the title. Browser Back returned to Settlement. `/game/` redirected to `/game/settlement` with the query kept. A 3px blue ring (`--color-nav-focus-ring`) is visible on charcoal.
- Landmarks and names: banner, `navigation "Game"`, region "Resources (preview)", main with one h1, footer. Navigation links are named by their full labels.
- Contrast from computed colours in light/dark × neutral/Arkazia/Veridor/Sylvara: every text pair is at least 4.5:1 (lowest 6.44:1); the current-item marker on charcoal at least 4.82:1; the current plate edge at least 3.28:1; the stage edge at least 3.62:1 on the canvas. The command-bar stripe is decorative and sits on charcoal.
- Overflow: no horizontal overflow at 375, 768, 1024 and 1280px. At 375px every link and button is at least 44px high and wide (tab items are 70×56px).
- Reduced motion: the shell's only motion is colour transitions on `--duration-fast`, which is 0 under `prefers-reduced-motion`.
- Screenshots (not committed; supplied with the implementation report): desktop light and dark, mobile light and dark, keyboard-focus states on desktop and mobile, all three kingdom accents, neutral signed-out, tablet, 1024px and the unknown-route page.

### Unresolved visual decisions

- Rail versus top navigation is a reasoned first proposal for owner review.
- In dark appearance the charcoal rail (`--color-nav-bg`) and the canvas are close in tone; only the rail's border separates them.
- The stage-frame ornament (corner fittings) and its colour, and whether real scenes keep the frame or run edge to edge.
- Resource-header density on phones (two rows for four slots) depends on the final resource count.
- The `--kingdom-on-nav` shades and the neutral grey stripe are first proposals.
- Signed-out game routes show the Log in action. Whether `/game` requires sign-in belongs to IDENTITY/MEMBERSHIP.

### Resumption on the LEGACY-TITLE result (2026-10-16)

The original commit d4f4bbe was rebased onto dev 2e295cf (LEGACY-TITLE merged and deployed). Conflicts and resolutions:

- `web/src/App.tsx`: both sides added routes. Merged order: `/game` redirect and `/game/*` shell first, then `/` title screen, then `/worlds` landing, then "Not available yet". The pathname comes from the History-API hook, so browser Back between `/game/*` and `/worlds` re-renders correctly; a trailing slash is ignored everywhere.
- `web/src/landing/Navbar.tsx`: kept the LEGACY-TITLE brand link to `/worlds` ("Weapons of Order main page") and added the shell's `children`/`wide` props.
- `web/src/main.tsx`: both stylesheets imported (`title.css`, then `game.css`).
- `tests/Woo.Server.Tests/HostTests.cs`: both direct-load cases kept (`/worlds`, `/game/settlement`).
- Documents (`tasks_roadmap.md`, `010-landing-page.md`, `design-system.md`, `screens-and-themes.md`): dev's renumbered versions taken whole, then this task's records re-applied. Git's rename detection paired the old `011-game-ui-shell.md` with `063-mobile-and-reliability.md`; that file was restored to dev's version and the content applied here instead.
- Planned-task labels in `web/src/game/views.ts` and references in `previewData.ts` follow the new numbering (VILLAGE-UI 013, CONSTRUCTION 041, MAP-UI 014, 045–048, FORGE-UI 016, FORGING 043, RECRUIT 042, COMBAT-UI 015, DUTIES 049, RESOURCES-ADMIN 025, RESOURCE-STATE 039).

Fix found during the resumption checks: on the Vite dev server (React StrictMode re-runs effects) a direct load moved focus to the page heading. The shell now moves focus only when the shown view changes; production and in-app navigation behaviour are unchanged.

Visual changes from the integration: none inside the shell. Its brand now opens `/worlds` instead of `/`, and the shell is reached directly at `/game/...` (the title and landing pages have no link to it yet).

## Dispatch history

1. 2026-10-16: dispatched as task 011 against dev 006be69e689180703c9acd3d826aab4476801d22 on branch `feat/011-game-shell`; local implementation commit d4f4bbe7f827aae208f611c514bcf7d84d84fd8a (game shell under `/game/`, screenshots supplied). Not reviewed, not pushed.
2. 2026-10-16: paused by the owner while the legacy title-screen migration (issue [#65](https://github.com/luisfpires18/woo/issues/65)) was inserted as [011 LEGACY-TITLE](011-legacy-title-screen.md); GAME-SHELL became 012 with its key and scope unchanged. The branch `feat/011-game-shell` is preserved unchanged at d4f4bbe for later review and rebase; its name is historical and is not reused. On resumption it must be rebased onto the LEGACY-TITLE result (the landing page moves to `/worlds`, and its own task records and roadmap edits were written before this renumbering).
3. 2026-10-16: resumed after LEGACY-TITLE was approved, merged and deployed at dev 2e295cfaacf4f56e22f8f9b51590bdeb8c935c52. `feat/011-game-shell` was rebased in place onto 2e295cf (conflicts listed under Result); the original commit d4f4bbe remains the recorded first implementation.

## Execution record

Task-spec commit: 006be69e689180703c9acd3d826aab4476801d22, dispatched 2026-10-16 as task 011; resumed 2026-10-16 as task 012
Expected dev base: 2e295cfaacf4f56e22f8f9b51590bdeb8c935c52 (resumption; verified equal to origin/dev, 2026-10-16). Original base 006be69e689180703c9acd3d826aab4476801d22
Implementation branch: feat/011-game-shell (local only; historical name kept)
Implementation commit: original d4f4bbe7f827aae208f611c514bcf7d84d84fd8a; rebased commit SHA given in the resumption report
Review: Pending ChatGPT review
Merged dev commit: None
Deployment: Deploys through DEV-DEPLOY after integration; not yet run
Issues: the historical branch name keeps 011; the task is now 012 after the LEGACY-TITLE insertion (see dispatch history). The frontend/design/accessibility skills in the inventory are not installed on this computer. No automated frontend tests.

## TLDR

Next: ChatGPT review of the resumed commit and the screenshots in screenshots/012/.
Done: game shell rebased onto the LEGACY-TITLE result; routes `/`, `/worlds` and `/game/*` coexist; full validation rerun.
Issues: see unresolved visual decisions.
