# Claude Code instructions

Read [AGENTS.md](AGENTS.md), [task roadmap](tasks_roadmap.md), [development workflow](docs/workflow/development-workflow.md), [prompt protocol](docs/workflow/implementation-protocol.md), [skills inventory](docs/workflow/skills-and-plugins.md) and the task file supplied by the owner.

Implement exactly one pinned task on feat/NNN-short-name based on the specified dev commit. Work locally and commit; do not push, merge, create PRs, touch master or deploy by default. Verify installed skills/plugins first, using approved sources for missing dependencies. Apply relevant skills rather than invoking the whole inventory.

Keep task WIP until ChatGPT review and remote merge/deployment verification. Report actual changes, requirement coverage, checks, limitations, base/spec/commit SHAs and Git state. End with TLDR: Next / Done / Issues.

LP WORK is a ChatGPT-side methodology, not a Claude installation requirement. No game implementation has yet been authorised. Azure credentials and resource names will be supplied through later step-by-step setup.

Current constraints: SQLite local/dev; PostgreSQL deferred. Exact dependency versions and architecture before scaffold. Version footer and image revisioning required. Follow early-page/deploy/UI then persisted accounts/world/membership roadmap; don't block first page on all future gameplay decisions. Keep mock UI honest. Update task details/queue during the cycle.
