package database

import (
	"database/sql"
	"kasir-api-go/internal/config"
	"kasir-api-go/internal/utils"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func NewPostgres(cfg *config.DatabaseConfig) (*sql.DB, func() error, error) {
	if cfg.URL == "" {
		log.Println("database url is empty")
		return nil, nil, utils.ErrEmptyDatabaseURL
	}

	// Open connection
	db, err := sql.Open("pgx", cfg.URL)
	if err != nil {
		log.Println("failed to open database connection", err)
		return nil, nil, err
	}

	// Set connection pool
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(time.Hour)

	// Ping database to ensure connection is alive
	if err := db.Ping(); err != nil {
		log.Printf("failed to ping database: %v", err)
		db.Close()
		return nil, nil, err
	}

	// Close function
	closeDB := func() error {
		log.Println("closing database connection")
		if err := db.Close(); err != nil {
			log.Printf("failed to close database: %v", err)
			return err
		}
		return nil
	}

	return db, closeDB, nil
}
