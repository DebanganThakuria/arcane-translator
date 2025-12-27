package sources

import (
	"fmt"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

type duopo struct{}

func NewDuopo() Source {
	return &duopo{}
}

func (d duopo) GetNovelId(url string) string {
	// The URL looks like this https://doupo.935666.xyz/b/5217/
	parts := strings.Split(url, "/")
	if len(parts) < 5 {
		return ""
	}

	return parts[4] // Return the part after the domain
}

func (d duopo) GetChapterId(chapterUrl string) string {
	// The URl looks like this https://doupo.935666.xyz/b/5217/3861075.html
	parts := strings.Split(chapterUrl, "/")
	if len(parts) < 6 {
		return ""
	}

	lastParts := strings.Split(parts[5], ".")
	if len(lastParts) < 2 {
		return ""
	}

	return lastParts[0] // Return the part before the ".html"
}

func (d duopo) GetNextChapterUrl(chapterContent, currentChapterUrl string) (string, error) {
	doc, err := html.Parse(strings.NewReader(chapterContent))
	if err != nil {
		return "", err
	}

	nextChapterHref, found := findNextChapterLink(doc)
	if !found {
		return "", fmt.Errorf("next chapter link not found")
	}

	// If the href is a relative path, construct the absolute URL
	if strings.HasPrefix(nextChapterHref, "/") {
		baseURL, err := url.Parse(currentChapterUrl)
		if err != nil {
			return "", fmt.Errorf("failed to parse current chapter URL: %w", err)
		}
		return fmt.Sprintf("%s://%s%s", baseURL.Scheme, baseURL.Host, nextChapterHref), nil
	}

	return nextChapterHref, nil
}

func (d duopo) GetNovelCoverImageUrl(pageContent string) (string, error) {
	doc, err := html.Parse(strings.NewReader(pageContent))
	if err != nil {
		return "", err
	}

	coverImageURL, found := findMetaOgImage(doc)
	if !found {
		return "", fmt.Errorf("cover image not found")
	}

	return coverImageURL, nil
}
