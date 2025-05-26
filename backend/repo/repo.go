package repo

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"backend/models"
	"backend/utils"

	"github.com/google/uuid"
)

var repository Repo

type Repo interface {
	GetStatus() (bool, error)

	// Novel methods
	GetAllNovels() ([]*models.Novel, error)
	GetNovelByID(id string) (*models.Novel, error)
	CreateNovel(novel *models.Novel) (*models.Novel, error)
	UpdateNovel(novel *models.Novel) error
	DeleteNovel(id string) error

	// Chapter methods
	GetNovelChapters(novelID string) ([]*models.Chapter, error)
	GetLastChapter(novelID string) (*models.Chapter, error)
	GetChapterByID(novelID string, chapterID string) (*models.Chapter, error)
	GetChapterByNumber(novelID string, chapterNumber int) (*models.Chapter, error)
	CreateChapter(chapter *models.Chapter) (*models.Chapter, error)
	UpdateChapter(chapter *models.Chapter) error
	DeleteChapter(id string) error
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

// Novel CRUD Operations

func (r *repo) GetAllNovels() ([]*models.Novel, error) {
	query := `
		SELECT id, title, original_title, cover, source, url, summary, author, status, 
		       genres, chapters_count, url_pattern, last_updated, date_added
		FROM novels
		ORDER BY title
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return models.ScanNovels(rows)
}

func (r *repo) GetNovelByID(id string) (*models.Novel, error) {
	query := `
		SELECT id, title, original_title, cover, source, url, summary, author, status, 
		       genres, chapters_count, url_pattern, last_updated, date_added
		FROM novels
		WHERE id = ?
	`

	row := r.db.QueryRow(query, id)
	return models.ScanNovel(row)
}

func (r *repo) CreateNovel(novel *models.Novel) (*models.Novel, error) {
	// Generate a new UUID if not provided
	if novel.ID == "" {
		novel.ID = uuid.New().String()
	}

	// Convert genres to JSON
	genresJSON, err := models.GenresToJSON(novel.Genres)
	if err != nil {
		return nil, err
	}

	query := `
		INSERT INTO novels (
			id, title, original_title, cover, source, url, summary, author, status, 
			genres, chapters_count, url_pattern, last_updated, date_added
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
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
		novel.URLPattern,
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
		    url_pattern = ?, last_updated = ?
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
		novel.URLPattern,
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

// Chapter CRUD Operations

func (r *repo) GetNovelChapters(novelID string) ([]*models.Chapter, error) {
	query := `
		SELECT id, novel_id, number, title, original_title, content, date_translated, word_count, url
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
		SELECT id, novel_id, number, title, original_title, content, date_translated, word_count, url
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
		SELECT id, novel_id, number, title, original_title, content, date_translated, word_count, url
		FROM chapters
		WHERE novel_id = ? AND id = ?
	`

	row := r.db.QueryRow(query, novelID, chapterID)
	return models.ScanChapter(row)
}

func (r *repo) GetChapterByNumber(novelID string, chapterNumber int) (*models.Chapter, error) {
	query := `
		SELECT id, novel_id, number, title, original_title, content, date_translated, word_count, url
		FROM chapters
		WHERE novel_id = ? AND number = ?
	`

	row := r.db.QueryRow(query, novelID, chapterNumber)
	return models.ScanChapter(row)
}

func (r *repo) CreateChapter(chapter *models.Chapter) (*models.Chapter, error) {
	// Generate a new UUID if not provided
	if chapter.ID == "" {
		chapter.ID = uuid.New().String()
	}

	query := `
		INSERT INTO chapters (
			id, novel_id, number, title, original_title, content, date_translated, word_count, url
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
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
	)
	if err != nil {
		return nil, err
	}

	// Update the novel's chapters count
	_, err = r.db.Exec(
		"UPDATE novels SET chapters_count = (SELECT COUNT(*) FROM chapters WHERE novel_id = ?), last_updated = ? WHERE id = ?",
		chapter.NovelID,
		time.Now().Unix(),
		chapter.NovelID,
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
		    date_translated = ?, word_count = ?, url = ?
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

	// Update the novel's last_updated timestamp
	_, err = r.db.Exec(
		"UPDATE novels SET last_updated = ? WHERE id = ?",
		time.Now().Unix(),
		chapter.NovelID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *repo) DeleteChapter(id string) error {
	// Get the novel ID before deleting the chapter
	var novelID string
	err := r.db.QueryRow("SELECT novel_id FROM chapters WHERE id = ?", id).Scan(&novelID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("chapter not found")
		}
		return err
	}

	// Delete the chapter
	result, err := r.db.Exec("DELETE FROM chapters WHERE id = ?", id)
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

	// Update the novel's chapters count and last_updated timestamp
	_, err = r.db.Exec(
		"UPDATE novels SET chapters_count = (SELECT COUNT(*) FROM chapters WHERE novel_id = ?), last_updated = ? WHERE id = ?",
		novelID,
		time.Now().Unix(),
		novelID,
	)
	if err != nil {
		return err
	}

	return nil
}
