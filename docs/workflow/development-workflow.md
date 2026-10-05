# Development workflow

Updated: 2026-10-05. Confirmed owner workflow; implementation has not started.

## Source and responsibilities

The owner requested LP AI WORK. The exact Windows file is not accessible here. An older reconstructed LoreX-focused LP WORK package was read as a methodology reference; it is not claimed to be the current project-agnostic skill. The owner's explicit WOO instructions take precedence. Do not copy or install that older ChatGPT-side skill into Claude or this repository. Obtain the current skill before claiming exact compliance.

ChatGPT plans, reads the latest GitHub state, prepares one prompt, reviews the implementation evidence and checks remote history/pipelines when accessible. Claude Code implements locally, validates behavior and returns evidence. The owner transfers prompts/reports and runs Git/Azure commands. Brainstorming chats edit task specifications; the main orchestration chat reads the committed revision.

## Task tracking

Use root [tasks_roadmap.md](../../tasks_roadmap.md) and detailed files in docs/tasks/. Statuses are TODO, WIP and DONE.

- TODO: planned; may still have prerequisites or unresolved decisions.
- WIP: one task handed to Claude, under implementation, awaiting review, follow-up or merge.
- DONE: accepted requirements met, ChatGPT approved, merged into dev and verified remotely; required deployment checked. Documentation-only tasks may state deployment not applicable.

Approval is not itself DONE. Track review, merge and deployment separately in the task detail. One active implementation task at a time. If pipeline remains running, record deployment pending; owner can explicitly proceed to another task without pretending the prior pipeline passed.

Numbers show current order, not permanent identity. Each task has an immutable key such as STACK, while filenames use the current number: 002-tech-stack.md. When inserting a refinement, renumber later rows and detail filenames atomically, update all links and dependencies, and preserve the immutable key. Never rewrite historical commits or rename an active implementation branch merely because its task moved. Keep its exact branch and original prompt in Execution record. Future branches use the current number. See insertion rules in tasks_roadmap.md. Use whole numbers; separately dispatched/reviewed work gets its own number rather than a letter suffix.

## Cycle

1. ChatGPT reads remote dev SHA, AGENTS.md, CLAUDE.md, tasks_roadmap.md, the current task file, the implementation protocol and relevant current code/docs. Resolve blocking requirements. Mark the selected task WIP only when issuing its implementation prompt.
2. Write one prompt using [implementation-protocol.md](implementation-protocol.md). Include exact task-file path, immutable task key, task-spec commit SHA, expected dev base SHA, named feat branch, required behavior, out-of-scope work, skill use and validation.
3. Owner sends it to Claude. Claude performs preflight, implements only that task, validates, commits locally and reports. Default restriction: no push, merge, PR or deployment.
4. Owner pastes report into ChatGPT. Review against the pinned task revision; a later brainstorming edit cannot silently expand an active task. Material changes require an explicit revised prompt.
5. Verdict is APPROVED, NEEDS CHANGES or REFUSAL. APPROVED gives guarded merge-to-dev commands. NEEDS CHANGES gives one focused fix prompt. REFUSAL explains fundamental scope/behavior mismatch and proposes a corrected task; preserve evidence, do not silently delete work.
6. Owner executes merge/push after approval. ChatGPT verifies remote dev and CI/deployment when available, updates execution records and status, then prepares the next task only when requested.

## Evidence policy

Claude's summary is evidence, not proof by itself. Review the changed code/diff when accessible. A local-only branch is invisible on GitHub: request the relevant diff/files/test output from the owner rather than claiming inspection. Missing essential evidence yields NEEDS CHANGES for evidence collection, not automatic approval.

Record exact commit, requirements addressed, changes, tests actually run/results, skipped validation, migrations/rollback implications, UI/browser checks if relevant and known limitations. No repeated full test suites solely for reassurance.

## Git

dev is the working integration branch and the owner's requested default branch. The repository setting remains pending verification; do not claim it changed based on this document. master is reserved for production promotion at alpha release. Feature branches originate at current dev and use feat/NNN-short-name. Fixes for an active task stay on its existing branch unless a concrete reason requires another branch.

Docs maintained by ChatGPT can be committed directly to dev within the owner's ongoing documentation authorisation. Game implementation goes through the review cycle. Do not touch master in normal implementation.

Merge commands must use verified remote base, check clean working tree and expected feature commit, pull with --ff-only and push normally. If dev moved, stop and inspect/revalidate instead of force pushing. Commands are tailored to the known shell and repository path.

## Commands

Always give executable, step-by-step Git/Azure commands with concrete verified values. No unresolved placeholders. If command 2 depends on command 1's output, give only command 1 and wait for pasted output. Independent commands may be batched. Authentication/account selection requires its own handoff. Explain destructive or paid actions briefly. Do not invent account/resource names, quotas, paths or command results.

## Closing report

ChatGPT and Claude end each task/prompt/review response with a short TLDR: what is to be done, what is done and issues encountered. Explicitly say when implementation, merge, deployment or remote verification has not happened.

## Continuous roadmap obligation

Update task roadmap/detail files after meaningful scope refinements, discovered prerequisites, implementation dispatch, review, merge and deployment. Keep task status honest and links coherent. Early runnable/deployed slices take precedence over one large all-game rules task. Current local/dev database is SQLite; production remains deferred.
