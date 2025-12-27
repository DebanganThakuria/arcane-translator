package sources

type Source interface {
	GetNovelId(url string) string
	GetChapterId(chapterUrl string) string
	GetNextChapterUrl(chapterContent, currentChapterUrl string) (string, error)
	GetNovelCoverImageUrl(pageContent string) (string, error)
}

func GetSource(sourceType string) Source {
	switch sourceType {
	case "69shuba":
		return NewShuba()
	case "69yue":
		return NewYue()
	case "shuhaige":
		return NewShuhaige()
	case "twkan":
		return NewTwkan()
	case "doupo":
		return NewDuopo()
	case "syosetu":
		return NewSyosetu()
	default:
		return nil
	}
}
