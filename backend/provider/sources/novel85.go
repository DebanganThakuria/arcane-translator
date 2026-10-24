package sources

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html"
)

const (
	novel85BaseURL = "https://www.85novel.com"

	novel85UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 " +
		"(KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36"

	// novel85MaxReply bounds the next-chapter reply; a real one is under 1 KB.
	novel85MaxReply = 64 << 10
)

var (
	// The reader page hands its XSRF token to the login widget, as in
	// Login(0, '2|09d74425|a7102c9b4f601850a3a6612ebd849e36|1792823901', '/').
	novel85TokenPattern = regexp.MustCompile(`Login\(\s*\d+\s*,\s*'([^']*)'`)
	// A token is version|mask|hex|timestamp. Anything else is refused before it
	// reaches a header or query string.
	novel85ValidToken = regexp.MustCompile(`^\d\|[0-9a-f]{8}\|[0-9a-f]{32}\|\d{1,12}$`)
	// Chapter pages live at /book/<novel>/<chapter>.html.
	novel85ChapterPath = regexp.MustCompile(`^/book/\d+/\d+\.html$`)
	// The site writes a doubled scheme: content="https:https://www.85novel.com/cover/2383736.jpg".
	novel85CoverPattern = regexp.MustCompile(`<meta\s+property="og:image"\s+content="([^"]+)"`)
)

// novel85 reads 85novel.com, which works against scrapers. A chapter page
// shuffles its paragraphs and restores their order in script, and its title
// and next-chapter link come only from a second request that the page's
// script makes. PrepareChapter restores the order and adds the title;
// GetNextChapterUrl asks for the next link.
type novel85 struct {
	baseURL string
	client  *http.Client
}

func NewNovel85() Source {
	return &novel85{baseURL: novel85BaseURL, client: &http.Client{Timeout: 25 * time.Second}}
}

func (n novel85) GetNovelId(novelURL string) string {
	// Example URL: https://www.85novel.com/book/2383736.html
	// Prefixed because the bare ids are numeric and could collide with another
	// source's novel ids.
	u, err := url.Parse(novelURL)
	if err != nil {
		return ""
	}

	id, ok := strings.CutPrefix(u.Path, "/book/")
	if !ok {
		return ""
	}
	id = strings.TrimSuffix(id, ".html")
	if _, err := strconv.Atoi(id); err != nil {
		return ""
	}

	return "85novel-" + id
}

func (n novel85) GetChapterId(chapterURL string) string {
	// Example URL: https://www.85novel.com/book/2383736/463436.html
	u, err := url.Parse(chapterURL)
	if err != nil || !novel85ChapterPath.MatchString(u.Path) {
		return ""
	}

	parts := strings.Split(strings.TrimSuffix(strings.TrimPrefix(u.Path, "/book/"), ".html"), "/")
	return "85novel-" + parts[0] + "-" + parts[1]
}

// novel85Chapter names a chapter in the site's reply.
type novel85Chapter struct {
	ChapterName string `json:"ChapterName"`
	ChapterURL  string `json:"ChapterUrl"`
}

// novel85Reply is the part of the site's reply to {"action": "auth"} this
// source reads: the current chapter and the one after it.
type novel85Reply struct {
	Code   int `json:"code"`
	Result struct {
		Current *novel85Chapter `json:"cInfo"`
		Next    *novel85Chapter `json:"nCInfo"`
	} `json:"result"`
}

// chapterInfo asks the site about the chapter at chapterURL, the same request
// its reader script makes. The request carries the page's own XSRF token as
// both cookie and parameter, which the site accepts. A pasted page carries a
// token too; without one the page is fetched for it.
func (n novel85) chapterInfo(chapterContent, chapterURL string) (*novel85Reply, error) {
	u, err := url.Parse(chapterURL)
	if err != nil || !novel85ChapterPath.MatchString(u.Path) {
		return nil, fmt.Errorf("85novel: not a chapter URL: %q", chapterURL)
	}
	// Only the validated path is kept: the request always goes to this site.
	pageURL := n.baseURL + u.Path

	token := novel85Token(chapterContent)
	if token == "" {
		page, err := n.get(pageURL)
		if err != nil {
			return nil, err
		}
		if token = novel85Token(page); token == "" {
			return nil, errors.New("85novel: chapter page has no XSRF token")
		}
	}

	req, err := http.NewRequest(http.MethodPost, pageURL+"?_xsrf="+url.QueryEscape(token), bytes.NewBufferString(`{"action":"auth"}`))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", novel85UserAgent)
	req.Header.Set("Referer", pageURL)
	req.Header.Set("Cookie", "_xsrf="+token)

	resp, err := n.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("85novel: chapter info: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("85novel: chapter info: status %d", resp.StatusCode)
	}

	var reply novel85Reply
	if err := json.NewDecoder(io.LimitReader(resp.Body, novel85MaxReply)).Decode(&reply); err != nil {
		return nil, fmt.Errorf("85novel: chapter info: %w", err)
	}
	if reply.Code != 0 {
		return nil, fmt.Errorf("85novel: chapter info: site returned code %d", reply.Code)
	}

	return &reply, nil
}

func (n novel85) GetNextChapterUrl(chapterContent, currentChapterURL string) (string, error) {
	reply, err := n.chapterInfo(chapterContent, currentChapterURL)
	if err != nil {
		return "", err
	}
	// The last chapter has no next chapter.
	if reply.Result.Next == nil || reply.Result.Next.ChapterURL == "" {
		return "", nil
	}

	next := "/book/" + strings.TrimPrefix(reply.Result.Next.ChapterURL, "/")
	if !novel85ChapterPath.MatchString(next) {
		return "", fmt.Errorf("85novel: unexpected next chapter %q", reply.Result.Next.ChapterURL)
	}

	return n.baseURL + next, nil
}

func (n novel85) get(pageURL string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, pageURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", novel85UserAgent)

	resp, err := n.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("85novel: fetch chapter: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("85novel: fetch chapter: status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
	return string(body), err
}

func novel85Token(page string) string {
	match := novel85TokenPattern.FindStringSubmatch(page)
	if match == nil || !novel85ValidToken.MatchString(match[1]) {
		return ""
	}

	return match[1]
}

func (n novel85) GetNovelCoverImageUrl(pageContent string) (string, error) {
	match := novel85CoverPattern.FindStringSubmatch(pageContent)
	if match == nil {
		return "", nil
	}

	cover := match[1]
	for strings.HasPrefix(cover, "https:https:") || strings.HasPrefix(cover, "http:https:") {
		cover = cover[strings.Index(cover, ":")+1:]
	}
	if strings.HasPrefix(cover, "/") && !strings.HasPrefix(cover, "//") {
		cover = n.baseURL + cover
	}

	return cover, nil
}

// PrepareChapter puts a chapter page's paragraphs back in reading order and
// adds the chapter title, which the page leaves for its script to fill in.
// A page without the site's content element comes back unchanged. If the
// title cannot be fetched, the chapter is still translated without one.
func (n novel85) PrepareChapter(chapterContent, chapterURL string) (string, error) {
	doc, err := html.Parse(strings.NewReader(chapterContent))
	if err != nil {
		return "", err
	}

	content := findHTMLNode(doc, func(node *html.Node) bool { return hasHTMLAttr(node, "data-chapter-content") })
	if content == nil {
		return chapterContent, nil
	}

	if err := restoreParagraphOrder(content); err != nil {
		return "", err
	}

	reply, err := n.chapterInfo(chapterContent, chapterURL)
	switch {
	case err != nil:
		log.Printf("85novel: chapter title: %v", err)
	case reply.Result.Current != nil && reply.Result.Current.ChapterName != "":
		title := &html.Node{Type: html.ElementNode, Data: "h1"}
		title.AppendChild(&html.Node{Type: html.TextNode, Data: reply.Result.Current.ChapterName})
		content.InsertBefore(title, content.FirstChild)
	}

	var out strings.Builder
	if err := html.Render(&out, doc); err != nil {
		return "", err
	}
	return out.String(), nil
}

// restoreParagraphOrder undoes the site's paragraph shuffle in place.
//
// The site sends paragraphs shuffled and lets its reader script reorder
// them, as the content element describes with data-content-order-strategy
// ("split" or "reverse") and data-content-order-mode (0 to 3). This follows
// the same steps as the site's restoreChapterContent.
func restoreParagraphOrder(content *html.Node) error {
	if htmlAttr(content, "data-content-restored") == "true" {
		return nil
	}

	mode, err := strconv.Atoi(htmlAttr(content, "data-content-order-mode"))
	strategy := htmlAttr(content, "data-content-order-strategy")
	if err != nil || mode < 0 || mode > 3 || (strategy != "split" && strategy != "reverse") {
		return fmt.Errorf("85novel: unknown paragraph order %q mode %q", strategy, htmlAttr(content, "data-content-order-mode"))
	}

	var blocks []*html.Node
	for child := content.FirstChild; child != nil; child = child.NextSibling {
		if hasHTMLAttr(child, "data-chapter-content-block") {
			blocks = append(blocks, child)
		}
	}

	for _, block := range blocks {
		var paragraphs []*html.Node
		collectParagraphs(block, &paragraphs)

		switch {
		case strategy == "split":
			paragraphs = restoreSplitOrder(paragraphs, mode)
		case mode == 1 || mode == 2:
			reverseNodes(paragraphs)
		}
		for _, p := range paragraphs {
			p.Parent.RemoveChild(p)
			block.AppendChild(p)
		}
	}

	reverseBlocks := mode == 1 || mode == 3
	if strategy == "split" {
		reverseBlocks = mode == 1 || mode == 2
	}
	if reverseBlocks && len(blocks) > 1 {
		// Swap the text blocks among their own positions; ad slots between them
		// stay put.
		markers := make([]*html.Node, len(blocks))
		for i, block := range blocks {
			markers[i] = &html.Node{Type: html.CommentNode}
			content.InsertBefore(markers[i], block)
			content.RemoveChild(block)
		}
		for i, marker := range markers {
			content.InsertBefore(blocks[len(blocks)-1-i], marker)
			content.RemoveChild(marker)
		}
	}

	setHTMLAttr(content, "data-content-restored", "true")
	return nil
}

// restoreSplitOrder undoes the "split" shuffle: one half of the block holds
// the even paragraphs and the other half the odd ones.
func restoreSplitOrder(items []*html.Node, mode int) []*html.Node {
	grouped := append([]*html.Node(nil), items...)
	if mode == 1 || mode == 2 {
		reverseNodes(grouped)
	}

	evenFirst := mode == 1 || mode == 3
	firstCount := len(grouped) / 2
	if evenFirst {
		firstCount = (len(grouped) + 1) / 2
	}
	even, odd := grouped[:firstCount], grouped[firstCount:]
	if !evenFirst {
		even, odd = odd, even
	}

	restored := make([]*html.Node, 0, len(grouped))
	for i := 0; i < len(even) || i < len(odd); i++ {
		if i < len(even) {
			restored = append(restored, even[i])
		}
		if i < len(odd) {
			restored = append(restored, odd[i])
		}
	}
	return restored
}

func collectParagraphs(node *html.Node, paragraphs *[]*html.Node) {
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if child.Type == html.ElementNode && child.Data == "p" {
			*paragraphs = append(*paragraphs, child)
			continue
		}
		collectParagraphs(child, paragraphs)
	}
}

func findHTMLNode(node *html.Node, match func(*html.Node) bool) *html.Node {
	if match(node) {
		return node
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if found := findHTMLNode(child, match); found != nil {
			return found
		}
	}
	return nil
}

func hasHTMLAttr(node *html.Node, key string) bool {
	if node.Type != html.ElementNode {
		return false
	}
	for _, a := range node.Attr {
		if a.Key == key {
			return true
		}
	}
	return false
}

func htmlAttr(node *html.Node, key string) string {
	for _, a := range node.Attr {
		if a.Key == key {
			return a.Val
		}
	}
	return ""
}

func setHTMLAttr(node *html.Node, key, value string) {
	for i := range node.Attr {
		if node.Attr[i].Key == key {
			node.Attr[i].Val = value
			return
		}
	}
	node.Attr = append(node.Attr, html.Attribute{Key: key, Val: value})
}

func reverseNodes(nodes []*html.Node) {
	for i, j := 0, len(nodes)-1; i < j; i, j = i+1, j-1 {
		nodes[i], nodes[j] = nodes[j], nodes[i]
	}
}
