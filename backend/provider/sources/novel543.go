package sources

import (
	"errors"
	"fmt"
	"path"
	"regexp"
	"slices"
	"strings"

	"golang.org/x/net/html"
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
	// A split chapter's title ends in its page, as in "第1章 葉長安 (1/2)".
	novel543PartMarker = regexp.MustCompile(`\s*[(（]\s*\d+\s*/\s*\d+\s*[)）]`)
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

// NextChapterPage returns the next page of the chapter that page belongs to,
// or "" on the chapter's last page. The site splits a chapter across pages
// named 8096_1.html, 8096_1_2.html, and so on; the last page's next link
// leads to the following chapter, 8096_2.html.
func (n novel543) NextChapterPage(page, pageURL string) string {
	match := novel543NextURLPattern.FindStringSubmatch(page)
	if match == nil || match[1] == "" {
		return ""
	}

	next, current := path.Base(match[1]), path.Base(pageURL)
	if next == current || novel543ChapterKey(next) != novel543ChapterKey(current) {
		return ""
	}

	if strings.HasPrefix(match[1], "http") {
		return match[1]
	}
	return novel543BaseURL + match[1]
}

// novel543ChapterKey names the chapter a page file belongs to: both
// 8096_1.html and 8096_1_2.html belong to 8096_1.
func novel543ChapterKey(file string) string {
	parts := strings.Split(strings.TrimSuffix(file, ".html"), "_")
	if len(parts) > 2 {
		parts = parts[:2]
	}
	return strings.Join(parts, "_")
}

// JoinChapterPages joins the pages of one chapter into its first page: the
// story paragraphs of every later page are appended to the first page's
// text, and the "(1/2)" page marker is dropped from the title.
func (n novel543) JoinChapterPages(pages []string) (string, error) {
	if len(pages) == 0 {
		return "", errors.New("novel543: no chapter pages to join")
	}

	doc, err := html.Parse(strings.NewReader(pages[0]))
	if err != nil {
		return "", err
	}
	content := novel543Content(doc)
	if content == nil {
		return "", errors.New("novel543: chapter page has no text block")
	}

	// The 繁體/简体 switch sits at the end of the text; drop it so the later
	// pages' text follows straight on.
	var paragraphs []*html.Node
	collectParagraphs(content, &paragraphs)
	for _, p := range paragraphs {
		if hasHTMLAttr(p, "class") {
			p.Parent.RemoveChild(p)
		}
	}

	for i, page := range pages[1:] {
		part, err := html.Parse(strings.NewReader(page))
		if err != nil {
			return "", err
		}
		partContent := novel543Content(part)
		if partContent == nil {
			return "", fmt.Errorf("novel543: page %d of the chapter has no text block", i+2)
		}

		var story []*html.Node
		collectParagraphs(partContent, &story)
		for _, p := range story {
			if hasHTMLAttr(p, "class") {
				continue
			}
			p.Parent.RemoveChild(p)
			content.AppendChild(p)
		}
	}

	novel543DropPartMarkers(doc)

	var out strings.Builder
	if err := html.Render(&out, doc); err != nil {
		return "", err
	}
	return out.String(), nil
}

// novel543Content finds the chapter text block, <div class="content py-5">.
func novel543Content(doc *html.Node) *html.Node {
	return findHTMLNode(doc, func(node *html.Node) bool {
		if node.Type != html.ElementNode || node.Data != "div" {
			return false
		}
		classes := strings.Fields(htmlAttr(node, "class"))
		return slices.Contains(classes, "content") && slices.Contains(classes, "py-5")
	})
}

func novel543DropPartMarkers(node *html.Node) {
	if node.Type == html.TextNode && node.Parent != nil && (node.Parent.Data == "title" || node.Parent.Data == "h1") {
		node.Data = novel543PartMarker.ReplaceAllString(node.Data, "")
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		novel543DropPartMarkers(child)
	}
}
