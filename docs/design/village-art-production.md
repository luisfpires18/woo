# Village art production runbook
Updated: 2026-10-05. Status: planned iterative workflow, not a guarantee of generated-art consistency.

## Roles and reference
Owner approves artistic direction and handles Git image uploads. ChatGPT helps generate/review assets and maintain specifications; Claude assembles authorised implementation tasks and records local evidence. The requested design skills must be checked for actual availability, not presumed installed.
Do not generate a complete kingdom catalogue before the small slice passes. Approved LF2 fighter style is not automatically the village environment standard.

## Step 1: prepare the style sheet
Use the preferred early village panel as layout/art direction. Mark what is illustrative. Approve camera, palette, lighting, a forge silhouette, material treatment, scale guide and composition footprint diagram. Keep UI independent from scene assets. The recent village mockup is not a source of independently usable buildings.

## Step 2: terrain first
Generate/edit a background without buildings: ground, roads, riverbanks, reserved plot foundations and distant scenery. If water is separate, prepare the exposed river bed/bank treatment and its mask. Avoid remnants from removed structures, strange foundations and scenic objects occupying plots.
Review at actual gameplay size, not just full resolution. Confirm the bridge location and building entrances before approving terrain.

## Step 3: one building and one upgrade
Generate forge against the approved reference with transparent background, required facing, ground anchor and clear margins. Review and repair only the failed feature rather than redesigning everything. Place on terrain before approval. Generate next visual level using the accepted forge as identity/perspective reference; keep entrance and plot alignment while permitting a taller silhouette.
Compare side by side and swap in place. Level art can be reused across multiple numeric levels later; do not require unique artwork for every level without a budget decision.

## Step 4: repeatability test
Produce two other medieval buildings using the same scene/style references. Record failures in angle, scale, roof shape, light and alpha. Add a separate tree and test overlap. Generate modular wall pieces separately, then inspect assembled endpoints. A beautiful isolated wall that fails to join is not accepted.
Test river/bridge contact and gate clearance. Effects are optional after the static scene passes; water animation must not hold up the first proof.

## Prompt templates
Fill approved values before use; these are art briefs, not commands or final asset dimensions.

Terrain brief: "Create an empty medieval village terrain using the approved camera and lighting reference. Preserve the agreed road, river and plot layout. Leave every building footprint clear. No buildings, labels, UI, characters or magic. Output the agreed scene dimensions."

Building brief: "Create one [building/visual level] as a transparent cutout matching reference [approved asset]. Match the camera, ground scale and light direction. Keep the entrance facing [approved facing] and ground anchor [approved location]. Fit the agreed footprint and height envelope. No environment, UI text, extra props or painted transparency grid. Follow the approved shadow convention."

Repair brief: "Change only [identified defect]. Preserve the accepted silhouette, perspective, facing, palette, dimensions and anchor. Inspect the repaired result in the assembled scene."

Wall brief: "Create [piece/orientation] matching approved wall reference. Match the specified thickness, height, lighting and connector endpoints. Transparent background, no terrain or text. This piece must join the approved neighbouring pieces."

## QA and production effort
For each asset record attempts, accepted output, manual repair steps, time and tools. Inspect on light/dark backgrounds and composited scene at normal and minimum zoom. Inspect junctions at high zoom. Check no baked symbols, textual labels or copyrighted/source-restricted material was accidentally imported.
Manual cleanup may be necessary using an image editor: alpha edges, ground contact, alignment and seams. Do not promise a designer-free workflow has zero manual work. Agree tools/budget only if measured repairs require them.
Stop after the slice to assess whether the effort scales to the alpha building set. Estimate from observed acceptance/repair rate, not assumed one-prompt perfection.

## Handoff and approval
Approval covers a named asset/revision and its intended usage. Preview-only generations are not automatically production assets. Owner uploads approved files to the agreed Git folder; documentation may be maintained remotely by ChatGPT. This supersedes earlier assistant-managed image commits.
The production asset reference should identify who may replace it and its geometry contract. Unapproved trial assets stay out of the delivery bucket. See [visual proof](village-visual-prototype.md) and [asset lifecycle](../technical/asset-lifecycle.md).
