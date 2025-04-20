package app // Or a middleware package

import (
	"log/slog"
	"net/http"
	"strings"
	"time"
)

// CacheControlMiddleware adds caching headers for static assets.
func CacheControlMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Apply cache headers only to GET requests for static assets
		if r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/static/") {
			// Cache aggressively: public, one year max-age, immutable
			// 'immutable' tells browsers the file content will NEVER change for this URL,
			// so it doesn't even need to revalidate. Use only if filenames
			// don't change without content changing (or use versioning/fingerprinting).
			// Since assets are embedded, this is generally safe.
			cacheDuration := time.Hour * 24 * 365 // One year
			w.Header().Set("Cache-Control", "public, max-age="+string(int(cacheDuration.Seconds()))+", immutable")

			// You could also set Expires header for older proxies/clients, though less common now
			// w.Header().Set("Expires", time.Now().Add(cacheDuration).Format(http.TimeFormat))
		}
		next.ServeHTTP(w, r)
	})
}

// Favicon specific handler with caching (update existing one slightly)
func FaviconHandler(w http.ResponseWriter, r *http.Request) {
	// Check method if needed, though usually just GET for favicon
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	file, err := staticFS.ReadFile("static/favicon.ico")
	if err != nil {
		slog.Error("Failed to read embedded favicon.ico", "error", err)
		http.NotFound(w, r) // Or Internal Server Error
		return
	}

	// Set aggressive cache headers for favicon too
	//cacheDuration := time.Hour * 24 * 365 // One year
	w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
	w.Header().Set("Content-Type", "image/x-icon") // Ensure correct content type

	_, err = w.Write(file)
	if err != nil {
		slog.Error("Failed to write favicon response", "error", err)
	}
}
