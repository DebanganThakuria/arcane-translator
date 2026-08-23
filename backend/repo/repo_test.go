package repo

import (
	"path/filepath"
	"testing"
	"time"

	"backend/models"
)

// newTestRepo builds a repository over a throwaway database, exercising the
// real schema and indexes rather than a stub.
func newTestRepo(t *testing.T) Repo {
	t.Helper()

	db, err := NewSQLiteDB(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	return NewRepo(db)
}

func seedNovel(t *testing.T, r Repo, novel *models.Novel) {
	t.Helper()
	if _, err := r.CreateNovel(novel); err != nil {
		t.Fatalf("seed novel %s: %v", novel.ID, err)
	}
}

func seedLibrary(t *testing.T, r Repo) {
	t.Helper()

	now := time.Now().Unix()
	seedNovel(t, r, &models.Novel{
		ID: "1", Title: "Azure Sword Immortal", Source: "69shuba", URL: "https://a",
		Summary: "s", Author: "Ling", Status: "Ongoing", Genres: []string{"Xianxia", "Action"},
		ChaptersCount: 300, LastReadTimestamp: now, LastUpdated: now, DateAdded: now - 900,
	})
	seedNovel(t, r, &models.Novel{
		ID: "2", Title: "beacon of the deep", Source: "syosetu", URL: "https://b",
		Summary: "s", Author: "Mori", Status: "Completed", Genres: []string{"Mystery"},
		ChaptersCount: 120, LastReadTimestamp: now - 100, LastUpdated: now - 50, DateAdded: now - 500,
	})
	seedNovel(t, r, &models.Novel{
		ID: "3", Title: "Crimson Ledger", Source: "royalroad", URL: "https://c",
		Summary: "s", Author: "Vale", Status: "Ongoing", Genres: []string{"Action"},
		ChaptersCount: 900, LastReadTimestamp: now - 200, LastUpdated: now - 10, DateAdded: now - 100,
	})
}

func TestQueryNovelsFilters(t *testing.T) {
	r := newTestRepo(t)
	seedLibrary(t, r)

	tests := []struct {
		name  string
		query models.NovelQuery
		want  []string
	}{
		{"no filter returns all", models.NovelQuery{}, []string{"1", "2", "3"}},
		{"by status", models.NovelQuery{Status: "Ongoing"}, []string{"1", "3"}},
		{
			"status is case-insensitive",
			models.NovelQuery{Status: "completed"},
			[]string{"2"},
		},
		{"by source", models.NovelQuery{SourceIDs: []string{"syosetu", "royalroad"}}, []string{"2", "3"}},
		{"by genre", models.NovelQuery{Genre: "Action"}, []string{"1", "3"}},
		{"search matches title", models.NovelQuery{Search: "ledger"}, []string{"3"}},
		{"search matches author", models.NovelQuery{Search: "Mori"}, []string{"2"}},
		{"search matches genre", models.NovelQuery{Search: "Xianxia"}, []string{"1"}},
		{"search is case-insensitive", models.NovelQuery{Search: "AZURE"}, []string{"1"}},
		{"no match", models.NovelQuery{Search: "nothing here"}, nil},
		{
			"filters compose with AND",
			models.NovelQuery{Status: "Ongoing", Genre: "Action", Search: "crimson"},
			[]string{"3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.query.Limit = 50
			tt.query.Sort = models.SortTitle
			tt.query.Ascending = true

			novels, count, err := r.QueryNovels(tt.query)
			if err != nil {
				t.Fatalf("QueryNovels: %v", err)
			}

			if count != len(tt.want) {
				t.Errorf("count = %d, want %d", count, len(tt.want))
			}

			got := make([]string, len(novels))
			for i, n := range novels {
				got[i] = n.ID
			}
			if len(got) != len(tt.want) {
				t.Fatalf("ids = %v, want %v", got, tt.want)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Fatalf("ids = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestQueryNovelsSorting(t *testing.T) {
	r := newTestRepo(t)
	seedLibrary(t, r)

	// Title ascending must be case-insensitive, or "beacon" sorts after "Crimson".
	novels, _, err := r.QueryNovels(models.NovelQuery{
		Sort: models.SortTitle, Ascending: true, Limit: 50,
	})
	if err != nil {
		t.Fatal(err)
	}
	if novels[0].ID != "1" || novels[1].ID != "2" || novels[2].ID != "3" {
		t.Errorf("title ascending gave %s, %s, %s; want 1, 2, 3 (case-insensitive)",
			novels[0].ID, novels[1].ID, novels[2].ID)
	}

	novels, _, err = r.QueryNovels(models.NovelQuery{Sort: models.SortChapterCount, Limit: 50})
	if err != nil {
		t.Fatal(err)
	}
	if novels[0].ID != "3" {
		t.Errorf("chapters descending put %s first, want 3 (900 chapters)", novels[0].ID)
	}
}

// The count must describe the whole filtered set, not the returned page.
func TestQueryNovelsPagination(t *testing.T) {
	r := newTestRepo(t)
	seedLibrary(t, r)

	first, count, err := r.QueryNovels(models.NovelQuery{
		Sort: models.SortTitle, Ascending: true, Limit: 2, Offset: 0,
	})
	if err != nil {
		t.Fatal(err)
	}
	if count != 3 {
		t.Errorf("total count = %d, want 3", count)
	}
	if len(first) != 2 {
		t.Fatalf("page size = %d, want 2", len(first))
	}

	second, _, err := r.QueryNovels(models.NovelQuery{
		Sort: models.SortTitle, Ascending: true, Limit: 2, Offset: 2,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(second) != 1 || second[0].ID == first[0].ID {
		t.Errorf("second page = %v, want the remaining novel", second)
	}
}

func TestQueryChapters(t *testing.T) {
	r := newTestRepo(t)
	seedNovel(t, r, &models.Novel{
		ID: "1", Title: "N", Source: "s", URL: "u", Summary: "s",
		ChaptersCount: 5, LastUpdated: 1, DateAdded: 1,
	})

	for i := 1; i <= 5; i++ {
		title := "Chapter " + string(rune('0'+i))
		if i == 3 {
			title = "Chapter 3: The Pill Furnace"
		}
		if _, err := r.CreateChapter(&models.Chapter{
			NovelID: "1", Number: i, Title: title, Content: "body",
			DateTranslated: int64(i), URL: "https://c/" + string(rune('0'+i)),
		}); err != nil {
			t.Fatalf("seed chapter %d: %v", i, err)
		}
	}

	t.Run("ascending page", func(t *testing.T) {
		chapters, count, err := r.QueryChapters("1", models.ChapterQuery{Ascending: true, Limit: 2})
		if err != nil {
			t.Fatal(err)
		}
		if count != 5 {
			t.Errorf("count = %d, want 5", count)
		}
		if len(chapters) != 2 || chapters[0].Number != 1 || chapters[1].Number != 2 {
			t.Errorf("got %d chapters starting at %d", len(chapters), chapters[0].Number)
		}
	})

	t.Run("descending page", func(t *testing.T) {
		chapters, _, err := r.QueryChapters("1", models.ChapterQuery{Limit: 2})
		if err != nil {
			t.Fatal(err)
		}
		if chapters[0].Number != 5 {
			t.Errorf("first chapter = %d, want 5", chapters[0].Number)
		}
	})

	t.Run("search narrows the count", func(t *testing.T) {
		chapters, count, err := r.QueryChapters("1", models.ChapterQuery{Search: "furnace", Limit: 10})
		if err != nil {
			t.Fatal(err)
		}
		if count != 1 || len(chapters) != 1 || chapters[0].Number != 3 {
			t.Errorf("search returned %d rows (count %d), want just chapter 3", len(chapters), count)
		}
	})

	t.Run("content is excluded from the list", func(t *testing.T) {
		chapters, _, err := r.QueryChapters("1", models.ChapterQuery{Limit: 1})
		if err != nil {
			t.Fatal(err)
		}
		if chapters[0].Content != "" {
			t.Error("list query returned chapter content, which is the largest column")
		}
	})

	t.Run("bounds describe the whole novel", func(t *testing.T) {
		first, last, err := r.ChapterNumberBounds("1")
		if err != nil {
			t.Fatal(err)
		}
		if first != 1 || last != 5 {
			t.Errorf("bounds = %d..%d, want 1..5", first, last)
		}
	})

	t.Run("bounds are zero for a novel with no chapters", func(t *testing.T) {
		first, last, err := r.ChapterNumberBounds("missing")
		if err != nil {
			t.Fatalf("bounds on an unknown novel should not error: %v", err)
		}
		if first != 0 || last != 0 {
			t.Errorf("bounds = %d..%d, want 0..0", first, last)
		}
	})
}

// Deleting a novel must take its chapters with it, or the chapter count drifts.
func TestDeleteNovelRemovesChapters(t *testing.T) {
	r := newTestRepo(t)
	seedNovel(t, r, &models.Novel{
		ID: "1", Title: "N", Source: "s", URL: "u", Summary: "s", LastUpdated: 1, DateAdded: 1,
	})
	if _, err := r.CreateChapter(&models.Chapter{
		NovelID: "1", Number: 1, Title: "C", Content: "body", DateTranslated: 1,
	}); err != nil {
		t.Fatal(err)
	}

	if err := r.DeleteNovel("1"); err != nil {
		t.Fatal(err)
	}

	_, count, err := r.QueryChapters("1", models.ChapterQuery{Limit: 10})
	if err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Errorf("%d chapters left behind after deleting the novel", count)
	}
}

// Chapters legitimately have no URL, so the unique index must not reject the
// second one that omits it.
func TestChaptersWithoutURLCoexist(t *testing.T) {
	r := newTestRepo(t)
	seedNovel(t, r, &models.Novel{
		ID: "1", Title: "N", Source: "s", URL: "u", Summary: "s", LastUpdated: 1, DateAdded: 1,
	})

	for i := 1; i <= 2; i++ {
		if _, err := r.CreateChapter(&models.Chapter{
			NovelID: "1", Number: i, Title: "C", Content: "body", DateTranslated: 1,
		}); err != nil {
			t.Fatalf("chapter %d with an empty url was rejected: %v", i, err)
		}
	}
}
