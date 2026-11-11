package sources

import "testing"

func TestCanghaishugeIds(t *testing.T) {
	s := NewCanghaishuge()

	novels := map[string]string{
		"https://www.canghaishuge.com/shenhuajiyuankaijusancunshihoujiyin/": "canghaishuge-shenhuajiyuankaijusancunshihoujiyin",
		"https://www.canghaishuge.com/shenhuajiyuankaijusancunshihoujiyin":  "canghaishuge-shenhuajiyuankaijusancunshihoujiyin",
		"https://www.canghaishuge.com/":                                     "",
	}
	for url, want := range novels {
		if got := s.GetNovelId(url); got != want {
			t.Errorf("GetNovelId(%q) = %q, want %q", url, got, want)
		}
	}

	chapters := map[string]string{
		"https://www.canghaishuge.com/shenhuajiyuankaijusancunshihoujiyin/1.html":   "canghaishuge-shenhuajiyuankaijusancunshihoujiyin-1",
		"https://www.canghaishuge.com/shenhuajiyuankaijusancunshihoujiyin/388.html": "canghaishuge-shenhuajiyuankaijusancunshihoujiyin-388",
		"https://www.canghaishuge.com/shenhuajiyuankaijusancunshihoujiyin/":         "",
	}
	for url, want := range chapters {
		if got := s.GetChapterId(url); got != want {
			t.Errorf("GetChapterId(%q) = %q, want %q", url, got, want)
		}
	}
}

func TestCanghaishugeNextChapterUrl(t *testing.T) {
	s := NewCanghaishuge()
	current := "https://www.canghaishuge.com/shenhuajiyuankaijusancunshihoujiyin/1.html"

	tests := map[string]struct {
		page string
		want string
	}{
		"next chapter": {
			page: `<a href="/shenhuajiyuankaijusancunshihoujiyin/">上一章</a><a href="/shenhuajiyuankaijusancunshihoujiyin/">章节目录</a><a href="/shenhuajiyuankaijusancunshihoujiyin/2.html">下一章</a>`,
			want: "https://www.canghaishuge.com/shenhuajiyuankaijusancunshihoujiyin/2.html",
		},
		"last chapter links back to the book": {
			page: `<a href="/shenhuajiyuankaijusancunshihoujiyin/">下一章</a>`,
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

func TestCanghaishugeCover(t *testing.T) {
	page := `<meta property="og:image" content="https://www.canghaishuge.com/storage/image/13/12784/shenhuajiyuankaijusancunshihoujiyin.webp">`
	got, err := NewCanghaishuge().GetNovelCoverImageUrl(page)
	if err != nil {
		t.Fatal(err)
	}
	if got != "https://www.canghaishuge.com/storage/image/13/12784/shenhuajiyuankaijusancunshihoujiyin.webp" {
		t.Errorf("cover = %q", got)
	}
}
