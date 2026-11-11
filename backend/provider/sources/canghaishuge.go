package sources

import (
	"net/url"
	"regexp"
	"strings"
)

var (
	// Every chapter page links the next one as <a href="/<book>/2.html">下一章</a>.
	canghaishugeNextLinkPattern = regexp.MustCompile(`<a[^>]+href="([^"]+)"[^>]*>\s*下一章\s*</a>`)
	// Chapter pages live at /<book>/<n>.html. On the last chapter the next
	// link points back at the book page instead.
	canghaishugeChapterPath  = regexp.MustCompile(`^/([A-Za-z0-9_-]+)/(\d+)\.html$`)
	canghaishugeNovelPath    = regexp.MustCompile(`^/([A-Za-z0-9_-]+)/?$`)
	canghaishugeCoverPattern = regexp.MustCompile(`<meta\s+property="og:image"\s+content="([^"]+)"`)
)

// canghaishuge reads www.canghaishuge.com (沧海书阁), a Simplified Chinese site
// with one page per chapter. Books are named by a pinyin slug, as in
// /shenhuajiyuankaijusancunshihoujiyin/.
type canghaishuge struct{}

func NewCanghaishuge() Source {
	return &canghaishuge{}
}

func (c canghaishuge) GetNovelId(novelURL string) string {
	// Example URL: https://www.canghaishuge.com/shenhuajiyuankaijusancunshihoujiyin/
	u, err := url.Parse(novelURL)
	if err != nil {
		return ""
	}

	match := canghaishugeNovelPath.FindStringSubmatch(u.Path)
	if match == nil {
		return ""
	}

	return "canghaishuge-" + match[1]
}

func (c canghaishuge) GetChapterId(chapterURL string) string {
	// Example URL: https://www.canghaishuge.com/shenhuajiyuankaijusancunshihoujiyin/1.html
	u, err := url.Parse(chapterURL)
	if err != nil {
		return ""
	}

	match := canghaishugeChapterPath.FindStringSubmatch(u.Path)
	if match == nil {
		return ""
	}

	return "canghaishuge-" + match[1] + "-" + match[2]
}

func (c canghaishuge) GetNextChapterUrl(chapterContent, currentChapterURL string) (string, error) {
	match := canghaishugeNextLinkPattern.FindStringSubmatch(chapterContent)
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

	// The last chapter's next link leads back to the book page.
	if !canghaishugeChapterPath.MatchString(next.Path) {
		return "", nil
	}

	return next.String(), nil
}

func (c canghaishuge) GetNovelCoverImageUrl(pageContent string) (string, error) {
	match := canghaishugeCoverPattern.FindStringSubmatch(pageContent)
	if match == nil {
		return "", nil
	}

	return match[1], nil
}
