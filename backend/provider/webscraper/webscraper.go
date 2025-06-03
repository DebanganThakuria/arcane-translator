package webscraper

import "github.com/gocolly/colly"

// ScraperService provides methods for scraping web content
type ScraperService interface {
	ScrapeWebPage(url string) (string, error)
}

type collyScraperService struct {
	collector *colly.Collector
}

var scraper *collyScraperService

func GetScraperService() ScraperService {
	return scraper
}

func init() {
	c := colly.NewCollector()
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36"
	scraper = &collyScraperService{c}
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
