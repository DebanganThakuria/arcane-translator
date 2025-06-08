package models

// Stats represents the statistics of the application
type Stats struct {
	NovelCount   int `json:"novel_count"`
	ChapterCount int `json:"chapter_count"`
}
