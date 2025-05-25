package webscraper

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"time"
)

// ScraperService provides methods for scraping web content
type ScraperService interface {
	ScrapeWebPage(url string) (string, error)
}

type scraperService struct {
	scraper *WebScraper
}

var scraperServiceInstance ScraperService

func init() {
	scraper, err := NewWebScraper()
	if err != nil {
		panic(fmt.Sprintf("Failed to create web scraper: %v", err))
	}

	scraperServiceInstance = NewScraperService(scraper)
}

// NewScraperService creates a new scraper service
func NewScraperService(scraper *WebScraper) ScraperService {
	return &scraperService{
		scraper: scraper,
	}
}

// GetScraperService returns the scraper service singleton
func GetScraperService() ScraperService {
	return scraperServiceInstance
}

// ScrapeWebPage scrapes the content of a web page
func (s *scraperService) ScrapeWebPage(url string) (string, error) {
	content, err := s.scraper.GetWithRetries(url, 3, 2*time.Second)
	if err != nil {
		return "", fmt.Errorf("failed to scrape web page: %w", err)
	}

	return string(content), nil
}

// Common user agents for popular browsers
var userAgents = []string{
	// Chrome
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 14_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36",
	// Firefox
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:124.0) Gecko/20100101 Firefox/124.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 14.5; rv:124.0) Gecko/20100101 Firefox/124.0",
	// Safari
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 14_5) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.4 Safari/605.1.15",
	// Edge
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36 Edg/125.0.0.0",
}

// WebScraper provides functionality to scrape web content
type WebScraper struct {
	client     *http.Client
	UserAgent  string
	Referer    string
	AcceptLang string
}

// NewWebScraper creates a new web scraper with default browser-like settings
func NewWebScraper() (*WebScraper, error) {
	// Create a cookie jar to handle cookies
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create cookie jar: %w", err)
	}

	// Initialize random seed
	rand.Seed(time.Now().UnixNano())

	// Create HTTP client with cookie support and reasonable timeouts
	client := &http.Client{
		Jar:     jar,
		Timeout: 30 * time.Second,
	}

	// Create scraper with random user agent
	scraper := &WebScraper{
		client:     client,
		UserAgent:  userAgents[rand.Intn(len(userAgents))],
		AcceptLang: "en-US,en;q=0.9",
	}

	return scraper, nil
}

// get performs a GET request with browser-like headers
func (s *WebScraper) get(url string) ([]byte, error) {
	return s.makeRequest("GET", url, nil, nil)
}

// makeRequest is a helper function to make HTTP requests with browser-like headers
func (s *WebScraper) makeRequest(method, urlStr string, body []byte, additionalHeaders map[string]string) ([]byte, error) {
	var bodyReader io.Reader
	if body != nil {
		bodyReader = bytes.NewReader(body)
	}

	// Create request
	req, err := http.NewRequest(method, urlStr, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set common browser headers
	req.Header.Set("User-Agent", s.UserAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", s.AcceptLang)
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Cache-Control", "max-age=0")

	// Set referer if specified
	if s.Referer != "" {
		req.Header.Set("Referer", s.Referer)
	}

	// Add any additional headers
	for key, value := range additionalHeaders {
		req.Header.Set(key, value)
	}

	// Perform request
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Update referer for subsequent requests
	s.Referer = urlStr

	// Check response status
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("request failed with status code: %d", resp.StatusCode)
	}

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return respBody, nil
}

// GetWithRetries performs a GET request with retries
func (s *WebScraper) GetWithRetries(url string, maxRetries int, retryDelay time.Duration) ([]byte, error) {
	var resp []byte
	var err error

	for i := 0; i <= maxRetries; i++ {
		resp, err = s.get(url)
		if err == nil {
			return resp, nil
		}

		if i < maxRetries {
			// Wait before retrying
			time.Sleep(retryDelay)
		}
	}

	return nil, fmt.Errorf("failed after %d retries: %w", maxRetries, err)
}
