package repo

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"backend/models"
	"backend/utils"

	"github.com/google/uuid"
)

var repository Repo

type Repo interface {
	GetStatus() (bool, error)

	// Stats methods
	GetStats() (int, int, error) // Returns (novelCount, chapterCount, error)

	// Novel methods
	GetAllNovels(offset, limit int) ([]*models.Novel, int, error)
	GetNovelByID(id string) (*models.Novel, error)
	GetNovelsBySourceIDs(source []string, offset, limit int) ([]*models.Novel, int, error)
	GetNovelsByGenre(genre string, offset, limit int) ([]*models.Novel, int, error)
	GetNovelsByRecentlyUpdated(count int) ([]*models.Novel, error)
	GetNovelsByRecentlyRead(count int) ([]*models.Novel, error)
	CreateNovel(novel *models.Novel) (*models.Novel, error)
	SearchNovel(query string) ([]*models.Novel, error)
	UpdateNovel(novel *models.Novel) error
	DeleteNovel(id string) error
	UpdateLastReadChapter(novelID string, chapterNumber int) error

	// Chapter methods
	GetNovelChapters(novelID string) ([]*models.Chapter, error)
	GetLastChapter(novelID string) (*models.Chapter, error)
	GetChapterByID(novelID string, chapterID string) (*models.Chapter, error)
	GetChapterByNumber(novelID string, chapterNumber int) (*models.Chapter, error)
	GetChapterByURL(url string) (*models.Chapter, error)
	CreateChapter(chapter *models.Chapter) (*models.Chapter, error)
	UpdateChapter(chapter *models.Chapter) error
}

type repo struct {
	db DB
}

func init() {
	// Set up the database file path
	dbPath := utils.GetDBPath()

	// Initialize SQLite database
	db, err := NewSQLiteDB(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize repository
	repository = NewRepo(db)
}

func NewRepo(db DB) Repo {
	return &repo{
		db: db,
	}
}

func GetRepo() Repo {
	return repository
}

func (r *repo) GetStatus() (bool, error) {
	err := r.db.Ping()
	return err == nil, err
}

// GetStats returns the total count of novels and chapters
func (r *repo) GetStats() (int, int, error) {
	// Get novel count
	var novelCount int
	err := r.db.QueryRow("SELECT COUNT(*) FROM novels").Scan(&novelCount)
	if err != nil {
		return 0, 0, err
	}

	// Get chapter count
	var chapterCount int
	err = r.db.QueryRow("SELECT COUNT(*) FROM chapters").Scan(&chapterCount)
	if err != nil {
		return 0, 0, err
	}

	return novelCount, chapterCount, nil
}

// Novel CRUD Operations

func (r *repo) GetAllNovels(offset, limit int) ([]*models.Novel, int, error) {
	query := `
		SELECT *
		FROM novels
		ORDER BY last_read_timestamp DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	novels, err := models.ScanNovels(rows)
	if err != nil {
		return nil, 0, err
	}

	var count int
	err = r.db.QueryRow("SELECT COUNT(*) FROM novels").Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	return novels, count, nil
}

func (r *repo) GetNovelByID(id string) (*models.Novel, error) {
	query := `
		SELECT *
		FROM novels
		WHERE id = ?
	`

	row := r.db.QueryRow(query, id)
	return models.ScanNovel(row)
}

func (r *repo) GetNovelsBySourceIDs(sources []string, offset, limit int) ([]*models.Novel, int, error) {
	if len(sources) == 0 {
		return []*models.Novel{}, 0, nil
	}

	// Create placeholders for each source
	placeholders := make([]string, len(sources))
	args := make([]interface{}, len(sources))
	for i, source := range sources {
		placeholders[i] = "?"
		args[i] = source
	}
	args = append(args, limit)
	args = append(args, offset)

	query := fmt.Sprintf(`
		SELECT *
		FROM novels
		WHERE source IN (%s)
		ORDER BY last_read_timestamp DESC
		LIMIT ? OFFSET ?
	`, strings.Join(placeholders, ","))

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	novels, err := models.ScanNovels(rows)
	if err != nil {
		return nil, 0, err
	}

	var count int
	err = r.db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM novels WHERE source IN (%s)", strings.Join(placeholders, ",")), args...).Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	return novels, count, nil
}

func (r *repo) GetNovelsByGenre(genre string, offset, limit int) ([]*models.Novel, int, error) {
	query := `
		SELECT *
		FROM novels
		WHERE genres LIKE ? COLLATE NOCASE
		ORDER BY last_read_timestamp DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.Query(query, "%"+genre+"%", limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	novels, err := models.ScanNovels(rows)
	if err != nil {
		return nil, 0, err
	}

	var count int
	err = r.db.QueryRow("SELECT COUNT(*) FROM novels WHERE genres LIKE ? COLLATE NOCASE", "%"+genre+"%").Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	return novels, count, nil
}

func (r *repo) GetNovelsByRecentlyUpdated(count int) ([]*models.Novel, error) {
	query := `
		SELECT *
		FROM novels
		ORDER BY last_updated DESC
		LIMIT ?
	`

	rows, err := r.db.Query(query, count)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return models.ScanNovels(rows)
}

func (r *repo) GetNovelsByRecentlyRead(count int) ([]*models.Novel, error) {
	query := `
		SELECT *
		FROM novels
		ORDER BY last_read_timestamp DESC
		LIMIT ?
	`

	rows, err := r.db.Query(query, count)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return models.ScanNovels(rows)
}

func (r *repo) SearchNovel(query string) ([]*models.Novel, error) {
	sqlQuery := `
		SELECT *
		FROM novels
		WHERE title LIKE ? COLLATE NOCASE
		ORDER BY last_read_timestamp DESC
		LIMIT 20
	`

	searchTerm := "%" + query + "%"

	rows, err := r.db.Query(sqlQuery, searchTerm)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return models.ScanNovels(rows)
}

func (r *repo) CreateNovel(novel *models.Novel) (*models.Novel, error) {
	// Convert genres to JSON
	genresJSON, err := models.GenresToJSON(novel.Genres)
	if err != nil {
		return nil, err
	}

	query := `
		INSERT INTO novels (id, title, original_title, cover, source, url, summary, author, status, genres, chapters_count, last_read_chapter_number, last_read_timestamp, last_updated, date_added)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err = r.db.Exec(
		query,
		novel.ID,
		novel.Title,
		novel.OriginalTitle,
		novel.Cover,
		novel.Source,
		novel.URL,
		novel.Summary,
		novel.Author,
		novel.Status,
		genresJSON,
		novel.ChaptersCount,
		novel.LastReadChapterNumber,
		novel.LastReadTimestamp,
		novel.LastUpdated,
		novel.DateAdded,
	)
	if err != nil {
		return nil, err
	}

	return novel, nil
}

func (r *repo) UpdateNovel(novel *models.Novel) error {
	// Convert genres to JSON
	genresJSON, err := models.GenresToJSON(novel.Genres)
	if err != nil {
		return err
	}

	query := `
		UPDATE novels
		SET title = ?, original_title = ?, cover = ?, source = ?, url = ?, 
		    summary = ?, author = ?, status = ?, genres = ?, chapters_count = ?, 
		    last_read_chapter_number = ?, last_read_timestamp = ?, last_updated = ?
		WHERE id = ?
	`

	result, err := r.db.Exec(
		query,
		novel.Title,
		novel.OriginalTitle,
		novel.Cover,
		novel.Source,
		novel.URL,
		novel.Summary,
		novel.Author,
		novel.Status,
		genresJSON,
		novel.ChaptersCount,
		novel.LastReadChapterNumber,
		novel.LastReadTimestamp,
		novel.LastUpdated,
		novel.ID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("novel not found")
	}

	return nil
}

func (r *repo) DeleteNovel(id string) error {
	// First delete all chapters for this novel
	_, err := r.db.Exec("DELETE FROM chapters WHERE novel_id = ?", id)
	if err != nil {
		return err
	}

	// Then delete the novel
	result, err := r.db.Exec("DELETE FROM novels WHERE id = ?", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("novel not found")
	}

	return nil
}

func (r *repo) UpdateLastReadChapter(novelID string, chapterNumber int) error {
	currentTime := time.Now().Unix()

	query := `
		UPDATE novels
		SET last_read_chapter_number = ?,
		    last_read_timestamp = ?
		WHERE id = ?
	`

	_, err := r.db.Exec(query, chapterNumber, currentTime, novelID)
	return err
}

// Chapter CRUD Operations

func (r *repo) GetNovelChapters(novelID string) ([]*models.Chapter, error) {
	// We are not loading the content of the chapters to decrease the memory usage
	query := `
		SELECT id, novel_id, number, title, original_title, date_translated, word_count, url, next_chapter_url
		FROM chapters
		WHERE novel_id = ?
		ORDER BY number
	`

	rows, err := r.db.Query(query, novelID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return models.ScanChapters(rows)
}

func (r *repo) GetLastChapter(novelID string) (*models.Chapter, error) {
	query := `
		SELECT *
		FROM chapters
		WHERE novel_id = ?
		ORDER BY number DESC
		LIMIT 1
	`

	row := r.db.QueryRow(query, novelID)
	return models.ScanChapter(row)
}

func (r *repo) GetChapterByID(novelID string, chapterID string) (*models.Chapter, error) {
	query := `
		SELECT *
		FROM chapters
		WHERE novel_id = ? AND id = ?
	`

	row := r.db.QueryRow(query, novelID, chapterID)
	return models.ScanChapter(row)
}

func (r *repo) GetChapterByNumber(novelID string, chapterNumber int) (*models.Chapter, error) {
	query := `
		SELECT *
		FROM chapters
		WHERE novel_id = ? AND number = ?
	`

	row := r.db.QueryRow(query, novelID, chapterNumber)
	return models.ScanChapter(row)
}

func (r *repo) GetChapterByURL(url string) (*models.Chapter, error) {
	query := `
		SELECT *
		FROM chapters
		WHERE url = ?
	`

	row := r.db.QueryRow(query, url)
	return models.ScanChapter(row)
}

func (r *repo) CreateChapter(chapter *models.Chapter) (*models.Chapter, error) {
	// Generate a new UUID if not provided
	if chapter.ID == "" {
		chapter.ID = uuid.New().String()
	}

	query := `
		INSERT INTO chapters (
			id, novel_id, number, title, original_title, content, date_translated, word_count, url, next_chapter_url
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.Exec(
		query,
		chapter.ID,
		chapter.NovelID,
		chapter.Number,
		chapter.Title,
		chapter.OriginalTitle,
		chapter.Content,
		chapter.DateTranslated,
		chapter.WordCount,
		chapter.URL,
		chapter.NextChapterURL,
	)
	if err != nil {
		return nil, err
	}

	return chapter, nil
}

func (r *repo) UpdateChapter(chapter *models.Chapter) error {
	query := `
		UPDATE chapters
		SET number = ?, title = ?, original_title = ?, content = ?, 
		    date_translated = ?, word_count = ?, url = ?, next_chapter_url = ?
		WHERE id = ? AND novel_id = ?
	`

	result, err := r.db.Exec(
		query,
		chapter.Number,
		chapter.Title,
		chapter.OriginalTitle,
		chapter.Content,
		chapter.DateTranslated,
		chapter.WordCount,
		chapter.URL,
		chapter.NextChapterURL,
		chapter.ID,
		chapter.NovelID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("chapter not found")
	}

	return nil
}
