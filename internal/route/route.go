package route

import (
	"net/http"

	"go-minio/internal/config"
	"go-minio/internal/handler"
	"go-minio/internal/middleware"
)

func Setup(cfg *config.Config, storageHandler *handler.StorageHandler) http.Handler {
	mux := http.NewServeMux()

	// Storage Routes
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	mux.HandleFunc("/upload", middleware.SecretAuth(cfg, storageHandler.Upload))
	mux.HandleFunc("/list", middleware.SecretAuth(cfg, storageHandler.List))
	mux.HandleFunc("/view", middleware.SecretAuth(cfg, storageHandler.View))
	mux.HandleFunc("/delete", middleware.SecretAuth(cfg, storageHandler.Delete))

	return mux
}
