## 2026-06-21 - Added Accessibility Attributes to Navbar Buttons
**Learning:** Found that the main navigation bar's icon-only buttons lacked essential ARIA attributes and focus styles, impacting screen reader and keyboard accessibility.
**Action:** Implemented `aria-label`, `aria-expanded`, and `aria-controls` for context and linkage, and added `focus-visible` styles to ensure keyboard accessibility. This pattern should be standard for all icon-only interactive elements across the platform.
