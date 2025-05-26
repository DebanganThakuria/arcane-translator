package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"backend/models"
	"backend/service"
)

// getAllNovels handles GET /novels
func getAllNovels(w http.ResponseWriter, r *http.Request) {
	novels, err := service.GetNovelService().GetAllNovels()
	if err != nil {
		http.Error(w, "Failed to retrieve novels: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// If no novels found, return an empty array
	if len(novels) == 0 {
		writeJSON(w, []models.Novel{}, http.StatusOK)
		return
	}

	writeJSON(w, novels, http.StatusOK)
}

// createNovel handles POST /novels
func createNovel(w http.ResponseWriter, r *http.Request) {
	var novel models.Novel

	// Decode request body
	err := json.NewDecoder(r.Body).Decode(&novel)
	if err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Create novel
	createdNovel, err := service.GetNovelService().CreateNovel(&novel)
	if err != nil {
		http.Error(w, "Failed to create novel: "+err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, createdNovel, http.StatusCreated)
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
		if errors.Is(err, errors.New("novel not found")) {
			http.Error(w, "Novel not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to retrieve novel: "+err.Error(), http.StatusInternalServerError)
		}
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
		if errors.Is(err, errors.New("novel not found")) {
			http.Error(w, "Novel not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to retrieve chapters: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// If no chapters found, return an empty array
	if len(chapters) == 0 {
		writeJSON(w, []models.Chapter{}, http.StatusOK)
		return
	}

	writeJSON(w, chapters, http.StatusOK)
}

// getNovelChapterByID handles GET /novels/{id}/chapters/{chapterId}
func getNovelChapterByID(w http.ResponseWriter, r *http.Request) {
	// Extract novel ID and chapter ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/novels/")
	pathParts := strings.Split(path, "/")
	if len(pathParts) < 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	novelID := pathParts[0]

	// The chapterId could be either an ID or a chapter number
	chapterParam := pathParts[2]

	// Try to parse as chapter number first
	if chapterNumber, err := strconv.Atoi(chapterParam); err == nil {
		// It's a number, so get chapter by number
		chapter, err := service.GetNovelService().GetChapterByNumber(novelID, chapterNumber)
		if err != nil {
			if errors.Is(err, errors.New("chapter not found")) {
				http.Error(w, "Chapter not found", http.StatusNotFound)
			} else {
				http.Error(w, "Failed to retrieve chapter: "+err.Error(), http.StatusInternalServerError)
			}
			return
		}

		writeJSON(w, chapter, http.StatusOK)
		return
	}

	// Not a number, so treat as chapter ID
	chapter, err := service.GetNovelService().GetChapterByID(novelID, chapterParam)
	if err != nil {
		if errors.Is(err, errors.New("chapter not found")) {
			http.Error(w, "Chapter not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to retrieve chapter: "+err.Error(), http.StatusInternalServerError)
		}
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
