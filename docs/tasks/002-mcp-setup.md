# 002: Microsoft Learn, Context7 and read-only GitHub integration

Task key: MCP-SETUP
Status: WIP
Updated: 2026-10-15
Milestone: Foundation
Dependencies: SKILLS

## Goal

Give Claude Code current Microsoft/library documentation and read-only GitHub issue access before continuing STACK. Audit existing connections first; preserve working tools and the official PixiJS skills.

## Acceptance criteria

Pinned in GitHub issue [#64](https://github.com/luisfpires18/woo/issues/64), revision 2026-10-15T10:24:43Z: audit before installing; Learn returns ASP.NET Core or EF Core documentation; Context7 resolves a frontend library and returns documentation without selecting WOO dependency versions; GitHub reads issue #2 and dev metadata with a verified read-only configuration; record sources, scope, authentication and limitations without secrets; keep browser tooling and the PixiJS skills; reconcile roadmap and links; run the checker.

## Required reading

- [Skills and plugins inventory](../workflow/skills-and-plugins.md#mcp-setup-2026-10-15)
- [Prompt protocol](../workflow/implementation-protocol.md)
- [Task roadmap](../../tasks_roadmap.md)

## Implementation boundary

Documentation and user-level Claude configuration only. No application scaffold, dependency selection, game code, artwork, cloud provisioning, deployment, pushes, merges, PRs or GitHub issue changes. STACK ([003](003-tech-stack.md)) is paused behind this task.

## Validation and administration

Results and limitations are in the [inventory](../workflow/skills-and-plugins.md#mcp-setup-2026-10-15). ChatGPT reviews the local diff and updates the remote issue status.

## Execution record

Task-spec commit: Issue #64, revision 2026-10-15T10:24:43Z
Expected dev base: e8703aa0c882c642dac8df0ad33cd6921f102ce7 (verified equal to origin/dev, 2026-10-15)
Implementation branch: feat/002-mcp-setup (local only)
Implementation commit: local commit on that branch; SHA given in the implementation report
Review: Pending ChatGPT review
Merged dev commit: None
Deployment: Not applicable (documentation and user-level configuration only)
Issues: see inventory limitations. GitHub issue: [#64](https://github.com/luisfpires18/woo/issues/64).

## TLDR

Next: ChatGPT review.
Done: Learn, Context7 and read-only GitHub installed and smoke-checked; PixiJS skills reinstalled on this computer.
Issues: see inventory limitations.
