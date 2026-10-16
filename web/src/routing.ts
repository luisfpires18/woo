// Minimal client routing (no router package): the current pathname as React state plus in-app navigation through
// the History API. Full page loads and reloads of any path still work because the server returns index.html for
// client paths (see src/Woo.Server/Program.cs).
import { useSyncExternalStore, type MouseEvent } from 'react'

const subscribe = (onChange: () => void) => {
  window.addEventListener('popstate', onChange)
  return () => window.removeEventListener('popstate', onChange)
}

export const usePathname = () => useSyncExternalStore(subscribe, () => window.location.pathname)

export function navigate(path: string, { replace = false } = {}) {
  if (replace) window.history.replaceState(null, '', path)
  else window.history.pushState(null, '', path)
  window.dispatchEvent(new PopStateEvent('popstate'))
}

/** onClick for in-app links: plain left clicks navigate without a reload; modified clicks (new tab etc.) are left alone. */
export function followLink(event: MouseEvent<HTMLAnchorElement>) {
  if (event.defaultPrevented || event.button !== 0 || event.metaKey || event.ctrlKey || event.shiftKey || event.altKey)
    return
  event.preventDefault()
  navigate(event.currentTarget.pathname)
}
