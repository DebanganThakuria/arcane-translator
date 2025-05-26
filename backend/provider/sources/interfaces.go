package sources

type Source interface {
	GetNovelId(url string) string
	GetChapterId(chapterUrl string) string
	GetChapterUrl(previousChapterUrl string) (string, error)
}

func GetSource(sourceType string) Source {
	switch sourceType {
	case "69shuba":
		return NewShuba()
	default:
		return nil
	}
}
