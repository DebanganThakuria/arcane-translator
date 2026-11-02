package sources

import (
	"regexp"
	"strings"
	"testing"
)

func TestNovel543Ids(t *testing.T) {
	s := NewNovel543()

	if got := s.GetNovelId("https://www.novel543.com/0725683934/"); got != "novel543-0725683934" {
		t.Errorf("GetNovelId = %q", got)
	}
	if got := s.GetNovelId("https://www.novel543.com/"); got != "" {
		t.Errorf("GetNovelId without a novel = %q, want empty", got)
	}

	chapters := map[string]string{
		"https://www.novel543.com/0725683934/8096_1.html":   "novel543-0725683934-8096_1",
		"https://www.novel543.com/0725683934/8096_1_2.html": "novel543-0725683934-8096_1_2",
	}
	for url, want := range chapters {
		if got := s.GetChapterId(url); got != want {
			t.Errorf("GetChapterId(%q) = %q, want %q", url, got, want)
		}
	}
}

func TestNovel543NextChapterUrl(t *testing.T) {
	s := NewNovel543()
	current := "https://www.novel543.com/0725683934/8096_1.html"

	tests := map[string]struct {
		page string
		want string
	}{
		"next page of the same chapter": {
			page: `<span class="button nextBtn"> 下一章 </span></section><script> var nextUrl = '/0725683934/8096_1_2.html'; var prevUrl = '/0725683934/8096_1.html';</script>`,
			want: "https://www.novel543.com/0725683934/8096_1_2.html",
		},
		"end of the novel": {
			page: `<script> var nextUrl = '/0725683934/end.html';</script>`,
			want: "",
		},
		"no pager": {
			page: `<html></html>`,
			want: "",
		},
	}
	for name, tt := range tests {
		got, err := s.GetNextChapterUrl(tt.page, current)
		if err != nil {
			t.Fatalf("%s: unexpected error: %v", name, err)
		}
		if got != tt.want {
			t.Errorf("%s: got %q, want %q", name, got, tt.want)
		}
	}
}

func TestNovel543Cover(t *testing.T) {
	page := `<meta name="og:image" content="https://i2.novel543.com/thumb/120x160/20260222/060358393433.jpg">`
	got, err := NewNovel543().GetNovelCoverImageUrl(page)
	if err != nil {
		t.Fatal(err)
	}
	if got != "https://i2.novel543.com/thumb/120x160/20260222/060358393433.jpg" {
		t.Errorf("cover = %q", got)
	}
}

func TestNovel543NextChapterPage(t *testing.T) {
	s := NewNovel543().(MultiPageChapter)

	tests := map[string]struct {
		page, url, want string
	}{
		"first page leads to the second": {
			page: `<script> var nextUrl = '/0312606674/8096_1_2.html';</script>`,
			url:  "https://www.novel543.com/0312606674/8096_1.html",
			want: "https://www.novel543.com/0312606674/8096_1_2.html",
		},
		"second page leads to the third": {
			page: `<script> var nextUrl = '/0312606674/8096_1_3.html';</script>`,
			url:  "https://www.novel543.com/0312606674/8096_1_2.html",
			want: "https://www.novel543.com/0312606674/8096_1_3.html",
		},
		"last page leads to the next chapter": {
			page: `<script> var nextUrl = '/0312606674/8096_2.html';</script>`,
			url:  "https://www.novel543.com/0312606674/8096_1_2.html",
			want: "",
		},
		"single-page chapter": {
			page: `<script> var nextUrl = '/0312606674/8096_11.html';</script>`,
			url:  "https://www.novel543.com/0312606674/8096_10.html",
			want: "",
		},
		"no pager": {
			page: `<html></html>`,
			url:  "https://www.novel543.com/0312606674/8096_1.html",
			want: "",
		},
	}
	for name, tt := range tests {
		if got := s.NextChapterPage(tt.page, tt.url); got != tt.want {
			t.Errorf("%s: got %q, want %q", name, got, tt.want)
		}
	}
}

// novel543Page builds a chapter page laid out like the site's: an ad slot,
// the story paragraphs, then the 繁體/简体 switch, all in one text block.
func novel543Page(title string, paragraphs ...string) string {
	var b strings.Builder
	b.WriteString(`<html><head><title>` + title + `</title></head><body><h1>` + title + `</h1>`)
	b.WriteString(`<div class="content py-5"><div class="gadBlock"><ins class="clickforceads"></ins></div>`)
	for _, p := range paragraphs {
		b.WriteString(`<p>` + p + `</p>`)
	}
	b.WriteString(`<p class="jft"><span class="button">繁體</span><span class="button">简体</span></p></div></body></html>`)
	return b.String()
}

func TestNovel543JoinChapterPages(t *testing.T) {
	s := NewNovel543().(MultiPageChapter)

	got, err := s.JoinChapterPages([]string{
		novel543Page("第1章 葉長安 (1/2)", "天荒大陸。", "東域一個小城。"),
		novel543Page("第1章 葉長安 (2/2)", "“這個太初神殿是什麼情況？”", "葉長安笑了。"),
	})
	if err != nil {
		t.Fatalf("JoinChapterPages() error = %v", err)
	}

	var text []string
	for _, m := range regexp.MustCompile(`<p>([^<]*)</p>`).FindAllStringSubmatch(got, -1) {
		text = append(text, m[1])
	}
	if want := "天荒大陸。|東域一個小城。|“這個太初神殿是什麼情況？”|葉長安笑了。"; strings.Join(text, "|") != want {
		t.Errorf("joined text = %q, want %q", strings.Join(text, "|"), want)
	}
	if strings.Contains(got, "繁體") {
		t.Error("the 繁體/简体 switch was kept")
	}
	if strings.Contains(got, "(1/2)") || !strings.Contains(got, "<title>第1章 葉長安</title>") || !strings.Contains(got, "<h1>第1章 葉長安</h1>") {
		t.Error("the page marker was not dropped from the title")
	}
}

func TestNovel543JoinChapterPagesWithoutText(t *testing.T) {
	s := NewNovel543().(MultiPageChapter)
	if _, err := s.JoinChapterPages([]string{novel543Page("第1章 (1/2)", "A"), "<html><body>blocked</body></html>"}); err == nil {
		t.Error("a later page without text should fail rather than translate half a chapter")
	}
}
