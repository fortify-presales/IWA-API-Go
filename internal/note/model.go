package note

import (
	"errors"
	"time"
)

//go:generate mockgen -destination=../mocks/mock_repository.go -package=mocks github.com/fortify-presales/insecure-go-api/model Repository

var (
	ErrNoteExists          = errors.New("note title exists")
	ErrNoteNotExists error = errors.New("note doesn't exist")
)

type Note struct {
	NoteID      string    `json:"noteid,omitempty"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedOn   time.Time `json:"createdon,omitempty"`
}

// CRUD interface
type Repository interface {
	Populate() error
	Create(Note) (string, error)
	Update(string, Note) error
	Delete(string) error
	GetById(string) (Note, error)
	GetAll(string) ([]Note, error)
}
