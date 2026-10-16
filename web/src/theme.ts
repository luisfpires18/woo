// Theme contract (docs/design/design-system.md): appearance and kingdom accent are independent.
// Appearance is applied to <html data-appearance>; kingdom accents are scoped with data-kingdom on any element.

export type Appearance = 'light' | 'dark'
/** 'system' follows the operating system and is the default until the Settings page stores a choice. */
export type AppearancePreference = Appearance | 'system'
export type Kingdom = 'arkazia' | 'veridor' | 'sylvara'

/** Player-facing kingdom names (confirmed); colours live only in tokens.css. */
export const kingdoms: readonly { id: Kingdom; name: string }[] = [
  { id: 'arkazia', name: 'Arkazia' },
  { id: 'veridor', name: 'Veridor' },
  { id: 'sylvara', name: 'Sylvara' },
]

const darkQuery = () => window.matchMedia('(prefers-color-scheme: dark)')

export const resolveAppearance = (preference: AppearancePreference): Appearance =>
  preference === 'system' ? (darkQuery().matches ? 'dark' : 'light') : preference

let stopFollowingSystem: (() => void) | undefined

/**
 * Applies an appearance preference to the document. With 'system' the document keeps following
 * operating-system changes until another preference is applied.
 */
export function applyAppearance(preference: AppearancePreference, root: HTMLElement = document.documentElement): void {
  stopFollowingSystem?.()
  stopFollowingSystem = undefined
  const set = () => {
    root.dataset.appearance = resolveAppearance(preference)
  }
  set()
  if (preference === 'system') {
    const query = darkQuery()
    query.addEventListener('change', set)
    stopFollowingSystem = () => query.removeEventListener('change', set)
  }
}

/** Sets or clears the document-wide kingdom accent (null = neutral). Prefer data-kingdom on a subtree for scoped accents. */
export function applyDocumentKingdom(kingdom: Kingdom | null, root: HTMLElement = document.documentElement): void {
  if (kingdom) root.dataset.kingdom = kingdom
  else delete root.dataset.kingdom
}
