package sources

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
)

func TestNovel85Ids(t *testing.T) {
	s := NewNovel85()

	novels := map[string]string{
		"https://www.85novel.com/book/2383736.html": "85novel-2383736",
		"https://www.85novel.com/book/abc.html":     "",
		"https://www.85novel.com/":                  "",
	}
	for url, want := range novels {
		if got := s.GetNovelId(url); got != want {
			t.Errorf("GetNovelId(%q) = %q, want %q", url, got, want)
		}
	}

	chapters := map[string]string{
		"https://www.85novel.com/book/2383736/463436.html": "85novel-2383736-463436",
		"https://www.85novel.com/book/2383736.html":        "",
	}
	for url, want := range chapters {
		if got := s.GetChapterId(url); got != want {
			t.Errorf("GetChapterId(%q) = %q, want %q", url, got, want)
		}
	}
}

func TestNovel85Cover(t *testing.T) {
	pages := map[string]string{
		// The site doubles the scheme.
		`<meta property="og:image" content="https:https://www.85novel.com/cover/2383736.jpg">`: "https://www.85novel.com/cover/2383736.jpg",
		`<meta property="og:image" content="https://www.85novel.com/cover/1.jpg">`:             "https://www.85novel.com/cover/1.jpg",
		`<meta property="og:image" content="/cover/2.jpg">`:                                    "https://www.85novel.com/cover/2.jpg",
		`<html></html>`: "",
	}
	for page, want := range pages {
		got, err := NewNovel85().GetNovelCoverImageUrl(page)
		if err != nil || got != want {
			t.Errorf("cover of %q = %q, %v; want %q", page, got, err, want)
		}
	}
}

// chapterPage builds a chapter page whose content element holds the given
// blocks of paragraphs, with an ad slot before each block as on the site.
func chapterPage(strategy string, mode int, blocks ...[]string) string {
	var b strings.Builder
	fmt.Fprintf(&b, `<html><body x-data="Login(0, '%s', '/')"><div data-chapter-content data-content-order-mode="%d" data-content-order-strategy="%s" data-content-restored="false">`, testNovel85Token, mode, strategy)
	for i, block := range blocks {
		fmt.Fprintf(&b, `<div class="ad">ad%d</div><div data-chapter-content-block style="display: contents">`, i)
		for _, p := range block {
			fmt.Fprintf(&b, `<p class="mt-3">%s</p>`, p)
		}
		b.WriteString(`</div>`)
	}
	b.WriteString(`</div></body></html>`)
	return b.String()
}

var paragraphOrAd = regexp.MustCompile(`<p class="mt-3">([^<]*)</p>|<div class="ad">([^<]*)</div>`)

// readingOrder lists the paragraphs and ad slots of a page in document order.
func readingOrder(page string) string {
	var order []string
	for _, m := range paragraphOrAd.FindAllStringSubmatch(page, -1) {
		order = append(order, m[1]+m[2])
	}
	return strings.Join(order, " ")
}

func TestNovel85PrepareChapter(t *testing.T) {
	tests := []struct {
		name string
		page string
		want string
	}{
		{
			// Mode 0 keeps the odd paragraphs in the first half: B D | A C E.
			name: "split mode 0",
			page: chapterPage("split", 0, []string{"B", "D", "A", "C", "E"}),
			want: "ad0 A B C D E",
		},
		{
			// Mode 3 keeps the even paragraphs first: A C E | B D.
			name: "split mode 3",
			page: chapterPage("split", 3, []string{"A", "C", "E", "B", "D"}),
			want: "ad0 A B C D E",
		},
		{
			// Mode 1 also reverses each block and the order of the blocks. The
			// ad slots stay where they are.
			name: "split mode 1 across two blocks",
			page: chapterPage("split", 1, []string{"H", "F", "G", "E"}, []string{"D", "B", "C", "A"}),
			want: "ad0 A B C D ad1 E F G H",
		},
		{
			name: "reverse mode 2 reverses paragraphs only",
			page: chapterPage("reverse", 2, []string{"C", "B", "A"}, []string{"F", "E", "D"}),
			want: "ad0 A B C ad1 D E F",
		},
		{
			name: "reverse mode 3 reverses blocks only",
			page: chapterPage("reverse", 3, []string{"D", "E", "F"}, []string{"A", "B", "C"}),
			want: "ad0 A B C ad1 D E F",
		},
		{
			name: "reverse mode 0 is already in order",
			page: chapterPage("reverse", 0, []string{"A", "B", "C"}),
			want: "ad0 A B C",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newNovel85Test(t, &novel85Site{reply: titleReply})
			got, err := s.PrepareChapter(tt.page, testChapterURL)
			if err != nil {
				t.Fatalf("PrepareChapter() error = %v", err)
			}
			if order := readingOrder(got); order != tt.want {
				t.Errorf("reading order = %q, want %q", order, tt.want)
			}
			if !strings.Contains(got, `data-content-restored="true"`) {
				t.Error("restored page is not marked restored")
			}
			if !strings.Contains(got, "<h1>第527章 斬殺</h1>") {
				t.Error("chapter title was not added")
			}
		})
	}
}

func TestNovel85PrepareChapterWithoutTitle(t *testing.T) {
	// The title request fails; the chapter is still restored and translated.
	s := newNovel85Test(t, &novel85Site{reply: `{"code":1,"result":null}`})
	got, err := s.PrepareChapter(chapterPage("split", 0, []string{"B", "A"}), testChapterURL)
	if err != nil {
		t.Fatalf("PrepareChapter() error = %v", err)
	}
	if order := readingOrder(got); order != "ad0 A B" || strings.Contains(got, "<h1>") {
		t.Errorf("got order %q and a title; want the restored text without one", order)
	}
}

func TestNovel85PrepareChapterLeavesOtherPagesAlone(t *testing.T) {
	s := newNovel85Test(t, &novel85Site{reply: titleReply})

	page := `<html><body><p>A</p><p>B</p></body></html>`
	if got, err := s.PrepareChapter(page, testChapterURL); err != nil || got != page {
		t.Errorf("a page without the content element changed or failed: %v", err)
	}

	restored := strings.Replace(chapterPage("split", 0, []string{"B", "A"}), `data-content-restored="false"`, `data-content-restored="true"`, 1)
	got, err := s.PrepareChapter(restored, testChapterURL)
	if err != nil || readingOrder(got) != "ad0 B A" {
		t.Errorf("an already restored page was reordered: %q, %v", readingOrder(got), err)
	}

	if _, err := s.PrepareChapter(chapterPage("shuffle", 0, []string{"A"}), testChapterURL); err == nil {
		t.Error("an unknown strategy should fail rather than translate shuffled text")
	}
}

func TestPrepareChapterForOtherSources(t *testing.T) {
	page := chapterPage("split", 0, []string{"B", "A"})
	if got, err := PrepareChapter("novel543", page, testChapterURL); err != nil || got != page {
		t.Errorf("PrepareChapter for a source without a preparer changed the page: %v", err)
	}
}

const (
	testNovel85Token = "2|09d74425|a7102c9b4f601850a3a6612ebd849e36|1792823901"
	testChapterURL   = "https://www.85novel.com/book/2383736/463436.html"
	titleReply       = `{"code":0,"result":{"cInfo":{"ChapterName":"第527章 斬殺","ChapterUrl":"2383736/463436.html"},"nCInfo":{"ChapterUrl":"2383736/463446.html"}}}`
)

// novel85Site fakes the site's chapter page and its next-chapter reply.
type novel85Site struct {
	reply    string
	requests []*http.Request
	bodies   []string
}

func (s *novel85Site) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	s.requests = append(s.requests, r)
	s.bodies = append(s.bodies, string(body))

	if r.Method == http.MethodGet {
		fmt.Fprintf(w, `<body x-data="Login(0, '%s', '/')"></body>`, testNovel85Token)
		return
	}
	if cookie, err := r.Cookie("_xsrf"); err != nil || cookie.Value != r.URL.Query().Get("_xsrf") {
		http.Error(w, "403: Forbidden", http.StatusForbidden)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, s.reply)
}

func newNovel85Test(t *testing.T, site *novel85Site) novel85 {
	t.Helper()
	srv := httptest.NewServer(site)
	t.Cleanup(srv.Close)
	return novel85{baseURL: srv.URL, client: srv.Client()}
}

func TestNovel85NextChapterUrl(t *testing.T) {
	site := &novel85Site{reply: `{"code":0,"result":{"nCInfo":{"ChapterName":"第528章 傳承","ChapterUrl":"2383736/463446.html","ChapterID":463446}}}`}
	s := newNovel85Test(t, site)
	page := fmt.Sprintf(`<body x-data="Login(0, '%s', '/')">`, testNovel85Token)

	got, err := s.GetNextChapterUrl(page, "https://www.85novel.com/book/2383736/463436.html")
	if err != nil {
		t.Fatalf("GetNextChapterUrl() error = %v", err)
	}
	if want := s.baseURL + "/book/2383736/463446.html"; got != want {
		t.Errorf("GetNextChapterUrl() = %q, want %q", got, want)
	}

	if len(site.requests) != 1 {
		t.Fatalf("got %d requests, want one POST using the page's token", len(site.requests))
	}
	req := site.requests[0]
	if req.Method != http.MethodPost || req.URL.Path != "/book/2383736/463436.html" {
		t.Errorf("request = %s %s, want POST to the chapter path", req.Method, req.URL.Path)
	}
	var body map[string]string
	if err := json.Unmarshal([]byte(site.bodies[0]), &body); err != nil || body["action"] != "auth" {
		t.Errorf("request body = %q, want {\"action\":\"auth\"}", site.bodies[0])
	}
}

func TestNovel85NextChapterUrlFetchesTokenForPastedPage(t *testing.T) {
	site := &novel85Site{reply: `{"code":0,"result":{"nCInfo":{"ChapterUrl":"2383736/463446.html"}}}`}
	s := newNovel85Test(t, site)

	got, err := s.GetNextChapterUrl("<html>no token here</html>", "https://www.85novel.com/book/2383736/463436.html")
	if err != nil || got != s.baseURL+"/book/2383736/463446.html" {
		t.Fatalf("GetNextChapterUrl() = %q, %v", got, err)
	}
	if len(site.requests) != 2 || site.requests[0].Method != http.MethodGet {
		t.Errorf("want a GET for the token, then the POST; got %d requests", len(site.requests))
	}
}

func TestNovel85NextChapterUrlEdges(t *testing.T) {
	page := fmt.Sprintf(`Login(0, '%s', '/')`, testNovel85Token)
	current := "https://www.85novel.com/book/2383736/463436.html"

	tests := []struct {
		name    string
		reply   string
		current string
		want    string
		wantErr bool
	}{
		{name: "last chapter", reply: `{"code":0,"result":{"nCInfo":null}}`, current: current},
		{name: "site error", reply: `{"code":1,"result":null}`, current: current, wantErr: true},
		{name: "next link outside the book pages", reply: `{"code":0,"result":{"nCInfo":{"ChapterUrl":"../../admin"}}}`, current: current, wantErr: true},
		// The request must only ever go to the site, whatever URL the reader enters.
		{name: "not a chapter URL", current: "http://169.254.169.254/latest/meta-data", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			site := &novel85Site{reply: tt.reply}
			s := newNovel85Test(t, site)

			got, err := s.GetNextChapterUrl(page, tt.current)
			if (err != nil) != tt.wantErr || got != tt.want {
				t.Errorf("GetNextChapterUrl() = %q, %v; want %q, error %v", got, err, tt.want, tt.wantErr)
			}
			if tt.name == "not a chapter URL" && len(site.requests) != 0 {
				t.Error("made a request for a URL that is not a chapter page")
			}
		})
	}
}

func TestNovel85RejectsMalformedToken(t *testing.T) {
	for _, page := range []string{
		`Login(0, '2|x|y|z', '/')`,
		"Login(0, '2|09d74425|a7102c9b4f601850a3a6612ebd849e36|1\r\nX-Evil: 1', '/')",
	} {
		if got := novel85Token(page); got != "" {
			t.Errorf("novel85Token(%q) = %q, want it refused", page, got)
		}
	}
}
