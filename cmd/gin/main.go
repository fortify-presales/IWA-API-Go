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

	"github.com/gin-gonic/gin"

	"github.com/fortify-presales/insecure-go-api/pkg/log"

	"github.com/fortify-presales/insecure-go-api/internal/config"
)

// Version indicates the current version of the application.
var Version = "1.0.0"

// Read configuration file path from command line argument, default is "./config/local.yml"
var flagConfig = flag.String("config", "./config/local.yml", "path to the config file")

// Entry point of the program
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
	// Build HTTP server
	address := fmt.Sprintf(":%v", cfg.ServerPort)
	srv := &http.Server{
		Addr:         address,
		Handler:      buildHandler(logger, cfg),
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
func buildHandler(logger log.Logger, cfg *config.Config) http.Handler {
	router := gin.Default()

	router.GET("/ping/:cmd", func(c *gin.Context) {
		logger.Debugf("Received request: %s", c.Request.URL.Path)
		//
		// gin c.Param - Not yet supported by Fortify
		//
		host := c.Param("cmd")
		cmd := exec.Command("ping", "-c", "4", host)
		output, err := cmd.CombinedOutput()
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error: %s", err))
		}
		c.String(http.StatusOK, fmt.Sprintf("<pre>%s</pre>", output))
	})

	return router
}
