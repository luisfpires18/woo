import { useEffect, useState } from 'react'
import { buildLabel, fetchServerBuild, isStale, pageBuild } from './build'

/** Shared footer for every layout: release identity plus a notice when a newer build is deployed. */
export function AppFooter() {
  const [stale, setStale] = useState(false)

  // Check on load and whenever the tab becomes visible again; no polling. Only a confirmed
  // newer server build shows the notice, and an unavailable server never clears or raises it.
  useEffect(() => {
    const controller = new AbortController()
    const check = () => {
      if (document.visibilityState !== 'visible') return
      fetchServerBuild(controller.signal).then((server) => {
        if (!controller.signal.aborted && isStale(pageBuild, server)) setStale(true)
      })
    }
    check()
    document.addEventListener('visibilitychange', check)
    return () => {
      document.removeEventListener('visibilitychange', check)
      controller.abort()
    }
  }, [])

  return (
    <footer className="app-footer">
      {/* The live region exists from the start so screen readers announce the notice when it appears. */}
      <div role="status">
        {stale && (
          <div className="update-notice">
            <span>A newer version of Weapons of Order is available.</span>
            <button type="button" onClick={() => window.location.reload()}>
              Reload page
            </button>
          </div>
        )}
      </div>
      <p>
        Weapons of Order <span title={pageBuild.builtAt ? `Built ${pageBuild.builtAt}` : undefined}>{buildLabel(pageBuild)}</span>
      </p>
    </footer>
  )
}
