# 001: Claude setup and skills health check

Task key: SKILLS
Status: DONE
Updated: 2026-10-15
Milestone: Foundation
Dependencies: None

## Goal

Verify the actual Claude environment; install missing approved skills/plugins from verified sources.

## Acceptance criteria

Record exact sources, versions, invocation and smoke checks; preserve working tools; clearly report unavailable sources. No game scaffold yet.

## Required reading

- [docs/workflow/skills-and-plugins.md](../../docs/workflow/skills-and-plugins.md)
- [Scope](../scope.md)
- [Prompt protocol](../workflow/implementation-protocol.md)

## Implementation boundary

A separately reviewable planning task. ChatGPT inspects current code and resolves this task's blockers before issuing one pinned Claude prompt. Preserve medieval v1. Split/refine based on actual development; do not implement successors automatically.

## Validation and administration

Run proportionate checks against acceptance criteria, including relevant failure/retry scenarios. Report actual local evidence. ChatGPT checks remote merge/deployment where applicable. Content/images follow admin configuration and cleanup rules; new mechanics remain code changes.

## Execution record

Task-spec commit: f8f6f20cb429c8f7615aed089da0c3718de74043
Expected dev base: f8f6f20cb429c8f7615aed089da0c3718de74043 (verified equal to origin/dev at preflight, 2026-10-13)
Implementation branch: feat/001-skills-preflight (local only)
Implementation commit: local commit on that branch; SHA given in the implementation report
Review: Accepted by owner; merged into dev (docs: record Claude environment and skills preflight)
Merged dev commit: e8703aa0c882c642dac8df0ad33cd6921f102ce7
Deployment: Not applicable (environment audit and documentation only)
Issues: see Audit outcome.

## Audit outcome

Evidence and per-item results: [skills inventory audit](../workflow/skills-and-plugins.md#audit-results-2026-10-13).

- Installed: entire official pixijs/pixijs-skills collection (26 skills, main @ 83760c6f53462ca9cecd68055041f5a8c94758ce) via `npx skills`, user scope. Discovery and a WebGL render smoke check passed.
- Verified without changes: Claude Code 2.1.296, Claude Browser MCP, Caveman, Ponytail, graphify CLI, RTK CLI.
- Installed, invocation not exercised: web-design-guidelines, impeccable, taste-skill (with image-to-code-skill), awesome-design-md, humanizer.
- Partial: graphify Claude skill, playwright-cli executable and RTK Claude hook are missing. Not installed because their sources are not owner-approved and the RTK hook changes global config.
- Unresolved, no approved source: aspnet-core-guidance, emilkowalski, huashu-design-guidelines, ui-ux-pro-max, phase-workflow, Understand Anything.
- Observed, left unchanged: Caveman SessionStart text injected twice (plugin plus standalone hook); humanizer user scope at 2.8.2 while its marketplace clone is at 3.0.0; Claude in Chrome has no connected browser.
- LP AI WORK is ChatGPT-side and was not installed.

## Follow-up 2026-10-15

The audit and PixiJS install were done on a different computer. On this computer (user profile `luisf`, Claude Code 2.1.291) the same 26-skill collection was absent and was reinstalled by task 002; see [MCP-SETUP](002-mcp-setup.md) and the [inventory](../workflow/skills-and-plugins.md#mcp-setup-2026-10-15). Other audit items remain as recorded for the original machine.

## Validation

- `python scripts/check_docs.py` and `git diff --check`; results are in the implementation report.
- Temporary smoke fixture was created outside the repository and removed. No game scaffold, app dependencies, artwork or infrastructure.

## TLDR

Next: ChatGPT review; owner supplies approved sources for unresolved items if still wanted.
Done: environment audited, PixiJS collection installed and verified, inventory recorded.
Issues: six names unresolved; three partial integrations; not merged or reviewed.
