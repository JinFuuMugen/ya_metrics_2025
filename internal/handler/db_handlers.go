package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/JinFuuMugen/ya_metrics_2025/internal/logger"
	"github.com/JinFuuMugen/ya_metrics_2025/internal/repository"
)

func PingDBHandler(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := repository.DB.PingContext(ctx); err != nil {
		logger.Errorf("error pinging database: %s", err)
		http.Error(w, "error pinging database", http.StatusInternalServerError)
		return
	} else {
		w.WriteHeader(http.StatusOK)
		return
	}
}
