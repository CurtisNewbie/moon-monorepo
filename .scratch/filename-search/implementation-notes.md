# Implementation notes

Where the code stands today, the facts that shaped the spec, and the bookkeeping the tickets depend on. The spec owns the behaviour; this file owns the anchors, and it goes stale as soon as the work lands.

## Where the code is

- The listing SQL is `listFilesSelective` in `backend/vfm/internal/vfm/core.go`. It holds the fulltext match to replace, the parent and type filters, and the order clauses. `ListFiles` above it resolves the effective sort first.
- `req.IsEmpty()` forces root scope when no parent, keyword or file key is present. Clear a term-less keyword before that check, or a keyword of separators only will turn a root listing into an account-wide one.
- The sort resolver is `resolveEffectiveSort`, `resolveEffectiveSortForDir` and `resolveFileTypeForSort` in the same file. Size coercion to File and the comic force already live there and need no change.
- `CalcFilePosition` already treats an absent parent as no parent filter, so the root size rule needs nothing from it.
- The second fulltext match, `ListVerFile` in `versioned_file.go`, stays exactly as it is.
- Tests sit beside the service code: `size_sort_test.go`, `core_test.go`, `resume_sort_test.go`. They run against a real MySQL database configured in `backend/vfm/conf.yml`.
- The frontend needs no change and is reference only. It sends `filename`, sends `parentFile` null at root, and already sets File plus descending when size is chosen.

## Facts that shaped the design

- `file_info.name` has no explicit collation, so it inherits utf8mb4 server default. That collation ends in `_ci`, which is why matching is case-insensitive with no `LOWER()`.
- `uploader_no_idx` narrows each search to one user's rows. The leading wildcard defeats `idx_parent_name`, so the scan is real but bounded.
- Nothing ever physically deletes `file_info` rows. Logical deletion is an update, so searching a stale projection is a non-issue and the fulltext index can go without leaving phantoms.
- A keyword at root is already account-wide, because the parent filter is only applied when the request names a parent.
- `github.com/m11s-io/wordninja-go` at `v0.0.0-20260710161850-848d63fa980c`, MIT, embeds a 528 KB gzipped English list, and adds no transitive requirements. `Split` returns a single value. Verified output: `myholiday` gives `my` and `holiday`, `MyHoliday` keeps case, `s01e02` gives `s`, `01`, `e`, `02`, and `therapist` stays whole.

## Choices the spec leaves to the ticket

- Deduplicate terms and cap the list at 16.
- No `ESCAPE` clause, because every surviving term is letters and digits only.
- Leave the deprecated `orderByName` flag alone.
- Keep the existing rule that a keyword under the time order sorts by row id.

## Release bookkeeping

- `docs/schema/v0.0.6/` does not exist yet, so the increment file has to be created. v0.0.6 is the in-flight release and its tag does not exist.
- `release/db-changes.md` has no v0.0.6 section yet. The section is added with the index drop.
- `release/v0.0.6.md` currently says no database migration and no config change. That line and its name-keyword note both need rewriting, and the root size scope belongs in the highlights.
- `docs/schema/current/vfm.sql` carries the index to drop.

## Build and verify

- `cd backend/vfm && go build ./... && go test ./...`. The service has `replace` directives to `../user-vault` and `../mini-fstore`.
- No frontend diff, so no frontend verification pass is needed.
- The full suite runs once, on the integration branch, after the writers are reconciled.

## Skills for the implementing session

- `miso`, mandatory before any backend Go code in this repo.
- `tdd`, red first at the two agreed seams: the listing operation and the pure splitter. TDD was not confirmed in the design session, so the implementing session confirms it with the user before the first test is written. The seams proposed there are the listing operation against the test database and the pure splitter as a table test.
- `code-review`, against `docs/agents/code-review-backend.md`.
