package handler

import (
	"net/http"

	"github.com/fortify-presales/insecure-go-api/pkg/log"

	httpSwagger "github.com/swaggo/http-swagger/v2"

	"github.com/fortify-presales/insecure-go-api/internal/config"
	"github.com/fortify-presales/insecure-go-api/internal/note"
	"github.com/fortify-presales/insecure-go-api/internal/site"

)

// BuildHandler sets up the HTTP routing and builds an HTTP handler.
func BuildHandler(logger log.Logger, cfg *config.Config, repo note.Repository) http.Handler {
	router := http.NewServeMux()

	router.Handle("GET /swagger/", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),                        // The url pointing to API definition
		httpSwagger.DefaultModelsExpandDepth(httpSwagger.HideModel), // Models will not be expanded
	))

	notesHandler := note.MakeHTTPHandler(repo)
	router.Handle("/api/v1/notes", notesHandler)
	router.Handle("/api/v1/notes/", notesHandler)

	siteHandler := site.MakeHTTPHandler(logger, cfg);
	router.Handle("/api/v1/site", siteHandler)
	router.Handle("/api/v1/site/", siteHandler)

	return router
}