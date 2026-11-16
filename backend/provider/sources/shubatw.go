package sources

import (
	"net/url"
	"regexp"
	"strings"
)

var (
	// Every chapter page links the next one as <a href="/read/88917/491715">下一章</a>.
	shubaTWNextLinkPattern = regexp.MustCompile(`<a[^>]+href="([^"]+)"[^>]*>\s*下一章\s*</a>`)
	// Chapter pages live at /read/<book>/<chapter>. The first chapter's
	// previous link and the last page's next link lead to the book page.
	shubaTWChapterPath  = regexp.MustCompile(`^/read/(\d+)/(\d+)/?$`)
	shubaTWNovelPath    = regexp.MustCompile(`^/book/(\d+)/?$`)
	shubaTWCoverPattern = regexp.MustCompile(`<meta\s+property="og:image"\s+content="([^"]+)"`)
	// Every chapter title ends in its page, as in "第1章 活命(1 / 1)". Only the
	// single-page marker is noise; a "(1 / 2)" would mean a split chapter,
	// which the page should keep showing.
	shubaTWSinglePageMarker = regexp.MustCompile(`\s*[(（]\s*1\s*/\s*1\s*[)）]`)
)

// shubaTW reads 69shuba.tw (69書吧), a Traditional Chinese site laid out
// differently from 69shuba.com: books at /book/<id>/ and chapters at
// /read/<id>/<chapter>.
type shubaTW struct{}

func NewShubaTW() Source {
	return &shubaTW{}
}

func (s shubaTW) GetNovelId(novelURL string) string {
	// Example URL: https://69shuba.tw/book/88917/
	// Prefixed so the numeric ids cannot collide with 69shuba.com's.
	u, err := url.Parse(novelURL)
	if err != nil {
		return ""
	}

	match := shubaTWNovelPath.FindStringSubmatch(u.Path)
	if match == nil {
		return ""
	}

	return "69shubatw-" + match[1]
}

func (s shubaTW) GetChapterId(chapterURL string) string {
	// Example URL: https://69shuba.tw/read/88917/491713
	u, err := url.Parse(chapterURL)
	if err != nil {
		return ""
	}

	match := shubaTWChapterPath.FindStringSubmatch(u.Path)
	if match == nil {
		return ""
	}

	return "69shubatw-" + match[1] + "-" + match[2]
}

func (s shubaTW) GetNextChapterUrl(chapterContent, currentChapterURL string) (string, error) {
	match := shubaTWNextLinkPattern.FindStringSubmatch(chapterContent)
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

	// The last page's next link leads back to the book page.
	if !shubaTWChapterPath.MatchString(next.Path) {
		return "", nil
	}

	return next.String(), nil
}

func (s shubaTW) GetNovelCoverImageUrl(pageContent string) (string, error) {
	match := shubaTWCoverPattern.FindStringSubmatch(pageContent)
	if match == nil {
		return "", nil
	}

	// The site writes the cover without a scheme: //p.69shuba.tw/88/88917/88917s.jpg.
	cover := match[1]
	if strings.HasPrefix(cover, "//") {
		cover = "https:" + cover
	}
	return cover, nil
}

// PrepareChapter drops the "(1 / 1)" page marker the site adds to every
// chapter title, so it does not end up in the translated title.
func (s shubaTW) PrepareChapter(chapterContent, _ string) (string, error) {
	return shubaTWSinglePageMarker.ReplaceAllString(chapterContent, ""), nil
}
