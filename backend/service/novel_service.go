package service

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"backend/models"
	"backend/repo"
	"backend/utils"
)

// NovelService provides business logic for novel operations
type NovelService interface {
	GetAllNovels(offset, limit int) (*models.NovelListResponse, error)
	GetNovelByID(id string) (*models.Novel, error)
	SearchNovel(query string) ([]*models.Novel, error)
	GetNovelsByFilter(filter, value string, offset, limit int) (*models.NovelListResponse, error)
	UpdateNovel(novel *models.Novel) error
	DeleteNovel(id string) error

	GetNovelChapters(novelID string) ([]*models.Chapter, error)
	GetChapterByID(novelID string, chapterID string) (*models.Chapter, error)
	GetChapterByNumber(novelID string, chapterNumber int) (*models.Chapter, error)
	CreateChapter(chapter *models.Chapter) (*models.Chapter, error)
	UpdateChapter(chapter *models.Chapter) error
	DeleteChapter(novelID string, chapterID string) error

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

func (s *novelService) GetAllNovels(offset, limit int) (*models.NovelListResponse, error) {
	if offset < 0 {
		offset = 0
	}
	if limit <= 0 {
		limit = 20 // Default limit
	} else if limit > 100 {
		limit = 100 // Max limit to prevent excessive data transfer
	}

	novels, totalCount, err := s.repo.GetAllNovels(offset, limit)
	if err != nil {
		return nil, err
	}

	return &models.NovelListResponse{
		Novels:      novels,
		TotalCount:  totalCount,
		CurrentPage: offset/limit + 1,
		TotalPages:  (totalCount + limit - 1) / limit,
	}, nil
}

func (s *novelService) GetNovelByID(id string) (*models.Novel, error) {
	if id == "" {
		return nil, errors.New("novel ID cannot be empty")
	}

	return s.repo.GetNovelByID(id)
}

func (s *novelService) SearchNovel(query string) ([]*models.Novel, error) {
	query = strings.TrimSpace(query)
	if query == "" || len(query) < 3 {
		return nil, errors.New("search query cannot be empty or less than 3 characters")
	}

	return s.repo.SearchNovel(query)
}

func (s *novelService) GetNovelsByFilter(filter, value string, offset, limit int) (*models.NovelListResponse, error) {
	switch filter {
	case "language":
		sources, err := s.GetAllSources()
		if err != nil {
			return nil, err
		}

		var sourceIDs []string
		for _, source := range sources {
			if source.Language == value {
				sourceIDs = append(sourceIDs, source.ID)
			}
		}

		novels, totalCount, err := s.repo.GetNovelsBySourceIDs(sourceIDs, offset, limit)
		if err != nil {
		}

		return &models.NovelListResponse{
			Novels:      novels,
			TotalCount:  totalCount,
			CurrentPage: offset/limit + 1,
			TotalPages:  (totalCount + limit - 1) / limit,
		}, nil
	case "genre":
		novels, totalCount, err := s.repo.GetNovelsByGenre(value, offset, limit)
		if err != nil {
			return nil, err
		}

		return &models.NovelListResponse{
			Novels:      novels,
			TotalCount:  totalCount,
			CurrentPage: offset/limit + 1,
			TotalPages:  (totalCount + limit - 1) / limit,
		}, nil
	case "recently_updated":
		count, err := strconv.Atoi(value)
		if err != nil {
			return nil, errors.New("invalid count")
		}
		if count > 100 {
			return nil, errors.New("count cannot be greater than 100")
		}
		novels, err := s.repo.GetNovelsByRecentlyUpdated(count)
		if err != nil {
			return nil, err
		}

		return &models.NovelListResponse{
			Novels:      novels,
			TotalCount:  len(novels),
			CurrentPage: 1,
			TotalPages:  1,
		}, nil
	case "recently_read":
		count, err := strconv.Atoi(value)
		if err != nil {
			return nil, errors.New("invalid count")
		}
		if count > 100 {
			return nil, errors.New("count cannot be greater than 100")
		}
		novels, err := s.repo.GetNovelsByRecentlyRead(count)
		if err != nil {
			return nil, err
		}

		return &models.NovelListResponse{
			Novels:      novels,
			TotalCount:  len(novels),
			CurrentPage: 1,
			TotalPages:  1,
		}, nil
	default:
		return nil, errors.New("invalid filter")
	}
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

	chapter, err := s.repo.GetChapterByID(novelID, chapterID)

	if chapter != nil {
		_ = s.UpdateLastReadChapter(novelID, chapter.Number)
	}

	return chapter, err
}

func (s *novelService) GetChapterByNumber(novelID string, chapterNumber int) (*models.Chapter, error) {
	if novelID == "" {
		return nil, errors.New("novel ID cannot be empty")
	}

	if chapterNumber <= 0 {
		return nil, errors.New("chapter number must be positive")
	}

	chapter, err := s.repo.GetChapterByNumber(novelID, chapterNumber)

	if chapter != nil {
		_ = s.UpdateLastReadChapter(novelID, chapter.Number)
	}

	return chapter, err
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

func (s *novelService) DeleteChapter(novelID string, chapterID string) error {
	if novelID == "" {
		return errors.New("novel ID cannot be empty")
	}

	if chapterID == "" {
		return errors.New("chapter ID cannot be empty")
	}

	// Check if novel exists
	_, err := s.repo.GetNovelByID(novelID)
	if err != nil {
		return err
	}

	// Check if chapter exists
	_, err = s.repo.GetChapterByID(novelID, chapterID)
	if err != nil {
		return err
	}

	return s.repo.DeleteChapter(novelID, chapterID)
}

func (s *novelService) UpdateLastReadChapter(novelID string, chapterNumber int) error {
	if novelID == "" {
		return errors.New("novel ID cannot be empty")
	}

	if chapterNumber <= 0 {
		return errors.New("chapter number must be positive")
	}

	return s.repo.UpdateLastReadChapter(novelID, chapterNumber)
}

// Source operations

var allSources = []*models.SourceSite{
	{
		ID:       "69shuba",
		Name:     "69shuba",
		URL:      "https://www.69shuba.com",
		Language: "chinese",
	},
	{
		ID:       "69yue",
		Name:     "69yue",
		URL:      "https://www.69yue.top",
		Language: "chinese",
	},
	{
		ID:       "shuhaige",
		Name:     "shuhaige",
		URL:      "https://m.shuhaige.net",
		Language: "chinese",
	},
	{
		ID:       "twkan",
		Name:     "twkan",
		URL:      "https://twkan.com",
		Language: "chinese",
	},
	{
		ID:       "doupo",
		Name:     "doupo",
		URL:      "https://doupo.935666.xyz",
		Language: "chinese",
	},
	{
		ID:       "syosetu",
		Name:     "syosetu",
		URL:      "https://syosetu.com/",
		Language: "japanese",
	},
}

func (s *novelService) GetAllSources() ([]*models.SourceSite, error) {
	return allSources, nil
}
