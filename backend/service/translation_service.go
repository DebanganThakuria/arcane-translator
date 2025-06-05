package service

import (
	"context"
	"errors"
	"log"
	"time"

	"backend/provider/webscraper"

	"backend/models"
	"backend/provider/gemini"
	"backend/provider/sources"
	"backend/repo"
	"backend/utils"
)

// TranslationService provides business logic for novel operations
type TranslationService interface {
	ExtractNovelDetails(ctx context.Context, request *models.NovelExtractionRequest) (*models.Novel, error)
	TranslateChapter(ctx context.Context, request *models.ChapterTranslationRequest) (*models.Chapter, error)
	TranslateFirstChapter(ctx context.Context, request *models.ChapterTranslationRequest) (*models.Chapter, error)
	RefreshNovel(ctx context.Context, novelId string) (*models.Novel, error)
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

	success := utils.Mutex.TryLock("extractNovelDetails"+request.URL, time.Millisecond)
	if !success {
		return nil, errors.New("another request is in progress")
	}
	defer utils.Mutex.Unlock("extractNovelDetails" + request.URL)

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

	// Get cover image URL
	coverUrl, err := sources.GetSource(request.Source).GetNovelCoverImageUrl(webpageContent)
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
		Cover:         coverUrl,
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

	success := utils.Mutex.TryLock("translateChapter"+request.NovelID, time.Millisecond)
	if !success {
		return nil, errors.New("another request is in progress")
	}
	defer utils.Mutex.Unlock("translateChapter" + request.NovelID)

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

	// Get the next chapter URL
	nextChapterURL, err := sources.GetSource(novel.Source).GetNextChapterUrl(chapterContent)
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
		log.Printf("Failed to update novel genres: %v", err)
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
		NextChapterURL: nextChapterURL,
	}

	return s.repo.CreateChapter(chapter)
}

func (s *translationService) TranslateChapter(ctx context.Context, request *models.ChapterTranslationRequest) (*models.Chapter, error) {
	if request == nil {
		return nil, errors.New("request cannot be nil")
	}

	success := utils.Mutex.TryLock("translateChapter"+request.NovelID, time.Millisecond)
	if !success {
		return nil, errors.New("another request is in progress")
	}
	defer utils.Mutex.Unlock("translateChapter" + request.NovelID)

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

	chapterContent, err := webscraper.GetScraperService().ScrapeWebPage(lastChapter.NextChapterURL)
	if err != nil {
		return nil, err
	}

	// Get the next chapter URL
	nextChapterUrl, err := sources.GetSource(novel.Source).GetNextChapterUrl(chapterContent)
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
		log.Printf("Failed to update novel genres: %v", err)
	}

	// Create a new chapter entry
	chapter := &models.Chapter{
		ID:             sources.GetSource(novel.Source).GetChapterId(lastChapter.NextChapterURL),
		NovelID:        novel.ID,
		Number:         lastChapter.Number + 1,
		Title:          translatedContent.TranslatedChapterTitle,
		OriginalTitle:  translatedContent.OriginalChapterTitle,
		Content:        translatedContent.TranslatedChapterContents,
		DateTranslated: time.Now().Unix(),
		WordCount:      utils.CountWords(translatedContent.TranslatedChapterContents),
		URL:            lastChapter.NextChapterURL, // Last Chapter's next chapter URL is the current chapter's URL
		NextChapterURL: nextChapterUrl,
	}

	return s.repo.CreateChapter(chapter)
}

func (s *translationService) RefreshNovel(ctx context.Context, novelId string) (*models.Novel, error) {
	if novelId == "" {
		return nil, errors.New("novel ID cannot be empty")
	}

	success := utils.Mutex.TryLock("refreshNovel"+novelId, time.Millisecond)
	if !success {
		return nil, errors.New("another request is in progress")
	}
	defer utils.Mutex.Unlock("refreshNovel" + novelId)

	// Get the novel by ID to ensure it exists
	novel, err := s.repo.GetNovelByID(novelId)
	if err != nil {
		return nil, err
	}

	// Scrape the webpage content for the novel
	webpageContent, err := webscraper.GetScraperService().ScrapeWebPage(novel.URL)
	if err != nil {
		return nil, err
	}

	// Cover image URL
	coverUrl, err := sources.GetSource(novel.Source).GetNovelCoverImageUrl(webpageContent)
	if err != nil {
		return nil, err
	}

	// Translate the novel details
	novelDetails, err := gemini.GetClient().TranslateNovelDetails(ctx, webpageContent)
	if err != nil {
		return nil, err
	}

	// Add the next chapter URL to the last chapter if there are new chapters
	if novelDetails.NumberOfChapters > novel.ChaptersCount {
		if err = s.addNextChapterUrlToLastChapter(novel.ID, novel.Source); err != nil {
			log.Printf("Failed to add next chapter URL to last chapter: %v\n", err)
		}
	}

	// Update the novel details in the database
	novel.Title = novelDetails.NovelTitleTranslated
	novel.OriginalTitle = novelDetails.NovelTitleOriginal
	novel.Summary = novelDetails.NovelSummaryTranslated
	novel.Cover = coverUrl
	novel.Author = novelDetails.NovelAuthorNameTranslated
	novel.Status = novelDetails.Status
	novel.ChaptersCount = novelDetails.NumberOfChapters
	novel.LastUpdated = time.Now().Unix()

	if err = s.repo.UpdateNovel(novel); err != nil {
		return nil, err
	}

	return s.repo.GetNovelByID(novelId)
}

func (s *translationService) addNextChapterUrlToLastChapter(novelId, source string) error {
	lastChapter, err := s.repo.GetLastChapter(novelId)
	if err != nil {
		return err
	}

	webpageContent, err := webscraper.GetScraperService().ScrapeWebPage(lastChapter.URL)
	if err != nil {
		return err
	}

	nextChapterUrl, err := sources.GetSource(source).GetNextChapterUrl(webpageContent)
	if err != nil {
		return err
	}

	lastChapter.NextChapterURL = nextChapterUrl
	return s.repo.UpdateChapter(lastChapter)
}
