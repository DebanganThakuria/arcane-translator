package sources

import (
	"strconv"
	"strings"
)

type ixdzs struct{}

func NewIxdzs() Source {
	return &ixdzs{}
}

func (s ixdzs) GetNovelId(url string) string {
	// The URL looks like this https://ixdzs.tw/read/526058/
	// The last number + ixdzs is the novel ID
	parts := strings.Split(url, "/")
	if len(parts) < 5 {
		return ""
	}

	return "ixdzs" + parts[4]
}

func (s ixdzs) GetChapterId(chapterUrl string) string {
	// The URL looks like this https://ixdzs.tw/read/526058/p1.html
	parts := strings.Split(chapterUrl, "/")
	if len(parts) < 6 {
		return ""
	}

	lastParts := strings.Split(parts[5], ".")
	if len(lastParts) < 2 {
		return ""
	}

	return parts[4] + lastParts[0]
}

func (s ixdzs) GetNextChapterUrl(chapterContent, currentChapterUrl string) (string, error) {
	// The current chapter URL looks like this https://ixdzs.tw/read/526058/p1.html
	// The next chapter URL looks like this https://ixdzs.tw/read/526058/p2.html
	// We can just increase the count to get the next chapter URL
	parts := strings.Split(currentChapterUrl, "/")
	if len(parts) < 6 {
		return "", nil
	}

	lastParts := strings.Split(parts[5], ".")
	if len(lastParts) < 2 {
		return "", nil
	}

	countStr := lastParts[0][1:] // Get the number part after "p"
	count, err := strconv.Atoi(countStr)
	if err != nil {
		return "", nil
	}

	nextCount := count + 1

	nextChapterUrl := strings.Join(parts[:5], "/") + "/p" + strconv.Itoa(nextCount) + ".html"
	return nextChapterUrl, nil
}

func (s ixdzs) GetNovelCoverImageUrl(pageContent string) (string, error) {
	// The url is present in <meta property="og:image" content="https://img22.ixdzs.com/e9/5c/e95c75a9eb2af8efa3e9cc00471700d8.jpg" />
	startIdx := strings.Index(pageContent, `<meta property="og:image" content="`)
	if startIdx == -1 {
		return "", nil
	}
	startIdx += len(`<meta property="og:image" content="`)

	endIdx := strings.Index(pageContent[startIdx:], `"`)
	if endIdx == -1 {
		return "", nil
	}
	endIdx += startIdx

	imageURL := pageContent[startIdx:endIdx]
	return imageURL, nil
}
