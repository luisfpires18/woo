# Azure versus Cloudflare hosting

Updated: 2026-10-05. Status: researched recommendation, no hosting change or provisioning.

## Decision context

Owner confirmed keeping C# / ASP.NET Core and React/TypeScript after the backend comparison. SQLite remains local/dev; PostgreSQL is deferred. Azure remains the previously selected host. This comparison revisits deployment choices without changing that decision.

Recommendation: initially serve the React build and ASP.NET Core API from one Azure App Service application, with Cloudflare R2 for admin-uploaded images. Validate F1 eligibility, runtime support and SQLite persistence before relying on it. This deployment layout is a proposal, not a completed implementation.

## Compare actual products

| Option | Fit for WOO | Main tradeoff |
|---|---|---|
| Azure App Service | Native ASP.NET Core app hosting; can serve React static build too | F1 CPU limits and sleeping; paid tier needed for dependable continuous processing |
| Cloudflare Workers static assets / Pages | Excellent React build delivery | Static hosting does not host our ordinary ASP.NET Core API |
| Cloudflare Workers | Strong edge/serverless runtime with managed services | Conventional ASP.NET Core server is not a standard Workers deployment; adopting its platform would require different backend integration |
| Cloudflare Containers | Can run a .NET container image | Worker/Durable Object orchestration, ephemeral default disk, separate database durability design and usage billing |

Do not claim Cloudflare cannot host .NET. Containers supports arbitrary languages/runtimes. Do not confuse Containers with uploading a React app to Pages.

## Initial Azure layout

One application origin simplifies secure-cookie authentication, deployment, routing and local/dev parity. React/PixiJS still runs in the browser, so choosing Azure does not constrain sprite rendering. Build assets and API can share an app, while uploaded images are served from R2 using the asset lifecycle/cache contract.

F1 is suitable for an initial running page and low-volume development subject to account eligibility. Its allowance is 60 CPU minutes per day, not 60 minutes of wall-clock availability. It has shared compute, 1 GB RAM, 1 GB storage and no Always On. Requests after idle unloading may encounter a cold start. Quota exhaustion is a different failure from ordinary idle sleeping.

Scheduled actions must be persisted with due timestamps, catch-up and idempotent processing. A sleeping app cannot guarantee that a battle executes exactly at its deadline. A dependable friends season needs a separately assessed always-running plan or independent scheduler. Basic and higher support Always On, but processes must still tolerate restarts.

SQLite durability remains a proof gate: correct writable persistent path, filesystem locking/journal compatibility, single instance, restart/redeploy survival and backup/restore. Azure hosting does not automatically make an arbitrary SQLite file safe.

## Optional split frontend

Later, host the React build on Cloudflare Workers static assets or Pages and keep the ASP.NET Core API on Azure. This improves static delivery and keeps the app shell available even if the API is sleeping. It does not remove API cold starts or execute game timers.

The split adds two deployment targets and an authentication/routing design: API URL, cookie attributes, credentialed requests, precise CORS policy, CSRF protection and version compatibility. A same-origin proxy is another design choice with its own routing/caching behavior. Never publicly cache personalized game API responses. Do not add this split before there is a useful reason.

## All-Cloudflare alternative

Cloudflare Containers can run our .NET image and sleep when idle. The standard Container class defaults to ten-minute inactivity sleep, configurable by the application. Host events can also stop instances.

Default disk is ephemeral: restarting after sleep starts with a fresh filesystem from the image. Filesystem snapshots exist for supported scheduling policies, but are point-in-time copies, not a live transactional database guarantee. They are tied to an image and require explicit lifecycle management. Do not treat R2 FUSE mounts as a proven active SQLite database volume.

An all-Cloudflare .NET architecture therefore needs a deliberately validated persistence strategy, container lifecycle/routing and scheduler design. D1 offers SQLite SQL semantics through Cloudflare interfaces; it is not a local SQLite file that EF Core's SQLite provider can simply open. Changing to D1 would be a separate data-access decision.

Containers is credible for a later experiment, particularly with external database hosting. It is less straightforward for our current single-instance SQLite dev plan.

## Costs and limits

- Azure F1: free within quotas and eligibility. B1 price depends on OS, region and subscription; inspect actual account/pricing before proposing a paid change. Paid App Service plan charges do not disappear merely by stopping an app.
- Workers static assets: static requests are free and unlimited under the documented asset model. Worker-first routing or dynamic API work can introduce metered invocations.
- Workers Free: 100,000 requests/day and 10 ms CPU per invocation. Workers Paid starts at USD 5/month with included usage; overages and other products can add cost.
- Containers: requires Workers Paid and adds container compute, provisioned memory/disk while running, and applicable orchestration/network charges. USD 5 is not a complete .NET hosting quote.
- R2 Standard: documented monthly allowance of 10 GB storage, one million Class A operations and ten million Class B operations; egress is free. Beyond allowances storage and operations are billed. R2 works with either backend host.

No estimated total is reliable until instance size, active time, traffic, image operations, database, backup and monitoring needs are known.

## Task implications

Keep STACK, AZURE-SETUP, DEV-DEPLOY and Cloudflare/R2 tasks. Record the backend confirmation in STACK; select exact supported versions there. Keep the initial deployment small and defer a frontend split, Containers and production PostgreSQL. No new implementation task or paid resource is created by this comparison.

## Primary references

Checked 2026-10-05:
- [Azure Windows App Service pricing and F1 limits](https://azure.microsoft.com/en-us/pricing/details/app-service/windows/)
- [Azure Always On guidance](https://learn.microsoft.com/en-us/troubleshoot/azure/app-service/troubleshoot-performance-slow-web-app)
- [Cloudflare Workers languages](https://developers.cloudflare.com/workers/languages/)
- [Workers static assets billing](https://developers.cloudflare.com/workers/static-assets/billing-and-limitations/)
- [Workers pricing](https://developers.cloudflare.com/workers/platform/pricing/)
- [Containers overview](https://developers.cloudflare.com/containers/)
- [Containers pricing](https://developers.cloudflare.com/containers/platform/pricing/)
- [Containers lifecycle and disk FAQ](https://developers.cloudflare.com/containers/faq/)
- [Container snapshots](https://developers.cloudflare.com/containers/guides/snapshots/)
- [D1 overview](https://developers.cloudflare.com/d1/)
- [R2 pricing](https://developers.cloudflare.com/r2/pricing/)

See also [Azure deployment constraints](azure-and-deployment.md), [architecture](architecture.md) and [backend comparison](backend-comparison.md).
