# 007: Azure dev provisioning through step-by-step commands

Task key: AZURE-SETUP
Status: WIP
Updated: 2026-10-16
Milestone: Deployment
Dependencies: BUILD-CI

## Goal

Identify actual subscription/region/quota, then establish an Azure dev app and deployment authentication.

## Acceptance criteria

Use concrete verified commands one dependent step at a time. Record runtime/plan/app settings and costs; prefer F1 and assess B1 only when justified and eligible. No production resources. DB not required for the initial landing deployment.

## Required reading

- [docs/technical/azure-and-deployment.md](../../docs/technical/azure-and-deployment.md)
- [Scope](../scope.md)
- [Prompt protocol](../workflow/implementation-protocol.md)

## Implementation boundary

A separately reviewable planning task. ChatGPT inspects current code and resolves this task's blockers before issuing one pinned Claude prompt. Preserve medieval v1. Split/refine based on actual development; do not implement successors automatically.

## Validation and administration

Run proportionate checks against acceptance criteria, including relevant failure/retry scenarios. Report actual local evidence. ChatGPT checks remote merge/deployment where applicable. Content/images follow admin configuration and cleanup rules; new mechanics remain code changes.

## Result

Provisioned on 2026-10-16 with owner approval at each decision point. Full resource record and limits: [Azure and deployment](../technical/azure-and-deployment.md#provisioned-dev-environment-task-007).

| Item | Value |
|---|---|
| Subscription | "Azure for Students", `7e1071a1-f3bb-4cda-9d0f-8299c1ee0d4d` (the only subscription; Enabled) |
| Tenant | Instituto Politécnico de Bragança, `dba3267e-25b4-4514-9a0d-7544f1d81724` |
| Region | West Europe (`westeurope`), approved by the owner |
| Resource group | `rg-woo-dev` (tags `project=woo`, `environment=dev`) |
| App Service plan | `plan-woo-dev`, Linux, SKU F1 (Free), 1 worker |
| Web app | `app-woo-dev-lfp18`, https://app-woo-dev-lfp18.azurewebsites.net |
| Runtime | `DOTNETCORE\|10.0` (Linux built-in stack, LTS) |
| Site settings | HTTPS only, minimum TLS 1.2, FTP disabled, HTTP/2 on, SCM and FTP basic-auth publishing disabled, Always On off (not available on F1), no app settings, no startup command |
| Deployment identity | user-assigned managed identity `id-woo-dev-deploy`, client ID `5d746d0d-ce9e-4ba2-94c2-5c3e783314fd` |
| Federated credential | `github-woo-dev`: issuer `https://token.actions.githubusercontent.com`, subject `repo:luisfpires18@29492601/woo@1405763514:environment:dev`, audience `api://AzureADTokenExchange` |
| Role | `Website Contributor` on `app-woo-dev-lfp18` only |
| Cost | €0: F1 is free and the managed identity has no charge. No paid resource was created |

Identifiers above are not secrets; there are no passwords, client secrets or publish profiles.

## Commands run (Azure CLI 2.91.0)

Discovery (read-only): `az account show`; `az account list --all`; `az rest --method get --url https://management.azure.com/subscriptions/7e1071a1-f3bb-4cda-9d0f-8299c1ee0d4d?api-version=2022-12-01` (spending limit); `az provider show --namespace Microsoft.Web`; `az group list`; `az appservice plan list`; `az webapp list`; `az resource list`; read-only inspection of the existing staging app, identity and role; `az appservice list-locations --sku F1 --linux-workers-enabled`; `checknameavailability` (POST, no resource created); `az webapp list-runtimes --os linux`; `az role assignment list` for the signed-in user; `gh api repos/luisfpires18/woo` and `.../actions/oidc/customization/sub` (read-only).

Provisioning, each verified before the next:

1. `az group create --name rg-woo-dev --location westeurope --tags project=woo environment=dev` → Succeeded.
2. `az appservice plan create --resource-group rg-woo-dev --name plan-woo-dev --location westeurope --is-linux --sku F1 --tags project=woo environment=dev` → verified `reserved: true`, SKU F1, tier Free, status Ready.
3. `az webapp create --resource-group rg-woo-dev --plan plan-woo-dev --name app-woo-dev-lfp18 --runtime "DOTNETCORE:10.0" --https-only true --basic-auth Disabled --tags project=woo environment=dev` → verified runtime, HTTPS-only, both basic-auth policies `allow: false`, Always On off; `https://` returned 200 (platform placeholder), `http://` returned 301 to HTTPS.
4. `az webapp config set -g rg-woo-dev -n app-woo-dev-lfp18 --ftps-state Disabled --http20-enabled true` → verified.
5. `az identity create --resource-group rg-woo-dev --name id-woo-dev-deploy --location westeurope --tags project=woo environment=dev` → verified.
6. `az identity federated-credential create --resource-group rg-woo-dev --identity-name id-woo-dev-deploy --name github-woo-dev --issuer https://token.actions.githubusercontent.com --subject "repo:luisfpires18@29492601/woo@1405763514:environment:dev" --audiences api://AzureADTokenExchange` → verified (only credential on the identity).
7. `az role assignment create --assignee-object-id b5e12f2d-1530-4656-8a3d-b0f1b8ff3769 --assignee-principal-type ServicePrincipal --role "Website Contributor" --scope /subscriptions/7e1071a1-f3bb-4cda-9d0f-8299c1ee0d4d/resourceGroups/rg-woo-dev/providers/Microsoft.Web/sites/app-woo-dev-lfp18` → verified it is the identity's only role assignment.

## Findings

- Spending limit is `Off` and the quota ID is `PayAsYouGo_2014-09-01`: despite its name, the subscription bills like pay-as-you-go, so a paid resource would be charged without a cap. Nothing paid was created and the spending limit was not changed.
- F1 quota: Linux F1 is listed for 49 regions and two F1 Linux plans already ran in West Europe, so West Europe was proven. Exact App Service quota numbers need the `Microsoft.Quota` provider registered, which was not done (a subscription change, not needed).
- The Microsoft.Web API reports the plan as `F1` / `Free`; the generic ARM resource index shows it as `U13` / `LinuxFree` (an internal Linux free SKU name; the older plans show `F1` there). Both are the free tier. Check the first billing/cost view to confirm €0.
- The subscription also holds unrelated resources that were not touched: `rg-lorex-dev`; `rg-weaponsoforder-staging` (an F1 app `app-woo-staging-le5jrs` for the earlier draft repository `weapons-of-order`, Application Insights, Log Analytics, Communication Services and its own OIDC identity); default resource groups; `queue-simulator-rg`.
- Azure CLI 2.91.0 was installed from Microsoft's official ZIP (`https://azcliprod.blob.core.windows.net/zip/azure-cli-2.91.0-x64.zip`, SHA-256 `c0c49639a96d8c0dfb9973cfe08c97c5eeb69be050102138ede806285e6fcf1f`; bundled `python.exe` has a valid Python Software Foundation signature) to `F:\Workspace\Tools\azure-cli-2.91.0`, user-local, with owner approval. The owner signed in with `az login`.

## Owner steps for DEV-DEPLOY (not done; GitHub settings are the owner's)

- Create the GitHub environment `dev` in `luisfpires18/woo`, with deployment branches limited to `dev`.
- Add environment variables `AZURE_CLIENT_ID` = `5d746d0d-ce9e-4ba2-94c2-5c3e783314fd`, `AZURE_TENANT_ID` = `dba3267e-25b4-4514-9a0d-7544f1d81724`, `AZURE_SUBSCRIPTION_ID` = `7e1071a1-f3bb-4cda-9d0f-8299c1ee0d4d`.
- The deploy job needs `environment: dev` and `permissions: id-token: write` (plus `contents: read`), signs in with `azure/login` using those variables, and must run only after the application build/test checks pass.

## Execution record

Task-spec commit: 0e69cada041971b634e735824544ea0fa60cd7e7 (owner dispatch, 2026-10-16)
Expected dev base: 0e69cada041971b634e735824544ea0fa60cd7e7 (verified equal to local dev and origin/dev, 2026-10-16)
Implementation branch: feat/007-azure-dev-setup (local only; documentation of the provisioning)
Implementation commit: local commit on that branch; SHA given in the implementation report
Review: Pending ChatGPT review
Merged dev commit: None
Deployment: Not applicable (provisioning only; nothing deployed; the app serves the platform placeholder page)
Issues: F1 quota numbers not read (needs `Microsoft.Quota` registration); ARM index shows the plan as `U13`/`LinuxFree`; spending limit off. GitHub environment and variables not created.

## TLDR

Next: ChatGPT review of the provisioning evidence; then DEV-DEPLOY when dispatched.
Done: rg-woo-dev, Linux F1 plan-woo-dev, app-woo-dev-lfp18 (.NET 10) and an app-scoped OIDC deployment identity in West Europe, at €0.
Issues: see execution record.
