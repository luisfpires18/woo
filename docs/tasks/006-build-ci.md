# 006: GitHub build and focused test Actions

Task key: BUILD-CI
Status: WIP
Updated: 2026-10-16
Milestone: Foundation
Dependencies: VERSION

## Goal

Add real application build/test checks to the existing documentation CI.

## Acceptance criteria

dev pushes and applicable PRs validate actual code; no Azure credentials needed. Failures block deployment; checks report what actually ran. Preserve working docs/task validation.

## Required reading

- [docs/workflow/development-workflow.md](../../docs/workflow/development-workflow.md)
- [Scope](../scope.md)
- [Prompt protocol](../workflow/implementation-protocol.md)

## Implementation boundary

A separately reviewable planning task. ChatGPT inspects current code and resolves this task's blockers before issuing one pinned Claude prompt. Preserve medieval v1. Split/refine based on actual development; do not implement successors automatically.

## Validation and administration

Run proportionate checks against acceptance criteria, including relevant failure/retry scenarios. Report actual local evidence. ChatGPT checks remote merge/deployment where applicable. Content/images follow admin configuration and cleanup rules; new mechanics remain code changes.

## Result

[`.github/workflows/app-ci.yml`](../../.github/workflows/app-ci.yml), "Application build and tests":

- Triggers: pushes to `dev` and `master`, pull requests targeting them, and manual `workflow_dispatch`, matching the documentation workflow. No path filters, so every change is built.
- Permissions: `contents: read` only; checkout does not persist credentials. No secrets, Azure credentials, artifacts or deployment steps.
- One `ubuntu-latest` job, 20-minute timeout. Concurrency: a pull request run uses the group `app-ci-pr-<number>` with `cancel-in-progress: true`, so a newer push to the same pull request cancels its older run. Every other run (pushes to `dev`/`master`, manual dispatch) uses its own group `app-ci-run-<run_id>`, so it is never cancelled or replaced while pending.
- Toolchain from the repository pins: `actions/setup-dotnet@v6` reads `global.json` (SDK 10.0.402), `actions/setup-node@v7` reads `.node-version` (Node 24.22.0, npm 11.19.1). NuGet and npm caches are keyed on the lock files and only speed up downloads; restores stay locked.
- Steps, each named for what it runs: show toolchain; `npm --prefix web ci`; `dotnet restore Woo.slnx --locked-mode`; `npm --prefix web run build` (typecheck plus standalone build); `dotnet build -c Release`; `dotnet test` (host tests); `node scripts/build.mjs` (the task 005 coordinated build, unchanged); `node scripts/verify-build.mjs`; `node scripts/check-stale-publish.mjs`.
- [`scripts/verify-build.mjs`](../../scripts/verify-build.mjs) runs the published app and fails unless `/api/build` reports `VERSION`, the checked-out commit, a UTC build time and the expected environment, the bundle loaded by the served `index.html` embeds the same commit and build time, entry HTML is `no-cache` and every manifest-listed hashed file is immutable. With `--expect-local` it instead requires the local fallback (empty commit and build time, environment `local`) and a served bundle without any commit. It reads the existing identity; it does not create metadata.
- [`scripts/check-stale-publish.mjs`](../../scripts/check-stale-publish.mjs) is the stale-publish regression: coordinated build A into `artifacts/publish`, standalone web build B, then (after deleting the server's `bin`/`obj`, so the SDK has no record of A) standalone `dotnet publish` B into the same folder. It fails if any A-only hashed file remains, if A's commit or build time appears in any file under the published `wwwroot`, or if `verify-build.mjs --expect-local` fails.
- The documentation workflow (`docs-check.yml`) is unchanged.
- Deployment gating: there is no deployment yet. DEV-DEPLOY must run only after this workflow succeeds for the same commit (for example by deploying from the same workflow after these steps, or by requiring this check), so a failing build blocks deployment.

Commands are in the [README](../../README.md#continuous-integration).

## Execution record

Task-spec commit: dc407cf325170ccee6fea19fab73ce15d42e06e9 (issue [#5](https://github.com/luisfpires18/woo/issues/5) dispatch)
Expected dev base: dc407cf325170ccee6fea19fab73ce15d42e06e9 (verified equal to local dev and origin/dev, 2026-10-16)
Implementation branch: feat/006-build-ci (local only)
Implementation commit: local commit on that branch; SHA given in the implementation report
Review: Pending ChatGPT review
Merged dev commit: None
Deployment: Not applicable (CI only; no deployment exists)
Issues: fixed in passing: a standalone `dotnet publish` into a folder holding a coordinated publish could keep coordinated web files. `PreserveNewest` kept the newer coordinated `index.html` and manifest, and obsolete hashed files were removed only when the SDK's `obj/` still recorded the earlier publish. The server project now deletes the published `wwwroot` before copying the current web build, and copies it `Always`; `check-stale-publish.mjs` fails without that deletion (verified by removing it in a scratch clone). Review fix (2026-10-16): push runs no longer share a concurrency group. The workflow has not run on GitHub yet; it was validated by running the same commands locally in a fresh clone on Windows, not on the Ubuntu runner. No workflow linter is installed and none was added.

## TLDR

Next: ChatGPT review of the local commit; the first real run happens after the owner pushes.
Done: application build/test workflow and published-build verification added; docs workflow preserved.
Issues: see execution record.
