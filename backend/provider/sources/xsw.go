package sources

import (
	"net/url"
	"regexp"
	"strings"
)

var (
	// Every chapter page links the next one twice, as <a href="/1758662/255859255.html">下一章</a>.
	xswNextLinkPattern = regexp.MustCompile(`<a[^>]+href="([^"]+)"[^>]*>\s*下一章\s*</a>`)
	// Chapter pages live at /<novel>/<chapter>.html. On the last chapter the
	// next link points at the chapter list instead.
	xswChapterPath = regexp.MustCompile(`^/(\d+)/(\d+)\.html$`)
	xswNovelPath   = regexp.MustCompile(`^/(\d+)(?:/|$)`)
	// The info page has no og:image; the cover is the only image on the
	// site's image host.
	xswCoverPattern = regexp.MustCompile(`<img[^>]+src="(https?://img\.xsw\.tw/[^"]+)"`)
)

// xsw reads m.xsw.tw (台灣小說網), a Traditional Chinese site. Its pages
// declare big5 but are served as UTF-8, which the scraper follows.
type xsw struct{}

func NewXsw() Source {
	return &xsw{}
}

func (x xsw) GetNovelId(novelURL string) string {
	// Example URL: https://m.xsw.tw/1758662/
	// Prefixed because the bare ids are numeric and could collide with another
	// source's novel ids.
	u, err := url.Parse(novelURL)
	if err != nil {
		return ""
	}

	match := xswNovelPath.FindStringSubmatch(u.Path)
	if match == nil {
		return ""
	}

	return "xsw-" + match[1]
}

func (x xsw) GetChapterId(chapterURL string) string {
	// Example URL: https://m.xsw.tw/1758662/255833825.html
	// Chapter ids are site-wide, but the novel id keeps them readable.
	u, err := url.Parse(chapterURL)
	if err != nil {
		return ""
	}

	match := xswChapterPath.FindStringSubmatch(u.Path)
	if match == nil {
		return ""
	}

	return "xsw-" + match[1] + "-" + match[2]
}

func (x xsw) GetNextChapterUrl(chapterContent, currentChapterURL string) (string, error) {
	match := xswNextLinkPattern.FindStringSubmatch(chapterContent)
	if match == nil {
		return "", nil
	}

	current, err := url.Parse(currentChapterURL)
	if err != nil {
		return "", err
	}
	next, err := current.Parse(strings.TrimSpace(match[1]))
	if err != nil {
		return "", err
	}

	// The last chapter's next link leads back to the chapter list.
	if !xswChapterPath.MatchString(next.Path) {
		return "", nil
	}

	return next.String(), nil
}

func (x xsw) GetNovelCoverImageUrl(pageContent string) (string, error) {
	match := xswCoverPattern.FindStringSubmatch(pageContent)
	if match == nil {
		return "", nil
	}

	return match[1], nil
}
