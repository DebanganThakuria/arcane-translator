package service

import (
	"backend/models"
	"backend/repo"
)

// StatsService provides business logic for statistics operations
type StatsService interface {
	GetNovelStats() (*models.Stats, error)
}

type statsService struct {
	repo repo.Repo
}

var statsServiceInstance StatsService

func init() {
	statsServiceInstance = NewStatsService(repo.GetRepo())
}

// NewStatsService creates a new instance of StatsService
func NewStatsService(repo repo.Repo) StatsService {
	return &statsService{
		repo: repo,
	}
}

// GetStatsService returns the singleton instance of StatsService
func GetStatsService() StatsService {
	return statsServiceInstance
}

// GetNovelStats returns the total count of novels and chapters
func (s *statsService) GetNovelStats() (*models.Stats, error) {
	novelCount, chapterCount, err := s.repo.GetStats()
	if err != nil {
		return nil, err
	}

	return &models.Stats{
		NovelCount:   novelCount,
		ChapterCount: chapterCount,
	}, nil
}
