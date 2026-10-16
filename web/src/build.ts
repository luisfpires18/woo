export type BuildIdentity = {
  version: string
  commit: string
  builtAt: string
  environment: string
}

declare const __WOO_BUILD__: BuildIdentity

/** Identity this page was built with (injected by vite.config.ts). */
export const pageBuild: BuildIdentity = __WOO_BUILD__

/** Coordinated builds carry a full commit; standalone and dev builds leave it empty. */
export const isCoordinated = (build: BuildIdentity) => /^[0-9a-f]{40}$/.test(build.commit)

export const buildLabel = (build: BuildIdentity) =>
  `${build.version} (${isCoordinated(build) ? build.commit.slice(0, 7) : 'local'})`

/** Server identity, or null when it is unavailable or malformed. */
export async function fetchServerBuild(signal: AbortSignal): Promise<BuildIdentity | null> {
  try {
    const response = await fetch('/api/build', { signal, cache: 'no-store', headers: { Accept: 'application/json' } })
    if (!response.ok) return null
    const body: unknown = await response.json()
    if (typeof body !== 'object' || body === null) return null
    const { version, commit, builtAt, environment } = body as Record<string, unknown>
    if ([version, commit, builtAt, environment].some((v) => typeof v !== 'string')) return null
    return body as BuildIdentity
  } catch {
    return null
  }
}

/** A newer deployment exists only when both sides are coordinated builds with different commits. */
export const isStale = (page: BuildIdentity, server: BuildIdentity | null) =>
  server !== null && isCoordinated(page) && isCoordinated(server) && page.commit !== server.commit
