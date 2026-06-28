## 2024-05-18 - Added ARIA labels and focus styles to icon-only buttons
**Learning:** Icon-only buttons like social media links or utility commands (Command+K) often lack explicit names, reducing screen reader accessibility. Mobile navigation toggles require `aria-expanded` attributes to communicate state changes effectively.
**Action:** Always provide `aria-label` attributes for semantic meaning in icon-only contexts. Add `aria-expanded` to mobile menu toggle triggers. Ensure keyboard interactivity is communicated visually via `focus-visible` styling using consistent focus rings across the design system.
## 2025-02-12 - Remove local scripts and dev logs before submitting PRs
**Learning:** Development artifacts like scripts (`replace_fonts.sh`) and logs (`docs/site/dev.log`) can accidentally be committed when making large sweeping changes across a repository, even when updating frontends.
**Action:** Always clean up temporary files created during visual or codebase refactoring before submitting a PR.
