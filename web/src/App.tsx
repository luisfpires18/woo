import { useCallback, useEffect, useState } from 'react'
import { AppFooter } from './AppFooter'
import { FoundationPreview } from './FoundationPreview'
import { Button, Page, PageHeader, Panel } from './ui'

type Health =
  | { state: 'loading' }
  | { state: 'ok'; status: string }
  | { state: 'error'; reason: string }

const TIMEOUT_MS = 5000

async function fetchHealth(signal: AbortSignal): Promise<string> {
  const response = await fetch('/api/health', { signal, headers: { Accept: 'application/json' } })
  if (!response.ok) throw new Error(`Server answered HTTP ${response.status}.`)
  const body: unknown = await response.json().catch(() => null)
  if (typeof body !== 'object' || body === null || typeof (body as { status?: unknown }).status !== 'string') {
    throw new Error('Server answered with an unexpected health response.')
  }
  return (body as { status: string }).status
}

export function App() {
  const [health, setHealth] = useState<Health>({ state: 'loading' })
  const [attempt, setAttempt] = useState(0)
  const retry = useCallback(() => setAttempt((n) => n + 1), [])

  useEffect(() => {
    const controller = new AbortController()
    const timer = setTimeout(() => controller.abort(new DOMException('timeout', 'TimeoutError')), TIMEOUT_MS)
    setHealth({ state: 'loading' })
    fetchHealth(controller.signal)
      .then((status) => setHealth({ state: 'ok', status }))
      .catch((error: unknown) => {
        if (controller.signal.aborted && controller.signal.reason?.name !== 'TimeoutError') return // unmounted
        const reason =
          controller.signal.reason?.name === 'TimeoutError'
            ? `No answer within ${TIMEOUT_MS / 1000} seconds.`
            : error instanceof TypeError
              ? 'The server could not be reached.'
              : error instanceof Error
                ? error.message
                : 'Unknown error.'
        setHealth({ state: 'error', reason })
      })
      .finally(() => clearTimeout(timer))
    return () => {
      clearTimeout(timer)
      controller.abort()
    }
  }, [attempt])

  return (
    <>
      <main>
        <Page>
          <PageHeader title="Weapons of Order" lead="Development skeleton. No game features are available yet." />

          <Panel title="Server status">
            <div role="status" aria-live="polite" aria-busy={health.state === 'loading'}>
              {health.state === 'loading' && <p>Checking the server…</p>}
              {health.state === 'ok' && (
                <p className="health--ok">
                  Server reachable. Health: <strong>{health.status}</strong>
                </p>
              )}
              {health.state === 'error' && <p className="health--error">Server unavailable. {health.reason}</p>}
            </div>
            {health.state === 'error' && (
              <div>
                <Button onClick={retry}>Try again</Button>
              </div>
            )}
          </Panel>

          <FoundationPreview />
        </Page>
      </main>
      <AppFooter />
    </>
  )
}
