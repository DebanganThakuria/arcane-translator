package service

import (
	"errors"
	"time"

	"backend/models"
	"backend/repo"
	"backend/utils"
)

// NovelService provides business logic for novel operations
type NovelService interface {
	GetAllNovels() ([]*models.Novel, error)
	GetNovelByID(id string) (*models.Novel, error)
	CreateNovel(novel *models.Novel) (*models.Novel, error)
	UpdateNovel(novel *models.Novel) error
	DeleteNovel(id string) error

	GetNovelChapters(novelID string) ([]*models.Chapter, error)
	GetChapterByID(novelID string, chapterID string) (*models.Chapter, error)
	GetChapterByNumber(novelID string, chapterNumber int) (*models.Chapter, error)
	CreateChapter(chapter *models.Chapter) (*models.Chapter, error)
	UpdateChapter(chapter *models.Chapter) error
	DeleteChapter(id string) error

	GetAllSources() ([]*models.SourceSite, error)
}

type novelService struct {
	repo repo.Repo
}

var novelServiceInstance NovelService

func init() {
	novelServiceInstance = NewNovelService(repo.GetRepo())
}

// NewNovelService creates a new novel service
func NewNovelService(r repo.Repo) NovelService {
	return &novelService{
		repo: r,
	}
}

// GetNovelService returns the novel service instance
func GetNovelService() NovelService {
	return novelServiceInstance
}

// Novel operations

func (s *novelService) GetAllNovels() ([]*models.Novel, error) {
	return s.repo.GetAllNovels()
}

func (s *novelService) GetNovelByID(id string) (*models.Novel, error) {
	if id == "" {
		return nil, errors.New("novel ID cannot be empty")
	}

	return s.repo.GetNovelByID(id)
}

func (s *novelService) CreateNovel(novel *models.Novel) (*models.Novel, error) {
	// Validate required fields
	if novel.Title == "" {
		return nil, errors.New("novel title cannot be empty")
	}

	if novel.Source == "" {
		novel.Source = "unknown"
	}

	if novel.URL == "" {
		novel.URL = ""
	}

	if novel.Summary == "" {
		novel.Summary = "No summary available"
	}

	// Set timestamps
	novel.DateAdded = time.Now().Unix()
	novel.LastUpdated = time.Now().Unix()

	// Set default values
	if novel.Status == "" {
		novel.Status = "Unknown"
	}

	return s.repo.CreateNovel(novel)
}

func (s *novelService) UpdateNovel(novel *models.Novel) error {
	if novel.ID == "" {
		return errors.New("novel ID cannot be empty")
	}

	// Check if novel exists
	_, err := s.repo.GetNovelByID(novel.ID)
	if err != nil {
		return err
	}

	// Validate required fields
	if novel.Title == "" {
		return errors.New("novel title cannot be empty")
	}

	return s.repo.UpdateNovel(novel)
}

func (s *novelService) DeleteNovel(id string) error {
	if id == "" {
		return errors.New("novel ID cannot be empty")
	}

	return s.repo.DeleteNovel(id)
}

// Chapter operations

func (s *novelService) GetNovelChapters(novelID string) ([]*models.Chapter, error) {
	if novelID == "" {
		return nil, errors.New("novel ID cannot be empty")
	}

	// Check if novel exists
	_, err := s.repo.GetNovelByID(novelID)
	if err != nil {
		return nil, err
	}

	return s.repo.GetNovelChapters(novelID)
}

func (s *novelService) GetChapterByID(novelID string, chapterID string) (*models.Chapter, error) {
	if novelID == "" {
		return nil, errors.New("novel ID cannot be empty")
	}

	if chapterID == "" {
		return nil, errors.New("chapter ID cannot be empty")
	}

	return s.repo.GetChapterByID(novelID, chapterID)
}

func (s *novelService) GetChapterByNumber(novelID string, chapterNumber int) (*models.Chapter, error) {
	if novelID == "" {
		return nil, errors.New("novel ID cannot be empty")
	}

	if chapterNumber <= 0 {
		return nil, errors.New("chapter number must be positive")
	}

	return s.repo.GetChapterByNumber(novelID, chapterNumber)
}

func (s *novelService) CreateChapter(chapter *models.Chapter) (*models.Chapter, error) {
	if chapter.NovelID == "" {
		return nil, errors.New("novel ID cannot be empty")
	}

	if chapter.Title == "" {
		return nil, errors.New("chapter title cannot be empty")
	}

	if chapter.Content == "" {
		return nil, errors.New("chapter content cannot be empty")
	}

	if chapter.Number <= 0 {
		return nil, errors.New("chapter number must be positive")
	}

	// Check if novel exists
	_, err := s.repo.GetNovelByID(chapter.NovelID)
	if err != nil {
		return nil, err
	}

	// Check if chapter with the same number already exists
	existingChapter, err := s.repo.GetChapterByNumber(chapter.NovelID, chapter.Number)
	if err == nil && existingChapter != nil {
		return nil, errors.New("chapter with this number already exists")
	}

	// Set default values
	chapter.DateTranslated = time.Now().Unix()

	// Calculate word count if not provided
	if chapter.WordCount <= 0 {
		chapter.WordCount = utils.CountWords(chapter.Content)
	}

	return s.repo.CreateChapter(chapter)
}

func (s *novelService) UpdateChapter(chapter *models.Chapter) error {
	if chapter.ID == "" {
		return errors.New("chapter ID cannot be empty")
	}

	if chapter.NovelID == "" {
		return errors.New("novel ID cannot be empty")
	}

	if chapter.Title == "" {
		return errors.New("chapter title cannot be empty")
	}

	if chapter.Content == "" {
		return errors.New("chapter content cannot be empty")
	}

	if chapter.Number <= 0 {
		return errors.New("chapter number must be positive")
	}

	// Check if novel exists
	_, err := s.repo.GetNovelByID(chapter.NovelID)
	if err != nil {
		return err
	}

	// Check if chapter exists
	existingChapter, err := s.repo.GetChapterByID(chapter.NovelID, chapter.ID)
	if err != nil {
		return err
	}

	// If changing chapter number, check if the new number is already taken
	if existingChapter.Number != chapter.Number {
		checkChapter, err := s.repo.GetChapterByNumber(chapter.NovelID, chapter.Number)
		if err == nil && checkChapter != nil && checkChapter.ID != chapter.ID {
			return errors.New("chapter with this number already exists")
		}
	}

	// Calculate word count if not provided
	if chapter.WordCount <= 0 {
		chapter.WordCount = utils.CountWords(chapter.Content)
	}

	return s.repo.UpdateChapter(chapter)
}

func (s *novelService) DeleteChapter(id string) error {
	if id == "" {
		return errors.New("chapter ID cannot be empty")
	}

	return s.repo.DeleteChapter(id)
}

// Source operations

var allSources = []*models.SourceSite{
	{
		ID:       "69shuba",
		Name:     "69shuba",
		URL:      "https://www.69shuba.com",
		Language: "Chinese",
	},
}

func (s *novelService) GetAllSources() ([]*models.SourceSite, error) {
	return allSources, nil
}
