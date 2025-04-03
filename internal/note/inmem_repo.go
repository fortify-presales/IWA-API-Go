package note

import (
	// internal
	"errors"
	"time"

	// external
	"github.com/gofrs/uuid"

	"github.com/fortify-presales/insecure-go-api/internal/models"
	"github.com/fortify-presales/insecure-go-api/pkg/log"

)

// inmemoryRepository provides concrete implementation for repository interface
type inmemoryRepository struct {
	noteStore map[string]Note
	logger log.Logger
}

func NewInmemoryRepository(logger log.Logger) (Repository, error) {
	return &inmemoryRepository{
		noteStore: make(map[string]Note),
		logger: logger,
	}, nil
}

func (i *inmemoryRepository) Populate() error{
	note1 := Note{
		NoteID:      "1",
		Title:       "slog",
		Description: "slog is a logging package",
		CreatedOn:   time.Now(),
	}
	note2 := Note{
		NoteID:      "2",
		Title:       "viper",
		Description: "viper is a configuration management package",
		CreatedOn:   time.Now(),
	}
	i.Create(note1)
	i.Create(note2)
	return nil
}

func (i *inmemoryRepository) isNoteTitleExists(title string) bool {
	for _, v := range i.noteStore {
		if v.Title == title {
			return true
		}
	}
	return false
}
func (i *inmemoryRepository) Create(n Note) (string, error) {
	if _, ok := i.noteStore[n.NoteID]; ok {
		return "", errors.New("NoteID exists")
	}
	if i.isNoteTitleExists(n.Title) {
		return "", ErrNoteExists
	}
	n.CreatedOn = time.Now()
	// Create a Version 4 UUID.
	uid, _ := uuid.NewV4()
	n.NoteID = uid.String()
	i.noteStore[n.NoteID] = n
	return n.NoteID, nil
}

func (i *inmemoryRepository) Update(id string, n Note) error {
	if _, ok := i.noteStore[id]; !ok {
		return ErrNoteNotExists
	}
	n.CreatedOn = time.Now()
	i.noteStore[id] = n
	return nil
}

func (i *inmemoryRepository) Delete(id string) error {
	if _, ok := i.noteStore[id]; !ok {
		return ErrNoteNotExists
	}
	delete(i.noteStore, id)
	return nil
}
func (i *inmemoryRepository) GetById(id string) (Note, error) {
	if v, ok := i.noteStore[id]; !ok {
		return Note{}, ErrNoteNotExists
	} else {
		return v, nil
	}

}

func (i *inmemoryRepository) GetAll(query string) ([]Note, error) {
	if len(i.noteStore) == 0 {
		return nil, model.ErrNotFound
	}
	notes := make([]Note, 0)
	for _, v := range i.noteStore {
		notes = append(notes, v)
	}
	return notes, nil
}
