# 03: Document the search semantics and ship in v0.0.6

**What to build:** The service API document describes what the keyword parameter now does: the keyword is split into terms, every term must appear in the name, order and separators do not matter, matching is case-insensitive, and the scope is the current directory or the whole account at root. It also notes that a size order at root spans every directory. The release note for v0.0.6 records the behaviour change and the dropped index.

**Blocked by:** 01, 02, 07

**Status:** ready-for-agent

- [ ] The listing endpoint's keyword parameter is documented with terms, any order, separator insensitivity, case insensitivity and scope
- [ ] The listing endpoint's sort enumeration notes that a size order at root spans every directory and stays files only
- [ ] No service document still describes name search as a fulltext match
- [ ] The v0.0.6 release note records the search change, the root size scope and the dropped fulltext index
