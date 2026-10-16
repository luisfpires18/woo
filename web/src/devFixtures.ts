// Development-only presentation fixtures. import.meta.env.DEV is false in production builds, so this code and the
// query flag are removed there; the deployed app never shows a signed-in user without a real session.
import type { Account } from './landing/Navbar'

/** `?fixture=signed-in` on the Vite dev server renders the signed-in navbar for visual review. */
export function signedInFixture(): Account | null {
  if (!import.meta.env.DEV) return null
  return new URLSearchParams(window.location.search).get('fixture') === 'signed-in'
    ? { nickname: 'Preview player' }
    : null
}
