# 008: Deploy the first page automatically from dev

Task key: DEV-DEPLOY
Status: WIP
Updated: 2026-10-16
Milestone: Deployment
Dependencies: AZURE-SETUP

## Goal

Create the dev-only application deployment Action and validate a complete deployed round trip.

## Acceptance criteria

A dev change builds/tests, deploys its exact artifact and serves a healthy versioned page. Deployment history identifies commit/version; failures visible; no prod workflow. Demonstrate recovery without inventing Azure resource identifiers.

## Required reading

- [docs/technical/azure-and-deployment.md](../../docs/technical/azure-and-deployment.md)
- [Scope](../scope.md)
- [Prompt protocol](../workflow/implementation-protocol.md)

## Implementation boundary

A separately reviewable planning task. ChatGPT inspects current code and resolves this task's blockers before issuing one pinned Claude prompt. Preserve medieval v1. Split/refine based on actual development; do not implement successors automatically.

## Validation and administration

Run proportionate checks against acceptance criteria, including relevant failure/retry scenarios. Report actual local evidence. ChatGPT checks remote merge/deployment where applicable. Content/images follow admin configuration and cleanup rules; new mechanics remain code changes.

## Result

Deployment is a second job, `deploy-dev`, in [`app-ci.yml`](../../.github/workflows/app-ci.yml), so it can only run after `build-test` succeeds in the same run and only ships that run's tested artifact.

Chain for a push to `dev` (or a manual run on `dev`):

1. `build-test` builds with `WOO_BUILD_ENVIRONMENT=dev` (every other branch/event stays `local`, so pull requests and `master` never produce a deployable artifact). All BUILD-CI steps are unchanged: locked restores, web typecheck/build, .NET build, host tests, `node scripts/build.mjs --environment dev`, `node scripts/verify-build.mjs --environment dev`, and the stale-publish regression.
2. Immediately after the verification, the verified `artifacts/publish` folder is uploaded as artifact `woo-dev-<commit>` (hidden `.vite/manifest.json` included, 14-day retention). The stale-publish regression runs afterwards and does not touch the uploaded artifact.
3. `deploy-dev` (`needs: build-test`) runs in GitHub environment `dev` with only `contents: read` and `id-token: write`. It checks out the same commit and:
   - skips with a notice, without deploying, if `dev` has already moved to a newer commit (`git ls-remote`), so a slow older run can never overwrite a newer deployment;
   - fails with a named error if any of the `dev` environment variables are missing;
   - downloads the artifact (download-artifact verifies its digest);
   - signs in with `azure/login@v3` using OIDC (`AZURE_CLIENT_ID`, `AZURE_TENANT_ID`, `AZURE_SUBSCRIPTION_ID` variables; no secret, no publish profile);
   - sets the startup command `dotnet Woo.Server.dll` (idempotent, keeps it in code);
   - zips the downloaded folder and runs `az webapp deploy --type zip --clean true --restart true --track-status false` to `app-woo-dev-lfp18` in `rg-woo-dev` (Microsoft Entra authentication, since basic auth is disabled; zip deployment does not rebuild);
   - runs `node scripts/verify-build.mjs --url https://app-woo-dev-lfp18.azurewebsites.net --environment dev --max-wait 300`, which waits with bounded backoff (5, 10, 15, 20, then 30 s; each request limited to 30 s; at most 300 s) until `/api/build` reports this commit, then checks `/api/health`, page `/` status, `/api/build` (`VERSION`, commit, `dev`, UTC build time, `no-store`), that the served bundle embeds the same commit and build time, `no-cache` entry HTML, immutable hashed files, and that every manifest-listed file served is byte-identical to the tested artifact;
   - writes the URL, version, commit and live `/api/build` to the run summary.
4. Concurrency: `build-test` keeps its own per-run group. `deploy-dev` uses one `deploy-dev` group without cancelling a running deployment; a newer waiting deployment replaces an older waiting one.

History: each run is a GitHub deployment to the `dev` environment, tied to the commit; the run summary and `/api/build` show the version and commit; App Service keeps its own deployment log.

Failure and recovery:

- Any failed step (build, tests, verification, missing variables, OIDC login, deploy, post-deploy verification) fails the workflow and blocks or marks the deployment failed; nothing uses `continue-on-error`.
- A failed or superseded deployment leaves the last deployed build running; the post-deploy verification reports which commit is actually served.
- Recover by re-running the failed jobs of the same run (same artifact), by running the workflow manually on `dev` (rebuilds and deploys the current `dev` head), or by reverting the bad commit on `dev` (the revert deploys normally). Older runs cannot be re-deployed over a newer `dev` head because of the head check.
- F1 constraints still apply: cold starts after idle or restart, 60 CPU minutes per day, and the SCM endpoint can return 503 under F1 resource limits; a re-run is the remedy, a paid tier only by explicit owner decision.

No production workflow: `master` and pull requests never deploy.

## Owner setup required before the first deployment

Checked 2026-10-16 (read-only): the repository has no GitHub environments, no environment or repository variables and no repository secrets.

1. Settings > Environments > New environment `dev`.
2. Deployment branches and tags: Selected branches, add `dev` only.
3. Environment variables (not secrets):
   - `AZURE_CLIENT_ID` = `5d746d0d-ce9e-4ba2-94c2-5c3e783314fd`
   - `AZURE_TENANT_ID` = `dba3267e-25b4-4514-9a0d-7544f1d81724`
   - `AZURE_SUBSCRIPTION_ID` = `7e1071a1-f3bb-4cda-9d0f-8299c1ee0d4d`
4. Optional: required reviewers on `dev` would pause every deployment for approval; not needed for dev.

The Azure federated credential only trusts tokens for `repo:luisfpires18@29492601/woo@1405763514:environment:dev`, so the environment name must be exactly `dev`.

## Execution record

Task-spec commit: dbb07cf17be0deffc3620212a5b475aa038d7486 (owner dispatch, 2026-10-16)
Expected dev base: dbb07cf17be0deffc3620212a5b475aa038d7486 (verified equal to local dev and origin/dev, 2026-10-16)
Implementation branch: feat/008-dev-deployment (local only)
Implementation commit: local commit on that branch; SHA given in the implementation report
Review: Pending ChatGPT review
Merged dev commit: None
Deployment: Not yet run. Requires the owner's `dev` environment setup and integration into `dev`; then the first run's deployment and verification must be checked remotely.
Issues: the deploy job cannot run locally. Local validation covered the coordinated `dev` build, the remote verification mode against a locally served copy (cold start, stale bytes, older commit) and the head-check logic, not GitHub OIDC or `az webapp deploy` against Azure.

## TLDR

Next: ChatGPT review; owner creates the `dev` environment and variables; first remote deployment verified after integration.
Done: dev-only deploy job shipping the exact tested artifact via OIDC, with bounded post-deploy verification.
Issues: see execution record.
