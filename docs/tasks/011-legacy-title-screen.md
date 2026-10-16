# 011: Port legacy title screen, then retire old staging Azure

Task key: LEGACY-TITLE
Status: WIP
Updated: 2026-10-16
Milestone: UI
Dependencies: LANDING

## Goal

Owner-requested intervening migration, issue [#65](https://github.com/luisfpires18/woo/issues/65). Make the full-screen forge/cavern title screen of the discontinued `luisfpires18/weapons-of-order` app (old staging `https://app-woo-staging-le5jrs.azurewebsites.net/`) the WOO entry screen, with ENTER WORLD leading to the existing landing page. After the new screen is deployed and verified, retire the old staging Azure resources under an owner-approved deletion list.

## Phases

1. **Migration** (this dispatch): port the title screen; move the landing page behind it. No Azure change.
2. **Azure retirement** (gated, not started): only after ChatGPT review, owner integration, successful deployment and live verification of the new title on WOO dev: inventory every Azure resource used only by the old staging app, confirm none is shared with WOO or another live project, show the owner the exact deletion list, delete only after explicit owner approval, then verify the old URL no longer serves the app and record what was removed. The old GitHub repository is kept unchanged as history.

## Acceptance criteria (migration phase)

`/` is the legacy title: original forge artwork full-viewport, centred WEAPONS / OF / ORDER wordmark with ember rules around OF, one ENTER WORLD action, no SETTINGS, no navbar, no footer. ENTER WORLD opens the task 010 landing page on its own route with navbar, world-preview cards, footer and behaviour unchanged. Direct loads of both routes work. No old authentication, session, world-guard, settings or routing code; no fake account or world state. Keyboard reachable with strong focus, decorative artwork, no horizontal scrolling, usable at 375px, sufficient contrast over the artwork.

## Required reading

- Issue [#65](https://github.com/luisfpires18/woo/issues/65)
- [Screens and themes](../design/screens-and-themes.md)
- [Design system](../design/design-system.md)
- [Landing task](010-landing-page.md)
- [Prompt protocol](../workflow/implementation-protocol.md)

## Result

### Provenance (old repository `luisfpires18/weapons-of-order`, main at f74ce3cb1b9b8c332812b89cd27e8bbaa6c7a210, read 2026-10-16)

| Old source | Used as | Notes |
|---|---|---|
| `art/backgrounds/forge-16x9.png` | `web/src/assets/title/forge-16x9.png` | Byte-identical copy, 1672×941 PNG, 1,952,386 bytes, SHA-256 4ab24d6eeffa377b666bf6950db550d94941f739262d4521d5f5c641148c3b75 |
| `web/src/screens/TitleScreen.tsx` | `web/src/title/TitleScreen.tsx`, `web/src/title.css` | Layout and scrim gradient (void 70% → 30% at 60% → 30%) ported; `useNavigate`, `WORLD_PATH`, `SettingsPanel` and the SETTINGS item dropped |
| `web/src/components/Wordmark.tsx` | same files | Three lines, per-line 0.06em tracking with trailing-space correction, OF at 0.55× with 1.5em × 2px ember rules and 0.5em gaps; visible text is aria-hidden and the heading name comes from a visually hidden "Weapons of Order" |
| `web/src/components/TitleMenu.tsx` | same files | Cinzel semibold 1.5rem, 0.10em tracking, bone text, ember on hover/focus, focus outline offset 6px. Arrow-key roving between items not ported (one item) |
| `web/src/theme/tokens.css` | `web/src/tokens.css` `--title-*` | `--color-void` #030A10, `--color-bone` #E6E1D3, `--color-ember` #FB5000, `--font-display` Cinzel, `--text-title` clamp(3rem, 7vw, 5.5rem), `--text-menu` 1.5rem |

The old app loaded Cinzel from Google Fonts at runtime. WOO uses no runtime font downloads, so the same font (Cinzel v26 latin variable woff2 served by Google Fonts, SHA-256 09941fb1c169c38fd414536b37690057fc01b3117a5c63dd6571186540c8f370) is self-hosted in `web/src/assets/fonts/` with its SIL Open Font License (`Cinzel-OFL.txt`, from google/fonts `ofl/cinzel/OFL.txt`). Nothing else from the old repository was copied: no Tailwind, react-router, zustand, React Query, auth/session, world guards, SettingsPanel or EmberRule.

### Routes and navigation

| Path | Screen |
|---|---|
| `/` | Legacy title screen (no navbar/footer) |
| `/worlds` | Task 010 landing page, unchanged content and behaviour; chosen because ENTER WORLD leads to "Choose a world" |
| any other path (e.g. `/login`, `/settings`, `/profile`) | Existing "Not available yet" page |

- ENTER WORLD is a real link to `/worlds` (full navigation, like the other links; no router added). A trailing slash is ignored.
- The navbar brand now goes to `/worlds` ("Weapons of Order main page") rather than to the title, and "Not available yet" offers "Back to the main page" (`/worlds`). The title is the entry screen; in-app "home" means the main page.
- Direct loads of `/` and `/worlds` use the existing server fallback; a host test now covers `/worlds`.

### Visual fidelity

Compared against the old staging app with headless Edge at 1280×800 and 375×812: same artwork crop, wordmark size/position, ember rules and menu typography. With SETTINGS removed, the centred column would drop about 30px onto the anvil's rune stones, so the content reserves the removed row's height (1.5 line + 1.5rem gap) and ENTER WORLD stays exactly where it was. The artwork is shown with `object-fit: cover` and never filtered; it is identical in light and dark appearance.

### Accessibility and responsive checks (2026-10-16, Vite dev server, built-in browser and headless Edge)

- Structure: `main` → `h1` (accessible name "Weapons of Order") → `nav "Main menu"` → one link; artwork `alt=""`, scrim `aria-hidden`.
- Keyboard: first Tab focuses ENTER WORLD (3px ember outline, 6px offset, text turns ember); Enter opens `/worlds`.
- Contrast measured against the brightest artwork pixel under each element after the scrim: ENTER WORLD bone text at least 13.98:1 and ember (focus ring and hover/focus text) at least 5.43:1 at 1280×800 and 375×812; wordmark at least 10.54:1. The ember rules are decorative.
- Layout: no horizontal overflow at 375×812, 812×375, 1280×800 and 1920×1080; ENTER WORLD is 44px high.
- Reduced motion: the only motion is the colour transition on `--duration-fast`, which is 0 under `prefers-reduced-motion`.
- `/worlds` direct load and after ENTER WORLD renders the landing page; `/settings` renders "Not available yet".
- Production bundle: no old auth/settings/session/router code (searched for SettingsPanel, SETTINGS, react-router, zustand, tanstack, useNavigate, WORLD_PATH, session); the only "Settings" strings are the task 010 account-menu link.
- Screenshots (not committed): old staging and new title at desktop and mobile, focus states, landscape, 1920px, `/worlds` after ENTER WORLD and direct, and the not-available page.

### Unresolved decisions

- The artwork is the original 1.95 MB PNG (served immutable with a content hash). A WebP/AVIF derivative would load faster; it was not made, to keep the original bytes.
- The title has no version footer or update notice (owner direction); both remain on `/worlds` and every other page. Confirm this exception to "all shared layouts include the app-version footer".
- The artwork is a fixed bundled asset, not an admin-controlled image slot; decide whether the title background joins the admin image-slot rules later.
- Whether the navbar brand should return to the title screen instead of `/worlds`.
- Cinzel is used only on the title; whether it replaces the local serif `--font-display` elsewhere is undecided.

## Validation and administration

Run proportionate checks against acceptance criteria. Report actual local evidence. ChatGPT checks remote merge/deployment. Azure retirement is destructive and needs the explicit owner approval described above; it is never part of an implementation commit.

## Execution record

Task-spec commit: Issue #65, created 2026-10-16T14:07:25Z; owner dispatch 2026-10-16 (migration phase only)
Expected dev base: 006be69e689180703c9acd3d826aab4476801d22 (verified equal to local dev and origin/dev, 2026-10-16)
Implementation branch: feat/011-legacy-title (local only)
Implementation commit: local commit on that branch; SHA given in the implementation report
Review: Pending ChatGPT review
Merged dev commit: None
Deployment: Deploys through DEV-DEPLOY after integration; not yet run
Azure retirement: Not started; gated on live verification and explicit owner approval. No old Azure resource was inspected, stopped, changed or deleted
Issues: the frontend/design/accessibility skills in the inventory are not installed on this computer; the port follows the old source directly. No automated frontend tests.

## TLDR

Next: ChatGPT review of the local commit and screenshots; after integration and live verification, the gated Azure inventory.
Done: legacy title at `/` with original artwork and wordmark, ENTER WORLD to the landing page at `/worlds`, provenance and checks recorded.
Issues: see unresolved decisions; Azure retirement not started.
