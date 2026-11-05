# 01: Cover comes from a visible child

**What to build:** `dir-thumbnail` and `dir-thumbnail/batch` pick a directory's cover from the same set of children the listing shows: not logically deleted, not physically deleted, not hidden, and holding a thumbnail. A directory whose children are all invisible answers with an empty token instead of a stale cover or a failed token call.

**Blocked by:** None

**Status:** done

- [x] The single-directory cover lookup skips children with `is_logic_deleted = 1`, `hidden = 1`, or `is_del = 1`
- [x] The batch cover lookup applies the same visibility predicate per directory
- [x] The batch lookup only considers a live directory row (`is_logic_deleted = 0` in addition to today's filters)
- [x] The oldest qualifying child still wins, and a directory with no qualifying child still answers with an empty token
- [x] One directory's trashed children cannot suppress another directory's cover in the same batch response
- [x] Backend tests cover: trashed child skipped for a live one, hidden child skipped, physically deleted child skipped, child without a thumbnail skipped, all-invisible children answer empty, oldest qualifying child wins
- [x] `go build ./...` and `go test ./...` pass in `backend/vfm`
- [x] No schema, endpoint, response-shape or ownership-check change
