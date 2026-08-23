package sources

import "strings"

type shuhaige struct{}

func NewShuhaige() Source {
	return &shuhaige{}
}

func (s shuhaige) GetNovelId(url string) string {
	// Example URL: https://m.shuhaige.net/345462
	parts := strings.Split(url, "/")
	if len(parts) < 4 {
		return ""
	}

	return parts[3] // Return the part after the domain
}

func (s shuhaige) GetChapterId(chapterUrl string) string {
	// Example URL: https://m.shuhaige.net/345462/121268892.html
	parts := strings.Split(chapterUrl, "/")
	if len(parts) < 5 {
		return ""
	}

	lastParts := strings.Split(parts[4], ".")
	if len(lastParts) < 2 {
		return ""
	}

	return lastParts[0] // Return the part before the ".html"
}

func (s shuhaige) GetNextChapterUrl(chapterContent, currentChapterUrl string) (string, error) {
	// Try to find both "下一章" and "下一页" as the next chapter/page link can use either
	keywords := []string{"下一章", "下一页"}
	startIdx := -1
	for _, kw := range keywords {
		startIdx = strings.Index(chapterContent, kw)
		if startIdx != -1 {
			break
		}
	}
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
	fullNextChapterURL := "https://m.shuhaige.net" + nextChapterURL
	return fullNextChapterURL, nil
}

func (s shuhaige) GetNovelCoverImageUrl(pageContent string) (string, error) {
	// Extract the cover image URL from the page content
	startIdx := strings.Index(pageContent, `<div class="detail"><img src="`)
	if startIdx == -1 {
		return "", nil
	}
	startIdx += len(`<div class="detail"><img src="`)

	endIdx := strings.Index(pageContent[startIdx:], `"`)
	if endIdx == -1 {
		return "", nil
	}
	endIdx += startIdx

	coverImageURL := pageContent[startIdx:endIdx]
	return coverImageURL, nil
}
