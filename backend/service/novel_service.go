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
	QueryNovels(query models.NovelQuery, language string) (*models.NovelListResponse, error)
	GetAllNovels(offset, limit int) (*models.NovelListResponse, error)
	GetNovelByID(id string) (*models.Novel, error)
	SearchNovel(query string) ([]*models.Novel, error)
	GetNovelsByFilter(filter, value string, offset, limit int) (*models.NovelListResponse, error)
	UpdateNovel(novel *models.Novel) error
	DeleteNovel(id string) error

	QueryChapters(novelID string, query models.ChapterQuery) (*models.ChapterListResponse, error)
	GetNovelChapters(novelID string) ([]*models.Chapter, error)
	GetChapterByID(novelID string, chapterID string) (*models.Chapter, error)
	GetChapterByNumber(novelID string, chapterNumber int) (*models.Chapter, error)
	CreateChapter(chapter *models.Chapter) (*models.Chapter, error)
	UpdateChapter(chapter *models.Chapter) error
	DeleteChapter(novelID string, chapterID string) error

	GetAllSources() ([]*models.SourceSite, error)
}

const (
	defaultPageSize = 20
	maxPageSize     = 100

	// Chapter rows are small once content is excluded, so the page can be
	// larger than the novel page without costing much.
	defaultChapterPageSize = 50
	maxChapterPageSize     = 200
)

type novelService struct {
	repo repo.Repo
}

// page wraps a result set with the paging metadata the API returns.
func page(novels []*models.Novel, totalCount, offset, limit int) *models.NovelListResponse {
	if novels == nil {
		novels = []*models.Novel{}
	}
	if limit <= 0 {
		limit = defaultPageSize
	}

	return &models.NovelListResponse{
		Novels:      novels,
		TotalCount:  totalCount,
		CurrentPage: offset/limit + 1,
		TotalPages:  max(1, (totalCount+limit-1)/limit),
	}
}

func emptyPage(limit int) *models.NovelListResponse {
	return page(nil, 0, 0, limit)
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

// sourceIDsForLanguage resolves a language name onto the sources that serve it.
// An unknown language yields no sources, which correctly produces no results
// rather than silently listing everything.
func (s *novelService) sourceIDsForLanguage(language string) ([]string, error) {
	sources, err := s.GetAllSources()
	if err != nil {
		return nil, err
	}

	ids := make([]string, 0, len(sources))
	for _, source := range sources {
		if strings.EqualFold(source.Language, language) {
			ids = append(ids, source.ID)
		}
	}
	return ids, nil
}

// QueryNovels applies filters, sorting and paging in the database. `language`
// is resolved to source ids here because the source list lives in this layer.
func (s *novelService) QueryNovels(query models.NovelQuery, language string) (*models.NovelListResponse, error) {
	if query.Offset < 0 {
		query.Offset = 0
	}
	if query.Limit <= 0 {
		query.Limit = defaultPageSize
	} else if query.Limit > maxPageSize {
		query.Limit = maxPageSize
	}

	if language != "" {
		ids, err := s.sourceIDsForLanguage(language)
		if err != nil {
			return nil, err
		}
		if len(ids) == 0 {
			return emptyPage(query.Limit), nil
		}
		query.SourceIDs = ids
	}

	novels, totalCount, err := s.repo.QueryNovels(query)
	if err != nil {
		return nil, err
	}

	return page(novels, totalCount, query.Offset, query.Limit), nil
}

func (s *novelService) GetAllNovels(offset, limit int) (*models.NovelListResponse, error) {
	return s.QueryNovels(models.NovelQuery{
		Sort:   models.SortLastRead,
		Offset: offset,
		Limit:  limit,
	}, "")
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

// GetNovelsByFilter is the older single-filter entry point, kept because the
// home page still uses the recently_* shortcuts. Language and genre delegate to
// QueryNovels so there is one query path.
func (s *novelService) GetNovelsByFilter(filter, value string, offset, limit int) (*models.NovelListResponse, error) {
	switch filter {
	case "language":
		return s.QueryNovels(models.NovelQuery{
			Sort:   models.SortLastRead,
			Offset: offset,
			Limit:  limit,
		}, value)
	case "genre":
		return s.QueryNovels(models.NovelQuery{
			Genre:  value,
			Sort:   models.SortLastRead,
			Offset: offset,
			Limit:  limit,
		}, "")
	case "recently_updated":
		count, err := strconv.Atoi(value)
		if err != nil {
			return nil, errors.New("invalid count")
		}
		if count <= 0 || count > maxPageSize {
			return nil, errors.New("count must be between 1 and 100")
		}
		novels, err := s.repo.GetNovelsByRecentlyUpdated(count)
		if err != nil {
			return nil, err
		}

		return page(novels, len(novels), 0, count), nil
	case "recently_read":
		count, err := strconv.Atoi(value)
		if err != nil {
			return nil, errors.New("invalid count")
		}
		if count <= 0 || count > maxPageSize {
			return nil, errors.New("count must be between 1 and 100")
		}
		novels, err := s.repo.GetNovelsByRecentlyRead(count)
		if err != nil {
			return nil, err
		}

		return page(novels, len(novels), 0, count), nil
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

// QueryChapters returns one page of a novel's chapters. The novel is not
// re-fetched first: the query is already scoped by novel_id, so a missing novel
// simply yields an empty page instead of costing an extra round trip.
func (s *novelService) QueryChapters(novelID string, query models.ChapterQuery) (*models.ChapterListResponse, error) {
	if novelID == "" {
		return nil, errors.New("novel ID cannot be empty")
	}

	if query.Offset < 0 {
		query.Offset = 0
	}
	if query.Limit <= 0 {
		query.Limit = defaultChapterPageSize
	} else if query.Limit > maxChapterPageSize {
		query.Limit = maxChapterPageSize
	}

	chapters, totalCount, err := s.repo.QueryChapters(novelID, query)
	if err != nil {
		return nil, err
	}
	if chapters == nil {
		chapters = []*models.Chapter{}
	}

	first, last, err := s.repo.ChapterNumberBounds(novelID)
	if err != nil {
		return nil, err
	}

	return &models.ChapterListResponse{
		Chapters:    chapters,
		TotalCount:  totalCount,
		CurrentPage: query.Offset/query.Limit + 1,
		TotalPages:  max(1, (totalCount+query.Limit-1)/query.Limit),
		FirstNumber: first,
		LastNumber:  last,
	}, nil
}

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
		ID:       "ixdzs",
		Name:     "ixdzs",
		URL:      "https://ixdzs.tw",
		Language: "chinese",
	},
	{
		ID:       "czbooks",
		Name:     "czbooks",
		URL:      "https://czbooks.net",
		Language: "chinese",
	},
	{
		ID:       "quanben",
		Name:     "quanben",
		URL:      "https://www.quanben.io",
		Language: "chinese",
	},
	{
		ID:       "sjks88",
		Name:     "sjks88",
		URL:      "https://www.sjks88.com",
		Language: "chinese",
	},
	{
		ID:       "huabenge",
		Name:     "huabenge",
		URL:      "https://www.huabenge.com",
		Language: "chinese",
	},
	{
		ID:       "ilwxs",
		Name:     "ilwxs",
		URL:      "https://m.ilwxs.com",
		Language: "chinese",
	},
	{
		ID:       "ffxs8",
		Name:     "ffxs8",
		URL:      "https://ffxs8.com",
		Language: "chinese",
	},
	{
		ID:       "syosetu",
		Name:     "syosetu",
		URL:      "https://syosetu.com/",
		Language: "japanese",
	},
	{
		ID:       "scribblehub",
		Name:     "scribblehub",
		URL:      "https://www.scribblehub.com",
		Language: "english",
	},
	{
		ID:       "royalroad",
		Name:     "royalroad",
		URL:      "https://www.royalroad.com",
		Language: "english",
	},
}

func (s *novelService) GetAllSources() ([]*models.SourceSite, error) {
	return allSources, nil
}
