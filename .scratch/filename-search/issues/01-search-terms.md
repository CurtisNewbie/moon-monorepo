# 01: Keyword search matches name fragments

**What to build:** The listing stops running a fulltext match on the file name. It splits the typed keyword into search terms inside `vfm`, then requires every term to appear somewhere in the name as a substring, in any order. Splitting happens in three tiers: at every non-alphanumeric character, at letter-to-digit, digit-to-letter and lower-to-upper changes, and, for a remaining run of Latin letters, through the bundled English word list in `github.com/m11s-io/wordninja-go`. A keyword that yields no terms is treated as absent, so it neither filters nor changes the scope of a root listing. The existing ordering rules and the existing keyword scope stay untouched.

**Blocked by:** none

**Status:** ready-for-agent

- [ ] A fragment of a name matches, so `holid` finds `holiday.mp4`
- [ ] Terms match in any order, so `trip holiday` finds `holiday-trip-2024.mp4` and `2024-holiday-my-trip.mp4` matches `my holiday`
- [ ] Spaces the user types stand for any separator on either side, so `my holiday`, `my-holiday` and `myholiday` all find `my_holiday.mp4`
- [ ] A short term narrows the list instead of being dropped, so `my` and `01` both filter
- [ ] A Chinese fragment matches, so `剧集` finds `剧集01.mp4`
- [ ] A letters-and-digits run matches, so `IMG2024` and `S01E02` find their files
- [ ] Every term is required, so a second term cannot widen the result
- [ ] A keyword containing `%` or `_` cannot widen the result
- [ ] A keyword of separators only leaves the listing unfiltered and leaves a root listing root-scoped
- [ ] Terms are deduplicated and capped
- [ ] Hidden, logically deleted and sensitive rows stay out exactly as before
- [ ] Keyword scope is unchanged, directory-scoped inside a directory and account-wide at root
- [ ] The time order under a keyword still orders by row id, and name, custom and size keep their own clauses
- [ ] The splitter has a table test covering the separator, transition and dictionary tiers, a Chinese run, and a run that is already a word
- [ ] Listing tests at the existing listing seam cover a fragment, word order, both separator directions, a name typed without separators, a short term, a Chinese name, a letters-and-digits name, all terms being required, and `%` and `_` staying literal
- [ ] The new dependency is declared in the service module with no transitive additions
