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
