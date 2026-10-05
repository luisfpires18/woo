# Asset lifecycle

Updated: 2026-10-05. Status: design documentation; no game implementation yet.

2026-10-05 explicit requirement: every visual entity (units, buildings, runes, weapons, armour, heroes and other relevant content) must have a configurable 2D sprite/image slot in the admin workspace. The slot and fallback support are mandatory; an uploaded image is not required. Use an emoji or symbol until artwork is supplied. Asset changes must work without code edits.

User intends Amazon S3-style storage through Cloudflare (interpreted as Cloudflare R2, exact service to confirm when the account is created). User will create the account later; no provisioning requested now. Store asset references in game data and bytes in the bucket.

Replacement policy: upload and validate the new image, save the entity reference successfully, then delete the previous image if no other entity uses it. Failed replacements preserve the working image and clean up the unsuccessful upload. Entity deletion or clearing its image also removes unreferenced objects. Shared assets require reference checks before deletion. No permanent historical image copies or unused uploads; this supersedes earlier asset-version retention proposals, while configuration/audit metadata may still be retained. Keep every retained object associated with content, with cleanup retries/reconciliation for interrupted operations. Brief overlap during a safe replacement is acceptable; steady-state orphan files are not.

Implementation guidance: use a new object key for each successful replacement so cached artwork updates reliably; expose replacement as one admin action. Store role-specific display geometry (dimensions, anchor and relevant hotspots) and review it when changing artwork. These are implementation recommendations serving the confirmed upload/cleanup requirements, not new gameplay rules.

## Proposed acceptance scenarios

An entity without an image renders its configured fallback. Admin upload produces a valid reference and correct preview. Failed upload leaves the previous image working. Successful replacement removes the previous unreferenced object. Shared files survive until their last reference is removed. Entity deletion clears unused objects. Cleanup retries after transient failures; reconciliation finds abandoned uploads. New object keys prevent stale-cache replacement issues.

## Open implementation details

Supported formats, dimensions, maximum size, image validation, delivery domain/cache policy, whether derivatives are needed, presigned versus backend upload, and reference accounting. User creates the account later. Do not provision services now.

## Cache requirement

Admin replacement refreshes the image independently of app deployments. Use a persisted asset revision/content hash/new key in its URL, keeping unchanged asset URLs stable. See [app/image versioning](versioning-and-cache.md). No timestamp per render or retained unused object histories.

Cloudflare setup is required roadmap work: verify account/CLI, create R2 bucket with commands and configure credentials/delivery/CORS before integrating uploads. See [setup runbook](cloudflare-r2.md).
