package sources

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

// txtChapterPath matches the chapter pages of 69shuba and twkan,
// /txt/<book>/<chapter>. On a book's latest chapter their 下一章 link leads
// elsewhere: twkan's to /txt/<book>/end.html, 69shuba's back to the book page.
var txtChapterPath = regexp.MustCompile(`^/txt/\d+/\d+/?$`)

// isTxtChapterURL reports whether rawURL is a 69shuba or twkan chapter page.
func isTxtChapterURL(rawURL string) bool {
	u, err := url.Parse(rawURL)
	return err == nil && txtChapterPath.MatchString(u.Path)
}

// txtNextChapterURL returns the page's 下一章 link resolved against
// currentChapterURL, or "" when it leads anywhere but a chapter page.
func txtNextChapterURL(chapterContent, currentChapterURL string) (string, error) {
	doc, err := html.Parse(strings.NewReader(chapterContent))
	if err != nil {
		return "", err
	}

	link, found := findNextChapterLink(doc)
	if !found {
		return "", fmt.Errorf("next chapter link not found")
	}

	current, err := url.Parse(currentChapterURL)
	if err != nil {
		return "", err
	}
	next, err := current.Parse(strings.TrimSpace(link))
	if err != nil {
		return "", err
	}
	if !txtChapterPath.MatchString(next.Path) {
		return "", nil
	}

	return next.String(), nil
}

type shuba struct{}

// NewShuba creates a new instance of the Shuba source
func NewShuba() Source {
	return &shuba{}
}

// GetNovelId extracts the novel ID from a Shuba URL
func (s *shuba) GetNovelId(url string) string {
	// Example URL: https://www.69shuba.com/book/36573.htm

	parts := strings.Split(url, "/")
	if len(parts) < 5 {
		return ""
	}

	lastParts := strings.Split(parts[len(parts)-1], ".")
	if len(lastParts) < 2 {
		return ""
	}

	return lastParts[0] // Return the part before the ".htm"
}

// GetChapterId extracts the chapter ID from a Shuba chapter URL
func (s *shuba) GetChapterId(chapterUrl string) string {
	// Example URL: https://www.69shuba.com/txt/36573/40179619
	// We need to extract the last part of the URL
	if !isTxtChapterURL(chapterUrl) {
		return ""
	}

	parts := strings.Split(chapterUrl, "/")
	if len(parts) < 6 {
		return ""
	}

	lastPart := parts[len(parts)-1]
	return lastPart // Return the last part as the chapter ID
}

// GetNextChapterUrl returns the page's 下一章 link, or "" on the latest chapter.
func (s *shuba) GetNextChapterUrl(chapterContent, currentChapterUrl string) (string, error) {
	return txtNextChapterURL(chapterContent, currentChapterUrl)
}

// findNextChapterLink traverses the HTML nodes to find the link for the next chapter.
// It looks for an <a> tag with the text "下一章".
func findNextChapterLink(n *html.Node) (string, bool) {
	if n.Type == html.ElementNode && n.Data == "a" {
		// Check if the link's text content is "下一章"
		var linkText strings.Builder
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.TextNode {
				linkText.WriteString(strings.TrimSpace(c.Data))
			}
		}

		if linkText.String() == "下一章" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					return attr.Val, true
				}
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if link, found := findNextChapterLink(c); found {
			return link, true
		}
	}

	return "", false
}

func (s *shuba) GetNovelCoverImageUrl(pageContent string) (string, error) {
	doc, err := html.Parse(strings.NewReader(pageContent))
	if err != nil {
		return "", err
	}

	bookImageURL, found := findMetaOgImage(doc)
	if found {
		return bookImageURL, nil
	}

	return "", fmt.Errorf("book image link not found")
}

// findMetaOgImage traverses the HTML nodes to find the content of <meta property="og:image">.
func findMetaOgImage(n *html.Node) (string, bool) {
	if n.Type == html.ElementNode && n.Data == "meta" {
		var property, content string
		for _, attr := range n.Attr {
			if attr.Key == "property" && attr.Val == "og:image" {
				property = attr.Val
			}
			if attr.Key == "content" {
				content = attr.Val
			}
		}
		if property == "og:image" && content != "" {
			return content, true
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if url, found := findMetaOgImage(c); found {
			return url, true
		}
	}
	return "", false
}
