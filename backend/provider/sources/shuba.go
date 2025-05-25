package sources

import (
	"strconv"
	"strings"
)

type shuba struct{}

// NewShuba creates a new instance of the Shuba source
func NewShuba() Source {
	return &shuba{}
}

// GetNovelId extracts the novel ID from a Shuba URL
func (s *shuba) GetNovelId(url string) string {
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
func (s *shuba) GetChapterId(chapterUrl string) string {
	// Example URL: https://www.69shuba.com/txt/36573/40179619
	// We need to extract the last part of the URL

	parts := strings.Split(chapterUrl, "/")
	if len(parts) < 6 {
		return ""
	}

	lastPart := parts[len(parts)-1]
	return lastPart // Return the last part as the chapter ID
}

// GetChapterUrl constructs the chapter URL based on the previous chapter URL
func (s *shuba) GetChapterUrl(previousChapterUrl string) (string, error) {
	// Example URL: https://www.69shuba.com/txt/36573/40179619
	// We have to add 1 to the last part of the URL to get the next chapter

	parts := strings.Split(previousChapterUrl, "/")
	if len(parts) < 6 {
		return "", nil // Invalid URL format
	}

	lastPart := parts[len(parts)-1]

	chapterNumber, err := strconv.ParseInt(lastPart, 10, 64)
	if err != nil {
		return "", err
	}
	chapterNumber++ // Increment the chapter number

	newChapterUrl := strings.Join(parts[:len(parts)-1], "/") + "/" + strconv.FormatInt(chapterNumber, 10)
	return newChapterUrl, nil
}
