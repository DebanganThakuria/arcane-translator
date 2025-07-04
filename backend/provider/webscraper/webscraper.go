package webscraper

import (
	"context"
	"fmt"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/go-rod/rod"
	"github.com/go-rod/stealth"
	"github.com/gocolly/colly"
)

// ScraperService provides methods for scraping web content
type ScraperService interface {
	ScrapeWebPage(url string) (string, error)
	ClosOpenCon()
}

var (
	collyScraper    ScraperService
	chromedpScraper ScraperService
	rodScraper      ScraperService
)

func GetScraperService() ScraperService {
	return collyScraper
}

// ------------------------- Colly Scraper Implementation -------------------------
type collyScraperService struct {
	collector *colly.Collector
}

func init() {
	c := colly.NewCollector(colly.AllowURLRevisit())
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36"
	collyScraper = &collyScraperService{c}
	//chromedpScraper = &chromedpScraperService{}
	//rodScraper = &rodScraperService{
	//	browser: rod.New().MustConnect(),
	//}
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

func (s *collyScraperService) ClosOpenCon() {}

// ------------------------- Chromedp Scraper Implementation -------------------------
type chromedpScraperService struct{}

func (s *chromedpScraperService) ScrapeWebPage(url string) (string, error) {
	var content string

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
		chromedp.Flag("disable-blink-features", "AutomationControlled"),
	)
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()
	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	// Increase timeout to allow Cloudflare challenge to complete
	ctx, cancel = context.WithTimeout(ctx, 40*time.Second)
	defer cancel()

	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitVisible("body", chromedp.ByQuery),
		chromedp.Sleep(10*time.Second),
		chromedp.OuterHTML("html", &content),
	)
	if err != nil {
		return "", err
	}
	return content, nil
}

func (s *chromedpScraperService) ClosOpenCon() {}

// ------------------------- Rod Scraper Implementation -------------------------
type rodScraperService struct {
	browser *rod.Browser
}

func (s *rodScraperService) ScrapeWebPage(url string) (string, error) {
	page := stealth.MustPage(s.browser)

	err := page.Navigate(url)
	if err != nil {
		return "", fmt.Errorf("failed to navigate to %s: %w", url, err)
	}

	// Wait until the JS challenge (if any) is done and page is fully loaded
	page.MustWaitLoad().MustWaitIdle()

	// Optional: wait a bit more to make sure everything is rendered
	time.Sleep(2 * time.Second)

	// Get page content
	html, err := page.HTML()
	if err != nil {
		return "", fmt.Errorf("failed to get page HTML: %w", err)
	}

	return html, nil
}

func (s *rodScraperService) ClosOpenCon() {
	if s.browser != nil {
		s.browser.MustClose()
	}
}
