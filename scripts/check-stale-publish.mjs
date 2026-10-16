// Regression check: a standalone publish into a folder that holds a coordinated publish must leave nothing of the
// coordinated web build behind.
//   A: node scripts/build.mjs                       -> artifacts/publish (coordinated)
//   B: npm --prefix web run build; dotnet publish    -> the same artifacts/publish (standalone, without obj/ history)
// Then: A-only hashed files are gone, A's commit and build time appear nowhere in the published web root, and
// verify-build.mjs --expect-local passes (local fallback from /api/build and the served page).
// Usage: node scripts/check-stale-publish.mjs   (needs committed tracked files, like build.mjs)
import { spawnSync } from 'node:child_process'
import { existsSync, readFileSync, readdirSync, rmSync } from 'node:fs'
import { join, relative } from 'node:path'
import { fileURLToPath } from 'node:url'

const root = fileURLToPath(new URL('..', import.meta.url))
const webRoot = join(root, 'artifacts', 'publish', 'wwwroot')
const fail = (message) => {
  console.error(`check-stale-publish: ${message}`)
  process.exit(1)
}
// One quoted command line through the shell: npm is a .cmd shim on Windows; no argument has spaces.
const run = (command) => {
  console.log(`> ${command}`)
  const result = spawnSync(command, { cwd: root, shell: true, encoding: 'utf8', stdio: ['ignore', 'pipe', 'inherit'] })
  process.stdout.write(result.stdout)
  if (result.status !== 0) fail(`"${command}" failed`)
  return result.stdout
}
const hashedFiles = () => {
  const manifest = JSON.parse(readFileSync(join(webRoot, '.vite', 'manifest.json'), 'utf8'))
  return Object.values(manifest).flatMap((c) => [c.file, ...(c.css ?? []), ...(c.assets ?? [])])
}
const allFiles = (dir) => readdirSync(dir, { withFileTypes: true, recursive: true }).filter((e) => e.isFile()).map((e) => join(e.parentPath, e.name))

const [, , commitA, builtAtA] = run('node scripts/build.mjs').match(/^build: (\S+) (\S+) (\S+) /m) ?? fail('no build line from build.mjs')
const filesA = hashedFiles()

run('npm --prefix web run build')
// The SDK deletes files from its own previous publish only when obj/ still records that publish. Drop that
// history so the check covers the case it cannot handle (a cleaned project, or another checkout publishing here).
rmSync(join(root, 'src', 'Woo.Server', 'bin'), { recursive: true, force: true })
rmSync(join(root, 'src', 'Woo.Server', 'obj'), { recursive: true, force: true })
run('dotnet publish src/Woo.Server -c Release -o artifacts/publish')
const filesB = hashedFiles()

const failures = []
const expect = (ok, message) => {
  console.log(`${ok ? 'ok  ' : 'FAIL'} ${message}`)
  if (!ok) failures.push(message)
}
const onlyA = filesA.filter((f) => !filesB.includes(f))
expect(onlyA.length > 0, `A had files that B does not (${onlyA.join(', ') || 'none'}), so the check is meaningful`)
for (const file of onlyA) expect(!existsSync(join(webRoot, file)), `A-only ${file} is gone`)
const leaks = allFiles(webRoot).filter((file) => {
  const text = readFileSync(file, 'latin1')
  return text.includes(commitA) || text.includes(builtAtA)
})
expect(leaks.length === 0, `A commit ${commitA} and build time ${builtAtA} appear in no published web file${leaks.length ? `: ${leaks.map((f) => relative(webRoot, f)).join(', ')}` : ''}`)

const verify = spawnSync(process.execPath, [join(root, 'scripts', 'verify-build.mjs'), '--expect-local'], { cwd: root, stdio: 'inherit' })
expect(verify.status === 0, 'verify-build.mjs --expect-local passes')

console.log(failures.length ? `check-stale-publish: ${failures.length} check(s) failed` : 'check-stale-publish: all checks passed')
process.exit(failures.length ? 1 : 0)
