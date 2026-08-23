package sources

import (
	"fmt"
	"strconv"
	"strings"
)

type syosetu struct{}

// NewSyosetu creates a new instance of the Syosetu source
func NewSyosetu() Source {
	return &syosetu{}
}

func (s syosetu) GetNovelId(url string) string {
	// Example URL: https://ncode.syosetu.com/n1514kj/
	parts := strings.Split(url, "/")
	for i, part := range parts {
		if part == "ncode.syosetu.com" && i+1 < len(parts) {
			// Return the part after ncode.syosetu.com, removing any trailing slashes
			return strings.TrimSuffix(parts[i+1], "/")
		}
	}
	return ""
}

func (s syosetu) GetChapterId(chapterUrl string) string {
	// Example URL: https://ncode.syosetu.com/n1514kj/1/
	parts := strings.Split(chapterUrl, "/")
	for i, part := range parts {
		if part == "ncode.syosetu.com" && i+2 < len(parts) {
			// Return the part after ncode.syosetu.com, removing any trailing slashes
			return s.GetNovelId(chapterUrl) + strings.TrimSuffix(parts[i+2], "/")
		}
	}
	return ""
}

func (s syosetu) GetNextChapterUrl(chapterContent, currentChapterUrl string) (string, error) {
	// Example currentChapterUrl: https://ncode.syosetu.com/n1514kj/1/
	parts := strings.Split(currentChapterUrl, "/")
	for i, part := range parts {
		if part == "ncode.syosetu.com" && i+2 < len(parts) {
			// Get the current chapter number
			chapterNum, err := strconv.Atoi(strings.TrimSuffix(parts[i+2], "/"))
			if err != nil {
				return "", fmt.Errorf("invalid chapter number in URL: %v", err)
			}

			// Construct next chapter URL
			nextChapter := strconv.Itoa(chapterNum + 1)
			return fmt.Sprintf("https://ncode.syosetu.com/%s/%s/",
				strings.TrimSuffix(parts[i+1], "/"),
				nextChapter), nil
		}
	}
	return "", nil
}

func (s syosetu) GetNovelCoverImageUrl(pageContent string) (string, error) {
	return "", nil
}
