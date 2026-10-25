# Filename Search

Status: ready-for-agent

## Problem Statement

The search box on the file manager misses files the user can see. It runs a MySQL fulltext match on the file name with the default parser and a three-character minimum token, so a two-letter word is dropped, a fragment of a word never matches, a Chinese name only matches when its whole character run is typed, and the parser, not the user, decides which characters separate one word from the next. Typing `my holiday` leaves `my-holiday.mp4` hidden.

The same parser decides what a word is, and file names are full of things a word parser does not recognise, such as `IMG_2024`, `S01E02` and `x265`. Each of those is one opaque token today.

A second problem sits beside it. Ordering a listing by size at root shows only the files that live at root. Files live inside directories, sizes belong to files, and directory rows have no size at all, so the only size ranking worth having at root covers every directory.

## Solution

Search matches a fragment of a name instead of a whole word. `vfm` splits the typed keyword into search terms itself, then requires every term to appear somewhere in the file name, in any order, whatever separates the words on either side. `my holiday`, `my-holiday` and `myholiday` all find `my_holiday.mp4`, and `holid` finds `holiday.mp4`.

At root, ordering by size spans every directory. The listing shows files only, largest first, because a directory has no size to rank.

## User Stories

1. As a moon user, I want a fragment of a word to find the file, so that typing `holid` finds `holiday.mp4` without me finishing the word.
2. As a moon user, I want the words I type to match in any order, so that `trip holiday` finds `holiday-trip-2024.mp4`.
3. As a moon user, I want the spaces I type to stand for any separator in the name, so that `my holiday` finds `my-holiday.mp4` and `my_holiday.mp4`.
4. As a moon user, I want a name typed without separators to match, so that `myholiday` finds `my_holiday.mp4`.
5. As a moon user, I want short words to count, so that `my` and `01` narrow the list instead of being ignored.
6. As a moon user, I want a fragment of a Chinese name to find it, so that `剧集` finds `剧集01.mp4`.
7. As a moon user, I want a name with letters and digits joined to match, so that `img2024` finds `IMG_2024_02.jpg` and `s01e02` finds `S01E02.mkv`.
8. As a moon user, I want case not to matter, so that `Holiday` and `holiday` find the same files.
9. As a moon user, I want every term I type to be required, so that a second word narrows the list rather than widening it.
10. As a moon user, I want a literal `%` or `_` in my keyword to stay literal, so that it cannot quietly match something I never typed.
11. As a moon user, I want keyword search inside a directory to stay inside it, so that browsing a folder and searching behaves the way it looks.
12. As a moon user, I want keyword search at root to cover every directory, so that I can find a file without remembering where I put it.
13. As a moon user, I want an empty or separator-only keyword to leave the listing as it was, so that typing a dash does not turn root browsing into a global search.
14. As a moon user, I want ordering by size at root to rank every file I own, so that the biggest space users are on page one without me opening directories one by one.
15. As a moon user, I want ordering by size inside a directory to stay scoped to that directory, so that the ranking matches what the address says I am looking at.
16. As a moon user, I want ordering by size to show files only, so that directories, which have no size, do not sit in a size ranking.
17. As a moon user, I want a keyword and a size order to work together at root, so that I can rank only the matches across every directory.
18. As a moon user, I want the largest file first and the direction control still working, so that I can start from the smallest files when I am cleaning up.
19. As a moon user, I want hidden, deleted and sensitive files to stay out of results exactly as they do now, so that a better search does not widen what I can see.
20. As a moon user, I want clicking a search result under a size order to land on the page that holds the file, so that I do not hunt for it.
21. As a release engineer, I want the now unused fulltext index removed with a recorded schema change, so that the database stops advertising a search path that no longer exists.

## Implementation Decisions

**Matching.** The listing stops using the fulltext match on the file name. The keyword is split into search terms, and each term becomes its own `LIKE` condition on the name, joined with `AND`. Case-insensitivity comes from the column collation, so it costs nothing.

**Term splitting, three tiers.** The first tier splits on every character that is neither a letter nor a digit, which makes separators interchangeable. The second splits a run at a letter-to-digit or digit-to-letter change and at a lower-to-upper change, so `IMG2024` and `MyHoliday` split without any dictionary. The third runs `github.com/m11s-io/wordninja-go` over a remaining Latin run, which splits `myholiday` into `my` and `holiday`. The dependency is MIT licensed, embeds its English word list, and has no transitive dependencies. Chinese text needs no segmentation because a Chinese term is already matched as a substring.

**Splitting is lossless.** Every term produced is a substring of what the user typed, so a split can only add matches and can never hide a file that the unsplit keyword would have found. A wrong split therefore costs a few extra rows, not a missing one.

**Keyword bounds.** Duplicate terms are dropped and the list is capped at a fixed small number. A keyword that produces no terms is treated as absent entirely, so it neither filters nor changes the scope of a root listing.

**No escaping needed.** Every term that survives splitting contains letters and digits only, because the first tier splits on everything else. `%`, `_` and the escape character cannot reach the query, so the keyword cannot widen a match.

**Ordering keeps today's rule.** A keyword with the time order still orders by row id, and the name, custom and size orders keep their own clauses. No relevance score is introduced anywhere.

**Root scope for size.** When the effective sort is size and the request names no directory, the parent filter is not applied, so the listing spans every directory of the requesting user. Size already coerces the type filter to files only, so the global listing is files only as well, ordered by size and then by row id in the requested direction. Inside a directory the size order stays scoped to that directory.

**Position endpoint.** The position calculation already treats an absent parent as no parent filter, which is the same rule the global size listing uses, so search-result jumps and deep links land on the page that listing shows. No change there.

**Search scope.** Keyword scope is unchanged. Inside a directory the keyword narrows that directory. At root it covers every directory of the user.

**Schema.** The fulltext index on the file name column of `file_info` is dropped in the in-flight release, with the current snapshot and the release manifest updated in the same change. The fulltext index on `versioned_file` stays, because versioned-file search keeps its current match.

**No client change.** The client already sends the keyword and an absent parent at root, so the root size behaviour needs no frontend work.

**Cost accepted.** A leading wildcard rules out an index, so each term scans the requesting user's rows, bounded by the existing user index. That is a few tens of thousands of rows at worst and the same order as the fulltext match it replaces.

## Testing Decisions

Good tests assert what the user gets. For this change that means which rows a listing returns, in what order, and which request the client sends, never how the query was built or which private helper ran.

Two seams cover the work. The listing operation itself is the backend seam, and it already runs against a test database. New tests there cover a fragment match, word order, separators on either side, a name typed without separators, a short term, a Chinese name, a letters-and-digits name, every term being required, a keyword containing `%` or `_` failing to widen the result, and that hidden and deleted rows stay out. The same seam covers root scope, so tests assert that a size order with no directory returns files from several directories, that a size order with a directory stays inside it, and that time and name orders at root stay root-level.

The term splitter is a pure function and a new seam. It gets a table test whose cases are the typed keyword and the expected terms, covering the separator, transition and dictionary tiers and their interaction, plus a Chinese run and a run that is already a word.

The existing sort and position tests are the prior art and stay as they are, apart from the cases the root rule changes.

## Out of Scope

- Versioned-file search, which keeps its fulltext match.
- Typo tolerance and relevance ranking.
- An embedded or external search index.
- A stored compacted name column.
- The bookmark and virtual-folder name filters, which keep their substring matches.
- Any change to the search box, its timing, or its state in the URL.
- Ordering a mixed listing, or directories, by size.

## Further Notes

Known limitation. The dictionary tier splits a concatenated run only when every piece is in the bundled word list. A run whose pieces are not in the bundled word list stays whole, so a user who types such a name without its separators will miss it. The separator and transition tiers still catch the common cases.

After this change the file name column has no fulltext index and no fallback search path. The listing with no keyword, the position calculation and every other query are unaffected.

Target release is v0.0.6.
