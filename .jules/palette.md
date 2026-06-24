## 2023-10-24 - Accessibility Enhancements in Navbar and Footer
**Learning:** Found multiple instances of icon-only links missing aria-labels and navigation elements lacking visible focus states. Ensure all interactive elements have keyboard focus visibility using Tailwind's `focus-visible` utility.
**Action:** Applied `focus-visible:ring-2` to buttons and links in the layout components to improve keyboard navigation accessibility. Added `aria-expanded` and `aria-controls` to mobile menu toggles.
