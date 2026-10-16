# 009: Shared UI foundations, themes and kingdom colours

Task key: DESIGN-SYSTEM
Status: WIP
Updated: 2026-10-16
Milestone: UI
Dependencies: DEV-DEPLOY

## Goal

Establish typography, spacing, layout, reusable controls and independent light/dark/kingdom themes.

## Acceptance criteria

Arkazia crimson/black, Sylvara green/gold and Veridor blue/silver accents work with readable surfaces. Forms, navigation, focus and responsive sizing consistent. No final gameplay data needed.

## Required reading

- [docs/design/screens-and-themes.md](../../docs/design/screens-and-themes.md)
- [Scope](../scope.md)
- [Prompt protocol](../workflow/implementation-protocol.md)

## Implementation boundary

A separately reviewable planning task. ChatGPT inspects current code and resolves this task's blockers before issuing one pinned Claude prompt. Preserve medieval v1. Split/refine based on actual development; do not implement successors automatically.

## Validation and administration

Run proportionate checks against acceptance criteria, including relevant failure/retry scenarios. Report actual local evidence. ChatGPT checks remote merge/deployment where applicable. Content/images follow admin configuration and cleanup rules; new mechanics remain code changes.

## Result

Token/theme contract, primitives, accessibility checks and open visual decisions: [design system foundation](../design/design-system.md).

- `web/src/tokens.css`: typography, spacing, sizing, shape, surfaces/text, navigation, brand action, focus, disabled, semantic and elevation tokens for light and dark; kingdom accent tokens for Arkazia, Veridor, Sylvara and the neutral state, scoped by `data-kingdom`.
- `web/src/theme.ts`: `applyAppearance('system' | 'light' | 'dark')` (system by default, follows OS changes), `applyDocumentKingdom`, kingdom ids and names.
- `web/src/ui.tsx`: `Page`, `PageHeader`, `Panel`, `Cluster`, `Button`, `TextField`, `StatusMessage`, `KingdomBadge`.
- The skeleton page and the shared footer/update notice use the primitives; a clearly labelled development preview shows button, field, status and kingdom accent states. No landing page, login, kingdom selection, navigation shell, session or game content.
- Dependencies unchanged: no UI, CSS, icon, font or test packages; system font stacks only.

## Execution record

Task-spec commit: 44d63f440a26f64d99422bce4611730d5844acbf (owner dispatch, 2026-10-16)
Expected dev base: 44d63f440a26f64d99422bce4611730d5844acbf (verified equal to local dev and origin/dev, 2026-10-16)
Implementation branch: feat/009-design-system (local only)
Implementation commit: local commit on that branch; SHA given in the implementation report
Review: Pending ChatGPT review
Merged dev commit: None
Deployment: Deploys through DEV-DEPLOY after integration; not yet run
Issues: contrast and keyboard/form/narrow-layout checks ran in the built-in browser against the dev server; no automated frontend test exists. Exact shades remain owner-adjustable (see open decisions).

## TLDR

Next: ChatGPT review of the local commit.
Done: design tokens, independent appearance and kingdom theming, shared primitives, footer on tokens, accessibility checks.
Issues: see execution record.
