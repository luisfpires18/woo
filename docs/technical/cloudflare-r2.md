# Cloudflare and R2 setup runbook

Updated: 2026-10-05. Required roadmap work; no account, bucket or token was created here.

## Task sequence

1. CLOUDFLARE: verify the owner account, R2 availability, exact CLI version and authentication.
2. R2-SETUP: select the actual dev bucket name/options, create through current supported commands and verify it exists.
3. R2-CONFIG: configure scoped S3-compatible upload/delete access, application settings, delivery URLs and CORS where applicable.
4. ASSETS and IMAGE-CACHE: connect admin uploads/replacements/deletion and stable asset revisions.

The owner creates/provides the account when ready. Do not infer access to an account from other projects. Account/service activation may require an interactive dashboard step; identify it honestly rather than invent a command for it.

## Command handoff rules

Commands are supplied when setup is reached, using verified account IDs, bucket names, shell and app/resource names. No unresolved placeholders. If a later command depends on output, give only the current command and wait. Independent verified commands may be grouped.

Read current official Wrangler/API docs before issuing commands. Record exact CLI version and authentication mode. Use owner sign-in or scoped credentials, never request pasted secret access keys in chat. Keep local credentials outside git and Azure server secrets outside browser bundles. Mark provisioning/configuration status with actual command output, not merely a prepared runbook.

## Configuration to establish

Dev bucket and app environment; account ID and S3 endpoint; bucket-scoped permissions; local/Azure settings; delivery base URL; allowed local/dev origins where browser requests require CORS; file types/sizes; cache headers and object revisions; cleanup/retry behavior.

Choose backend-mediated upload versus presigned browser upload explicitly. S3 API endpoints, public development URLs and custom-domain delivery have different roles. Presigned requests use the S3 API endpoint; a custom image domain is not a substitute for that endpoint. CORS controls browser access, not account authorisation.

Use a dev delivery URL initially if suitable; custom domain and Cloudflare cache can be configured as a separately verified step. No production bucket/domain is required now. Direct API calls/CLI are the intended setup path where supported.

## Smoke check and cleanup

Upload one approved temporary fixture through the configured path, fetch it via intended delivery, replace it using a distinct revision and remove all temporary objects. Confirm upload/delete denied without the necessary credentials. Then connect real admin entity assets. Every steady-state retained image must be referenced; no abandoned test files or retired asset histories.

## Official references

Checked 2026-10-05; recheck before live setup.

- [Create buckets](https://developers.cloudflare.com/r2/buckets/create-buckets/)
- [Wrangler R2 commands](https://developers.cloudflare.com/r2/reference/wrangler-commands/)
- [S3 integration](https://developers.cloudflare.com/r2/get-started/s3/)
- [CORS](https://developers.cloudflare.com/r2/buckets/cors/)
- [Public/custom-domain delivery](https://developers.cloudflare.com/r2/buckets/public-buckets/)
- [Presigned URLs](https://developers.cloudflare.com/r2/api/s3/presigned-urls/)

## TLDR

Next: configure account/tool access, create bucket, configure uploads/delivery with concrete step-by-step commands.
Done: documented sequence only.
Issues: account identifiers, permissions, bucket naming and actual CLI/runtime configuration not yet supplied.
