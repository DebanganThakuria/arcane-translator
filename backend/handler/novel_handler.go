package handler

import (
	"database/sql"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"backend/models"
	"backend/service"
)

// respondError maps a service error onto a status code. A missing row is a 404
// rather than the 500 every lookup used to return.
func respondError(w http.ResponseWriter, message string, err error) {
	status := http.StatusInternalServerError
	if errors.Is(err, sql.ErrNoRows) {
		status = http.StatusNotFound
	}
	http.Error(w, message+": "+err.Error(), status)
}

// paginationParams reads page and limit, clamping limit to maxLimit.
func paginationParams(query url.Values, defaultLimit, maxLimit int) (page, limit int) {
	page = 1
	if parsed, err := strconv.Atoi(query.Get("page")); err == nil && parsed > 0 {
		page = parsed
	}

	limit = defaultLimit
	if parsed, err := strconv.Atoi(query.Get("limit")); err == nil && parsed > 0 {
		limit = min(parsed, maxLimit)
	}

	return page, limit
}

// getNovelsUsingFilter handles GET /novels.
//
// Two shapes are supported:
//
//	?filter=recently_read&value=5              legacy single-filter shortcuts
//	?q=&language=&genre=&status=&sort=&dir=    composable filters
//
// Both accept page and limit. Filtering, sorting and paging all happen in the
// database, so total_count reflects the filters rather than the whole table.
func getNovelsUsingFilter(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	page, limit := paginationParams(query, 20, 100)
	offset := (page - 1) * limit

	var (
		response *models.NovelListResponse
		err      error
	)

	if filter := query.Get("filter"); filter != "" {
		response, err = service.GetNovelService().GetNovelsByFilter(
			filter, query.Get("value"), offset, limit,
		)
	} else {
		response, err = service.GetNovelService().QueryNovels(models.NovelQuery{
			Search:    strings.TrimSpace(query.Get("q")),
			Genre:     strings.TrimSpace(query.Get("genre")),
			Status:    strings.TrimSpace(query.Get("status")),
			Sort:      models.ParseNovelSort(query.Get("sort")),
			Ascending: strings.EqualFold(query.Get("dir"), "asc"),
			Offset:    offset,
			Limit:     limit,
		}, strings.TrimSpace(query.Get("language")))
	}

	if err != nil {
		http.Error(w, "Failed to retrieve novels: "+err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, response, http.StatusOK)
}

// searchNovel handles GET /search/novels/{query}
func searchNovel(w http.ResponseWriter, r *http.Request) {
	query, err := url.PathUnescape(r.PathValue("query"))
	if err != nil {
		http.Error(w, "Invalid search query", http.StatusBadRequest)
		return
	}

	foundNovels, err := service.GetNovelService().SearchNovel(query)
	if err != nil {
		http.Error(w, "Failed to search novels: "+err.Error(), http.StatusBadRequest)
		return
	}

	if foundNovels == nil {
		foundNovels = []*models.Novel{}
	}

	writeJSON(w, foundNovels, http.StatusOK)
}

// getNovelByID handles GET /novels/{id}
func getNovelByID(w http.ResponseWriter, r *http.Request) {
	novel, err := service.GetNovelService().GetNovelByID(r.PathValue("id"))
	if err != nil {
		respondError(w, "Failed to retrieve novel", err)
		return
	}

	writeJSON(w, novel, http.StatusOK)
}

// deleteNovel handles DELETE /novels/{id}
func deleteNovel(w http.ResponseWriter, r *http.Request) {
	if err := service.GetNovelService().DeleteNovel(r.PathValue("id")); err != nil {
		respondError(w, "Failed to delete novel", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// getNovelChapters handles GET /novels/{id}/chapters.
//
// Without paging parameters it returns the full array, which is what the older
// clients expect. Pass page or limit to get a ChapterListResponse envelope
// instead; novels here run to several hundred chapters, so the paged form is
// the one worth using.
func getNovelChapters(w http.ResponseWriter, r *http.Request) {
	novelID := r.PathValue("id")
	query := r.URL.Query()

	paged := query.Has("page") || query.Has("limit") || query.Has("q") || query.Has("dir")
	if !paged {
		chapters, err := service.GetNovelService().GetNovelChapters(novelID)
		if err != nil {
			respondError(w, "Failed to retrieve chapters", err)
			return
		}
		if chapters == nil {
			chapters = []*models.Chapter{}
		}
		writeJSON(w, chapters, http.StatusOK)
		return
	}

	page, limit := paginationParams(query, 50, 200)

	response, err := service.GetNovelService().QueryChapters(novelID, models.ChapterQuery{
		Search:    strings.TrimSpace(query.Get("q")),
		Ascending: !strings.EqualFold(query.Get("dir"), "desc"),
		Offset:    (page - 1) * limit,
		Limit:     limit,
	})
	if err != nil {
		respondError(w, "Failed to retrieve chapters", err)
		return
	}

	writeJSON(w, response, http.StatusOK)
}

// getNovelChapterByNumber handles GET /novels/{id}/chapters/num/{chapterNumber}
func getNovelChapterByNumber(w http.ResponseWriter, r *http.Request) {
	chapterNumber, err := strconv.Atoi(r.PathValue("chapterNumber"))
	if err != nil {
		http.Error(w, "Invalid chapter number", http.StatusBadRequest)
		return
	}

	chapter, err := service.GetNovelService().GetChapterByNumber(r.PathValue("id"), chapterNumber)
	if err != nil {
		respondError(w, "Failed to retrieve chapter", err)
		return
	}

	writeJSON(w, chapter, http.StatusOK)
}

// deleteChapter handles DELETE /novels/{id}/chapters/{chapterId}
func deleteChapter(w http.ResponseWriter, r *http.Request) {
	err := service.GetNovelService().DeleteChapter(r.PathValue("id"), r.PathValue("chapterId"))
	if err != nil {
		respondError(w, "Failed to delete chapter", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// getAllSources handles GET /sources
func getAllSources(w http.ResponseWriter, r *http.Request) {
	sources, err := service.GetNovelService().GetAllSources()
	if err != nil {
		respondError(w, "Failed to retrieve sources", err)
		return
	}

	writeJSON(w, sources, http.StatusOK)
}
