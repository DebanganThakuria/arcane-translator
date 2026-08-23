package sources

import (
	"fmt"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

type huabenge struct{}

func NewHuaBenGe() Source {
	return &huabenge{}
}

func (h *huabenge) GetNovelId(rawURL string) string {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}

	parts := strings.Split(strings.Trim(parsedURL.Path, "/"), "/")
	for i := 0; i < len(parts)-1; i++ {
		if parts[i] == "book" {
			return parts[i+1]
		}
	}

	return ""
}

func (h *huabenge) GetChapterId(chapterURL string) string {
	parsedURL, err := url.Parse(chapterURL)
	if err != nil {
		return ""
	}

	parts := strings.Split(strings.Trim(parsedURL.Path, "/"), "/")
	if len(parts) == 0 {
		return ""
	}

	return strings.TrimSuffix(parts[len(parts)-1], ".html")
}

func (h *huabenge) GetNextChapterUrl(chapterContent, currentChapterURL string) (string, error) {
	doc, err := html.Parse(strings.NewReader(chapterContent))
	if err != nil {
		return "", err
	}

	nextChapterURL, found := findHuabengeLinkByID(doc, "linkNext")
	if !found {
		nextChapterURL, found = findNextChapterLink(doc)
	}
	if !found {
		nextChapterURL, found = findHuabengeChapterLink(doc, "下一页")
	}
	if !found {
		return "", fmt.Errorf("next chapter link not found")
	}

	if strings.HasPrefix(nextChapterURL, "http") {
		return nextChapterURL, nil
	}

	parsedCurrentURL, err := url.Parse(currentChapterURL)
	if err != nil {
		return "", err
	}

	parsedNextURL, err := url.Parse(nextChapterURL)
	if err != nil {
		return "", err
	}

	return parsedCurrentURL.ResolveReference(parsedNextURL).String(), nil
}

func findHuabengeLinkByID(n *html.Node, id string) (string, bool) {
	if n.Type == html.ElementNode && n.Data == "a" {
		var href string
		for _, attr := range n.Attr {
			if attr.Key == "id" && attr.Val != id {
				goto children
			}
			if attr.Key == "href" {
				href = attr.Val
			}
		}
		if href != "" {
			return href, true
		}
	}

children:
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if link, found := findHuabengeLinkByID(c, id); found {
			return link, true
		}
	}

	return "", false
}

func (h *huabenge) GetNovelCoverImageUrl(pageContent string) (string, error) {
	doc, err := html.Parse(strings.NewReader(pageContent))
	if err != nil {
		return "", err
	}

	bookImageURL, found := findMetaOgImage(doc)
	if found {
		return bookImageURL, nil
	}

	return "", nil
}

func findHuabengeChapterLink(n *html.Node, text string) (string, bool) {
	if n.Type == html.ElementNode && n.Data == "a" {
		var linkText strings.Builder
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.TextNode {
				linkText.WriteString(strings.TrimSpace(c.Data))
			}
		}

		if linkText.String() == text {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					return attr.Val, true
				}
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if link, found := findHuabengeChapterLink(c, text); found {
			return link, true
		}
	}

	return "", false
}
