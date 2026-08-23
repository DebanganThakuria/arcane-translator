package sources

import (
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

type royalroad struct{}

func NewRoyalroad() *royalroad {
	return &royalroad{}
}

func (r *royalroad) GetNovelId(url string) string {
	// URL: https://www.royalroad.com/fiction/158686/incarnation-of-sorcery-op-mc-litrpg
	url = strings.TrimSuffix(url, "/")
	parts := strings.Split(url, "/")
	// Find the "fiction" part and get the next part which is the ID
	for i, part := range parts {
		if part == "fiction" && i+1 < len(parts) {
			return parts[i+1]
		}
	}
	return ""
}

func (r *royalroad) GetChapterId(chapterUrl string) string {
	// URL: https://www.royalroad.com/fiction/150039/new-life-of-the-dark-archmage-op-mc-isekai-litrpg/chapter/2978405/chapter-1-the-dark-archmage
	chapterUrl = strings.TrimSuffix(chapterUrl, "/")
	parts := strings.Split(chapterUrl, "/")
	// Find the "chapter" part and get the next part which is the chapter ID
	for i, part := range parts {
		if part == "chapter" && i+1 < len(parts) {
			return "royalroad_" + parts[i+1]
		}
	}
	return ""
}

func (r *royalroad) GetNextChapterUrl(chapterContent, currentChapterUrl string) (string, error) {
	doc, err := html.Parse(strings.NewReader(chapterContent))
	if err != nil {
		return "", err
	}

	nextURL, found := findRoyalRoadNextChapterLink(doc)
	if found {
		// Handle relative URLs
		if strings.HasPrefix(nextURL, "/") {
			return "https://www.royalroad.com" + nextURL, nil
		}
		return nextURL, nil
	}

	return "", fmt.Errorf("next chapter link not found")
}

// findRoyalRoadNextChapterLink looks for <link rel="next" href="...">
func findRoyalRoadNextChapterLink(n *html.Node) (string, bool) {
	if n.Type == html.ElementNode && n.Data == "link" {
		var href string
		isNext := false
		for _, attr := range n.Attr {
			if attr.Key == "rel" && attr.Val == "next" {
				isNext = true
			}
			if attr.Key == "href" {
				href = attr.Val
			}
		}
		if isNext && href != "" {
			return href, true
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if link, found := findRoyalRoadNextChapterLink(c); found {
			return link, true
		}
	}
	return "", false
}

func (r *royalroad) GetNovelCoverImageUrl(pageContent string) (string, error) {
	// RoyalRoad has <meta property="og:image" content="..."> in the page head
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
