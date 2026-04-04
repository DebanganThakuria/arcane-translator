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
	case "ixdzs":
		return NewIxdzs()
	case "czbooks":
		return NewCzbooks()
	case "quanben":
		return NewQuanben()
	case "sjks88":
		return NewSjks88()
	default:
		return nil
	}
}
