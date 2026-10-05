# Implementation prompting protocol

Updated: 2026-10-05. This is the shared file that brainstorming chats may refine and the orchestration chat must read before every prompt.

## Prompt contract

Prepare exactly one implementation prompt at a time. Populate the following from current evidence:

- Task title, immutable key and detailed task-file path.
- Task-spec commit SHA and current dev base SHA, with verified facts separated from assumptions.
- Feature branch name and suggested commit message.
- Player/product goal and existing behavior.
- Required behavior, invariants, scope boundaries and explicit exclusions.
- Relevant source files/docs and decisions; unresolved blocker means no implementation prompt yet.
- Applicable verified Claude skills/plugins and preflight instructions.
- Proportionate local build/test/browser validation and migration/backup implications.
- Documentation updates, task execution record and final report requirements.
- Git restrictions: local commit only; no push/merge/PR/master edits/deploy unless explicitly authorised.

The owner can iterate a task file in another chat. The orchestration chat reads the latest committed version before prompting. After dispatch, preserve the pinned revision. New requirements get a focused amendment; Claude must not silently chase moving scope.

## Claude preflight

Read AGENTS.md, CLAUDE.md and the named task. Inspect local branch/worktree and verify the expected base. Do not overwrite unrelated changes. Check installed skills/plugins against the inventory, install missing approved sources when available, and report unavailable sources instead of guessing packages. Apply only relevant skills. No automatic installation of every design skill for every task.

## Required implementation report

Task key and current number; specification SHA; base SHA; branch; final local commit SHA; files changed; requirement-by-requirement result; actual commands/check results; skipped checks and why; migration/data effects; UI/browser evidence where applicable; unresolved issues; working-tree status; suggested readiness.

End with TLDR: Next / Done / Issues. This report does not authorise merge.

## Review template

Verdict: APPROVED / NEEDS CHANGES / REFUSAL.
Evidence examined: task revision, report, diff/files and checks actually inspected.
Reasons: concrete requirement outcomes, defects or scope mismatches.
Next action: guarded merge commands for approval, one correction prompt for changes, or scope reset proposal for refusal.
TLDR: Next / Done / Issues.

Remote verification belongs to ChatGPT when available. Local behavior validation belongs to Claude. Do not waste Claude tokens polling Actions.
