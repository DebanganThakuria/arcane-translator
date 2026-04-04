package sources

import "strings"

type quanben struct{}

func NewQuanben() *quanben {
	return &quanben{}
}

func (q *quanben) GetNovelId(url string) string {
	// Example URL: https://www.quanben.io/n/gaowu-wodetianfuwuxianshengji/list.html
	parts := strings.Split(url, "/")
	if len(parts) < 5 {
		return ""
	}
	return parts[4]
}

func (q *quanben) GetChapterId(chapterUrl string) string {
	// Example URL: https://www.quanben.io/n/gaowu-wodetianfuwuxianshengji/1.html
	// Chapter ID: novel_id + "_" + chapter_number
	parts := strings.Split(chapterUrl, "/")
	if len(parts) < 6 {
		return ""
	}
	novelId := parts[4]
	chapterFile := strings.Split(parts[5], ".")
	if len(chapterFile) < 2 {
		return ""
	}
	return novelId + "_" + chapterFile[0]
}

func (q *quanben) GetNextChapterUrl(chapterContent, currentChapterUrl string) (string, error) {
	// Next chapter link: <a href="/n/gaowu-wodetianfuwuxianshengji/2.html" itemprop="url" rel="next">下一页</a>
	relNext := `rel="next"`
	relIdx := strings.Index(chapterContent, relNext)
	if relIdx == -1 {
		return "", nil
	}

	// Search backward from rel="next" to find href="
	segment := chapterContent[:relIdx]
	hrefMarker := `href="`
	hrefIdx := strings.LastIndex(segment, hrefMarker)
	if hrefIdx == -1 {
		return "", nil
	}
	hrefIdx += len(hrefMarker)

	endIdx := strings.Index(chapterContent[hrefIdx:], `"`)
	if endIdx == -1 {
		return "", nil
	}

	href := chapterContent[hrefIdx : hrefIdx+endIdx]
	if strings.HasPrefix(href, "/") {
		href = "https://www.quanben.io" + href
	}
	return href, nil
}

func (q *quanben) GetNovelCoverImageUrl(pageContent string) (string, error) {
	// Cover image: <img src="https://img.c0m.io/quanben.io/..." alt="..." itemprop="image" />
	marker := `itemprop="image"`
	imgIdx := strings.Index(pageContent, marker)
	if imgIdx == -1 {
		return "", nil
	}

	// Search backward from itemprop="image" to find src="
	segment := pageContent[:imgIdx]
	srcMarker := `src="`
	srcIdx := strings.LastIndex(segment, srcMarker)
	if srcIdx == -1 {
		return "", nil
	}
	srcIdx += len(srcMarker)

	endIdx := strings.Index(pageContent[srcIdx:], `"`)
	if endIdx == -1 {
		return "", nil
	}

	return pageContent[srcIdx : srcIdx+endIdx], nil
}
