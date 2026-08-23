package sources

import "strings"

type ilwxs struct{}

func NewIlwxs() Source {
	return &ilwxs{}
}

func (i ilwxs) GetNovelId(url string) string {
	// Example URL: https://m.ilwxs.com/shu/113687/
	parts := strings.Split(strings.TrimSuffix(url, "/"), "/")
	for idx := 0; idx < len(parts)-1; idx++ {
		if parts[idx] == "shu" {
			return parts[idx+1]
		}
	}

	return ""
}

func (i ilwxs) GetChapterId(chapterUrl string) string {
	// Example URL: https://m.ilwxs.com/shu/113687/158872011.html
	parts := strings.Split(strings.TrimSuffix(chapterUrl, "/"), "/")
	if len(parts) == 0 {
		return ""
	}

	lastParts := strings.Split(parts[len(parts)-1], ".")
	return lastParts[0] // Return the part before the ".html"
}

func (i ilwxs) GetNextChapterUrl(chapterContent, currentChapterUrl string) (string, error) {
	// The pager uses "下一章" for the next chapter link
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
	if strings.HasPrefix(nextChapterURL, "http") {
		return nextChapterURL, nil
	}

	return "https://m.ilwxs.com" + nextChapterURL, nil
}

func (i ilwxs) GetNovelCoverImageUrl(pageContent string) (string, error) {
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
