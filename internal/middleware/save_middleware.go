package middleware

import (
	"net/http"

	"github.com/JinFuuMugen/ya_metrics_2025/internal/config"
	"github.com/JinFuuMugen/ya_metrics_2025/internal/filestorage"
	"github.com/JinFuuMugen/ya_metrics_2025/internal/logger"
)

func SyncSaveMiddleware(cfg *config.ServerConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)

			if cfg.StoreInterval == 0 {
				syncSaver, err := filestorage.NewSaver(cfg.FileStoragePath)
				if err != nil {
					http.Error(w, "internal server error", http.StatusInternalServerError)
					logger.Errorf("cannot create sync saver: %w", err)
				}

				syncSaver.SaveMetrics()
				syncSaver.Close()
			}
		})
	}
}
