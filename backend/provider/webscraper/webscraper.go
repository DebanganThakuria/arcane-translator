package webscraper

import (
	"github.com/gocolly/colly"
)

// ScraperService provides methods for scraping web content
type ScraperService interface {
	ScrapeWebPage(url string) (string, error)
	ScrapeWebPageWithContent(url string, htmlContent string) (string, error)
}

var collyScraper ScraperService

func GetScraperService() ScraperService {
	return collyScraper
}

type collyScraperService struct {
	collector *colly.Collector
}

func init() {
	c := colly.NewCollector(colly.AllowURLRevisit())
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36"
	collyScraper = &collyScraperService{c}
}

func (s *collyScraperService) ScrapeWebPage(url string) (string, error) {
	var content []byte
	s.collector.OnResponse(func(r *colly.Response) {
		content = r.Body
	})

	err := s.collector.Visit(url)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func (s *collyScraperService) ScrapeWebPageWithContent(url string, htmlContent string) (string, error) {
	// For browser-based scraping, we already have the HTML content
	if htmlContent != "" {
		return htmlContent, nil
	}
	// Fallback to normal scraping
	return s.ScrapeWebPage(url)
}
