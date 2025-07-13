package handler

import (
	"encoding/json"
	"net/http"

	"backend/repo"
)

// RegisterHandlers registers all API handlers to the provided router
func RegisterHandlers(mux *http.ServeMux) {
	// Register health check endpoint
	mux.HandleFunc("GET /health", healthCheckHandler)

	// Stats endpoint
	mux.HandleFunc("GET /stats/novels", getNovelStats)

	// Sources CRUD APIs
	mux.HandleFunc("GET /sources", getAllSources)

	// Novels CRUD APIs
	mux.HandleFunc("GET /novels", getNovelsUsingFilter)
	mux.HandleFunc("GET /search/novels/{query}", searchNovel)
	mux.HandleFunc("GET /novels/{id}", getNovelByID)
	mux.HandleFunc("GET /novels/{id}/chapters", getNovelChapters)
	mux.HandleFunc("GET /novels/{id}/chapters/num/{chapterNumber}", getNovelChapterByNumber)

	// Translation APIs
	mux.HandleFunc("POST /novels/translate", extractNovelDetails)
	mux.HandleFunc("POST /novels/translate/chapter", translateNovelChapter)
	mux.HandleFunc("POST /novels/translate/first_chapter", translateFirstChapter)
	mux.HandleFunc("POST /novels/refresh", refreshNovel)
}

// healthCheckHandler provides a simple health check endpoint
func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	up, err := repo.GetRepo().GetStatus()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if !up {
		http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
		return
	}

	resp := map[string]string{
		"status": "ok",
	}

	writeJSON(w, resp, http.StatusOK)
}

// writeJSON is a helper function to write JSON responses
func writeJSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
