package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os/exec"
	"strings"

	//"github.com/rs/cors"

	"github.com/fortify-presales/insecure-go-api/internal/handlers"
	"github.com/fortify-presales/insecure-go-api/internal/memstore"
	"github.com/fortify-presales/insecure-go-api/internal/middleware"
)

// Entry point of the program
func main() {
	// Initialize storage
	repo, err := memstore.NewInmemoryRepository() // with in-memory database
	if err != nil {
		log.Fatal("Error:", err)
	}
	repo.Populate() // Populate the in-memory database

	// Initialize middleware stack
	logger := slog.Default()
	stack := middleware.MiddlewareStack(
		middleware.RateLimiter(200),
		middleware.PanicRecovery(logger),
	)

	// Iniitialize handlers
	noteHandler := &handlers.NoteHandler{
		Repository: repo, // Injecting dependency
	}

	// Setup router
	serverMux := http.NewServeMux()

	// Register the routes
	serverMux.HandleFunc("GET /api/notes", noteHandler.GetAll)
	serverMux.HandleFunc("GET /api/notes/{id}", noteHandler.Get)
	serverMux.HandleFunc("POST /api/notes", noteHandler.Post)
	serverMux.HandleFunc("PUT /api/notes/{id}", noteHandler.Put)
	serverMux.HandleFunc("DELETE /api/notes/{id}", noteHandler.Delete)

    serverMux.HandleFunc("/ping/{cmd}", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("handling ping at %s\n", r.URL.Path)

        // Get the 'host' parameter from the query string
        //host := r.URL.Query().Get("host")
		// Get the 'host' parameter from the URL path
		host := strings.TrimPrefix(r.URL.Path, "/ping/")
		// Uncomment the following line when support is added for PathValue
		//host := r.PathValue("cmd")
		log.Printf("host is %s\n", host)

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

	// CORS middleware
	//serverMux = cors.Default().Handler(router)
	
	// Start server
	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", stack(serverMux)); err != nil {
		log.Fatal(err)
	}

}

