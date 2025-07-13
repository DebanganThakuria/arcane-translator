package utils

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

// GetDBPath returns the path to the SQLite database file
func GetDBPath() string {
	// Get the absolute path to the project root (arcane-translator)
	cwd, err := os.Getwd()
	if err != nil {
		log.Printf("Warning: Failed to get current working directory: %v", err)
		cwd = "."
	}
	// Find the index of "arcane-translator" in the path using strings.Index
	projectRoot := cwd
	if idx := strings.Index(cwd, "arcane-translator"); idx != -1 {
		projectRoot = cwd[:idx+len("arcane-translator")]
	}
	dataDir := filepath.Join(projectRoot, "data")

	// Create a data directory if it doesn't exist
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		if err := os.MkdirAll(dataDir, 0o755); err != nil {
			log.Printf("Warning: Failed to create data directory: %v", err)
			// Fall back to the current directory
			return "data.db"
		}
	}

	return filepath.Join(dataDir, "data.db")
}

func CountWords(text string) int {
	if text == "" {
		return 0
	}

	wordCount := 0
	inWord := false

	// Iterate through each rune in the text
	for _, r := range text {
		// Check if the current rune is part of a word
		// Letters, numbers, apostrophes, and hyphens within words are considered part of a word
		isWordChar := unicode.IsLetter(r) || unicode.IsNumber(r) ||
			(inWord && (r == '\'' || r == '-') &&
				len(text) > 1)

		// Start of a new word
		if isWordChar && !inWord {
			wordCount++
			inWord = true
		} else if !isWordChar {
			// End of a word
			inWord = false
		}
	}

	return wordCount
}

func RemoveDuplicatesFromSlice(array []string) []string {
	uniqueWords := make(map[string]struct{})
	var result []string

	for _, word := range array {
		if _, exists := uniqueWords[word]; !exists {
			uniqueWords[word] = struct{}{}
			result = append(result, word)
		}
	}

	return result
}
