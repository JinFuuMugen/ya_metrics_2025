package repository

import (
	"database/sql"
	"fmt"

	"github.com/JinFuuMugen/ya_metrics_2025/internal/config"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var DB *sql.DB

func InitDB(cfg *config.ServerConfig) error {
	var err error

	DB, err = sql.Open("pgx", cfg.DatabaseDSN)
	if err != nil {
		return fmt.Errorf("cannot connect to database: %w", err)
	}

	return nil
}
