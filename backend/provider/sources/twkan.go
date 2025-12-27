package sources

import (
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

type twkan struct{}

// NewTwkan creates a new instance of the Shuba source
func NewTwkan() Source {
	return &twkan{}
}

// GetNovelId extracts the novel ID from a Shuba URL
func (s *twkan) GetNovelId(url string) string {
	// Example URL: https://www.69shuba.com/book/36573.htm

	parts := strings.Split(url, "/")
	if len(parts) < 5 {
		return ""
	}

	lastParts := strings.Split(parts[len(parts)-1], ".")
	if len(lastParts) < 2 {
		return ""
	}

	return lastParts[0] // Return the part before the ".htm"
}

// GetChapterId extracts the chapter ID from a Shuba chapter URL
func (s *twkan) GetChapterId(chapterUrl string) string {
	// Example URL: https://www.69shuba.com/txt/36573/40179619
	// We need to extract the last part of the URL

	parts := strings.Split(chapterUrl, "/")
	if len(parts) < 6 {
		return ""
	}

	lastPart := parts[len(parts)-1]
	return lastPart // Return the last part as the chapter ID
}

// GetNextChapterUrl constructs the chapter URL based on the previous chapter URL
func (s *twkan) GetNextChapterUrl(chapterContent, currentChapterUrl string) (string, error) {
	doc, err := html.Parse(strings.NewReader(chapterContent))
	if err != nil {
		return "", err
	}

	nextChapterURL, found := findNextChapterLink(doc)
	if found {
		return nextChapterURL, nil
	}

	return "", fmt.Errorf("next chapter link not found")
}

func (s *twkan) GetNovelCoverImageUrl(pageContent string) (string, error) {
	doc, err := html.Parse(strings.NewReader(pageContent))
	if err != nil {
		return "", err
	}

	bookImageURL, found := findMetaOgImage(doc)
	if found {
		return bookImageURL, nil
	}

	return "", fmt.Errorf("book image link not found")
}
