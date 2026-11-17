package sources

import (
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

var (
	// A chapter runs over several pages, linked as <a href="71653685-2.html">下一页</a>.
	// On its last page 下一页 is plain text, and only 下一章 leads on.
	mingshutangNextPagePattern    = regexp.MustCompile(`<a[^>]+href="([^"]+)"[^>]*>\s*下一页\s*</a>`)
	mingshutangNextChapterPattern = regexp.MustCompile(`<a[^>]+href="([^"]+)"[^>]*>\s*下一章\s*</a>`)
	// Chapter pages live at /mingshut/<book>/<chapter>.html, and the pages
	// after the first at /mingshut/<book>/<chapter>-<page>.html. On the last
	// chapter the next link leads back to the book page instead.
	mingshutangChapterPath  = regexp.MustCompile(`^/mingshut/(\d+)/(\d+)(?:-(\d+))?\.html$`)
	mingshutangNovelPath    = regexp.MustCompile(`^/mingshu/(\d+)/?$`)
	mingshutangCoverPattern = regexp.MustCompile(`<meta\s+property="og:image"\s+content="([^"]+)"`)
	mingshutangTitlePattern = regexp.MustCompile(`<h1>([^<]+)</h1>`)
	// Each page ends its text with filler and a notice: "…。。<br>本章未完，请点击下一页继续阅读！"
	// (chapter not finished, read on) or, on a chapter's last page,
	// "39195113。。<br>一寸欢喜提醒您：看完记得收藏【名书堂】…" (bookmark the site).
	mingshutangPageNotice = regexp.MustCompile(`(?:…|\d+)?。。<br\s*/?>[\s\x{3000}]*(?:本章未完|[^<]*名书堂)[^<]*`)
)

// mingshutang reads www.mingshutang.net (名书堂), a Simplified Chinese site.
// Books live at /mingshu/<book>/ and chapters at /mingshut/<book>/<chapter>.html.
//
// A chapter is split over pages of about a thousand characters, and the server
// gets a captcha form instead of any of them, so each page is pasted in and
// translated as its own chapter; see manualHTMLSources.
type mingshutang struct{}

func NewMingshutang() Source {
	return &mingshutang{}
}

func (m mingshutang) GetNovelId(novelURL string) string {
	// Example URL: https://www.mingshutang.net/mingshu/51416484288/
	u, err := url.Parse(novelURL)
	if err != nil {
		return ""
	}

	match := mingshutangNovelPath.FindStringSubmatch(u.Path)
	if match == nil {
		return ""
	}

	return "mingshutang-" + match[1]
}

func (m mingshutang) GetChapterId(chapterURL string) string {
	// Example URLs: https://www.mingshutang.net/mingshut/51416484288/71653685.html
	// and its second page, .../71653685-2.html.
	u, err := url.Parse(chapterURL)
	if err != nil {
		return ""
	}

	match := mingshutangChapterPath.FindStringSubmatch(u.Path)
	if match == nil {
		return ""
	}

	id := "mingshutang-" + match[1] + "-" + match[2]
	if match[3] != "" {
		id += "-" + match[3]
	}
	return id
}

// GetNextChapterUrl returns the chapter's next page while it has one, and the
// next chapter after its last page.
func (m mingshutang) GetNextChapterUrl(chapterContent, currentChapterURL string) (string, error) {
	current, err := url.Parse(currentChapterURL)
	if err != nil {
		return "", err
	}

	for _, pattern := range []*regexp.Regexp{mingshutangNextPagePattern, mingshutangNextChapterPattern} {
		match := pattern.FindStringSubmatch(chapterContent)
		if match == nil {
			continue
		}

		next, err := current.Parse(strings.TrimSpace(match[1]))
		if err != nil {
			return "", err
		}
		// The last chapter's next link leads back to the book page.
		if mingshutangChapterPath.MatchString(next.Path) {
			return next.String(), nil
		}
	}

	return "", nil
}

func (m mingshutang) GetNovelCoverImageUrl(pageContent string) (string, error) {
	match := mingshutangCoverPattern.FindStringSubmatch(pageContent)
	if match == nil {
		return "", nil
	}

	return match[1], nil
}

// PrepareChapter drops the notice the site ends each page with, and numbers
// the title of every page after a chapter's first, as in "第1章 催收（2）", so
// the pages of one chapter can be told apart.
func (m mingshutang) PrepareChapter(chapterContent, chapterURL string) (string, error) {
	page := mingshutangPageNotice.ReplaceAllString(chapterContent, "")

	u, err := url.Parse(chapterURL)
	if err != nil {
		return "", err
	}
	match := mingshutangChapterPath.FindStringSubmatch(u.Path)
	if match == nil || match[3] == "" {
		return page, nil
	}
	if n, err := strconv.Atoi(match[3]); err != nil || n < 2 {
		return page, nil
	}

	title := mingshutangTitlePattern.FindStringSubmatch(page)
	if title == nil {
		return page, nil
	}
	// The <title> and breadcrumb repeat the heading; number them all alike.
	return strings.ReplaceAll(page, title[1], title[1]+"（"+match[3]+"）"), nil
}
