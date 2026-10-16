// Development-only presentation fixtures. import.meta.env.DEV is false in production builds, so this code and the
// query flag are removed there; the deployed app never shows a signed-in user without a real session.
import type { Account } from './landing/Navbar'
import { kingdoms, type Kingdom } from './theme'

/** `?fixture=signed-in` on the Vite dev server renders the signed-in navbar for visual review. */
export function signedInFixture(): Account | null {
  if (!import.meta.env.DEV) return null
  return new URLSearchParams(window.location.search).get('fixture') === 'signed-in'
    ? { nickname: 'Preview player' }
    : null
}

/**
 * `?kingdom=arkazia|veridor|sylvara` on the Vite dev server previews the game shell's kingdom accent. It is a visual
 * fixture, not membership: the shell labels it "Accent preview" and production builds always stay neutral.
 */
export function kingdomFixture(): Kingdom | null {
  if (!import.meta.env.DEV) return null
  const value = new URLSearchParams(window.location.search).get('kingdom')
  return kingdoms.find((k) => k.id === value)?.id ?? null
}

/** `?scene=fail` on the Vite dev server simulates a renderer that cannot start, to check the Settlement fallback. */
export function sceneFailureFixture(): boolean {
  return import.meta.env.DEV && new URLSearchParams(window.location.search).get('scene') === 'fail'
}
