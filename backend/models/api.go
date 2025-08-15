package models

// NovelListResponse represents the response structure for paginated novels
// @Description Paginated list of novels with metadata
type NovelListResponse struct {
	Novels      []*Novel `json:"novels"`
	TotalCount  int      `json:"total_count"`
	CurrentPage int      `json:"current_page"`
	TotalPages  int      `json:"total_pages"`
}

// NovelDetails represents details about a novel extracted and translated from a source
type NovelDetails struct {
	NovelTitleOriginal        string   `json:"novel_title_original"`
	NovelTitleTranslated      string   `json:"novel_title_translated"`
	NovelSummaryTranslated    string   `json:"novel_summary_translated"`
	NovelAuthorNameTranslated string   `json:"novel_author_name_translated,omitempty"`
	PossibleNovelGenres       []string `json:"possible_novel_genres,omitempty"`
	NumberOfChapters          int      `json:"number_of_chapters"`
	Status                    string   `json:"status,omitempty"`
}

// TranslatedChapter represents a chapter that has been translated
type TranslatedChapter struct {
	TranslatedChapterTitle    string   `json:"translated_chapter_title"`
	OriginalChapterTitle      string   `json:"original_chapter_title,omitempty"`
	TranslatedChapterContents string   `json:"translated_chapter_contents"`
	PossibleNewGenres         []string `json:"possible_new_genres,omitempty"`
}

// NovelExtractionRequest represents a request to extract novel details from a URL
type NovelExtractionRequest struct {
	URL         string  `json:"url"`
	Source      string  `json:"source"`
	HTMLContent *string `json:"html_content"`
}

// ChapterTranslationRequest represents a request to translate a chapter
type ChapterTranslationRequest struct {
	NovelID       string  `json:"novel_id"`
	ChapterNumber int     `json:"chapter_number"`
	ChapterURL    string  `json:"chapter_url,omitempty"`
	HTMLContent   *string `json:"html_content"`
}

// NovelRefreshRequest represents a request to refresh a novel's details
type NovelRefreshRequest struct {
	NovelID     string  `json:"novel_id"`
	HTMLContent *string `json:"html_content"`
}

// NovelRefreshResponse represents the response when refreshing a novel
type NovelRefreshResponse struct {
	Success          bool   `json:"success"`
	NewChaptersCount int    `json:"newChaptersCount"`
	Message          string `json:"message"`
	UpdatedDetails   struct {
		Summary string   `json:"summary"`
		Author  string   `json:"author"`
		Status  string   `json:"status"`
		Genres  []string `json:"genres"`
	} `json:"updatedDetails"`
}

// FirstChapterURLResponse represents the response when setting a first chapter URL
type FirstChapterURLResponse struct {
	Success           bool   `json:"success"`
	Message           string `json:"message"`
	FirstChapterURL   string `json:"firstChapterUrl,omitempty"`
	URLPattern        string `json:"urlPattern,omitempty"`
	ChapterTranslated bool   `json:"chapterTranslated,omitempty"`
}

// ChapterURLResponse represents the response when getting a chapter URL
type ChapterURLResponse struct {
	URL string `json:"url"`
}
