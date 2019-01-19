package rest_backend

import (
	"net/http"
)

func DisableBrowserCache(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache-Control", "private, max-age=0")

		next.ServeHTTP(w, r)
	}
}
