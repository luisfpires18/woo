---
description: "Use when writing or modifying CSS Module files. Covers responsive design, kingdom theming, typography, spacing, and z-index conventions."
applyTo: "client/**/*.module.css"
---

# CSS Module Conventions

Full reference: `docs/04-frontend/styling-guide.md`

## CSS Modules Only

No inline styles. No global CSS for components. Every component gets a co-located `.module.css`.

## Responsive Design (Mandatory)

- **Desktop-first**: Base styles are desktop layout.
- **Every component** MUST have `@media (max-width: 768px)` overrides in the same `.module.css` file.
- No separate `.mobile.css` files.
- **CSS variables CANNOT be used inside `@media` queries** — write pixel values directly: `@media (max-width: 768px)`, not `@media (max-width: var(--breakpoint-mobile))`.

## Kingdom Theming

Theme is driven by `[data-kingdom="..."]` on `<html>`. Use CSS custom properties:

```css
.container {
  background: var(--bg-primary);
  color: var(--text-primary);
  border-color: var(--border-primary);
}
```

Available variable groups: `--bg-*`, `--text-*`, `--accent-*`, `--border-*`, `--shadow-*`.

## Typography

- **Headings**: `font-family: 'Cinzel', serif`
- **Body text**: `font-family: 'EB Garamond', serif`
- Never use system fonts for game UI.

## Spacing

Use CSS custom properties: `--spacing-xs` (4px) through `--spacing-3xl` (64px).

## Z-Index Scale

| Layer | Value |
|-------|-------|
| Base | 1 |
| Dropdown | 100 |
| Sticky | 200 |
| Overlay | 300 |
| Modal | 400 |
| Tooltip | 500 |
| Toast | 600 |

## Colors

Never hardcode colors — use `var(--accent-primary)`, `var(--bg-secondary)`, etc.
Each kingdom overrides these variables in `themes.css`.
