package models

import "testing"

// The sort key reaches SQL as an ORDER BY fragment, so anything outside the
// whitelist must fall back rather than be interpolated.
func TestParseNovelSort(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  NovelSort
	}{
		{"known key", "title", SortTitle},
		{"uppercase", "TITLE", SortTitle},
		{"padded", "  chapters  ", SortChapterCount},
		{"empty falls back", "", SortLastRead},
		{"unknown falls back", "not_a_column", SortLastRead},
		{"injection attempt falls back", "title; DROP TABLE novels--", SortLastRead},
		{"column not in whitelist falls back", "summary", SortLastRead},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseNovelSort(tt.input); got != tt.want {
				t.Errorf("ParseNovelSort(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestNovelSortExpression(t *testing.T) {
	// Every declared sort must map to SQL, or ORDER BY would be built empty.
	for _, sort := range []NovelSort{
		SortLastRead, SortLastUpdated, SortDateAdded, SortTitle, SortChapterCount,
	} {
		if sort.Expression() == "" {
			t.Errorf("sort %q has no expression", sort)
		}
	}

	if got := NovelSort("bogus").Expression(); got != sortExpressions[SortLastRead] {
		t.Errorf("unknown sort produced %q, want the last_read expression", got)
	}

	if got := SortTitle.Expression(); got != "title COLLATE NOCASE" {
		t.Errorf("title sort = %q, want a case-insensitive collation", got)
	}
}
