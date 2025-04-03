package note

import (
	// internal
	"database/sql"
	"errors"

	"github.com/gofrs/uuid"
	"github.com/mattn/go-sqlite3"

	"github.com/fortify-presales/insecure-go-api/pkg/log"
)

// SQLiteRepository  provides concrete implementation for repository interface
type SQLiteRepository struct {
	db     *sql.DB
	logger log.Logger
}

func NewSQLiteRepository(db *sql.DB, logger log.Logger) (Repository, error) {
	return &SQLiteRepository{
		db:     db,
		logger: logger,
	}, nil
}

func (r *SQLiteRepository) Populate() error {
	r.logger.Info("Populating SQLite database with initial data")
	query := `
    CREATE TABLE IF NOT EXISTS notes (
        id TEXT PRIMARY KEY, -- Storing UUID as text
        title TEXT NOT NULL UNIQUE,
        description TEXT NOT NULL,
        created_on DATETIME NOT NULL
    );
    `

	_, err := r.db.Exec(query)

	note1 := Note{
		NoteID:      "1",
		Title:       "slog",
		Description: "slog is a logging package",
	}
	note2 := Note{
		NoteID:      "2",
		Title:       "viper",
		Description: "viper is a configuration management package",
	}
	created1, err := r.Create(note1)
	if created1 != "" {
	}
	if err != nil {
		r.logger.Error(err)
		return err
	}
	created2, err := r.Create(note2)
	if created2 != "" {
	}
	if err != nil {
		r.logger.Error(err)
		return err
	}
	return nil
}

/*
func (r *SQLiteRepository) isNoteTitleExists(title string) bool {
	for _, v := range i.noteStore {
		if v.Title == title {
			return true
		}
	}
	return false
}*/

func (r *SQLiteRepository) Create(n Note) (string, error) {
	r.logger.Infof("Creating a new note with title: %s", n.Title)
	//if _, ok := i.noteStore[n.NoteID]; ok {
	//	return "", errors.New("NoteID exists")
	//}
	//if i.isNoteTitleExists(n.Title) {
	//	return "", model.ErrNoteExists
	//}
	// Create a Version 4 UUID.
	uid, _ := uuid.NewV4()
	res, err := r.db.Exec("INSERT INTO notes(id, title, description, created_on) values(?,?,?, datetime('now'))",
		uid, n.Title, n.Description)
	if res != nil {
		r.logger.Info(res)
		r.logger.Infof("Created note with ID: %s", uid)
	}
	if err != nil {
		r.logger.Info(err)
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
				return "", ErrNoteExists
			}
		}
		return "", err
	}

	//id, err := res.LastInsertId()
	//if err != nil {
	//    return "", err
	//}
	//n.NoteID = id

	return n.NoteID, nil
}

func (r *SQLiteRepository) Update(id string, n Note) error {
	if id == "" {
		return errors.New("invalid NoteID")
	}
	// check if note with id exists
	res, err := r.db.Exec("UPDATE notes SET title = ?, description = ? WHERE id = ?",
		n.Title, n.Description, n.NoteID)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("update failed")
	}

	return nil
}

func (r *SQLiteRepository) Delete(id string) error {
	if id == "" {
		return errors.New("invalid NoteID")
	}
	// check if note with id exists
	res, err := r.db.Exec("DELETE FROM notes WHERE noteid = ?", id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("delete failed")
	}

	return err
}

func (r *SQLiteRepository) GetById(id string) (Note, error) {
	if id == "" {
		return Note{}, errors.New("invalid NoteID")
	}
	var note Note
	row := r.db.QueryRow("SELECT * FROM notes WHERE noteid = ?", id)
	if err := row.Scan(&note.NoteID, &note.Title, &note.Description, &note.CreatedOn); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Note{}, ErrNoteNotExists
		}
		return Note{}, err
	}
	return note, nil
}

func (r *SQLiteRepository) GetAll(keywords string) ([]Note, error) {
	if keywords == "" {
		r.logger.Info("Retrieving all notes")
	} else {
		r.logger.Infof("Retrieving notes using keywords: %s", keywords)
	}
	rows, err := r.db.Query("SELECT * FROM notes WHERE title LIKE ? OR description LIKE ?", "%"+keywords+"%", "%"+keywords+"%")

	if err != nil {
		r.logger.Error(err)
		return nil, err
	}
	defer rows.Close()

	r.logger.Infof("Found %d notes", rows)
	var all []Note
	for rows.Next() {
		var note Note
		r.logger.Info(note)

		if err := rows.Scan(&note.NoteID, &note.Title, &note.Description, &note.CreatedOn); err != nil {
			return nil, err
		}
		all = append(all, note)
	}
	return all, nil
}
