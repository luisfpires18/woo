// Coordinated build: capture build identity once and give the same values to the web build and the
// server publish (contract: docs/technical/stack.md#build-identity-contract-001-dev).
// Usage: node scripts/build.mjs [--environment local|dev|production]
// Output: artifacts/publish (single-origin app). The web build goes to artifacts/web, never web/dist,
// so standalone builds cannot pick up this build's metadata afterwards.
import { execFileSync, spawnSync } from 'node:child_process'
import { readFileSync, rmSync } from 'node:fs'
import { fileURLToPath } from 'node:url'
import { join } from 'node:path'

const root = fileURLToPath(new URL('..', import.meta.url))
const fail = (message) => {
  console.error(`build: ${message}`)
  process.exit(1)
}

const args = process.argv.slice(2)
const envIndex = args.indexOf('--environment')
const environment = envIndex === -1 ? 'local' : args[envIndex + 1]
if (!['local', 'dev', 'production'].includes(environment)) fail(`unknown environment "${environment}"`)

const version = readFileSync(join(root, 'VERSION'), 'utf8').trim()
if (!/^\d+\.\d+\.\d+(-[0-9A-Za-z.-]+)?$/.test(version)) fail(`VERSION "${version}" is not a release version`)

const git = (...gitArgs) => execFileSync('git', gitArgs, { cwd: root, encoding: 'utf8' }).trim()
// A commit only identifies the build if tracked files match it.
if (git('status', '--porcelain', '--untracked-files=no')) fail('tracked files have uncommitted changes; commit or stash them first')
const commit = git('rev-parse', 'HEAD')
const builtAt = new Date().toISOString()

const webOut = join(root, 'artifacts', 'web')
const publishOut = join(root, 'artifacts', 'publish')
rmSync(webOut, { recursive: true, force: true })
rmSync(publishOut, { recursive: true, force: true })

const check = (result, name) => {
  if (result.status !== 0) fail(`${name} failed`)
}

// npm is a .cmd shim on Windows and needs a shell, so it gets one quoted command line.
const npmCommand = `npm --prefix web run build -- --outDir "${webOut}" --emptyOutDir`
console.log(`> ${npmCommand}`)
check(spawnSync(npmCommand, {
  cwd: root, stdio: 'inherit', shell: true,
  env: { ...process.env, WOO_COMMIT: commit, WOO_BUILT_AT: builtAt, WOO_ENVIRONMENT: environment },
}), 'web build')

const run = (command, commandArgs) => {
  console.log(`> ${command} ${commandArgs.join(' ')}`)
  check(spawnSync(command, commandArgs, { cwd: root, stdio: 'inherit' }), command)
}

run('dotnet', [
  'publish', 'src/Woo.Server', '-c', 'Release', '-o', publishOut,
  `-p:WebDist=${webOut}${process.platform === 'win32' ? '\\' : '/'}`,
  `-p:WooCommit=${commit}`, `-p:WooBuiltAt=${builtAt}`, `-p:WooEnvironment=${environment}`,
])

console.log(`build: ${version} ${commit} ${builtAt} ${environment} -> ${publishOut}`)
