// Smoke-checks a coordinated publish (node scripts/build.mjs) by running it:
// - /api/health answers;
// - /api/build reports VERSION, the current Git commit, a UTC build time and the expected environment;
// - the published page bundle embeds the same commit and build time (one shared identity);
// - entry HTML revalidates and every manifest-listed hashed file is immutable.
// Usage: node scripts/verify-build.mjs [--environment local|dev|production]   (default local)
import { execFileSync, spawn } from 'node:child_process'
import { readFileSync, readdirSync } from 'node:fs'
import { join } from 'node:path'
import { fileURLToPath } from 'node:url'

const root = fileURLToPath(new URL('..', import.meta.url))
const publish = join(root, 'artifacts', 'publish')
const args = process.argv.slice(2)
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
  expect(build.commit === commit, `commit is HEAD ${commit}`)
  expect(build.environment === environment, `environment is ${environment}`)
  expect(/^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(\.\d+)?Z$/.test(build.builtAt), 'builtAt is a UTC timestamp')

  const assets = join(publish, 'wwwroot', 'assets')
  const bundle = readdirSync(assets).filter((f) => f.endsWith('.js')).map((f) => readFileSync(join(assets, f), 'utf8')).join('\n')
  expect(bundle.includes(build.commit) && bundle.includes(build.builtAt), 'page bundle embeds the same commit and builtAt')

  const entry = await fetch(`${base}/`)
  expect(entry.headers.get('cache-control') === 'no-cache', 'entry HTML is no-cache')
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
