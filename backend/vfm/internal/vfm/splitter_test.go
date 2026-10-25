package vfm

import (
	"reflect"
	"testing"
	"unicode"
)

// TestSplitSearchTerms pins the keyword splitter: the terms a typed keyword
// yields before each becomes a required name substring.
func TestSplitSearchTerms(t *testing.T) {
	cases := []struct {
		name    string
		keyword string
		want    []string
	}{
		{"empty", "", nil},
		{"spaces only", "    ", nil},
		{"separators only", "---___%", nil},

		// tier 1: separators are interchangeable
		{"space separator", "my holiday", []string{"my", "holiday"}},
		{"dash separator", "my-holiday", []string{"my", "holiday"}},
		{"underscore separator", "my_holiday", []string{"my", "holiday"}},
		{"punctuation separator", "%my.holiday!", []string{"my", "holiday"}},

		// tier 3: dictionary splits a concatenated Latin run, but only when the
		// split is a real decomposition into word-list entries
		{"dictionary tier", "myholiday", []string{"my", "holiday"}},
		{"one-letter word a is a real piece", "aholiday", []string{"a", "holiday"}},
		{"pinyin decomposes into word-list entries", "zhangsan", []string{"zhang", "san"}},
		{"nonsense run stays whole", "qwzzx", []string{"qwzzx"}},
		{"brand run stays whole", "sennheiser", []string{"sennheiser"}},

		// tier 2: letter/digit and lower/upper transitions
		{"letter to digit", "IMG2024", []string{"IMG", "2024"}},
		{"digit to letter", "2024holiday", []string{"2024", "holiday"}},
		{"lower to upper", "MyHoliday", []string{"My", "Holiday"}},
		{"all transitions", "S01E02", []string{"S", "01", "E", "02"}},

		// tiers interacting
		{"separator then transitions", "my-holiday_IMG2024", []string{"my", "holiday", "IMG", "2024"}},

		// non-Latin runs are matched as substrings, never segmented
		{"chinese run", "剧集", []string{"剧集"}},
		{"chinese then digits", "剧集01", []string{"剧集", "01"}},

		// a run that is already a word stays whole
		{"already a word", "holiday", []string{"holiday"}},
		{"dictionary does not over-split", "therapist", []string{"therapist"}},

		// dedupe and cap
		{"dedupe", "my my holiday", []string{"my", "holiday"}},
		{"dedupe ignores case", "My my", []string{"My"}},
		{
			"cap at 16",
			"1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17",
			[]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16"},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := splitSearchTerms(c.keyword)
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("splitSearchTerms(%q) = %#v, want %#v", c.keyword, got, c.want)
			}
		})
	}
}

// TestSplitSearchTermsProducesOnlyAlphanumerics is the no-escaping rule: a
// surviving term is letters and digits only, so % and _ can never reach LIKE.
func TestSplitSearchTermsProducesOnlyAlphanumerics(t *testing.T) {
	for _, term := range splitSearchTerms(`a%b_c\d'e"f/g 100%`) {
		for _, r := range term {
			if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
				t.Fatalf("term %q contains non-alphanumeric rune %q", term, r)
			}
		}
	}
}
