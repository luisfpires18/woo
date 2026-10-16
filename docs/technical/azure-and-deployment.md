# Azure hosting and deployment

Updated: 2026-10-05.

## Confirmed owner direction

Host on Azure. Owner reports an Azure for Students account with its annual credit already used. Prefer an F1 App Service development environment that can sleep when unused. If the initial deployment is blocked by quotas, consider temporary B1, then downgrade to F1 when compatible and ready. This records the desired path; subscription eligibility, available quota and upgrade billing are not verified.

Development deployments follow dev. Production deployments follow master, with production workflow added only at alpha release. No Azure resources or deployment credentials have been provisioned by this documentation task.

## Verified platform constraints

F1 is shared compute with CPU quotas and no Always On. It may unload while idle, and exceeding quotas can stop the app until reset. It is not a guarantee of zero-cost always-running simulation. B1 uses paid dedicated compute; availability and subscription spending limits must be checked before changing tiers. A temporary tier change is not a universal fix for region/subscription quota problems.

References checked 2026-10-05:
- [Plans](https://learn.microsoft.com/en-us/azure/app-service/overview-hosting-plans)
- [Limits](https://learn.microsoft.com/en-us/azure/azure-resource-manager/management/azure-subscription-service-limits)
- [Quotas](https://learn.microsoft.com/en-us/azure/app-service/web-sites-monitor)
- [Scaling](https://learn.microsoft.com/en-us/azure/app-service/manage-scale-up)

## Architecture implication

Persist scheduled actions and due timestamps. For sleeping development, request-triggered catch-up can demonstrate offline completion, but does not execute attacks at their exact deadline while asleep. Choose an independently scheduled processor or always-running tier when punctual shared-world actions matter. Do not silently place a permanent .NET worker in F1 and assume it runs continuously.

SQLite is the confirmed local/dev database. Establish its durable file location and backup behavior in persistence tasks. PostgreSQL hosting is deferred to future production. Asset delivery and any independent scheduler are separate decisions. F1 suitability for chosen runtime/deployment format must be checked during stack selection.

## Actions plan

The repository has documentation/task validation CI and, since task 006 BUILD-CI, application build/test CI (`app-ci.yml`) on dev/master pushes and pull requests; it needs no Azure credentials. Dev deployment Actions will be added in AZURE-SETUP and DEV-DEPLOY once subscription/resources and authentication are known, and must deploy only commits that passed the application build/test checks. This is intentionally not a dummy deployment workflow that claims success without deploying.

Dev deploy should build/test an exact dev revision, authenticate with Azure using a selected mechanism (OIDC preferred subject to support), deploy that artifact to the named dev app and check a health endpoint. Do not deploy PR branches or arbitrary refs. Define environment-level configuration, migration handling, serial deployment and failure reporting. Docs-only changes need not redeploy the app once path filtering is defined.

Production Action stays absent until alpha-release approval. Then add master-only deployment with production configuration and checks; never reuse dev database/assets accidentally.

## Command handoff

First collect actual subscription, region, runtime, resource group/app/plan names and current SKU through read-only commands, one dependent step at a time. Assess quota failure before selecting a remedy. Give concrete commands only after outputs are known. Before B1 explain the paid change and check account eligibility. After validation, downgrade with compatible settings (including disabling Always On when required) and verify the resulting SKU. No paid changes or credential setup are executed now.

## Initial deploy versus database deploy

Deploy the first versioned page before database/account/game implementation. Subsequently verify SQLite survives both Azure restart and redeploy before real player data is introduced. Do not store the DB inside a replacement/readonly deployment artifact.

Persistent App Service storage can be backed by a shared filesystem, depending on platform/layout. SQLite WAL is not generally supported over network filesystems; journal/locking behavior must be checked on the actual configuration, not selected from a generic tutorial. Keep dev on a small, deliberately limited single-instance architecture. If storage compatibility fails, report evidence and adjust the Azure layout with the owner while retaining SQLite dev direction; do not silently switch to PostgreSQL.

References for this gate (checked 2026-10-05): [App Service filesystem](https://learn.microsoft.com/en-us/azure/app-service/operating-system-functionality), [container storage](https://learn.microsoft.com/en-us/azure/app-service/configure-custom-container), [SQLite WAL](https://www.sqlite.org/wal.html) and [network caveats](https://www.sqlite.org/useovernet.html).

## Hosting reassessment, 2026-10-05
See [Azure versus Cloudflare hosting](hosting-comparison.md). With the now-confirmed ASP.NET Core backend, the initial recommendation is one Azure application serving API and React build, plus R2 asset delivery. A Cloudflare static frontend is optional later. Cloudflare Containers can run .NET but needs a distinct durability design for SQLite; it is not a drop-in replacement for the planned dev app. No host change or paid provisioning is implied.

## STACK selection, 2026-10-15
The [stack decision record](stack.md#azure-application-and-local-development) selects .NET 10 (built-in App Service stack), one app serving API and React build, Vite-proxied same-origin local development and a GitHub Actions zip deployment from dev. F1 subscription eligibility and SQLite storage compatibility remain explicit verification gates.

## Provisioned dev environment (task 007)

Created 2026-10-16 in subscription "Azure for Students" (`7e1071a1-f3bb-4cda-9d0f-8299c1ee0d4d`, tenant `dba3267e-25b4-4514-9a0d-7544f1d81724`), region West Europe. Commands and evidence: [task 007](../tasks/007-azure-dev-setup.md).

- `rg-woo-dev`: resource group for WOO dev only.
- `plan-woo-dev`: Linux App Service plan, F1 (Free), 1 worker. Free tier limits apply: shared compute, 60 CPU minutes per day, 1 GB memory, 1 GB storage, no Always On, idle unloading with cold starts, and the app stops when a quota is exhausted until it resets. Nothing runs while the app is unloaded, so timed actions still need the persisted catch-up design.
- `app-woo-dev-lfp18` (https://app-woo-dev-lfp18.azurewebsites.net): `DOTNETCORE|10.0`, HTTPS only, TLS 1.2 minimum, FTP disabled, HTTP/2, basic-auth publishing disabled, no app settings yet. `ASPNETCORE_ENVIRONMENT` is left at its Production default; the build identity's `environment` (`dev`) comes from the coordinated build, not from app settings.
- `id-woo-dev-deploy`: user-assigned managed identity for GitHub Actions OIDC. Its single federated credential trusts only jobs in the `dev` GitHub environment of `luisfpires18/woo` (ID-based subject `repo:luisfpires18@29492601/woo@1405763514:environment:dev`), and its single role is `Website Contributor` on `app-woo-dev-lfp18`. It cannot create resources, change the plan's tier or reach other resource groups. No client secret or publish profile exists.
- Cost: €0. The subscription's spending limit is `Off` (pay-as-you-go behaviour), so any future paid tier (for example B1) would be billed immediately; it requires an explicit owner decision with the expected cost.
- Not provisioned: database, storage for SQLite, Application Insights, production resources, R2 changes, deployment workflow.

## Dev deployment (task 008)

`app-ci.yml` job `deploy-dev` deploys only from `dev`, after `build-test` passes in the same run, using the artifact that run verified (no rebuild). It signs in with `id-woo-dev-deploy` through GitHub OIDC in environment `dev`, sets the startup command `dotnet Woo.Server.dll`, runs `az webapp deploy --type zip --clean true --restart true`, and verifies `/api/health`, the page, `/api/build` and byte-identical files with bounded cold-start backoff. Details, recovery and the owner's GitHub setup: [task 008](../tasks/008-dev-deployment.md).
