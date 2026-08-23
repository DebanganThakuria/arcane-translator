package handler

import (
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// RegisterStatic serves the built frontend from webDir, so a packaged install
// runs as one process on one port instead of needing a separate static server.
//
// It is a no-op when the directory is absent, which is the case during
// development: the Vite dev server owns the frontend then.
func RegisterStatic(mux *http.ServeMux, webDir string) bool {
	if webDir == "" {
		return false
	}

	index := filepath.Join(webDir, "index.html")
	if _, err := os.Stat(index); err != nil {
		return false
	}

	files := http.FileServer(http.Dir(webDir))

	// Registered on "/" so it acts as the fallback: Go's router prefers the
	// more specific API patterns, and everything else lands here.
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		clean := filepath.Clean(strings.TrimPrefix(r.URL.Path, "/"))

		// Hashed build assets are immutable; the shell must never be cached.
		if strings.HasPrefix(r.URL.Path, "/_app/immutable/") {
			w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		}

		if clean != "." && exists(webDir, clean) {
			files.ServeHTTP(w, r)
			return
		}

		// Unknown paths are client-side routes, so hand back the SPA shell.
		w.Header().Set("Cache-Control", "no-cache")
		http.ServeFile(w, r, index)
	})

	return true
}

// exists reports whether name is a regular file inside root, refusing any path
// that escapes it.
func exists(root, name string) bool {
	if filepath.IsAbs(name) || strings.HasPrefix(name, "..") {
		return false
	}

	info, err := fs.Stat(os.DirFS(root), filepath.ToSlash(name))
	return err == nil && info.Mode().IsRegular()
}
