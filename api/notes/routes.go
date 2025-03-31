package notes

import (
	"net/http"

	"github.com/fortify-presales/insecure-go-api/internal/handlers"
	model "github.com/fortify-presales/insecure-go-api/internal/models"
)

func MakeHTTPHandler(repo model.Repository) http.Handler {

	// Iniitialize handlers
	noteHandler := &handlers.NoteHandler{
		Repository: repo, // Injecting dependency
	}

	router := http.NewServeMux()
	router.HandleFunc("GET /api/v1/notes", noteHandler.GetAll)
	router.HandleFunc("GET /api/v1/notes/{id}", noteHandler.Get)
	router.HandleFunc("POST /api/v1/notes", noteHandler.Post)
	router.HandleFunc("PUT /api/v1/notes/{id}", noteHandler.Put)
	router.HandleFunc("DELETE /api/v1/notes/{id}", noteHandler.Delete)

	return router
}
