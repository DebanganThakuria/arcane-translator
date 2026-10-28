package sources

import "testing"

func TestXswIds(t *testing.T) {
	s := NewXsw()

	novels := map[string]string{
		"https://m.xsw.tw/1758662/":            "xsw-1758662",
		"https://m.xsw.tw/1758662/page-1.html": "xsw-1758662",
		"https://m.xsw.tw/":                    "",
	}
	for url, want := range novels {
		if got := s.GetNovelId(url); got != want {
			t.Errorf("GetNovelId(%q) = %q, want %q", url, got, want)
		}
	}

	chapters := map[string]string{
		"https://m.xsw.tw/1758662/255833825.html": "xsw-1758662-255833825",
		"https://m.xsw.tw/1758662/page-1.html":    "",
	}
	for url, want := range chapters {
		if got := s.GetChapterId(url); got != want {
			t.Errorf("GetChapterId(%q) = %q, want %q", url, got, want)
		}
	}
}

func TestXswNextChapterUrl(t *testing.T) {
	s := NewXsw()
	current := "https://m.xsw.tw/1758662/255833825.html"

	tests := map[string]struct {
		page string
		want string
	}{
		"next chapter": {
			page: `<a href="/1758662/255822258.html">上一章</a><a href="/1758662/page-1.html">目錄</a><a href="/1758662/255859255.html">下一章</a>`,
			want: "https://m.xsw.tw/1758662/255859255.html",
		},
		"last chapter links back to the list": {
			page: `<a href="/1758662/page-1.html">下一章</a>`,
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

func TestXswCover(t *testing.T) {
	page := `<img src="/images/logo.png"><img src="https://img.xsw.tw/1758/1758662/1758662s.jpg" alt="苟聖！">`
	got, err := NewXsw().GetNovelCoverImageUrl(page)
	if err != nil {
		t.Fatal(err)
	}
	if got != "https://img.xsw.tw/1758/1758662/1758662s.jpg" {
		t.Errorf("cover = %q", got)
	}
}
