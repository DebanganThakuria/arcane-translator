package models

import (
	"database/sql"
	"encoding/json"
)

// Novel represents a novel in the database
type Novel struct {
	ID                    string   `json:"id"`
	Title                 string   `json:"title"`
	OriginalTitle         string   `json:"original_title,omitempty"`
	Cover                 string   `json:"cover,omitempty"`
	Source                string   `json:"source"`
	URL                   string   `json:"url"`
	Summary               string   `json:"summary"`
	Author                string   `json:"author,omitempty"`
	Status                string   `json:"status,omitempty"`
	Genres                []string `json:"genres,omitempty"`
	ChaptersCount         int      `json:"chapters_count"`
	LastReadChapterNumber int      `json:"last_read_chapter_number,omitempty"`
	LastReadTimestamp     int64    `json:"last_read_timestamp,omitempty"`
	LastUpdated           int64    `json:"last_updated"`
	DateAdded             int64    `json:"date_added"`
}

// Chapter represents a chapter in the database
type Chapter struct {
	ID             string `json:"id"`
	NovelID        string `json:"novel_id"`
	Number         int    `json:"number"`
	Title          string `json:"title"`
	OriginalTitle  string `json:"original_title,omitempty"`
	Content        string `json:"content"`
	DateTranslated int64  `json:"date_translated"`
	WordCount      int    `json:"word_count,omitempty"`
	URL            string `json:"url,omitempty"`
	NextChapterURL string `json:"next_chapter_url,omitempty"`
}

// SourceSite represents a source site for novels
type SourceSite struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	URL      string `json:"url"`
	Language string `json:"language"` // "Chinese", "Korean", "Japanese", "Other"
	Icon     string `json:"icon,omitempty"`
}

// ScanNovel scans a novel from a SQL row
func ScanNovel(row *sql.Row) (*Novel, error) {
	var novel Novel
	var genresJSON string
	var lastUpdatedUnix, dateAddedUnix int64

	err := row.Scan(
		&novel.ID,
		&novel.Title,
		&novel.OriginalTitle,
		&novel.Cover,
		&novel.Source,
		&novel.URL,
		&novel.Summary,
		&novel.Author,
		&novel.Status,
		&genresJSON,
		&novel.ChaptersCount,
		&novel.LastReadChapterNumber,
		&novel.LastReadTimestamp,
		&lastUpdatedUnix,
		&dateAddedUnix,
	)
	if err != nil {
		return nil, err
	}

	// Parse genres from JSON string
	if genresJSON != "" {
		if err := json.Unmarshal([]byte(genresJSON), &novel.Genres); err != nil {
			// If there's an error, set an empty genres slice
			novel.Genres = []string{}
		}
	}

	novel.LastUpdated = lastUpdatedUnix
	novel.DateAdded = lastUpdatedUnix

	return &novel, nil
}

// ScanNovels scans multiple novels from SQL rows
func ScanNovels(rows *sql.Rows) ([]*Novel, error) {
	var novels []*Novel

	for rows.Next() {
		var novel Novel
		var genresJSON string
		var lastUpdatedUnix, dateAddedUnix int64

		err := rows.Scan(
			&novel.ID,
			&novel.Title,
			&novel.OriginalTitle,
			&novel.Cover,
			&novel.Source,
			&novel.URL,
			&novel.Summary,
			&novel.Author,
			&novel.Status,
			&genresJSON,
			&novel.ChaptersCount,
			&novel.LastReadChapterNumber,
			&novel.LastReadTimestamp,
			&lastUpdatedUnix,
			&dateAddedUnix,
		)
		if err != nil {
			return nil, err
		}

		// Parse genres from JSON string
		if genresJSON != "" {
			if err := json.Unmarshal([]byte(genresJSON), &novel.Genres); err != nil {
				// If there's an error, set empty genres slice
				novel.Genres = []string{}
			}
		}

		novel.LastUpdated = lastUpdatedUnix
		novel.DateAdded = lastUpdatedUnix

		novels = append(novels, &novel)
	}

	return novels, nil
}

// ScanChapter scans a chapter from a SQL row
func ScanChapter(row *sql.Row) (*Chapter, error) {
	var chapter Chapter
	var dateTranslatedUnix int64

	err := row.Scan(
		&chapter.ID,
		&chapter.NovelID,
		&chapter.Number,
		&chapter.Title,
		&chapter.OriginalTitle,
		&chapter.Content,
		&dateTranslatedUnix,
		&chapter.WordCount,
		&chapter.URL,
		&chapter.NextChapterURL,
	)
	if err != nil {
		return nil, err
	}

	chapter.DateTranslated = dateTranslatedUnix

	return &chapter, nil
}

// ScanChapters scans multiple chapters from SQL rows
func ScanChapters(rows *sql.Rows) ([]*Chapter, error) {
	var chapters []*Chapter

	for rows.Next() {
		var chapter Chapter
		var dateTranslatedUnix int64

		err := rows.Scan(
			&chapter.ID,
			&chapter.NovelID,
			&chapter.Number,
			&chapter.Title,
			&chapter.OriginalTitle,
			&dateTranslatedUnix,
			&chapter.WordCount,
			&chapter.URL,
			&chapter.NextChapterURL,
		)
		if err != nil {
			return nil, err
		}

		chapter.DateTranslated = dateTranslatedUnix

		chapters = append(chapters, &chapter)
	}

	return chapters, nil
}

// GenresToJSON converts a slice of genre strings to a JSON string for storage
func GenresToJSON(genres []string) (string, error) {
	if len(genres) == 0 {
		return "[]", nil
	}

	bytes, err := json.Marshal(genres)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func GenresToString(genres []string) string {
	if len(genres) == 0 {
		return ""
	}

	bytes, err := json.Marshal(genres)
	if err != nil {
		return ""
	}

	return string(bytes)
}
