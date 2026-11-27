# Project card layout with long CMS content

Project cards use a flex column and keep their links in normal flow. Technology
tags can wrap to any number of lines without links covering them. An auto top
margin keeps link rows aligned when sibling cards stretch in the grid; top
padding preserves a readable gap. Existing colors, typography, content, placement,
optional links and semantic CMS proof fields are unchanged.

The four-card editorial preview exposed the previous absolute-positioned links
overlapping real 12–16 technology labels. A local browser baseline reproduced it.
The permanent compiled proof now temporarily clones a card with long bilingual
copy and 20 tags, verifies at least 16px clearance and containment in both themes
at desktop/mobile widths, and removes the clone before content/accessibility proof.
This exercises the real CSS instead of merely asserting a class or style string.

Validation results are recorded with the matching private CMS editorial review.
The normal public source checks are `pnpm validate` and `pnpm test:cms`. This
change introduces no dependency, client framework, CMS contract or provider write.
It is local source until committed/pushed and admitted by the normal CI/source-pin
rollout. An older candidate cannot attest these changed stylesheet bytes.

Validation on27 November2026: `pnpm validate` and `pnpm test:cms` passed, including
the permanent long-copy/tag regression and full-page themed accessibility proof.
The real four-card EN/ID preview passed eight desktop/mobile light/dark scenarios
with exact text/tags/links, at least16px CTA clearance, successful trial-click hit
targets, no horizontal overflow and no serious/critical project-section axe findings.
The initial geometry assertion measured the padded container edge; it was corrected
to measure the actual clickable anchor, with the old negative-clearance baseline
retained and the full checks rerun. Screenshots hide the fixed navigation only for
clear section captures; the keyboard/theme checks run against the unmodified page.
