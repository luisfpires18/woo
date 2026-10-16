import { useEffect, useState } from 'react'
import { buildLabel, confirmStale, fetchServerBuild, pageBuild } from './build'

/** Shared footer for every layout: release identity plus a notice when a newer build is deployed. */
export function AppFooter() {
  const [stale, setStale] = useState(false)

  // Check on load and whenever the tab becomes visible again; no polling. Each check cancels the one
  // before it, so an older response can never override a newer one. Only a valid coordinated server
  // identity changes the notice; unavailable or unconfirmable responses keep the last confirmed state.
  useEffect(() => {
    let current: AbortController | undefined
    const check = () => {
      if (document.visibilityState !== 'visible') return
      current?.abort()
      const controller = (current = new AbortController())
      fetchServerBuild(controller.signal).then((server) => {
        if (controller.signal.aborted) return
        const confirmed = confirmStale(pageBuild, server)
        if (confirmed !== undefined) setStale(confirmed)
      })
    }
    check()
    document.addEventListener('visibilitychange', check)
    return () => {
      document.removeEventListener('visibilitychange', check)
      current?.abort()
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
