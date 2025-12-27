package sources

import "strings"

type yue struct{}

func NewYue() Source {
	return &yue{}
}

func (y yue) GetNovelId(url string) string {
	// Example URL: https://www.69yue.top/articlecategroy/1xvh.html
	parts := strings.Split(url, "/")
	if len(parts) < 5 {
		return ""
	}

	lastParts := strings.Split(parts[len(parts)-1], ".")
	if len(lastParts) < 2 {
		return ""
	}

	return lastParts[0] // Return the part before the ".html"
}

func (y yue) GetChapterId(chapterUrl string) string {
	// Example URL: https://www.69yue.top/article/15189329271721020.html
	parts := strings.Split(chapterUrl, "/")
	if len(parts) < 5 {
		return ""
	}

	lastParts := strings.Split(parts[len(parts)-1], ".")
	if len(lastParts) < 2 {
		return ""
	}

	return lastParts[0] // Return the part before the ".html"
}

func (y yue) GetNextChapterUrl(chapterContent, currentChapterUrl string) (string, error) {
	startIdx := strings.Index(chapterContent, `下一章`)
	if startIdx == -1 {
		return "", nil
	}

	aTagStart := strings.LastIndex(chapterContent[:startIdx], `<a `)
	if aTagStart == -1 {
		return "", nil
	}

	aTagEnd := strings.Index(chapterContent[aTagStart:], `>`)
	if aTagEnd == -1 {
		return "", nil
	}
	aTagEnd += aTagStart

	hrefStart := strings.Index(chapterContent[aTagStart:aTagEnd], `href="`)
	if hrefStart == -1 {
		return "", nil
	}
	hrefStart += aTagStart + len(`href="`)

	hrefEnd := strings.Index(chapterContent[hrefStart:], `"`)
	if hrefEnd == -1 {
		return "", nil
	}
	hrefEnd += hrefStart

	nextChapterURL := chapterContent[hrefStart:hrefEnd]
	if strings.HasPrefix(nextChapterURL, "http") {
		return nextChapterURL, nil
	}

	// Ensure nextChapterURL starts with a slash
	if !strings.HasPrefix(nextChapterURL, "/") {
		nextChapterURL = "/" + nextChapterURL
	}

	baseDomain := "https://www.69yue.top"

	fullNextChapterURL := baseDomain + nextChapterURL
	return fullNextChapterURL, nil
}

func (y yue) GetNovelCoverImageUrl(pageContent string) (string, error) {
	return "", nil
}
