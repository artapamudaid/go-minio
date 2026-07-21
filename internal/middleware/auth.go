package middleware

import (
	"net/http"

	"go-minio/internal/config"
	"go-minio/pkg/response"
)

func SecretAuth(cfg *config.Config, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		secret := r.FormValue("secret")
		if secret == "" {
			secret = r.URL.Query().Get("secret")
		}

		if secret != cfg.ServerSecretKey {
			response.Error(w, http.StatusUnauthorized, "Unauthorized: Invalid server secret key")
			return
		}

		next(w, r)
	}
}
