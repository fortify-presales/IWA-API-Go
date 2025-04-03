package model

import (
	"errors"
)

//go:generate mockgen -destination=../mocks/mock_repository.go -package=mocks github.com/fortify-presales/insecure-go-api/model Repository

var (
	ErrNotFound            = errors.New("no records found")
	ErrUpdateFailed	error  = errors.New("update failed")
)

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
