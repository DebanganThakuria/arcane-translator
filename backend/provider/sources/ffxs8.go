package sources

import "strings"

type ffxs8 struct{}

func NewFfxs8() Source {
	return &ffxs8{}
}

func (f ffxs8) GetNovelId(url string) string {
	// Example URL: https://ffxs8.com/xhmf/22537/
	parts := strings.Split(strings.TrimSuffix(url, "/"), "/")
	for idx := 1; idx < len(parts)-1; idx++ {
		if parts[idx] == "index" {
			return parts[idx-1]
		}
	}

	if len(parts) > 0 {
		return parts[len(parts)-1]
	}

	return ""
}

func (f ffxs8) GetChapterId(chapterUrl string) string {
	// Example URL: https://ffxs8.com/xhmf/22537/index/1.html
	// Chapter numbers alone are not unique, so prefix with the source id
	parts := strings.Split(strings.TrimSuffix(chapterUrl, "/"), "/")
	if len(parts) == 0 {
		return ""
	}

	chapterNumber := strings.TrimSuffix(parts[len(parts)-1], ".html")
	return "ffxs8-" + chapterNumber
}

func (f ffxs8) GetNextChapterUrl(chapterContent, currentChapterUrl string) (string, error) {
	// The pager uses "下一节" for the next chapter link
	before, _, ok := strings.Cut(chapterContent, "下一节")
	if !ok {
		return "", nil
	}

	aTagStart := strings.LastIndex(before, `<a `)
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

	return "https://ffxs8.com" + nextChapterURL, nil
}

func (f ffxs8) GetNovelCoverImageUrl(pageContent string) (string, error) {
	// The novel page does not include a cover image
	return "", nil
}
