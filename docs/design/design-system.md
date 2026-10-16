# Design system foundation

Updated: 2026-10-16. Status: task 009 DESIGN-SYSTEM, approved and merged (ed1f3c2); extended by task 010 LANDING and task 011 LEGACY-TITLE. Applies the confirmed direction in [screens and themes](screens-and-themes.md); exact shades below are the first implemented proposal, open to owner adjustment.

## Theme contract

Appearance and kingdom accent are independent. Components read CSS custom properties only (`web/src/tokens.css`); none hardcodes a colour or a kingdom.

| Concept | Mechanism | Default |
|---|---|---|
| Appearance | `<html data-appearance="light\|dark">`, set by `applyAppearance(preference)` in `web/src/theme.ts`. Preference is `system`, `light` or `dark`; `system` follows the operating system live (media query listener) | `system`, applied before the first render in `main.tsx`; light is the CSS fallback if the attribute is missing |
| Kingdom accent | `data-kingdom="arkazia\|veridor\|sylvara"` on any element scopes `--kingdom-*` to that subtree; `applyDocumentKingdom(kingdom \| null)` sets it document-wide | none: neutral accent |
| Brand action | `--color-action-*`: crimson in both appearances, independent of the kingdom | always |
| Semantic feedback | `--color-success-*`, `--color-warning-*`, `--color-danger-*`, `--color-info-*` (`text`, `surface`, `border`) | always; never replaced by kingdom colours |

There is no Settings page yet: a later task stores the player's explicit choice and calls `applyAppearance('light' | 'dark')`. Theme and authentication state stay independent.

## Tokens

- Typography: `--font-sans` (system UI stack), `--font-display` (local serif stack for headings, matching the approved landing direction), `--font-mono` (version label). No remote fonts. Sizes `--text-xs` to `--text-3xl` (the two largest scale with the viewport), `--leading-*`, `--weight-*`.
- Spacing `--space-1` (4px) to `--space-8` (64px); `--page-gutter` scales from 16px to 32px.
- Sizing: `--control-height` 44px (minimum touch target for buttons and inputs), `--content-max` 72rem, `--content-narrow` 40rem.
- Shape: `--radius-sm|md|lg|pill`, `--border-width`, `--border-width-strong`.
- Surfaces and text: `--color-canvas` (warm ivory / charcoal page), `--color-surface`, `--color-surface-raised`, `--color-surface-sunken`, `--color-text`, `--color-text-muted`, `--color-text-inverse`, `--color-border` (decorative), `--color-border-strong` (control boundaries).
- Navigation: `--color-nav-bg` (dark charcoal in both appearances), `--color-nav-text`, `--color-nav-border`, `--color-nav-border-strong`, `--color-nav-focus-ring`, `--color-nav-action-border` (3:1 edge for crimson actions on charcoal), `--color-nav-avatar-bg`; used by the landing navbar and ready for the game shell (task 012).
- Title banner (task 010): `--banner-art` (`none` until the approved asset exists), `--banner-stone`, `--banner-stone-light`, `--banner-lava`, `--banner-ember`, `--banner-rune-light`, `--banner-title`, `--banner-title-shadow`, `--banner-height`; identical in both appearances.
- Title screen (task 011): `--title-void`, `--title-bone`, `--title-ember` (legacy weapons-of-order theme values, identical in both appearances because they sit on artwork), `--font-title` (Cinzel, self-hosted under the SIL OFL in `web/src/assets/fonts/`), `--text-title`, `--text-title-menu`.
- States: `--color-disabled-*`, `--color-focus-ring` (blue, distinct from crimson), `--focus-ring-width`, `--focus-ring-offset`.
- Elevation: `--shadow-1`, `--shadow-2`. Motion: `--duration-fast`, `--duration-slow` (both 0 with `prefers-reduced-motion`).
- Kingdom accent (per appearance): `--kingdom-primary`, `--kingdom-on-primary`, `--kingdom-primary-hover`, `--kingdom-secondary` (Arkazia black, Veridor silver, Sylvara gold), `--kingdom-text` (accent text on surfaces), `--kingdom-tint` (subtle accent surface), `--kingdom-border`.

Artwork is never filtered or darkened for dark appearance; the same images are used in both.

## Primitives (`web/src/ui.tsx`, styles in `web/src/styles.css`)

| Primitive | Purpose and accessibility |
|---|---|
| `Page`, `PageHeader` | Centred content column with responsive gutters; one `h1` and optional lead |
| `Panel` | Bordered surface; with `title` it is a labelled region (`aria-labelledby`); `kingdom` adds a scoped accent top border |
| `Cluster` | Wrapping horizontal group for buttons/badges |
| `Button` | Variants `primary` (brand crimson), `secondary` (neutral outline), `kingdom` (nearest kingdom accent). 44px minimum height and width, hover, visible focus ring, `disabled`, and `aria-pressed` selected state shown with a border and check mark |
| `TextField` | Visible label (with "(required)"), optional hint, error message; `aria-invalid` and `aria-describedby` link hint and error; invalid border is thicker, and the error text has an icon, not colour alone |
| `StatusMessage` | Success / warning / error / note with colour, icon and a text label |
| `KingdomBadge` | Kingdom ownership marker: name text plus a two-colour swatch, so colour is never the only indicator |

The shared `AppFooter` (version label and update notice) uses these tokens; the notice's reload control is a `Button`. Utilities: `.visually-hidden`, `.skip-link`; links can use `ui-button` classes (`a.ui-button`). The task 009 development preview was replaced by the landing page in task 010.

## Accessibility and responsive checks (2026-10-16, dev server, built-in browser)

- Contrast measured from the computed tokens for every pair the primitives use, 100 pairs across light/dark and neutral/Arkazia/Veridor/Sylvara: all text pairs at least 4.84:1 (AA 4.5), all control boundaries and focus rings at least 3.93:1 (AA 3). Disabled text is also at least 4.5:1 although WCAG exempts it. Filled dark buttons carry a lighter border token so their edge also reaches 3:1 against dark surfaces.
- Keyboard: Tab reaches controls in order; the focused control shows a 3px blue outline with 2px offset.
- Form: submitting the empty sample field sets `aria-invalid="true"`, the input's description reads the hint and the error, and the border becomes 2px danger colour.
- 375px viewport: no horizontal overflow; smallest button/input height 44px.
- Appearance: changing the emulated system scheme switched `data-appearance` from light to dark without a reload.

## Unresolved visual decisions

- Exact shades are a first proposal. In particular the brand crimson action and Arkazia's crimson accent are close by design (Arkazia is crimson); a screen that shows both side by side may need a clearer Arkazia shade or a non-crimson action there.
- Dark appearance direction beyond "dark charcoal": the ivory reference is not the final dark design.
- Display serif: a local system serif stack stands in until the owner chooses whether to bundle a specific font.
- Kingdom secondary colours (black/silver/gold) are used only as swatch rings and stripes so far; their use in larger areas is undecided.
- Focus-on-first-invalid-field and live error announcement are left to the screens that own forms (login, settings).
