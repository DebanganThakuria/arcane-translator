package sources

type Source interface {
	GetNovelId(url string) string
	GetChapterId(chapterUrl string) string
	GetNextChapterUrl(chapterContent, currentChapterUrl string) (string, error)
	GetNovelCoverImageUrl(pageContent string) (string, error)
}

// ChapterPreparer is implemented by a source whose chapter pages need fixing
// before translation, such as one that shuffles its paragraphs.
type ChapterPreparer interface {
	PrepareChapter(chapterContent, chapterURL string) (string, error)
}

// PrepareChapter returns the chapter page at chapterURL as the translator
// should read it.
func PrepareChapter(sourceType, chapterContent, chapterURL string) (string, error) {
	if preparer, ok := GetSource(sourceType).(ChapterPreparer); ok {
		return preparer.PrepareChapter(chapterContent, chapterURL)
	}
	return chapterContent, nil
}

// manualHTMLSources are sources the server must never scrape. Their pages come
// back with a normal 200 but carry the wrong details, so the scrape "succeeds"
// and nothing downstream notices. The reader supplies the HTML from their own
// browser instead.
var manualHTMLSources = map[string]bool{
	"ixdzs": true,
}

// RequiresManualHTML reports whether pages from sourceType must be pasted in by
// the reader rather than scraped by the server.
func RequiresManualHTML(sourceType string) bool {
	return manualHTMLSources[sourceType]
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
	case "scribblehub":
		return NewScribbleHub()
	case "royalroad":
		return NewRoyalroad()
	case "huabenge":
		return NewHuaBenGe()
	case "ilwxs":
		return NewIlwxs()
	case "ffxs8":
		return NewFfxs8()
	case "novel543":
		return NewNovel543()
	case "85novel":
		return NewNovel85()
	default:
		return nil
	}
}
