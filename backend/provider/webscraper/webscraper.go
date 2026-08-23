package webscraper

import (
	"errors"
	"time"

	"github.com/gocolly/colly"
)

const (
	// Source sites are slow and occasionally hang. Without this a stuck request
	// occupies a handler until the server's own write timeout fires.
	requestTimeout = 25 * time.Second

	userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 " +
		"(KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36"
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

type collyScraperService struct{}

func init() {
	collyScraper = &collyScraperService{}
}

// newCollector builds a collector for a single scrape.
//
// A collector must not be shared: colly *appends* response callbacks rather
// than replacing them, so reusing one package-level collector meant every
// scrape added another handler that stayed registered forever. That leaked
// memory, made each request do more work than the last, and let a later
// response overwrite an earlier caller's buffer.
func newCollector() *colly.Collector {
	c := colly.NewCollector(colly.AllowURLRevisit(), colly.DetectCharset())
	c.UserAgent = userAgent
	c.SetRequestTimeout(requestTimeout)
	return c
}

func (s *collyScraperService) ScrapeWebPage(url string) (string, error) {
	if url == "" {
		return "", errors.New("url cannot be empty")
	}

	c := newCollector()

	var content []byte
	c.OnResponse(func(r *colly.Response) {
		content = r.Body
	})

	if err := c.Visit(url); err != nil {
		return "", err
	}

	// Visit returns nil for a response that produced no body, so guard here
	// rather than handing the translator an empty document.
	if len(content) == 0 {
		return "", errors.New("source returned an empty page")
	}

	return string(content), nil
}

func (s *collyScraperService) ScrapeWebPageWithContent(url string, htmlContent string) (string, error) {
	// The caller already fetched the page in their own browser.
	if htmlContent != "" {
		return htmlContent, nil
	}

	return s.ScrapeWebPage(url)
}
