package service

import (
	"context"
	"errors"
	"time"

	"backend/models"
	"backend/provider/gemini"
	"backend/provider/sources"
	"backend/provider/webscraper"
	"backend/repo"
	"backend/utils"
)

// TranslationService provides business logic for novel operations
type TranslationService interface {
	ExtractNovelDetails(ctx context.Context, request *models.NovelExtractionRequest) (*models.Novel, error)
	TranslateChapter(ctx context.Context, request *models.ChapterTranslationRequest) (*models.Chapter, error)
	TranslateFirstChapter(ctx context.Context, request *models.ChapterTranslationRequest) (*models.Chapter, error)
	RefreshNovel(ctx context.Context, request *models.NovelRefreshRequest) (*models.Novel, error)
}

type translationService struct {
	repo repo.Repo
}

var translationServiceInstance TranslationService

func init() {
	translationServiceInstance = NewTranslationService(repo.GetRepo())
}

// NewTranslationService creates a new novel service
func NewTranslationService(r repo.Repo) TranslationService {
	return &translationService{
		repo: r,
	}
}

// GetTranslationService returns the novel service instance
func GetTranslationService() TranslationService {
	return translationServiceInstance
}

func (s *translationService) ExtractNovelDetails(ctx context.Context, request *models.NovelExtractionRequest) (*models.Novel, error) {
	if request == nil {
		return nil, errors.New("request cannot be nil")
	}

	// Validate required fields
	if request.URL == "" {
		return nil, errors.New("URL cannot be empty")
	}
	if request.Source == "" {
		return nil, errors.New("source cannot be empty")
	}

	// Scrape the webpage content
	webpageContent, err := webscraper.GetScraperService().ScrapeWebPage(request.URL)
	if err != nil {
		return nil, err
	}

	// Translate the novel details
	novelDetails, err := gemini.GetClient().TranslateNovelDetails(ctx, webpageContent)
	if err != nil {
		return nil, err
	}

	// Create a new novel entry in the database
	newNovel := &models.Novel{
		ID:            sources.GetSource(request.Source).GetNovelId(request.URL),
		Title:         novelDetails.NovelTitleTranslated,
		OriginalTitle: novelDetails.NovelTitleOriginal,
		Cover:         "",
		Source:        request.Source,
		URL:           request.URL,
		Summary:       novelDetails.NovelSummaryTranslated,
		Author:        novelDetails.NovelAuthorNameTranslated,
		Status:        novelDetails.Status,
		Genres:        novelDetails.PossibleNovelGenres,
		ChaptersCount: novelDetails.NumberOfChapters,
		LastUpdated:   time.Now().Unix(),
		DateAdded:     time.Now().Unix(),
	}

	return s.repo.CreateNovel(newNovel)
}

func (s *translationService) TranslateFirstChapter(ctx context.Context, request *models.ChapterTranslationRequest) (*models.Chapter, error) {
	if request == nil {
		return nil, errors.New("request cannot be nil")
	}

	// Validate required fields
	if request.NovelID == "" {
		return nil, errors.New("novel ID cannot be empty")
	}

	// Get the novel by ID to ensure it exists
	novel, err := s.repo.GetNovelByID(request.NovelID)
	if err != nil {
		return nil, err
	}

	// Scrape the chapter content
	chapterContent, err := webscraper.GetScraperService().ScrapeWebPage(request.ChapterURL)
	if err != nil {
		return nil, err
	}

	// Translate the chapter content
	translatedContent, err := gemini.GetClient().TranslateNovelChapter(ctx, novel.Genres, chapterContent)
	if err != nil {
		return nil, err
	}

	novel.Genres = append(novel.Genres, translatedContent.PossibleNewGenres...)
	novel.Genres = utils.RemoveDuplicatesFromSlice(novel.Genres)
	if err = s.repo.UpdateNovel(novel); err != nil {
		return nil, err
	}

	// Create a new chapter entry
	chapter := &models.Chapter{
		ID:             sources.GetSource(novel.Source).GetChapterId(request.ChapterURL),
		NovelID:        novel.ID,
		Number:         1,
		Title:          translatedContent.TranslatedChapterTitle,
		OriginalTitle:  translatedContent.OriginalChapterTitle,
		Content:        translatedContent.TranslatedChapterContents,
		DateTranslated: time.Now().Unix(),
		WordCount:      utils.CountWords(translatedContent.TranslatedChapterContents),
		URL:            request.ChapterURL,
	}

	return s.repo.CreateChapter(chapter)
}

func (s *translationService) TranslateChapter(ctx context.Context, request *models.ChapterTranslationRequest) (*models.Chapter, error) {
	if request == nil {
		return nil, errors.New("request cannot be nil")
	}

	// Validate required fields
	if request.NovelID == "" {
		return nil, errors.New("novel ID cannot be empty")
	}

	// Get the novel by ID to ensure it exists
	novel, err := s.repo.GetNovelByID(request.NovelID)
	if err != nil {
		return nil, err
	}

	lastChapter, err := s.repo.GetLastChapter(novel.ID)
	if err != nil {
		return nil, err
	}

	nextChapterUrl, err := sources.GetSource(novel.Source).GetChapterUrl(lastChapter.URL)
	if err != nil {
		return nil, err
	}

	chapterContent, err := webscraper.GetScraperService().ScrapeWebPage(nextChapterUrl)
	if err != nil {
		return nil, err
	}

	// Translate the chapter content
	translatedContent, err := gemini.GetClient().TranslateNovelChapter(ctx, novel.Genres, chapterContent)
	if err != nil {
		return nil, err
	}

	novel.Genres = append(novel.Genres, translatedContent.PossibleNewGenres...)
	novel.Genres = utils.RemoveDuplicatesFromSlice(novel.Genres)
	if err = s.repo.UpdateNovel(novel); err != nil {
		return nil, err
	}

	// Create a new chapter entry
	chapter := &models.Chapter{
		ID:             sources.GetSource(novel.Source).GetChapterId(nextChapterUrl),
		NovelID:        novel.ID,
		Number:         lastChapter.Number + 1,
		Title:          translatedContent.TranslatedChapterTitle,
		OriginalTitle:  translatedContent.OriginalChapterTitle,
		Content:        translatedContent.TranslatedChapterContents,
		DateTranslated: time.Now().Unix(),
		URL:            nextChapterUrl,
	}

	return s.repo.CreateChapter(chapter)
}

func (s *translationService) RefreshNovel(ctx context.Context, request *models.NovelRefreshRequest) (*models.Novel, error) {
	if request.NovelID == "" {
		return nil, errors.New("novel ID cannot be empty")
	}

	// Get the novel by ID to ensure it exists
	novel, err := s.repo.GetNovelByID(request.NovelID)
	if err != nil {
		return nil, err
	}

	// Scrape the webpage content for the novel
	webpageContent, err := webscraper.GetScraperService().ScrapeWebPage(novel.URL)
	if err != nil {
		return nil, err
	}

	// Translate the novel details
	novelDetails, err := gemini.GetClient().TranslateNovelDetails(ctx, webpageContent)
	if err != nil {
		return nil, err
	}

	// Update the novel details in the database
	novel.Title = novelDetails.NovelTitleTranslated
	novel.OriginalTitle = novelDetails.NovelTitleOriginal
	novel.Summary = novelDetails.NovelSummaryTranslated
	novel.Author = novelDetails.NovelAuthorNameTranslated
	novel.Status = novelDetails.Status
	novel.ChaptersCount = novelDetails.NumberOfChapters
	novel.LastUpdated = time.Now().Unix()

	if err = s.repo.UpdateNovel(novel); err != nil {
		return nil, err
	}

	return s.repo.GetNovelByID(request.NovelID)
}
