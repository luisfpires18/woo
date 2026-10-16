# App versions and image cache revisions

Updated: 2026-10-05.

## Confirmed requirements

Versioning is a first-class foundation task. Display an application version such as 0.0.1-dev in the shared footer across landing, game and admin pages. Choose exact dependency versions and architecture early. Image replacement must refresh reliably even when browser/CDN caches exist.

## Proposed app version contract

Maintain one build-time release-version source shared by frontend/backend, with environment suffix and commit identity. Start at 0.0.1-dev when the first application build is created; this documentation commit does not claim an application release. Define increments during task VERSION. CI injects/validates metadata so manually hardcoded footer copies cannot diverge.

Hashed JS/CSS assets may be cached long-term; entry HTML/version metadata must permit the client to discover a deployed build. Decide headers and stale-client handling in the version task. Show commit/build details where useful without cluttering the footer. Do not append a random timestamp to every asset request.

## Independent image revision

Each uploaded replacement gets a new immutable object key, content hash or persisted revision in its delivered URL. UI updates the stored reference after successful upload. A single image change does not require bumping the app release version. Unchanged images keep their stable URL/cache benefit.

Only delete the old object after the new reference succeeds and no remaining entity uses it. Cleanup/storage rules still prohibit unused version histories. An app version query parameter alone is insufficient for admin image changes between deployments. Browser/CDN invalidation and reference refresh must be verified with real replacements.

## Acceptance examples

First page and deployed health/build information identify the same release. Landing/game/admin footers agree. A new app build is discoverable by an existing browser session. Replacing one unit image shows the new art without changing app version or clearing every cached image. Failed replacement preserves the old displayed image and leaves no abandoned object.

Exact release-number policy, CDN headers, response cache policy and stale-client prompt remain task-level decisions.

## Initial build identity contract, 2026-10-15
Task 003 STACK defines the [0.0.1-dev build identity contract](stack.md#build-identity-contract-001-dev): a root `VERSION` file, with commit and build time captured once per coordinated build and passed to both the web build and the server publish; local development reports `local` with no invented commit or timestamp. SKELETON provides only the page and health seam; VERSION implements the contract, footer and build-identity response.

## Implementation (task 005)

- `VERSION` (repository root) holds the release version, currently `0.0.1-dev`. Both builds read it; nothing else hardcodes it.
- `node scripts/build.mjs [--environment local|dev|production]` is the coordinated entry point. It refuses uncommitted tracked changes (a commit would not describe the build), reads `VERSION`, captures `git rev-parse HEAD` and one UTC timestamp, then runs the web build with `WOO_COMMIT`, `WOO_BUILT_AT` and `WOO_ENVIRONMENT` in that child process only and publishes the server with the same values as MSBuild properties. The coordinated web build is written to `artifacts/web`, not `web/dist`, and both output folders are emptied first, so a later standalone build or publish cannot reuse coordinated metadata. Generated outputs stay in Git-ignored folders.
- Standalone builds (`npm run build`, `npm run dev`, `dotnet build`, `dotnet run`, plain `dotnet publish`) embed version plus environment `local` and empty commit and build time. MSBuild regenerates the assembly attributes whenever these values change, so a standalone build after a coordinated one does not keep the old commit. The server's informational version is set to `VERSION` and the SDK's automatic `+<git commit>` suffix is disabled, so a standalone DLL contains no commit at all.
- `GET /api/build` returns `{version, commit, builtAt, environment}` with `Cache-Control: no-store`. `/api/health` and the 404 problem response for unknown `/api/...` routes are unchanged.
- Footer (`web/src/AppFooter.tsx`, reusable by later landing/game/admin layouts): `Weapons of Order 0.0.1-dev (abc1234)` for coordinated builds, `Weapons of Order 0.0.1-dev (local)` otherwise.
- Update notice: on load and whenever the tab becomes visible again, the page fetches `/api/build`; each check cancels the previous one, so an older response cannot override a newer one. The notice follows the latest valid coordinated server identity: a different full 40-character commit shows "A newer version of Weapons of Order is available." with a "Reload page" button in a polite status region, and a matching commit clears it. Unavailable, malformed or non-coordinated (local) responses leave the last confirmed state unchanged. A page from a standalone build never shows the notice. No automatic reload or polling.
- Cache headers: the Vite build writes `.vite/manifest.json` (`build.manifest`), and at startup the server loads it from `wwwroot` and marks exactly the files it lists (entry chunks, CSS and imported assets, all content-hashed) `public, max-age=31536000, immutable`. Every other static file, including any unlisted file under `/assets/`, the manifest itself and entry HTML (`/`, client routes, `/index.html`), gets `Cache-Control: no-cache` (store but revalidate). Without a manifest (development, tests without one) nothing is immutable. `/api/build` gets `no-store`. The manifest is publicly readable; it only names files that `index.html` already references.

## Version bump policy

- `VERSION` is the only release version and follows SemVer 2.0. Until the alpha release it stays on `0.0.x-dev`; the `-dev` suffix is removed only for production releases from master.
- Bump the patch number (for example `0.0.1-dev` to `0.0.2-dev`) in the task commit that changes what testers can see or do in the deployed application. The task prompt states whether a bump is expected.
- Minor and major bumps are owner decisions at milestones (for example `0.1.0` for the friends alpha).
- Do not bump for documentation-only, CI-only or test-only changes, for dependency updates that do not change application behaviour, or for image replacements: images refresh through their own revisioned URLs, and the commit already identifies the exact build.
