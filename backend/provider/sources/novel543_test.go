package sources

import "testing"

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
