## 2026-06-16 - Added Keyboard Focus to Links and Buttons
**Learning:** Found that buttons/links missed keyboard-friendly ':focus-visible' outline making navigation difficult for users reliant on keyboard. Also lacked click feedback.
**Action:** Always add ':focus-visible' (or ':focus') paired with an outline and outline-offset for interactive elements, and consider adding ':active' states to provide tactile visual feedback to click interactions.
