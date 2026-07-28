package middleware

import "net/http"

func ParseBody(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseMultipartForm(50 << 20)
		next(w, r)
	}
}
