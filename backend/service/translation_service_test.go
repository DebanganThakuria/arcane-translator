package service

import (
	"errors"
	"testing"

	"backend/models"
	"backend/repo"
)

func TestPageHTMLUsesPastedHTML(t *testing.T) {
	pasted := "<html>pasted</html>"

	got, err := pageHTML("ixdzs", "https://ixdzs.tw/read/526058/p1.html", &pasted)
	if err != nil {
		t.Fatalf("pageHTML returned error: %v", err)
	}
	if got != pasted {
		t.Fatalf("pageHTML = %q, want the pasted HTML", got)
	}
}

func TestPageHTMLRefusesToScrapeManualSources(t *testing.T) {
	// An unroutable URL: if the refusal regressed, the test fails on the
	// error type rather than quietly reaching the real site.
	_, err := pageHTML("ixdzs", "http://127.0.0.1:1/read/526058/", nil)
	if !errors.Is(err, ErrManualHTMLRequired) {
		t.Fatalf("pageHTML error = %v, want ErrManualHTMLRequired", err)
	}
}

// chapterRepo serves one novel and its last chapter, and records the update.
type chapterRepo struct {
	repo.Repo
	novel   *models.Novel
	last    *models.Chapter
	updated *models.Chapter
}

func (r *chapterRepo) GetNovelByID(string) (*models.Novel, error) { return r.novel, nil }

func (r *chapterRepo) GetLastChapter(string) (*models.Chapter, error) {
	// A copy, so a test sees only what the service saved.
	last := *r.last
	return &last, nil
}

func (r *chapterRepo) UpdateChapter(chapter *models.Chapter) error {
	r.updated = chapter
	return nil
}

func TestRefreshNextChapterURL(t *testing.T) {
	const (
		lastURL = "https://69shuba.tw/read/88917/491713"
		nextURL = "https://69shuba.tw/read/88917/491715"
	)
	pasted := func(page string) *string { return &page }

	tests := map[string]struct {
		source      string
		linked      string
		page        *string
		wantNext    string
		wantUpdated bool
	}{
		"a new chapter is out": {
			source:      "69shubatw",
			page:        pasted(`<a href="/read/88917/491715">下一章</a>`),
			wantNext:    nextURL,
			wantUpdated: true,
		},
		"no new chapter yet": {
			source: "69shubatw",
			page:   pasted(`<a href="/book/88917/">下一章</a>`),
		},
		// ixdzs refuses server fetches, so reaching for the page would fail.
		"a chapter that already links onward is not fetched": {
			source:   "ixdzs",
			linked:   nextURL,
			wantNext: nextURL,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			r := &chapterRepo{
				novel: &models.Novel{ID: "novel", Source: tt.source},
				last:  &models.Chapter{Number: 10, URL: lastURL, NextChapterURL: tt.linked},
			}

			got, err := NewTranslationService(r).RefreshNextChapterURL(&models.NextChapterRefreshRequest{NovelID: "novel", HTMLContent: tt.page})
			if err != nil {
				t.Fatalf("RefreshNextChapterURL: %v", err)
			}
			if got.NextChapterURL != tt.wantNext {
				t.Errorf("NextChapterURL = %q, want %q", got.NextChapterURL, tt.wantNext)
			}
			if updated := r.updated != nil; updated != tt.wantUpdated {
				t.Errorf("chapter saved = %v, want %v", updated, tt.wantUpdated)
			}
		})
	}
}

func TestRefreshNextChapterURLAsksForThePageItCannotFetch(t *testing.T) {
	r := &chapterRepo{
		novel: &models.Novel{ID: "novel", Source: "ixdzs"},
		last:  &models.Chapter{Number: 10, URL: "http://127.0.0.1:1/read/526058/p10.html"},
	}

	_, err := NewTranslationService(r).RefreshNextChapterURL(&models.NextChapterRefreshRequest{NovelID: "novel"})
	if !errors.Is(err, ErrManualHTMLRequired) {
		t.Fatalf("error = %v, want ErrManualHTMLRequired", err)
	}
}
