package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/cors"

	"github.com/fortify-presales/insecure-go-api/api/notes"
	"github.com/fortify-presales/insecure-go-api/internal/config"
	"github.com/fortify-presales/insecure-go-api/internal/memstore"
	"github.com/fortify-presales/insecure-go-api/internal/middleware"
	model "github.com/fortify-presales/insecure-go-api/internal/models"
	"github.com/fortify-presales/insecure-go-api/pkg/log"
)

// Version indicates the current version of the application.
var Version = "1.0.0"

// Read configuration file path from command line argument, default is "./config/local.yml"
var flagConfig = flag.String("config", "./config/local.yml", "path to the config file")

func main() {
	// Parse command line flags
	flag.Parse()
	// Create root logger tagged with server version
	logger := log.New().With(nil, "version", Version)
	// Load application configurations
	cfg, err := config.Load(*flagConfig, logger)
	if err != nil {
		logger.Errorf("failed to load application configuration: %s", err)
		os.Exit(-1)
	}
	// Initialize storage
	repo, err := memstore.NewInmemoryRepository() // with in-memory database
	if err != nil {
		logger.Error("Error:", err)
		os.Exit(-1)
	}
	repo.Populate() // Populate the in-memory database

	// Initialize middleware stack
	stack := middleware.MiddlewareStack(
		middleware.RateLimiter(200),
		middleware.PanicRecovery(logger),
	)
	// Initialize CORS
	serverMux := cors.Default().Handler(buildHandler(logger, cfg, repo))
	// Build HTTP server
	address := fmt.Sprintf(":%v", cfg.ServerPort)
	srv := &http.Server{
		Addr:         address,
		Handler:      stack(serverMux),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}
	// Start the server in a separate Goroutine.
	go func() {
		logger.Infof("Starting the server on :%d", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Errorf("Failed to start server: %s", err)
			os.Exit(-1)
		}
	}()
	// Implement graceful shutdown.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Infof("Shutting down the server...")
	// Set a timeout for shutdown (for example, 5 seconds).
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// Shutdown the server gracefully.
	if err := srv.Shutdown(ctx); err != nil {
		logger.Errorf("Server shutdown error: %v", err)
	}
	logger.Infof("Server gracefully stopped")

}

// buildHandler sets up the HTTP routing and builds an HTTP handler.
func buildHandler(logger log.Logger, cfg *config.Config, repo model.Repository) http.Handler {
	router := http.NewServeMux()

	notesHandler := notes.MakeHTTPHandler(repo)
	router.Handle("/api/v1/notes", notesHandler)
	router.Handle("/api/v1/notes/", notesHandler)

	router.HandleFunc("/api/v1/ping/{cmd}", func(w http.ResponseWriter, r *http.Request) {
		logger.Debugf("handling ping at %s\n", r.URL.Path)
		//
		// servemux r.PathValue - Not yet supported by Fortify without custom rules
		//
		host := r.PathValue("cmd")
		//
		// instead we use following
		//host := strings.TrimPrefix(r.URL.Path, "/api/v1/ping/")
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

	return router
}
