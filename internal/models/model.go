package model

import (
	"errors"
	"time"
)

//go:generate mockgen -destination=../mocks/mock_repository.go -package=mocks github.com/fortify-presales/insecure-go-api/model Repository

var (
	ErrNotFound            = errors.New("no records found")
	ErrNoteExists          = errors.New("note title exists")
	ErrNoteNotExists error = errors.New("note doesn't exist")
)

type Note struct {
	NoteID      string    `json:"noteid,omitempty"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedOn   time.Time `json:"createdon,omitempty"`
}

// APIError
type APIMessage struct {
	Message string
	//CreatedAt    time.Time
}

// APIError
type APIError struct {
	ErrorCode    int
	ErrorMessage string
	//CreatedAt    time.Time
}

// CRUD interface
type Repository interface {
	Populate()
	Create(Note) (string, error)
	Update(string, Note) error
	Delete(string) error
	GetById(string) (Note, error)
	GetAll() ([]Note, error)
}
