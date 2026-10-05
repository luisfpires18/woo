# 00: Landing page

Status: Approved visual reference, with mandatory correction below.
Approved by owner: 2026-10-05.
Implementation task: LANDING (currently 009). Approval of artwork does not mark the implementation task DONE.

## Files and states

- [00-landing-page.png](00-landing-page.png): primary light-theme reference, logged in; generic avatar and example nickname Unreally, profile dropdown open.
- [00-landing-page-dark.png](00-landing-page-dark.png): companion dark-theme reference, logged out; navbar Log in button.

Both images are 1448 × 1086. Theme and authentication are independent: support logged-in dark and logged-out light too. These two screenshots illustrate different states, not a rule coupling login to theme.

## Mandatory correction: Settings is a separate page

The light screenshot includes Appearance controls inside an expanded Settings dropdown. DO NOT implement those controls there.

The profile dropdown contains only:
1. Profile
2. Settings
3. Log out

Clicking Settings navigates to a dedicated Settings page. Appearance choices live on that page, not in the navbar or dropdown. Light/Dark/System in the screenshot are provisional choice labels; the required confirmed appearance support is light/dark. The standalone Settings page will receive its own design/specification.

This written correction takes precedence over the screenshot. No regeneration was requested for this minor change.

## Approved landing structure

Slim navbar: anvil + WOO at left; logged out shows Log in at right; logged in shows avatar + nickname + dropdown at right. No Overview, Worlds or standalone theme controls.

The owner's forge/cavern/lava artwork with embedded Weapons of Order title is directly below the navbar and above content. The original New Game, Continue and Exit menu text must not appear. Use the same banner artwork in both themes. Final clean banner asset must be supplied/extracted as a separately approved asset; a screenshot of the whole page is not a background asset.

Below banner: brief game introduction, three feature labels, world/player metrics, then Choose a world with server cards. No email/password/login form on the landing page. Log in opens an independent authentication page/layout.

World cards contain availability, name, pace, player/online counts, kingdom availability and selection action. Selected world should survive the authentication flow where appropriate; concrete join/session behavior belongs to account/membership specifications.

Footer displays the shared app version, illustrated as v0.0.1-dev. Build it from actual app metadata, not a literal screenshot string.

## Data and art limitations

Friends Alpha I/II, 2 worlds, 48 registered players, 12 online, and all per-world numbers are illustrative preview data. They do not establish real servers, balanced stats or a second alpha commitment. Use explicit mocks at the UI task stage and actual world-scoped data once APIs exist. Do not ship misleading hardcoded activity counts.

Unreally and the portrait demonstrate the profile layout; display the signed-in account's actual data. Keep the same anvil brand icon in both themes.

Banner runestones are supplied visual branding, not an instruction to implement runeforging in the medieval first version. Text on ordinary UI controls is rendered as real HTML; do not flatten the page into one clickable image.

## Claude image-to-code handoff

Read this file and both screenshots before the landing task. Use the relevant installed image-to-code/design/accessibility skills. Recreate spacing, hierarchy, typography, bordered cards and restrained crimson actions as reusable React components. Respect written corrections over pixels.

Create responsive equivalents rather than force the desktop screenshot onto mobile: stack server cards, keep title/banner legible, retain readable metrics and touch targets. Implement accessible account menu behavior, focus handling and keyboard navigation. Exact mobile composition is not yet an approved screenshot.

Separate shared navbar/footer, banner, introduction, metrics and server cards. Light/dark UI changes surfaces/text, not banner identity. Login and Settings are separate routes/layouts; no inline landing authentication or dropdown appearance editor.

## Approval and file policy

Only owner-approved images enter mockup/. Use NN-screen-name.png plus a matching NN-screen-name.md, with descriptive theme/state suffixes for companions. Future mockups require approval before Git upload. Preserve an existing approved reference unless replacement is explicitly authorised.

## TLDR

Next: use these references during the separately authorised landing UI task.
Done: light and dark mockups approved and documented.
Issues: remove Appearance from dropdown; standalone Login/Settings and mobile mockups are not supplied yet; banner source asset still needs its own clean final handoff.
