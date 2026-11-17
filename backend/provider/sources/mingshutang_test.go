package sources

import "testing"

func TestMingshutangIds(t *testing.T) {
	s := NewMingshutang()

	novels := map[string]string{
		"https://www.mingshutang.net/mingshu/51416484288/": "mingshutang-51416484288",
		"https://www.mingshutang.net/mingshu/51416484288":  "mingshutang-51416484288",
		"https://www.mingshutang.net/":                     "",
	}
	for url, want := range novels {
		if got := s.GetNovelId(url); got != want {
			t.Errorf("GetNovelId(%q) = %q, want %q", url, got, want)
		}
	}

	chapters := map[string]string{
		"https://www.mingshutang.net/mingshut/51416484288/71653685.html":   "mingshutang-51416484288-71653685",
		"https://www.mingshutang.net/mingshut/51416484288/71653685-2.html": "mingshutang-51416484288-71653685-2",
		"https://www.mingshutang.net/mingshu/51416484288/":                 "",
	}
	for url, want := range chapters {
		if got := s.GetChapterId(url); got != want {
			t.Errorf("GetChapterId(%q) = %q, want %q", url, got, want)
		}
	}
}

func TestMingshutangNextChapterUrl(t *testing.T) {
	s := NewMingshutang()
	const book = "https://www.mingshutang.net/mingshut/51416484288/"

	tests := map[string]struct {
		current, page, want string
	}{
		"first page goes to the second": {
			current: book + "71653685.html",
			page:    `<a href="71653686.html">下一章</a><div id="PageSet"><span>1</span><a href="71653685-2.html">2</a><a href="71653685-2.html">下一页</a></div><a href="71653686.html">下一章</a>`,
			want:    book + "71653685-2.html",
		},
		"last page goes to the next chapter": {
			current: book + "71653685-3.html",
			page:    `<a href="71653686.html">下一章</a><div id="PageSet"><a href="71653685-2.html">上一页</a><a href="71653685.html">1</a><span>3</span><span class="zw">下一页</span></div><a href="71653686.html">下一章</a>`,
			want:    book + "71653686.html",
		},
		"last chapter links back to the book": {
			current: book + "76825881.html",
			page:    `<span class="zw">下一页</span><a href="/mingshu/51416484288/">下一章</a>`,
			want:    "",
		},
		"captcha page": {
			current: book + "71653685.html",
			page:    `<form name="frmss" method="post"><b>页面安全验证</b><input type="text" name="yzm"></form>`,
			want:    "",
		},
	}
	for name, tt := range tests {
		got, err := s.GetNextChapterUrl(tt.page, tt.current)
		if err != nil {
			t.Fatalf("%s: unexpected error: %v", name, err)
		}
		if got != tt.want {
			t.Errorf("%s: got %q, want %q", name, got, tt.want)
		}
	}
}

func TestMingshutangCover(t *testing.T) {
	page := `<meta property="og:image" content="https://www.mingshutang.net/mingshuimg/0/487780.jpg" />`
	got, err := NewMingshutang().GetNovelCoverImageUrl(page)
	if err != nil {
		t.Fatal(err)
	}
	if got != "https://www.mingshutang.net/mingshuimg/0/487780.jpg" {
		t.Errorf("cover = %q", got)
	}
}

func TestMingshutangPrepareChapter(t *testing.T) {
	const book = "https://www.mingshutang.net/mingshut/51416484288/"

	tests := map[string]struct {
		url, page, want string
	}{
		"first page keeps its title and loses the read-on notice": {
			url:  book + "71653685.html",
			page: `<title>第1章 催收_原来我是反派天骄_名书堂</title><h1>第1章 催收</h1><div id="content"><p>　　“少东家。”…。。<br>　　本章未完，请点击下一页继续阅读！</p></div>`,
			want: `<title>第1章 催收_原来我是反派天骄_名书堂</title><h1>第1章 催收</h1><div id="content"><p>　　“少东家。”</p></div>`,
		},
		"a later page is numbered and loses the bookmark notice": {
			url:  book + "71653685-3.html",
			page: `<title>第1章 催收_原来我是反派天骄_名书堂</title><h1>第1章 催收</h1><div id="content"><p>　　身上。39195113。。<br />　　一寸欢喜提醒您：看完记得收藏【名书堂】 www.mingshutang.net，随时随地都可以畅阅无阻...</p></div>`,
			want: `<title>第1章 催收（3）_原来我是反派天骄_名书堂</title><h1>第1章 催收（3）</h1><div id="content"><p>　　身上。</p></div>`,
		},
		"ordinary text is left alone": {
			url:  book + "71653686.html",
			page: `<h1>第2章 水源</h1><p>　　“好。。”<br>他说。</p>`,
			want: `<h1>第2章 水源</h1><p>　　“好。。”<br>他说。</p>`,
		},
	}
	for name, tt := range tests {
		got, err := PrepareChapter("mingshutang", tt.page, tt.url)
		if err != nil {
			t.Fatalf("%s: unexpected error: %v", name, err)
		}
		if got != tt.want {
			t.Errorf("%s:\n got %q\nwant %q", name, got, tt.want)
		}
	}
}

func TestMingshutangRequiresManualHTML(t *testing.T) {
	if !RequiresManualHTML("mingshutang") {
		t.Error("mingshutang pages must be pasted in: the server gets a captcha form")
	}
}
