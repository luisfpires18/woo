---
description: "Scaffold a new React component with TypeScript props, CSS Module with mobile overrides, and accessibility built in."
agent: "agent"
argument-hint: "Component name and purpose (e.g., 'ResourceBar — displays resource counts with icons')"
---

# New Component Scaffold

Generate a reusable React component following WOO frontend conventions.

## What to Generate

### 1. Component File

Location depends on scope:
- **Reusable** (2+ places): `client/src/components/{Name}/{Name}.tsx`
- **Feature-specific**: `client/src/features/{feature}/components/{Name}.tsx`

```typescript
import styles from './{Name}.module.css';

interface {Name}Props {
  // Typed props — no `any`
  className?: string;
  onClick?: () => void;
}

export function {Name}({ className, onClick }: {Name}Props) {
  // If clickable: add role="button", tabIndex={0}, Enter/Space handlers
  return (
    <div className={`${styles.container} ${className ?? ''}`}>
      {/* content */}
    </div>
  );
}
```

### 2. CSS Module File

Co-located with component: `{Name}.module.css`

```css
.container {
  /* Desktop-first styles */
  /* Use CSS variables: var(--bg-primary), var(--spacing-md), etc. */
  /* Fonts: 'Cinzel' for headings, 'EB Garamond' for body */
}

@media (max-width: 768px) {
  .container {
    /* Mobile overrides — MANDATORY */
  }
}
```

## Rules

- Use CSS custom properties for all colors, spacing, fonts
- CSS variables CANNOT be used inside `@media` — write pixel values
- Kingdom theming via `data-kingdom` — colors adapt automatically via CSS vars
- If component uses game config: import from `client/src/config/` (never hardcode values)
- If component shows building names: use `useBuildingDisplayNames()` hook

## Reference Pattern

Study `client/src/components/Card/Card.tsx` for accessibility and structure patterns.
