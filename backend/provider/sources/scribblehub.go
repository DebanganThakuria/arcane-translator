package sources

import (
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

type scribbleHub struct{}

func NewScribbleHub() *scribbleHub {
	return &scribbleHub{}
}

func (s scribbleHub) GetNovelId(url string) string {
	// URL: https://www.scribblehub.com/series/131121/transmigrated-into-a-nobles-beaten-son/
	url = strings.TrimSuffix(url, "/")
	parts := strings.Split(url, "/")
	return parts[len(parts)-1]
}

func (s scribbleHub) GetChapterId(chapterUrl string) string {
	// URL: https://www.scribblehub.com/read/131121-transmigrated-into-a-nobles-beaten-son/chapter/131126/
	chapterUrl = strings.TrimSuffix(chapterUrl, "/")
	parts := strings.Split(chapterUrl, "/")
	return "scribblehub" + parts[len(parts)-1]
}

func (s scribbleHub) GetNextChapterUrl(chapterContent, currentChapterUrl string) (string, error) {
	doc, err := html.Parse(strings.NewReader(chapterContent))
	if err != nil {
		return "", err
	}

	nextURL, found := findScribbleHubNextChapterLink(doc)
	if found {
		return nextURL, nil
	}

	return "", fmt.Errorf("next chapter link not found")
}

// findScribbleHubNextChapterLink looks for <a class="btn-wi btn-next" href="...">
func findScribbleHubNextChapterLink(n *html.Node) (string, bool) {
	if n.Type == html.ElementNode && n.Data == "a" {
		var href string
		isNextBtn := false
		for _, attr := range n.Attr {
			if attr.Key == "class" && strings.Contains(attr.Val, "btn-next") {
				isNextBtn = true
			}
			if attr.Key == "href" {
				href = attr.Val
			}
		}
		if isNextBtn && href != "" && href != "#" {
			return href, true
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if link, found := findScribbleHubNextChapterLink(c); found {
			return link, true
		}
	}
	return "", false
}

func (s scribbleHub) GetNovelCoverImageUrl(pageContent string) (string, error) {
	// ScribbleHub has <meta property="og:image" content="..."> in the page head
	doc, err := html.Parse(strings.NewReader(pageContent))
	if err != nil {
		return "", err
	}

	imageURL, found := findMetaOgImage(doc)
	if found {
		return imageURL, nil
	}

	return "", fmt.Errorf("cover image not found")
}
