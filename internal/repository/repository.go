package repository

import (
	"database/sql"
	"os"

	"github.com/fortify-presales/insecure-go-api/pkg/log"

	"github.com/fortify-presales/insecure-go-api/internal/config"
	"github.com/fortify-presales/insecure-go-api/internal/note"
)

const fileName = "sqlite.db"

func BuildRepository(logger log.Logger, cfg *config.Config) note.Repository {
	// Initialize SQLite3 database
	os.Remove(fileName)
	db, err := sql.Open("sqlite3", fileName)
	if err != nil {
		logger.Error(err)
		os.Exit(-1)
	}
	repo, err := note.NewSQLiteRepository(db, logger)
	if err != nil {
		logger.Error("Error:", err)
		os.Exit(-1)
	}
	repo.Populate() // Populate the database

	return repo
}