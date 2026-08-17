package database

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"subscriptions-api-postgres/internal/config"
)

func DatabaseBuild(cfg config.DatabaseConfig) string {
	databaseURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
		cfg.SSLMode,
	)
	return databaseURL
}

func Connect(cfg config.DatabaseConfig) (*sql.DB, error) {
	databaseURL := DatabaseBuild(cfg)

	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
