# 010: First landing page UI polish

Task key: LANDING
Status: WIP
Updated: 2026-10-16
Milestone: UI
Dependencies: DESIGN-SYSTEM

## Goal

Design and implement a clear WOO landing page using the approved visual direction.

## Acceptance criteria

Desktop/mobile page introduces the medieval game, offers clear entry actions and includes version footer. No claims of implemented future features; review actual screenshots/browser behavior. Placeholder art allowed. Follow the approved mockups and their mandatory written correction: no inline login, no Overview/Worlds/theme navbar links, profile menu Settings opens a separate page and never embeds appearance controls. All screenshot activity/world data remains illustrative.

## Required reading

- [docs/design/screens-and-themes.md](../../docs/design/screens-and-themes.md)
- [Approved landing mockup and corrections](../../mockup/00-landing-page.md)
- [Scope](../scope.md)
- [Prompt protocol](../workflow/implementation-protocol.md)

## Implementation boundary

A separately reviewable planning task. ChatGPT inspects current code and resolves this task's blockers before issuing one pinned Claude prompt. Preserve medieval v1. Split/refine based on actual development; do not implement successors automatically.

## Validation and administration

Run proportionate checks against acceptance criteria, including relevant failure/retry scenarios. Report actual local evidence. ChatGPT checks remote merge/deployment where applicable. Content/images follow admin configuration and cleanup rules; new mechanics remain code changes.

## Result

The development skeleton and design-foundation preview are replaced by the landing page (`web/src/landing/`), built on the task 009 tokens and primitives.

- Slim navbar on dark charcoal in both appearances: anvil + WOO at left (link home); signed out shows a crimson **Log in** link to `/login`; signed in shows avatar initials + nickname with an account menu containing only Profile (`/profile`), Settings (`/settings`) and Log out. No section links, no theme toggle, no appearance controls.
- Account menu: disclosure button (`aria-expanded`, `aria-controls`, label "<nickname>: account menu"); Enter/Space/ArrowDown open on the first item, ArrowUp on the last; arrows, Home and End move; Escape closes and returns focus; Tab out or a click outside closes. Log out has no behaviour until real sessions exist.
- Title banner: an intentional CSS placeholder (stone, forge light, cool rune light, lava floor, real HTML "Weapons of Order" title, visible "Placeholder banner · final artwork pending" note). It is identical in light and dark. `--banner-art` in `tokens.css` takes the final approved asset later without redesigning the page; the screenshot mockup is not used as artwork.
- Introduction: "Forge your kingdom’s future", tagline, and the three features with inline SVG icons (Settlement building, Weapons & armour, Kingdom warfare).
- Metrics band (Open worlds, Registered players, Online now) showing dashes and a dashed "Preview · live figures appear when worlds open" tag; screen readers hear "Not available yet".
- Choose a world: two "Example world" cards marked Preview (era, pace, dashed player/online counts, kingdom badges with names, a disabled **Select** button explained by the note "World selection opens once worlds are available. These cards show the layout only.").
- Shared footer: "Weapons of Order" at left, the real build label at right (`0.0.1-dev (local)` in development), with the update notice preserved.
- Routing: no router package. `/` renders the landing page; every other path renders a "Not available yet" page with the shared navbar and footer, so `/login`, `/settings` and `/profile` stay separate page boundaries without implementing those screens. The server already returns `index.html` for client paths and 404 for unknown `/api` routes.
- Development fixture: on the Vite dev server only, `?fixture=signed-in` renders the signed-in navbar ("Preview player") for visual review. `import.meta.env.DEV` removes it from production builds (verified: the production bundle contains neither the flag nor the fixture name).

Preview data is isolated in `web/src/landing/previewData.ts` and labelled in the UI; no real world, player or activity count is shown.

### Visual decisions and deliberate differences from the mockups

- Navbar is dark charcoal in light appearance too (latest written direction), instead of the light navbar in the light screenshot.
- "Village building" and "Build your village" became "Settlement building" and "Build your settlement" (settlements is the preferred player-facing term).
- Friends Alpha I/II, Open badges, counts and the profile name/portrait are not reproduced: preview cards say "Example world 1/2", "Preview", and use dashes; the signed-in state exists only in the dev fixture with generic initials.
- Select buttons are disabled with an explanation rather than active.
- Footer omits About/Help links and "Concept preview" text: no such pages exist.
- Mobile is a re-composition, not a squeeze: features stack, metrics become labelled rows, world cards stack with kingdom badges wrapping, the account name truncates.

### Accessibility and responsive checks (2026-10-16, dev server, built-in browser)

- Headings: h1 Weapons of Order (banner), h2 introduction / visually hidden "World activity" / Choose a world, h3 per world, h4 Kingdoms. Landmarks: header, main (skip link target), labelled sections, footer.
- Keyboard: account menu sequence above verified with real key presses; visible 3px focus ring on charcoal (`--color-nav-focus-ring`) and in the menu.
- Contrast from computed colours, both appearances: all landing text pairs at least 7.1:1 except body pairs already covered by task 009; Log in text 6.26–6.75:1 and edge on charcoal at least 4.82:1; disabled Select text at least 4.86:1; feature icons now use the stronger action edge token.
- 1280px and 375px: no horizontal overflow (`scrollWidth` equals viewport); every button/link control at least 44px high.
- Reduced motion: the only animation (banner light flicker) is disabled under `prefers-reduced-motion`.
- `/login` and `/settings` render "Not available yet" with no form; no console errors.

## Execution record

Task-spec commit: ed1f3c2975e2041dc7375f1ae900e2d797bce121 (owner dispatch, 2026-10-16)
Expected dev base: ed1f3c2975e2041dc7375f1ae900e2d797bce121 (verified equal to local dev and origin/dev, 2026-10-16)
Implementation branch: feat/010-landing-page (local only)
Implementation commit: local commit on that branch; SHA given in the implementation report
Review: Pending ChatGPT review
Merged dev commit: None
Deployment: Deploys through DEV-DEPLOY after integration; not yet run
Issues: the requested frontend/image-to-code/design/accessibility skills are not installed on this computer (only the PixiJS collection), so the work follows the mockups, written corrections and design records directly. Final banner asset and real Login/Settings/Profile pages remain open. No automated frontend tests.

## TLDR

Next: ChatGPT review of the local commit and screenshots.
Done: landing page with navbar/account menu, placeholder banner, introduction, preview metrics and world cards, shared footer; responsive and accessibility checks.
Issues: see execution record.
