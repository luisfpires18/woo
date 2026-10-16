// Smoke-checks the app published in artifacts/publish by running it:
// - /api/health answers;
// - coordinated publish (default): /api/build reports VERSION, the current Git commit, a UTC build time and the
//   expected environment, and the bundle the served page loads embeds the same commit and build time;
// - standalone publish (--expect-local): /api/build is the local fallback and the served bundle carries no commit;
// - entry HTML revalidates and every manifest-listed hashed file is immutable.
// Usage: node scripts/verify-build.mjs [--environment local|dev|production | --expect-local]
import { execFileSync, spawn } from 'node:child_process'
import { readFileSync } from 'node:fs'
import { join } from 'node:path'
import { fileURLToPath } from 'node:url'

const root = fileURLToPath(new URL('..', import.meta.url))
const publish = join(root, 'artifacts', 'publish')
const args = process.argv.slice(2)
const expectLocal = args.includes('--expect-local')
const envIndex = args.indexOf('--environment')
const environment = envIndex === -1 ? 'local' : args[envIndex + 1]
const version = readFileSync(join(root, 'VERSION'), 'utf8').trim()
const commit = execFileSync('git', ['rev-parse', 'HEAD'], { cwd: root, encoding: 'utf8' }).trim()
const base = 'http://127.0.0.1:5099'

const failures = []
const expect = (ok, message) => {
  console.log(`${ok ? 'ok  ' : 'FAIL'} ${message}`)
  if (!ok) failures.push(message)
}

const server = spawn('dotnet', ['Woo.Server.dll'], { cwd: publish, env: { ...process.env, ASPNETCORE_URLS: base }, stdio: 'inherit' })
try {
  let healthy = false
  for (let i = 0; i < 30 && !healthy; i++) {
    healthy = await fetch(`${base}/api/health`).then((r) => r.ok, () => false)
    if (!healthy) await new Promise((r) => setTimeout(r, 1000))
  }
  expect(healthy, '/api/health answers within 30 s')
  if (!healthy) throw new Error('server did not start')

  const response = await fetch(`${base}/api/build`)
  const build = await response.json()
  console.log(`     /api/build: ${JSON.stringify(build)}`)
  expect(response.headers.get('cache-control') === 'no-store', '/api/build is no-store')
  expect(build.version === version, `version is ${version}`)

  // The bundle the active page actually loads, as referenced by the served index.html.
  const entry = await fetch(`${base}/`)
  expect(entry.headers.get('cache-control') === 'no-cache', 'entry HTML is no-cache')
  const scripts = [...(await entry.text()).matchAll(/<script[^>]+src="([^"]+)"/g)].map((m) => m[1])
  expect(scripts.length > 0, `served page loads ${scripts.join(', ') || 'no script'}`)
  let bundle = ''
  for (const src of scripts) {
    const js = await fetch(`${base}${src}`)
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
    const asset = await fetch(`${base}/${file}`)
    expect(asset.ok && asset.headers.get('cache-control') === 'public, max-age=31536000, immutable', `${file} is immutable`)
  }
} catch (error) {
  expect(false, String(error))
} finally {
  server.kill()
}

console.log(failures.length ? `verify-build: ${failures.length} check(s) failed` : 'verify-build: all checks passed')
process.exit(failures.length ? 1 : 0)
