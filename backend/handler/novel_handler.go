package handler

import (
	"net/http"
	"strconv"
	"strings"

	"backend/models"
	"backend/service"
)

// getNovelsUsingFilter handles GET /novels?filter={filter}&value={value}&page={page}&limit={limit}
func getNovelsUsingFilter(w http.ResponseWriter, r *http.Request) {
	filter, value, page, limit := getAllRequestParams(r)

	offset := (page - 1) * limit

	var err error
	var response *models.NovelListResponse

	if filter != "" {
		response, err = service.GetNovelService().GetNovelsByFilter(filter, value, offset, limit)
	} else {
		response, err = service.GetNovelService().GetAllNovels(offset, limit)
	}

	if err != nil {
		http.Error(w, "Failed to retrieve novels: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if response.Novels == nil {
		response.Novels = []*models.Novel{}
	}

	writeJSON(w, response, http.StatusOK)
}

func getAllRequestParams(r *http.Request) (filter, value string, page, limit int) {
	query := r.URL.Query()

	filter = query.Get("filter")

	value = query.Get("value")

	page = 1
	if p := query.Get("page"); p != "" {
		if pInt, err := strconv.Atoi(p); err == nil && pInt > 0 {
			page = pInt
		}
	}

	limit = 20
	if l := query.Get("limit"); l != "" {
		if lInt, err := strconv.Atoi(l); err == nil && lInt > 0 {
			if lInt > 100 {
				lInt = 100
			}
			limit = lInt
		}
	}

	return
}

// searchNovel handles GET /search/novels/{query}
func searchNovel(w http.ResponseWriter, r *http.Request) {
	// Extract query from URL path
	query := strings.TrimPrefix(r.URL.Path, "/search/novels/")
	if idx := strings.Index(query, "/"); idx > -1 {
		query = query[:idx]
	}

	// Search novel
	foundNovels, err := service.GetNovelService().SearchNovel(query)
	if err != nil {
		http.Error(w, "Failed to create novel: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if len(foundNovels) == 0 {
		foundNovels = []*models.Novel{}
	}

	writeJSON(w, foundNovels, http.StatusCreated)
}

// getNovelByID handles GET /novels/{id}
func getNovelByID(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL path
	id := strings.TrimPrefix(r.URL.Path, "/novels/")
	if idx := strings.Index(id, "/"); idx > -1 {
		id = id[:idx]
	}

	novel, err := service.GetNovelService().GetNovelByID(id)
	if err != nil {
		http.Error(w, "Failed to retrieve novel: "+err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, novel, http.StatusOK)
}

// getNovelChapters handles GET /novels/{id}/chapters
func getNovelChapters(w http.ResponseWriter, r *http.Request) {
	// Extract novel ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/novels/")
	pathParts := strings.Split(path, "/")
	if len(pathParts) < 2 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	novelID := pathParts[0]

	chapters, err := service.GetNovelService().GetNovelChapters(novelID)
	if err != nil {
		http.Error(w, "Failed to retrieve chapters: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// If no chapters found, return an empty array
	if len(chapters) == 0 {
		writeJSON(w, []models.Chapter{}, http.StatusOK)
		return
	}

	writeJSON(w, chapters, http.StatusOK)
}

// getNovelChapterByNumber handles GET /novels/{id}/chapters/num/{chapterNumber}
func getNovelChapterByNumber(w http.ResponseWriter, r *http.Request) {
	// Extract novel ID and chapter number from URL path
	path := strings.TrimPrefix(r.URL.Path, "/novels/")
	pathParts := strings.Split(path, "/")
	if len(pathParts) < 4 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	novelID := pathParts[0]
	chapterNumber, err := strconv.Atoi(pathParts[3])
	if err != nil {
		http.Error(w, "Invalid chapter number", http.StatusBadRequest)
		return
	}

	chapter, err := service.GetNovelService().GetChapterByNumber(novelID, chapterNumber)
	if err != nil {
		http.Error(w, "Failed to retrieve chapter: "+err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, chapter, http.StatusOK)
}

// getAllSources handles GET /sources
func getAllSources(w http.ResponseWriter, r *http.Request) {
	sources, err := service.GetNovelService().GetAllSources()
	if err != nil {
		http.Error(w, "Failed to retrieve sources: "+err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, sources, http.StatusOK)
}
