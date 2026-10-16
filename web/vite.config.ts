import { readFileSync } from 'node:fs'
import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// Build identity (docs/technical/stack.md). Only scripts/build.mjs sets WOO_COMMIT, WOO_BUILT_AT and
// WOO_ENVIRONMENT, for its own child process; standalone builds and the dev server get the local fallback.
const build = {
  version: readFileSync(new URL('../VERSION', import.meta.url), 'utf8').trim(),
  commit: process.env.WOO_COMMIT ?? '',
  builtAt: process.env.WOO_BUILT_AT ?? '',
  environment: process.env.WOO_ENVIRONMENT || 'local',
}

// The dev server proxies /api to the ASP.NET Core server (launchSettings "http" profile),
// so the browser sees one origin in development as in production.
export default defineConfig({
  plugins: [react()],
  define: {
    __WOO_BUILD__: JSON.stringify(build),
  },
  server: {
    proxy: {
      '/api': 'http://localhost:5080',
    },
  },
})
