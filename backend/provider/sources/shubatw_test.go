package sources

import "testing"

func TestShubaTWIds(t *testing.T) {
	s := NewShubaTW()

	novels := map[string]string{
		"https://69shuba.tw/book/88917/": "69shubatw-88917",
		"https://69shuba.tw/book/88917":  "69shubatw-88917",
		"https://69shuba.tw/":            "",
	}
	for url, want := range novels {
		if got := s.GetNovelId(url); got != want {
			t.Errorf("GetNovelId(%q) = %q, want %q", url, got, want)
		}
	}

	chapters := map[string]string{
		"https://69shuba.tw/read/88917/491713": "69shubatw-88917-491713",
		"https://69shuba.tw/book/88917/":       "",
	}
	for url, want := range chapters {
		if got := s.GetChapterId(url); got != want {
			t.Errorf("GetChapterId(%q) = %q, want %q", url, got, want)
		}
	}
}

func TestShubaTWNextChapterUrl(t *testing.T) {
	s := NewShubaTW()
	current := "https://69shuba.tw/read/88917/491713"

	tests := map[string]struct {
		page string
		want string
	}{
		"next chapter": {
			page: `<a href="/book/88917/">上一章</a><a href="/indexlist/88917/">目錄</a><a href="/read/88917/491715">下一章</a>`,
			want: "https://69shuba.tw/read/88917/491715",
		},
		"last page links back to the book": {
			page: `<a href="/book/88917/">下一章</a>`,
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

func TestShubaTWCover(t *testing.T) {
	page := `<meta property="og:image" content="//p.69shuba.tw/88/88917/88917s.jpg"/>`
	got, err := NewShubaTW().GetNovelCoverImageUrl(page)
	if err != nil {
		t.Fatal(err)
	}
	if got != "https://p.69shuba.tw/88/88917/88917s.jpg" {
		t.Errorf("cover = %q", got)
	}
}

func TestShubaTWPrepareChapter(t *testing.T) {
	tests := map[string]struct {
		page, want string
	}{
		"single-page marker is dropped": {
			page: `<title>斬妖除魔(鴿一刀)_第1章 活命（1 / 1）_斬妖除魔</title><h1>第1章 活命(1 / 1)</h1>`,
			want: `<title>斬妖除魔(鴿一刀)_第1章 活命_斬妖除魔</title><h1>第1章 活命</h1>`,
		},
		"a split chapter's marker is kept": {
			page: `<h1>第2章 長章(1 / 2)</h1>`,
			want: `<h1>第2章 長章(1 / 2)</h1>`,
		},
	}
	for name, tt := range tests {
		got, err := PrepareChapter("69shubatw", tt.page, "https://69shuba.tw/read/88917/491713")
		if err != nil {
			t.Fatalf("%s: unexpected error: %v", name, err)
		}
		if got != tt.want {
			t.Errorf("%s: got %q, want %q", name, got, tt.want)
		}
	}
}
