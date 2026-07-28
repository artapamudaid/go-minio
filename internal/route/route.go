package route

import (
	"net/http"
	"os"
	"strconv"

	"go-minio/internal/config"
	"go-minio/internal/handler"
	"go-minio/internal/middleware"

	"golang.org/x/time/rate"
)

func Setup(cfg *config.Config, storageHandler *handler.StorageHandler) http.Handler {
	mux := http.NewServeMux()

	rps, _ := strconv.Atoi(os.Getenv("RATE_LIMIT_RPS"))
	burst, _ := strconv.Atoi(os.Getenv("RATE_LIMIT_BURST"))
	globalLimiter := middleware.NewGlobalLimiter(rate.Limit(float64(rps)), burst)

	protect := func(h http.HandlerFunc) http.HandlerFunc {
		return middleware.ParseBody(middleware.GlobalRateLimit(globalLimiter, middleware.SecretAuth(cfg, h)))
	}

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Routes
	mux.HandleFunc("/upload", protect(storageHandler.Upload))
	mux.HandleFunc("/list", protect(storageHandler.List))
	mux.HandleFunc("/view", protect(storageHandler.View))
	mux.HandleFunc("/delete", protect(storageHandler.Delete))

	return mux
}
