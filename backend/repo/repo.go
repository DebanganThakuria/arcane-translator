package repo

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"backend/config"
	"backend/models"

	"github.com/google/uuid"
)

var repository Repo

type Repo interface {
	GetStatus() (bool, error)

	// Stats methods
	GetStats() (int, int, error) // Returns (novelCount, chapterCount, error)

	// Novel methods
	QueryNovels(query models.NovelQuery) ([]*models.Novel, int, error)
	GetNovelByID(id string) (*models.Novel, error)
	GetNovelsByRecentlyUpdated(count int) ([]*models.Novel, error)
	GetNovelsByRecentlyRead(count int) ([]*models.Novel, error)
	CreateNovel(novel *models.Novel) (*models.Novel, error)
	SearchNovel(query string) ([]*models.Novel, error)
	UpdateNovel(novel *models.Novel) error
	DeleteNovel(id string) error
	UpdateLastReadChapter(novelID string, chapterNumber int) error

	// Chapter methods
	QueryChapters(novelID string, query models.ChapterQuery) ([]*models.Chapter, int, error)
	ChapterNumberBounds(novelID string) (first, last int, err error)
	GetNovelChapters(novelID string) ([]*models.Chapter, error)
	GetLastChapter(novelID string) (*models.Chapter, error)
	GetChapterByID(novelID string, chapterID string) (*models.Chapter, error)
	GetChapterByNumber(novelID string, chapterNumber int) (*models.Chapter, error)
	GetChapterByURL(url string) (*models.Chapter, error)
	CreateChapter(chapter *models.Chapter) (*models.Chapter, error)
	UpdateChapter(chapter *models.Chapter) error
	DeleteChapter(novelID string, chapterID string) error
}

type repo struct {
	db DB
}

func init() {
	// Set up the database file path
	dbPath := config.Get().DBPath

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

// QueryNovels is the single filtered/sorted/paginated novel lookup. Filters
// compose with AND, so language, genre, status and a search term can be applied
// together, and the count reflects the same filters as the page.
//
// The ORDER BY fragment comes from models.NovelSort, never from raw input.
func (r *repo) QueryNovels(query models.NovelQuery) ([]*models.Novel, int, error) {
	var conditions []string
	var args []any

	if len(query.SourceIDs) > 0 {
		placeholders := make([]string, len(query.SourceIDs))
		for i, source := range query.SourceIDs {
			placeholders[i] = "?"
			args = append(args, source)
		}
		conditions = append(conditions, "source IN ("+strings.Join(placeholders, ",")+")")
	}

	if query.Genre != "" {
		conditions = append(conditions, "genres LIKE ? COLLATE NOCASE")
		args = append(args, "%"+query.Genre+"%")
	}

	if query.Status != "" {
		conditions = append(conditions, "status = ? COLLATE NOCASE")
		args = append(args, query.Status)
	}

	if query.Search != "" {
		conditions = append(conditions, `(
			title LIKE ? COLLATE NOCASE
			OR original_title LIKE ? COLLATE NOCASE
			OR author LIKE ? COLLATE NOCASE
			OR genres LIKE ? COLLATE NOCASE
		)`)
		term := "%" + query.Search + "%"
		args = append(args, term, term, term, term)
	}

	where := ""
	if len(conditions) > 0 {
		where = " WHERE " + strings.Join(conditions, " AND ")
	}

	var count int
	if err := r.db.QueryRow("SELECT COUNT(*) FROM novels"+where, args...).Scan(&count); err != nil {
		return nil, 0, err
	}

	direction := "DESC"
	if query.Ascending {
		direction = "ASC"
	}

	listSQL := fmt.Sprintf(
		"SELECT * FROM novels%s ORDER BY %s %s LIMIT ? OFFSET ?",
		where, query.Sort.Expression(), direction,
	)

	rows, err := r.db.Query(listSQL, append(append([]any{}, args...), query.Limit, query.Offset)...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	novels, err := models.ScanNovels(rows)
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

// chapterListColumns omits `content`, which is by far the largest column and is
// only ever needed when reading one chapter.
const chapterListColumns = `id, novel_id, number, title, original_title,
	date_translated, word_count, url, next_chapter_url`

// QueryChapters returns one page of a novel's chapters plus the total matching
// the same filter. Novels here reach several hundred chapters, so the list is
// paged in SQL rather than loaded whole and sliced in the client.
func (r *repo) QueryChapters(novelID string, query models.ChapterQuery) ([]*models.Chapter, int, error) {
	conditions := []string{"novel_id = ?"}
	args := []any{novelID}

	if query.Search != "" {
		conditions = append(conditions,
			"(title LIKE ? COLLATE NOCASE OR original_title LIKE ? COLLATE NOCASE)")
		term := "%" + query.Search + "%"
		args = append(args, term, term)
	}

	where := " WHERE " + strings.Join(conditions, " AND ")

	var count int
	if err := r.db.QueryRow("SELECT COUNT(*) FROM chapters"+where, args...).Scan(&count); err != nil {
		return nil, 0, err
	}

	direction := "DESC"
	if query.Ascending {
		direction = "ASC"
	}

	listSQL := fmt.Sprintf(
		"SELECT %s FROM chapters%s ORDER BY number %s LIMIT ? OFFSET ?",
		chapterListColumns, where, direction,
	)

	rows, err := r.db.Query(listSQL, append(append([]any{}, args...), query.Limit, query.Offset)...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	chapters, err := models.ScanChapters(rows)
	if err != nil {
		return nil, 0, err
	}

	return chapters, count, nil
}

// ChapterNumberBounds returns the lowest and highest chapter numbers a novel
// has, or zeroes when it has none.
func (r *repo) ChapterNumberBounds(novelID string) (int, int, error) {
	var first, last sql.NullInt64

	err := r.db.
		QueryRow("SELECT MIN(number), MAX(number) FROM chapters WHERE novel_id = ?", novelID).
		Scan(&first, &last)
	if err != nil {
		return 0, 0, err
	}

	return int(first.Int64), int(last.Int64), nil
}

func (r *repo) GetNovelChapters(novelID string) ([]*models.Chapter, error) {
	// We are not loading the content of the chapters to decrease the memory usage
	query := `
		SELECT ` + chapterListColumns + `
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

func (r *repo) DeleteChapter(novelID string, chapterID string) error {
	query := `
		DELETE FROM chapters
		WHERE novel_id = ? AND id = ?
	`

	result, err := r.db.Exec(query, novelID, chapterID)
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
