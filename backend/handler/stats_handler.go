package handler

import (
	"net/http"

	"backend/service"
)

// getNovelStats handles the GET endpoint
func getNovelStats(w http.ResponseWriter, r *http.Request) {
	// Get stats from service
	stats, err := service.GetStatsService().GetNovelStats()
	if err != nil {
		http.Error(w, "Failed to get stats: "+err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, stats, http.StatusOK)
}
