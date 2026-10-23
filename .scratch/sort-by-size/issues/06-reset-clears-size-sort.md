# 06: Reset clears a size sort to time

**What to build:** With size selected, clicking the reset button returns the listing to the normal time-descending order. The type filter, name search, and file-key search clear as they do today. Reset does not write an order-by preference; it only changes the sort in effect for the current listing.

**Blocked by:** 03

**Status:** ready-for-agent

- [ ] Reset while size is active sets the sort to time and the direction to descending
- [ ] The resulting list request carries orderBy time and orderDir desc, with the type filter cleared
- [ ] Reset sends nothing to the order-by preference endpoint
- [ ] Reset while a non-size sort is active behaves as it does today
- [ ] A component test covers the size-to-time reset
