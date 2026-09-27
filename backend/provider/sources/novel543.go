package sources

import (
	"regexp"
	"strings"
)

const novel543BaseURL = "https://www.novel543.com"

var (
	// Every chapter page sets the pager target in a script variable, which is
	// less ambiguous than the "下一章" text that also appears in the settings
	// modal.
	novel543NextURLPattern = regexp.MustCompile(`var\s+nextUrl\s*=\s*'([^']*)'`)
	// The site writes og tags with name= rather than property=, so the shared
	// findMetaOgImage helper does not see them.
	novel543CoverPattern = regexp.MustCompile(`<meta\s+(?:name|property)="og:image"\s+content="([^"]+)"`)
)

type novel543 struct{}

func NewNovel543() Source {
	return &novel543{}
}

func (n novel543) GetNovelId(url string) string {
	// Example URL: https://www.novel543.com/0725683934/
	// Prefixed because the bare ids are numeric and could collide with another
	// source's novel ids.
	path := strings.TrimPrefix(url, "https://")
	path = strings.TrimPrefix(path, "http://")
	parts := strings.Split(path, "/")
	if len(parts) < 2 || parts[1] == "" {
		return ""
	}

	return "novel543-" + parts[1]
}

func (n novel543) GetChapterId(chapterUrl string) string {
	// Example URLs: https://www.novel543.com/0725683934/8096_1.html
	//               https://www.novel543.com/0725683934/8096_1_2.html
	// A chapter is split across several pages and each page is stored as its
	// own entry. Page names repeat across novels, so include the novel id.
	parts := strings.Split(strings.TrimSuffix(chapterUrl, "/"), "/")
	if len(parts) < 2 {
		return ""
	}

	page := strings.TrimSuffix(parts[len(parts)-1], ".html")
	return "novel543-" + parts[len(parts)-2] + "-" + page
}

func (n novel543) GetNextChapterUrl(chapterContent, currentChapterUrl string) (string, error) {
	match := novel543NextURLPattern.FindStringSubmatch(chapterContent)
	if match == nil || match[1] == "" {
		return "", nil
	}

	nextChapterURL := match[1]
	// The last page of the final chapter points at an "end" page, not a chapter.
	if strings.HasSuffix(nextChapterURL, "/end.html") {
		return "", nil
	}

	if strings.HasPrefix(nextChapterURL, "http") {
		return nextChapterURL, nil
	}

	return novel543BaseURL + nextChapterURL, nil
}

func (n novel543) GetNovelCoverImageUrl(pageContent string) (string, error) {
	match := novel543CoverPattern.FindStringSubmatch(pageContent)
	if match == nil {
		return "", nil
	}

	return match[1], nil
}
