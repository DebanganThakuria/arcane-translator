package sources

import (
	"strconv"
	"strings"
)

type sjks88 struct{}

func NewSjks88() *sjks88 {
	return &sjks88{}
}

func (s *sjks88) GetNovelId(url string) string {
	// Example URL: https://www.sjks88.com/xuanhuan/51964.html
	parts := strings.Split(url, "/")
	if len(parts) < 5 {
		return ""
	}
	return strings.Split(parts[4], ".")[0]
}

func (s *sjks88) GetChapterId(chapterUrl string) string {
	// Example URL: https://www.sjks88.com/xuanhuan/51964/1.html
	parts := strings.Split(chapterUrl, "/")
	if len(parts) < 6 {
		return ""
	}
	novelId := parts[4]
	chapterNum := strings.Split(parts[5], ".")[0]
	return novelId + "_" + chapterNum
}

func (s *sjks88) GetNextChapterUrl(chapterContent, currentChapterUrl string) (string, error) {
	// URLs are sequential: /xuanhuan/51964/1.html -> /xuanhuan/51964/2.html
	parts := strings.Split(currentChapterUrl, "/")
	if len(parts) < 6 {
		return "", nil
	}
	currentNum, err := strconv.Atoi(strings.Split(parts[5], ".")[0])
	if err != nil {
		return "", nil
	}
	return strings.Join(parts[:5], "/") + "/" + strconv.Itoa(currentNum+1) + ".html", nil
}

func (s *sjks88) GetNovelCoverImageUrl(pageContent string) (string, error) {
	// Cover image is in <meta property="og:image" content="...">
	marker := `property="og:image"`
	idx := strings.Index(pageContent, marker)
	if idx == -1 {
		return "", nil
	}

	contentMarker := `content="`
	contentIdx := strings.Index(pageContent[idx:], contentMarker)
	if contentIdx == -1 {
		return "", nil
	}
	contentIdx += idx + len(contentMarker)

	endIdx := strings.Index(pageContent[contentIdx:], `"`)
	if endIdx == -1 {
		return "", nil
	}

	return pageContent[contentIdx : contentIdx+endIdx], nil
}
