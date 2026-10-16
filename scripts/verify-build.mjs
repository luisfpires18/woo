// Smoke-checks the app published in artifacts/publish, either by running it locally or, with --url, against a
// deployment of that same folder:
// - /api/health answers;
// - coordinated publish (default): /api/build reports VERSION, the current Git commit, a UTC build time and the
//   expected environment, and the bundle the served page loads embeds the same commit and build time;
// - standalone publish (--expect-local): /api/build is the local fallback and the served bundle carries no commit;
// - entry HTML revalidates and every manifest-listed hashed file is immutable;
// - with --url: every manifest-listed file served by the deployment is byte-identical to artifacts/publish.
// Usage: node scripts/verify-build.mjs [--environment local|dev|production | --expect-local]
//                                      [--url https://host [--max-wait seconds]]
// --url waits (bounded backoff, default 300 s) until the deployment answers with the expected commit, so a cold
// start or a restart still serving the previous build is not treated as a failure until the wait runs out.
import { execFileSync, spawn } from 'node:child_process'
import { readFileSync } from 'node:fs'
import { join } from 'node:path'
import { fileURLToPath } from 'node:url'

const root = fileURLToPath(new URL('..', import.meta.url))
const publish = join(root, 'artifacts', 'publish')
const args = process.argv.slice(2)
const option = (name, fallback) => (args.includes(name) ? args[args.indexOf(name) + 1] : fallback)
const expectLocal = args.includes('--expect-local')
const environment = option('--environment', 'local')
const remote = option('--url', null)?.replace(/\/+$/, '') ?? null
const maxWaitMs = Number(option('--max-wait', '300')) * 1000
const version = readFileSync(join(root, 'VERSION'), 'utf8').trim()
const commit = execFileSync('git', ['rev-parse', 'HEAD'], { cwd: root, encoding: 'utf8' }).trim()
const base = remote ?? 'http://127.0.0.1:5099'

const failures = []
const expect = (ok, message) => {
  console.log(`${ok ? 'ok  ' : 'FAIL'} ${message}`)
  if (!ok) failures.push(message)
}
const sleep = (ms) => new Promise((r) => setTimeout(r, ms))
// Each request is bounded too, so a hanging cold start cannot stall the check forever.
const get = (path) => fetch(`${base}${path}`, { signal: AbortSignal.timeout(30_000), cache: 'no-store' })

// Waits until /api/build answers (locally: health only) and, remotely, reports the expected commit.
async function waitForBuild() {
  const started = Date.now()
  const delays = remote ? [5, 10, 15, 20, 30] : [1]
  for (let attempt = 1; ; attempt++) {
    let seen = 'no answer'
    try {
      const response = await get('/api/build')
      if (response.ok) {
        const build = await response.json()
        seen = `commit ${build.commit || '(none)'}`
        if (!remote || expectLocal || build.commit === commit) return true
      } else seen = `HTTP ${response.status}`
    } catch (error) {
      seen = error.name
    }
    const elapsed = Date.now() - started
    const limit = remote ? maxWaitMs : 30_000
    if (elapsed >= limit) {
      console.log(`     gave up after ${Math.round(elapsed / 1000)} s (last: ${seen})`)
      return false
    }
    const delay = Math.min(delays[Math.min(attempt - 1, delays.length - 1)] * 1000, limit - elapsed)
    if (remote) console.log(`     attempt ${attempt}: ${seen}; retrying in ${Math.round(delay / 1000)} s`)
    await sleep(delay)
  }
}

const server = remote ? null : spawn('dotnet', ['Woo.Server.dll'], { cwd: publish, env: { ...process.env, ASPNETCORE_URLS: base }, stdio: 'inherit' })
try {
  const ready = await waitForBuild()
  expect(ready, remote ? `${base} serves commit ${commit} within ${maxWaitMs / 1000} s` : '/api/build answers within 30 s')
  if (!ready) throw new Error('application not ready')

  const health = await get('/api/health')
  expect(health.ok, '/api/health answers')

  const response = await get('/api/build')
  const build = await response.json()
  console.log(`     /api/build: ${JSON.stringify(build)}`)
  expect(response.headers.get('cache-control') === 'no-store', '/api/build is no-store')
  expect(build.version === version, `version is ${version}`)

  // The bundle the active page actually loads, as referenced by the served index.html.
  const entry = await get('/')
  expect(entry.ok, `page / answers HTTP ${entry.status}`)
  expect(entry.headers.get('cache-control') === 'no-cache', 'entry HTML is no-cache')
  const scripts = [...(await entry.text()).matchAll(/<script[^>]+src="([^"]+)"/g)].map((m) => m[1])
  expect(scripts.length > 0, `served page loads ${scripts.join(', ') || 'no script'}`)
  let bundle = ''
  for (const src of scripts) {
    const js = await get(src)
    expect(js.ok, `${src} is served`)
    bundle += await js.text()
  }

  if (expectLocal) {
    expect(build.commit === '' && build.builtAt === '' && build.environment === 'local', '/api/build is the local fallback')
    expect(bundle.includes(version) && !/[0-9a-f]{40}/.test(bundle), 'served page bundle has the version and no commit')
  } else {
    expect(build.commit === commit, `commit is HEAD ${commit}`)
    expect(build.environment === environment, `environment is ${environment}`)
    expect(/^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(\.\d+)?Z$/.test(build.builtAt), 'builtAt is a UTC timestamp')
    expect(bundle.includes(build.commit) && bundle.includes(build.builtAt), 'served page bundle embeds the same commit and builtAt')
  }

  const manifest = JSON.parse(readFileSync(join(publish, 'wwwroot', '.vite', 'manifest.json'), 'utf8'))
  const hashed = Object.values(manifest).flatMap((c) => [c.file, ...(c.css ?? []), ...(c.assets ?? [])])
  for (const file of hashed) {
    const asset = await get(`/${file}`)
    expect(asset.ok && asset.headers.get('cache-control') === 'public, max-age=31536000, immutable', `${file} is immutable`)
    if (remote) {
      const served = Buffer.from(await asset.arrayBuffer())
      expect(served.equals(readFileSync(join(publish, 'wwwroot', file))), `${file} is byte-identical to the tested artifact`)
    }
  }
} catch (error) {
  expect(false, String(error))
} finally {
  server?.kill()
}

console.log(failures.length ? `verify-build: ${failures.length} check(s) failed` : 'verify-build: all checks passed')
process.exit(failures.length ? 1 : 0)
