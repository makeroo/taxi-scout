package rest

import (
	"net/http"
)

// DisableBrowserCache add HTTP header Cache-Control to disable browser cache.
func DisableBrowserCache(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache-Control", "private, max-age=0")

		next.ServeHTTP(w, r)
	}
}
