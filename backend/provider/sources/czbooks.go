package sources

import "strings"

type czbooks struct{}

func NewCzbooks() *czbooks {
	return &czbooks{}
}

func (c czbooks) GetNovelId(url string) string {
	// URL looks like this https://czbooks.net/n/s6g12fn4
	parts := strings.Split(url, "/")
	if len(parts) < 5 {
		return ""
	}
	return parts[4]
}

func (c czbooks) GetChapterId(chapterUrl string) string {
	// URL looks like this https://czbooks.net/n/s6g12fn4/s6n6j52hc?chapterNumber=0
	chapterUrl = strings.Split(chapterUrl, "?")[0]
	parts := strings.Split(chapterUrl, "/")
	if len(parts) < 6 {
		return ""
	}
	return parts[5]
}

func (c czbooks) GetNextChapterUrl(chapterContent, currentChapterUrl string) (string, error) {
	// The next chapter URL is in: <a class = "next-chapter" href="//czbooks.net/n/s6g12fn4/s6n6j52dk?chapterNumber=1">下一章</a>
	marker := `class = "next-chapter" href="`
	startIdx := strings.Index(chapterContent, marker)
	if startIdx == -1 {
		return "", nil
	}
	startIdx += len(marker)

	endIdx := strings.Index(chapterContent[startIdx:], `"`)
	if endIdx == -1 {
		return "", nil
	}

	href := chapterContent[startIdx : startIdx+endIdx]
	if strings.HasPrefix(href, "//") {
		href = "https:" + href
	}
	return href, nil
}

func (c czbooks) GetNovelCoverImageUrl(pageContent string) (string, error) {
	// The cover image is in: <div class = "thumbnail"><script...></script><img src="URL" referrerpolicy="no-referrer" ...>
	thumbMarker := `<div class = "thumbnail">`
	thumbIdx := strings.Index(pageContent, thumbMarker)
	if thumbIdx == -1 {
		return "", nil
	}

	imgMarker := `<img src="`
	imgIdx := strings.Index(pageContent[thumbIdx:], imgMarker)
	if imgIdx == -1 {
		return "", nil
	}
	imgIdx += thumbIdx + len(imgMarker)

	endIdx := strings.Index(pageContent[imgIdx:], `"`)
	if endIdx == -1 {
		return "", nil
	}

	return pageContent[imgIdx : imgIdx+endIdx], nil
}