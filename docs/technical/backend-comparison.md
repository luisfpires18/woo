# Backend technology comparison
Updated: 2026-10-05. Status: researched recommendation, not a selected backend or implementation authorisation.

## Project workload
React/TypeScript is the owner-preferred frontend. Village/map rendering runs in the browser; backend language does not determine scene artwork or sprite quality.
WOO needs account/owner authorisation, editable content, World/Kingdom isolation, village resources, durable construction/recruitment/travel/forging jobs, transactional spending, automatic battle calculation, reports and R2 image lifecycle. This is asynchronous strategy; no continuous real-time physics simulation is currently scoped.

## Comparison
These tradeoffs are project-specific engineering judgements, not benchmark rankings.

| Candidate | Advantages for WOO | Costs / qualifications | Recommendation |
|---|---|---|---|
| C# / ASP.NET Core | Owner expertise; Identity account management; EF SQLite provider; hosted-service infrastructure; direct Azure stack support; typed domain modelling | Separate frontend/backend languages; EF/provider behaviour needs review; custom React admin still required; easy to overbuild abstractions | First recommendation for this owner/project |
| Java / Spring Boot | Strong typed domain and transaction ecosystem; Spring Security and scheduling; direct Azure stack support | Different review/maintenance stack; conventional Spring runtime startup/memory must be measured; Hibernate SQLite uses community dialect rather than core support | Fully viable; choose for Java expertise/preference, not an assumed game limitation in C# |
| Go | Compiled executable deployment, concurrency primitives, explicit service code, SQL transactions and SQLite/Postgres drivers | More choices/integration for accounts, validation, migrations and jobs; no single framework's full account-management stack; chosen SQLite driver may change build/deploy dependencies; Azure deployment route must be verified | Good if compact service/runtime footprint is the priority and integration work is accepted |
| TypeScript / Node.js | Same source language as React; potential shared DTO/schema tooling; strong fit for asynchronous HTTP APIs | Runtime validation still required; CPU-heavy simulation must not block API event loop; worker/process separation if measurements justify it; ecosystem choices for ORM/auth/jobs | Best alternative for one-language development |
| Python / Django | Built-in model-based internal admin is valuable for rapid catalogue editing | Custom scene editor/player UI still needed; owner maintenance learning; simulation workload and typing practices need evaluation | Strong if internal admin speed becomes the overriding priority |
| Rust | Strong compile-time safety and systems control | Learning/integration cost for this owner; benefits not shown necessary by current workload | Do not choose solely for speculative future performance |

Go/Java/.NET/Node can all implement the current workload. No player-capacity, throughput, memory or hosting-cost figure has been measured; do not invent one.

## Why .NET was initially recommended
The owner is an experienced C#/.NET engineer and must review AI-produced implementations and debug production behaviour. Familiarity is part of maintainability, not the only consideration.
ASP.NET Core Identity handles users/passwords/roles/claims/tokens; the React account screens and game-specific membership/permissions remain our implementation.
EF Core has a Microsoft SQLite provider, but SQLite migration/type/concurrency limitations still matter. Future PostgreSQL is not a connection-string-only migration; provider-specific migrations and behaviour need a dedicated future task.
Typed domain logic is useful for forging, counter rules, village ownership and replayable calculations. Java offers comparable strengths. Neither language requires maximal DDD or microservices.

## Scheduling and hosting dominate the early decision
Persist due times and sufficient command/state data, then process jobs transactionally and idempotently. Simultaneous spending must not race; a retry/restart must not complete an action twice. Resource catch-up can compute elapsed production; cross-dependent battle/travel events need deterministic chronological processing.
A BackgroundService, goroutine, Spring scheduler or Node timer does not become durable merely by existing. Process-local timers alone cannot survive restart or sleep.
Azure App Service offers built-in .NET, Java, Node, Python and PHP stacks; Go requires verifying an alternative deployment route rather than assuming an equivalent built-in stack.
Always On is unavailable on the planned Free tier and idle apps can unload. While sleeping, exact-time events cannot be guaranteed by an in-process scheduler. Either accept catch-up for dev or explicitly authorise an always-running host/independent durable scheduling approach when exact timing is required.
SQLite serialises writes regardless of language. Keep transactions short, design retry/idempotency, validate actual Azure file persistence and locking, and keep one application instance initially. Do not assume goroutines/threads increase SQLite write concurrency.

## Recommended initial architecture, still a proposal
React/TypeScript/Vite, proposed PixiJS scene renderer, ASP.NET Core modular monolith with ordinary modules, Identity with same-origin secure cookie configuration, EF Core/SQLite for local/dev, R2 behind a storage boundary.
One deployable application initially. Keep job-processing logic separable; run in-process when hosting allows and reconcile persisted due work when necessary. Do not provision a second always-running worker merely because a previous architecture table proposed it.
Use REST and straightforward refresh/polling initially; SignalR is optional later if player experience requires push. Generate frontend API contracts from an explicit schema rather than manually duplicating them. No Redis, broker, Kubernetes or PostgreSQL until a task justifies it.
Final supported versions, Azure OS/runtime route, SQLite persistence proof and authentication details are chosen in STACK and subsequent tasks.

## Selection validation
Before finalising STACK, record owner preference and supported versions. Validate the leading candidate with the early skeleton/deployment and later SQLite/Identity/durable-order tasks. Representative evidence: deploy/restart, data survival, two simultaneous spending requests, duplicate-job retry and catch-up after downtime. Do not build six full comparative backends.
Prefer Go if measured resource/deployment constraints and team expertise outweigh integrated feature delivery. Prefer Java if the owner/team prefers maintaining Spring. Prefer Node if frontend/backend language unity has greater value than the .NET maintenance advantage. Revisiting the backend now is cheap; no game code exists.

## Official sources reviewed
Checked 2026-10-05. Sources support framework/platform facts; recommendation and tradeoffs are engineering judgements.
- [ASP.NET Core Identity](https://learn.microsoft.com/en-us/aspnet/core/security/authentication/identity).
- [EF SQLite provider](https://learn.microsoft.com/en-us/ef/core/providers/sqlite/) and [limitations](https://learn.microsoft.com/en-us/ef/core/providers/sqlite/limitations).
- [EF migrations across providers](https://learn.microsoft.com/en-us/ef/core/managing-schemas/migrations/providers).
- [Hosted services](https://learn.microsoft.com/en-us/aspnet/core/fundamentals/host/hosted-services).
- [Spring Security](https://spring.io/projects/spring-security/) and [scheduling](https://docs.spring.io/spring-boot/reference/features/task-execution-and-scheduling.html).
- [Hibernate dialect support](https://docs.hibernate.org/stable/orm/dialect/).
- [Spring native-image benefits and limitations](https://docs.spring.io/spring-boot/reference/packaging/native-image/introducing-graalvm-native-images.html).
- [Go database access](https://go.dev/doc/database/) and [FAQ](https://go.dev/doc/faq).
- [Node event-loop considerations](https://nodejs.org/en/learn/asynchronous-work/dont-block-the-event-loop) and [worker threads](https://nodejs.org/api/worker_threads.html).
- [Django admin](https://docs.djangoproject.com/en/stable/ref/contrib/admin/).
- [Azure stacks](https://learn.microsoft.com/en-us/azure/app-service/overview), [configuration](https://learn.microsoft.com/en-us/azure/app-service/configure-common) and [Always On tier](https://learn.microsoft.com/en-us/troubleshoot/azure/app-service/troubleshoot-performance-slow-web-app).
- [SQLite isolation](https://www.sqlite.org/isolation.html).
