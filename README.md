# Weapons of Order

Updated: 2026-10-05. Status: design documentation; no game implementation yet.

A seasonal browser strategy game set in Bellum. The player is a blacksmith who controls village development and troops, serving one kingdom. Forging should shape military decisions.

The first playable version is simple and medieval. Runes and advanced fantasy systems belong to the evolving game, with release timing still open. The closed alpha is for the owner and friends, with Arkazia, Veridor and Sylvara playable.

Start with [the documentation index](docs/INDEX.md), [first-version scope](docs/scope.md), [open questions](docs/open-questions.md) and [development roadmap](docs/roadmap.md). Coding agents must read [AGENTS.md](AGENTS.md).

This repository is now the primary game-design record. Update the relevant documents and decision log as discussions evolve. Documentation authorisation does not authorise starting game implementation.

## Development orchestration

The working integration branch is **dev**. Read [tasks_roadmap.md](tasks_roadmap.md), [workflow](docs/workflow/development-workflow.md) and [prompt protocol](docs/workflow/implementation-protocol.md). Claude Code implements one reviewed task at a time; ChatGPT prepares prompts and reviews; the owner runs merge/Azure commands. Azure is the selected host direction. Production deployment from master is deferred until alpha release.

GitHub's repository default branch may still be master; explicitly select dev to read the latest workflow. Changing the default branch is an outstanding repository setting, not claimed complete.
