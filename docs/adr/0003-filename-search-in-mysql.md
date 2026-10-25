# Keep filename search in MySQL, split the keyword in vfm

Status: accepted

The manage-files search box narrows a listing by file name. It ran a MySQL fulltext match on the name column with the index's default parser, which imposes a three character minimum token, matches whole tokens only, and keeps a run of Chinese characters as one token. The parser therefore decided what a word is, and file names are full of runs it does not recognise, such as `IMG_2024`, `S01E02` and `x265`.

Search stays in MySQL. `vfm` splits the typed keyword into terms itself, at every separator, at letter and digit transitions, and, for a remaining run of Latin letters, into dictionary words through a bundled English word list. Each term becomes a `LIKE` condition on the name column, and the conditions are joined with `AND`. Case-insensitivity comes from the column collation.

A leading wildcard rules out an index, so each term scans the requesting user's rows, bounded by the existing user index. For a corpus in the tens of thousands per user that is a few milliseconds, the same order as the fulltext match it replaces.

## Considered Options

- **Embedded search index (bleve, or bluge).** Gives prefix, fuzzy and n-gram matching with scoring, but adds an index store that must stay in sync with `file_info`, a rebuild path, a persistent volume, and an exclusive lock that makes the index a per-process resource. With more than one `vfm` instance each process would need its own queue on the event exchange, because a stream queue delivers each event to exactly one consumer.
- **External search service (Meilisearch, Typesense).** Solves the multi-instance sync problem by making the index shared, and brings better matching for free, but adds a stateful service, a volume, a key, a sync worker, a page enrichment path and a new failure mode. Worth it when ranking or typo tolerance is wanted, or when another client needs search. Neither was wanted here.
- **MySQL fulltext with the ngram parser.** Accepts substrings and Chinese, but phrase matching is still exact, so the user must still type the name accurately, and the index grows considerably.
- **A stored compacted name column.** Separators become irrelevant in both directions, but it adds a derived column, a backfill, and a second copy of every name to keep in sync, and a leading wildcard still gets no index out of it.

## Consequences

- Search quality is bounded by substring matching. Typo tolerance and relevance ranking are absent by decision, not by oversight.
- A concatenated run splits only when its parts are English dictionary words. Pinyin, brand names and model numbers stay whole.
- The fulltext index on the file name column is dropped with the change, so there is no fulltext fallback on that table. The versioned-file table keeps its own index and its current match.
- Revisiting this is additive. A compacted column, or an external search service, can be added without changing the request or response shape.
