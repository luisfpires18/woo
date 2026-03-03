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

Use CSS custom properties for breakpoints:

```css
:root {
  --breakpoint-mobile: 768px;
  --breakpoint-tablet: 1024px;
}
```

---

## Theming

### CSS Custom Properties

All colors, spacing, and typography values use CSS custom properties defined in `:root`. Theme switching swaps these variables.

### Dark Mode (Default)

```css
:root,
[data-theme="dark"] {
  /* Colors */
  --bg-primary: #0a0a0a;
  --bg-secondary: #1a1a1a;
  --bg-tertiary: #2a2a2a;
  --bg-elevated: #1e1e1e;

  --text-primary: #f0f0f0;
  --text-secondary: #b0b0b0;
  --text-muted: #707070;

  --accent: #DC143C;          /* Crimson */
  --accent-hover: #B22222;    /* Firebrick */
  --accent-light: #FF6B6B;    /* Lighter crimson for subtle highlights */

  --border: #333333;
  --border-light: #444444;

  /* Status Colors */
  --success: #2ECC71;
  --warning: #F39C12;
  --error: #E74C3C;
  --info: #3498DB;

  /* Shadows */
  --shadow-sm: 0 1px 2px rgba(0, 0, 0, 0.5);
  --shadow-md: 0 4px 6px rgba(0, 0, 0, 0.5);
  --shadow-lg: 0 10px 15px rgba(0, 0, 0, 0.5);
}
```

### Light Mode

```css
[data-theme="light"] {
  --bg-primary: #FAFAFA;
  --bg-secondary: #F0F0F0;
  --bg-tertiary: #E5E5E5;
  --bg-elevated: #FFFFFF;

  --text-primary: #1a1a1a;
  --text-secondary: #4a4a4a;
  --text-muted: #8a8a8a;

  --accent: #001F5B;          /* Navy Blue */
  --accent-hover: #003087;    /* Brighter navy */
  --accent-light: #1A4F8B;    /* Lighter navy for subtle highlights */

  --border: #D0D0D0;
  --border-light: #E0E0E0;

  --success: #27AE60;
  --warning: #E67E22;
  --error: #C0392B;
  --info: #2980B9;

  --shadow-sm: 0 1px 2px rgba(0, 0, 0, 0.1);
  --shadow-md: 0 4px 6px rgba(0, 0, 0, 0.1);
  --shadow-lg: 0 10px 15px rgba(0, 0, 0, 0.1);
}
```

### Theme Switching

The theme is set via a `data-theme` attribute on `<html>`:

```tsx
// In a theme hook or store
document.documentElement.setAttribute('data-theme', theme); // 'dark' or 'light'
```

---

## Kingdom Theming (Future)

> Planned for Phase 7. Document for reference.

Each kingdom will have additional CSS custom properties loaded dynamically:

```css
[data-kingdom="veridor"] {
  --kingdom-primary: #001F5B;    /* Navy Blue */
  --kingdom-secondary: #4A90D9;  /* Light Blue */
  --kingdom-accent: #C0C0C0;    /* Silver */
}

[data-kingdom="sylvara"] {
  --kingdom-primary: #1B5E20;    /* Deep Green */
  --kingdom-secondary: #FF8F00;  /* Amber */
  --kingdom-accent: #8D6E63;    /* Brown */
}

[data-kingdom="arkazia"] {
  --kingdom-primary: #424242;    /* Iron Grey */
  --kingdom-secondary: #DC143C;  /* Crimson */
  --kingdom-accent: #B8860B;    /* Dark Gold */
}
```

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
- Dark mode variables (`:root` and `[data-theme="dark"]`)
- Light mode variables (`[data-theme="light"]`)
- Kingdom theme variables (future)

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
| 2026-03-03 | Simplified mobile strategy: removed separate .mobile.css files, all responsive overrides now live inside .module.css via @media queries |
