# Styling Guide

> CSS architecture, theming, fonts, and responsive design conventions. Read before writing any CSS.

---

## Architecture: CSS Modules

All component styles use **CSS Modules** (`.module.css` extension). This provides:
- Automatic class name scoping (no global conflicts)
- Co-location with components
- TypeScript support via Vite's built-in CSS module handling

### File Convention

Every component that has styles **MUST** have this structure:

```
ComponentName/
├── ComponentName.tsx
└── ComponentName.module.css     # Desktop-first styles + mobile overrides via @media
```

### How to Use

```tsx
// ComponentName.tsx
import styles from './ComponentName.module.css';

export function ComponentName() {
  return <div className={styles.container}>...</div>;
}
```

```css
/* ComponentName.module.css */
.container {
  display: flex;
  padding: var(--spacing-md);
  background: var(--bg-primary);
}

.title {
  font-family: var(--font-heading);
  color: var(--text-primary);
}

/* Mobile overrides */
@media (max-width: 768px) {
  .container {
    flex-direction: column;
    padding: var(--spacing-sm);
  }
}
```

All responsive overrides live **inside the same `.module.css` file** using `@media` queries. This keeps CSS Module scoping intact and simplifies the file structure.

---

## Mobile Strategy

### Web-First, Mobile-Second

All styles are written for **desktop first**. Mobile adaptations are overrides written as `@media` queries **inside the same `.module.css` file**.

### Why Not Separate `.mobile.css` Files?

- Separate non-module CSS files break CSS Module scoping
- Developers must maintain two files per component instead of one
- Media queries inside modules are automatically scoped — no `data-component` attributes or `:global` hacks needed
- Standard CSS pattern that all developers are familiar with

### Mobile CSS Pattern

```css
/* ComponentName.module.css */

/* Desktop-first styles */
.container {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: var(--spacing-lg);
  padding: var(--spacing-xl);
}

.sidebar {
  width: 300px;
}

/* Tablet adjustments */
@media (max-width: 1024px) {
  .container {
    grid-template-columns: 1fr;
  }
  .sidebar {
    width: 100%;
  }
}

/* Mobile adjustments */
@media (max-width: 768px) {
  .container {
    padding: var(--spacing-sm);
    gap: var(--spacing-sm);
  }
}
```

### Breakpoints

| Breakpoint | Width | Target |
|-----------|-------|--------|
| **Desktop** | > 1024px | Primary design target |
| **Tablet** | 769px – 1024px | Adjusted layout |
| **Mobile** | ≤ 768px | Simplified, stacked layout |

Document breakpoint values as reference constants. Note that CSS custom properties **cannot** be used inside `@media` queries — use the raw pixel values directly in `@media`:

```css
/* Reference only — these vars CANNOT be used in @media queries.
   Always write @media (max-width: 768px) or @media (max-width: 1024px) directly. */
:root {
  --breakpoint-mobile: 768px;   /* reference */
  --breakpoint-tablet: 1024px;  /* reference */
}
```

---

## Theming

### CSS Custom Properties

All colors, spacing, and typography values use CSS custom properties defined in `:root`. Kingdom theming overrides appearance variables per-kingdom via `[data-kingdom]` selectors.

### Default (Pre-auth / Neutral)

`:root` defines a dark neutral default used before a kingdom is known (login/register pages). Once the player is authenticated and a kingdom is assigned, the `data-kingdom` attribute on `<html>` activates the appropriate kingdom theme overriding `--bg-*`, `--text-*`, `--accent*`, `--border*`, and `--shadow-*` variables.

See `client/src/styles/themes.css` for the full variable definitions.

### Kingdom Theming

The light/dark theme system has been replaced with **kingdom-based theming**. Each of the 8 kingdoms fully overrides all appearance CSS variables via a `data-kingdom` attribute on `<html>`.

**How it works:**

1. `:root` defines the neutral dark default (used on login/register pages before a kingdom is known).
2. After login, `themeStore` reads `player.kingdom` from `authStore` and applies `data-kingdom="<kingdom>"` on `<html>`.
3. Each `[data-kingdom="..."]` selector overrides `--bg-*`, `--text-*`, `--accent*`, `--border*`, and `--shadow-*` variables.
4. The theme is **derived from player data**, not localStorage. No manual toggle exists.

```tsx
// Applied automatically by themeStore
document.documentElement.setAttribute('data-kingdom', kingdom); // e.g. 'veridor'
```

**Kingdom color palettes:**

| Kingdom | Accent (Main) | Background tone | Text | `--text-on-accent` |
|---------|--------------|----------------|------|--------------------|
| **Arkazia** | `#DC143C` Crimson Red | Dark (black) | White | `#FFFFFF` |
| **Draxys** | `#F9A825` Yellow | Dark (black) | White | `#000000` |
| **Drakanith** | `#FF6D00` Orange | Dark (black) | White | `#FFFFFF` |
| **Zandres** | `#795548` Brown | Dark (black) | White | `#FFFFFF` |
| **Veridor** | `#2196F3` Blue | Light (white/blue tint) | Black | `#FFFFFF` |
| **Nordalh** | `#7B1FA2` Purple | Light (white/purple tint) | Black | `#FFFFFF` |
| **Lumus** | `#FBC02D` Golden Yellow | Light (white/yellow tint) | Black | `#000000` |
| **Sylvara** | `#2E7D32` Forest Green | Warm parchment (golden) | Black | `#FFFFFF` |

See `client/src/styles/themes.css` for the full variable definitions per kingdom.

---

## Typography

### Font Stack

```css
:root {
  --font-heading: 'Cinzel', 'Georgia', serif;
  --font-body: 'EB Garamond', 'Georgia', serif;
  --font-mono: 'JetBrains Mono', 'Fira Code', monospace;
}
```

### Font Loading

Load fonts via Google Fonts in `index.html`:

```html
<link rel="preconnect" href="https://fonts.googleapis.com">
<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
<link href="https://fonts.googleapis.com/css2?family=Cinzel:wght@400;600;700&family=EB+Garamond:ital,wght@0,400;0,500;0,600;1,400&display=swap" rel="stylesheet">
```

### Typography Scale

```css
h1 {
  font-family: var(--font-heading);
  font-size: 2.5rem;
  font-weight: 700;
  letter-spacing: 0.02em;
}

h2 {
  font-family: var(--font-heading);
  font-size: 2rem;
  font-weight: 600;
}

h3 {
  font-family: var(--font-heading);
  font-size: 1.5rem;
  font-weight: 600;
}

h4 {
  font-family: var(--font-heading);
  font-size: 1.25rem;
  font-weight: 400;
}

body, p, li, td {
  font-family: var(--font-body);
  font-size: 1.1rem;
  font-weight: 400;
  line-height: 1.6;
}

small, .text-sm {
  font-size: 0.875rem;
}

code, pre {
  font-family: var(--font-mono);
  font-size: 0.9rem;
}
```

### Font Weight Usage

- **Cinzel**: 400 (normal headings), 600 (emphasized headings), 700 (page titles)
- **EB Garamond**: 400 (body text), 500 (emphasized body), 600 (bold body), 400 italic (quotes/lore text)

---

## Spacing System

Use a consistent spacing scale via CSS variables:

```css
:root {
  --spacing-xs: 0.25rem;   /* 4px */
  --spacing-sm: 0.5rem;    /* 8px */
  --spacing-md: 1rem;       /* 16px */
  --spacing-lg: 1.5rem;     /* 24px */
  --spacing-xl: 2rem;       /* 32px */
  --spacing-2xl: 3rem;      /* 48px */
  --spacing-3xl: 4rem;      /* 64px */
}
```

Always use these variables for padding, margin, and gap. Never hardcode pixel values.

---

## Border Radius

```css
:root {
  --radius-sm: 4px;
  --radius-md: 8px;
  --radius-lg: 12px;
  --radius-xl: 16px;
  --radius-full: 9999px;    /* Pill shape */
}
```

---

## Z-Index Scale

```css
:root {
  --z-base: 0;
  --z-dropdown: 100;
  --z-sticky: 200;
  --z-overlay: 300;
  --z-modal: 400;
  --z-tooltip: 500;
  --z-toast: 600;
}
```

---

## CSS Anti-Patterns (DO NOT)

- **No inline styles**: Always use CSS Modules + variables.
- **No `!important`**: Fix specificity instead.
- **No magic numbers**: Use spacing/radius/z-index variables.
- **No global class names**: Use CSS Modules for scoping. Global styles only in `styles/`.
- **No vendor prefixes**: Vite handles autoprefixing via PostCSS.
- **No px for font sizes**: Use `rem` for scalability.
- **No color literals in components**: Always reference `var(--color-name)`.

---

## Global CSS Files

### globals.css

Contains:
- CSS reset (box-sizing, margin removal)
- Root CSS variables (spacing, radius, z-index)
- Base element styles (html, body)

### themes.css

Contains:
- Neutral dark default variables (`:root`)
- 8 kingdom theme overrides (`[data-kingdom="..."]`) for all appearance variables

### typography.css

Contains:
- Font import declarations
- Heading styles (h1–h4)
- Body text styles
- Utility text classes (.text-sm, .text-muted, etc.)

---

## Changelog

| Date | Change |
|------|--------|
| 2026-03-03 | Initial creation of styling guide |
| 2026-03-03 | Added note that CSS custom properties cannot be used in @media queries — use raw pixel values |
| 2026-03-03 | Simplified mobile strategy: removed separate .mobile.css files, all responsive overrides now live inside .module.css via @media queries |
| 2026-03-06 | Replaced light/dark theme toggle with 8 kingdom-based themes. Removed ThemeToggle component. Theme derived from player data. |
