package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"backend/models"
	"backend/service"
)

// extractNovelDetails handles the POST request to extract and translate novel details from a URL
// It creates a new novel entry in the database with the translated information
func extractNovelDetails(w http.ResponseWriter, r *http.Request) {
	// Parse the incoming request
	var request models.NovelExtractionRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	createdNovel, err := service.GetTranslationService().ExtractNovelDetails(r.Context(), &request)
	if err != nil {
		http.Error(w, "Failed to extract novel details: "+err.Error(), http.StatusInternalServerError)
	}

	// Send the response
	writeJSON(w, createdNovel, http.StatusCreated)
}

// translateNovelChapter handles the POST request to translate a novel chapter
func translateNovelChapter(w http.ResponseWriter, r *http.Request) {
	// Parse the incoming request
	var request models.ChapterTranslationRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Call the translation service to translate the chapter
	translatedChapter, err := service.GetTranslationService().TranslateChapter(r.Context(), &request)
	if err != nil {
		http.Error(w, "Failed to translate chapter: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Send the response
	writeJSON(w, translatedChapter, http.StatusOK)
}

// translateFirstChapter handles the POST request to translate the first chapter of a novel
func translateFirstChapter(w http.ResponseWriter, r *http.Request) {
	// Parse the incoming request
	var request models.ChapterTranslationRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Call the translation service to translate the first chapter
	response, err := service.GetTranslationService().TranslateFirstChapter(r.Context(), &request)
	if err != nil {
		http.Error(w, "Failed to translate first chapter: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Send the response
	writeJSON(w, response, http.StatusOK)
}

// refreshNovel handles the POST request to refresh a novel's details
// endpoint /novels/{id}/refresh
func refreshNovel(w http.ResponseWriter, r *http.Request) {
	// parse the id from the URL path
	var id string
	urlParts := strings.TrimPrefix(r.URL.Path, "/novels/")
	if parts := strings.Split(urlParts, "/"); len(parts) == 2 {
		id = parts[0]
	}

	response, err := service.GetTranslationService().RefreshNovel(r.Context(), id)
	if err != nil {
		http.Error(w, "Failed to refresh novel: "+err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, response, http.StatusOK)
}
