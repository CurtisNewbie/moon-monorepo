package vfm

import (
	"strings"
	"unicode"

	"github.com/m11s-io/wordninja-go"
)

// maxSearchTerms bounds how many required name fragments one keyword adds.
const maxSearchTerms = 16

// splitSearchTerms turns a typed keyword into the name substrings the listing
// requires, every one of them. Splitting happens in three tiers:
//
//  1. every character that is neither a letter nor a digit separates runs, so
//     spaces, dashes and underscores are interchangeable;
//  2. a run splits at letter/digit and digit/letter changes and at a
//     lower-to-upper change, which handles runs without any dictionary;
//  3. a remaining Latin-letter run goes through the bundled English dictionary,
//     which turns a concatenated name into words. The dictionary tier keeps a
//     run whole when its pieces are not all in the bundled word list.
//
// Every term is a substring of the keyword, so splitting can only widen a match,
// never hide a row. The surviving terms are letters and digits only, so no LIKE
// wildcard or escape character can reach the query. Terms are deduplicated
// case-insensitively and capped.
func splitSearchTerms(keyword string) []string {
	terms := make([]string, 0, 4)
	seen := make(map[string]struct{}, 4)
	addTerm := func(term string) {
		if term == "" {
			return
		}
		lower := strings.ToLower(term)
		if _, ok := seen[lower]; ok {
			return
		}
		seen[lower] = struct{}{}
		terms = append(terms, term)
	}

	run := make([]rune, 0, len(keyword))
	flushRun := func() {
		if len(run) == 0 {
			return
		}
		for _, chunk := range splitRunTransitions(string(run)) {
			for _, term := range splitLatinRun(chunk) {
				addTerm(term)
			}
		}
		run = run[:0]
	}

	for _, r := range keyword {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			run = append(run, r)
			continue
		}
		flushRun()
	}
	flushRun()

	if len(terms) == 0 {
		return nil
	}
	if len(terms) > maxSearchTerms {
		terms = terms[:maxSearchTerms]
	}
	return terms
}

// splitRunTransitions splits one alphanumeric run at letter/digit and
// digit/letter changes and at lower-to-upper changes.
func splitRunTransitions(run string) []string {
	runes := []rune(run)
	chunks := make([]string, 0, 4)
	start := 0
	for i := 1; i < len(runes); i++ {
		prev, cur := runes[i-1], runes[i]
		letterToDigit := unicode.IsLetter(prev) && unicode.IsDigit(cur)
		digitToLetter := unicode.IsDigit(prev) && unicode.IsLetter(cur)
		lowerToUpper := unicode.IsLower(prev) && unicode.IsUpper(cur)
		if letterToDigit || digitToLetter || lowerToUpper {
			chunks = append(chunks, string(runes[start:i]))
			start = i
		}
	}
	return append(chunks, string(runes[start:]))
}

// splitLatinRun runs the dictionary tier over a run of ASCII letters. Any other
// run (digits, Chinese, other scripts) is kept whole, because a non-Latin term
// already matches as a substring.
//
// The dictionary split is accepted only when it is a real decomposition. The
// model prices an entry it does not know as infinite, so a run it cannot
// decompose comes back as one character per token; a single-letter piece that is
// not one of the English one-letter words therefore means the pieces are not all
// in the bundled word list, and the run stays whole.
func splitLatinRun(chunk string) []string {
	if len(chunk) < 2 || !asciiLettersOnly(chunk) {
		return []string{chunk}
	}
	parts := wordninja.Split(chunk)
	if len(parts) == 0 || !realWordSplit(parts) {
		return []string{chunk}
	}
	return parts
}

// realWordSplit reports whether a dictionary split decomposes the run into
// words. A one-letter piece outside a/A/i/I is the model's one-character
// fallback for an unknown run, not a word. The input is ASCII letters only, so
// a one-byte piece is a one-letter piece.
func realWordSplit(parts []string) bool {
	for _, part := range parts {
		if len(part) == 1 && !oneLetterWord(part) {
			return false
		}
	}
	return true
}

// oneLetterWord reports whether s is one of the one-letter English words, the
// only single-letter pieces a real dictionary split can contain.
func oneLetterWord(s string) bool {
	switch s {
	case "a", "A", "i", "I":
		return true
	}
	return false
}

func asciiLettersOnly(s string) bool {
	for _, r := range s {
		if r < 'A' || (r > 'Z' && r < 'a') || r > 'z' {
			return false
		}
	}
	return true
}
