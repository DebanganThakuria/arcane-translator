package sources

import "testing"

// 69shuba and twkan share one layout: chapters at /txt/<book>/<chapter>,
// linked by 下一章, which on a book's latest chapter leads to a placeholder.
func TestTxtSourcesNextChapterUrl(t *testing.T) {
	tests := map[string]struct {
		source, current, page, want string
	}{
		"twkan next chapter": {
			source:  "twkan",
			current: "https://twkan.com/txt/60148/58840351",
			page:    `<a href="https://twkan.com/txt/60148/58925199">下一章</a>`,
			want:    "https://twkan.com/txt/60148/58925199",
		},
		"twkan latest chapter links to end.html": {
			source:  "twkan",
			current: "https://twkan.com/txt/60148/58925199",
			page:    `<a href="https://twkan.com/txt/60148/end.html">下一章</a>`,
			want:    "",
		},
		"69shuba next chapter": {
			source:  "69shuba",
			current: "https://www.69shuba.com/txt/89532/40178501",
			page:    `<a href="https://www.69shuba.com/txt/89532/40178502">下一章</a>`,
			want:    "https://www.69shuba.com/txt/89532/40178502",
		},
		"69shuba latest chapter links to the book": {
			source:  "69shuba",
			current: "https://www.69shuba.com/txt/89532/40178501",
			page:    `<a href="https://www.69shuba.com/book/89532.htm">下一章</a>`,
			want:    "",
		},
	}
	for name, tt := range tests {
		got, err := GetSource(tt.source).GetNextChapterUrl(tt.page, tt.current)
		if err != nil {
			t.Fatalf("%s: unexpected error: %v", name, err)
		}
		if got != tt.want {
			t.Errorf("%s: got %q, want %q", name, got, tt.want)
		}
	}
}

func TestTxtSourcesChapterIds(t *testing.T) {
	tests := map[string]struct {
		source, url, want string
	}{
		"twkan chapter":         {"twkan", "https://twkan.com/txt/60148/58925199", "58925199_60148"},
		"twkan end placeholder": {"twkan", "https://twkan.com/txt/60148/end.html", ""},
		"69shuba chapter":       {"69shuba", "https://www.69shuba.com/txt/89532/40178501", "40178501"},
		"69shuba book page":     {"69shuba", "https://www.69shuba.com/book/89532.htm", ""},
	}
	for name, tt := range tests {
		if got := GetSource(tt.source).GetChapterId(tt.url); got != tt.want {
			t.Errorf("%s: GetChapterId(%q) = %q, want %q", name, tt.url, got, tt.want)
		}
	}
}
