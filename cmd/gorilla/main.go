package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os/exec"

	"github.com/rs/cors"
	"github.com/gorilla/mux"

	apphttp "github.com/fortify-presales/insecure-go-api/http/gorilla"
	"github.com/fortify-presales/insecure-go-api/internal/memstore"
	"github.com/fortify-presales/insecure-go-api/internal/middleware"
)

// Entry point of the program
func main() {
	repo, err := memstore.NewInmemoryRepository() // With in-memory database
	if err != nil {
		log.Fatal("Error:", err)
	}
	repo.Populate() // Populate the in-memory database

	h := &apphttp.NoteHandler{
		Repository: repo, // Injecting dependency
	}
	router := initializeRoutes(h) // configure routes

	logger := slog.Default()
	// Adding middleware http
	router = middleware.Apply(router,
		middleware.RateLimiter(200),
		middleware.PanicRecovery(logger),
	)
	// CORS middleware
	router = cors.Default().Handler(router)
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	log.Println("Listening...")
	server.ListenAndServe() // Run the http server
}

func initializeRoutes(h *apphttp.NoteHandler) http.Handler {
	mux := mux.NewRouter()
	mux.HandleFunc("/api/notes", h.GetAll).Methods("GET")
	mux.HandleFunc("/api/notes/{id}", h.Get).Methods(("GET"))
	mux.HandleFunc("/api/notes", h.Post).Methods("POST")
	mux.HandleFunc("/api/notes/{id}", h.Put).Methods("PUT")
	mux.HandleFunc("/api/notes/{id}", h.Delete).Methods("DELETE")

    mux.HandleFunc("/ping/{cmd}", func(w http.ResponseWriter, r *http.Request) {
        // Get the 'host' parameter from the query string
        //host := r.URL.Query().Get("host")
		host := mux.Vars(r)["cmd"]

        // Directly using user input in a shell command
        cmd := exec.Command("ping", "-c", "4", host)
        output, err := cmd.CombinedOutput()
        if err != nil {
            http.Error(w, fmt.Sprintf("Error: %s", err), http.StatusInternalServerError)
            return
        }

        // Return the command output to the user
        w.Header().Set("Content-Type", "text/plain")
        w.Write(output)
    })
	
	return mux
}
