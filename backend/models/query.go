package models

import "strings"

// NovelSort names a column the novel list may be ordered by. The API accepts
// only these values so the ORDER BY clause is never built from user input.
type NovelSort string

const (
	SortLastRead     NovelSort = "last_read"
	SortLastUpdated  NovelSort = "last_updated"
	SortDateAdded    NovelSort = "date_added"
	SortTitle        NovelSort = "title"
	SortChapterCount NovelSort = "chapters"
)

// SortExpression maps a sort key onto the SQL it is allowed to produce.
var sortExpressions = map[NovelSort]string{
	SortLastRead:     "last_read_timestamp",
	SortLastUpdated:  "last_updated",
	SortDateAdded:    "date_added",
	SortTitle:        "title COLLATE NOCASE",
	SortChapterCount: "chapters_count",
}

// ParseNovelSort returns the matching sort, falling back to last read.
func ParseNovelSort(value string) NovelSort {
	sort := NovelSort(strings.TrimSpace(strings.ToLower(value)))
	if _, ok := sortExpressions[sort]; ok {
		return sort
	}
	return SortLastRead
}

// Expression returns the whitelisted SQL fragment for this sort.
func (s NovelSort) Expression() string {
	if expression, ok := sortExpressions[s]; ok {
		return expression
	}
	return sortExpressions[SortLastRead]
}

// NovelQuery describes a filtered, sorted, paginated novel lookup. Zero values
// mean "no constraint", so an empty query lists everything.
type NovelQuery struct {
	// Search matches title, original title, author and genres.
	Search string
	// SourceIDs restricts to a set of sources; a language resolves to these.
	SourceIDs []string
	Genre     string
	// Status matches the stored value case-insensitively, e.g. "Ongoing".
	Status string

	Sort      NovelSort
	Ascending bool

	Offset int
	Limit  int
}
