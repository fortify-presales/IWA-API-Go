package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	apphttp "github.com/fortify-presales/insecure-go-api/http/echo"
	"github.com/fortify-presales/insecure-go-api/internal/memstore"
)

func main() {
	repo, err := memstore.NewInmemoryRepository() // With in-memory database
	if err != nil {
		log.Fatal("Error:", err)
	}
	h := &apphttp.NoteHandler{
		Repository: repo, // Injecting dependency
	}
	e := echo.New()
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))
	// Routes
	e.GET("/api/notes", h.GetAll)
	e.GET("/api/notes/:id", h.Get)
	e.POST("/api/notes", h.Post)
	e.PUT("/api/notes/:id", h.Put)
	e.DELETE("/api/notes/:id", h.Delete)
	e.GET("/ping", func(c echo.Context) error {
		 // Get the 'host' parameter from the query string
		 r := c.Request()
		 host := r.URL.Query().Get("host")

		 // Directly using user input in a shell command
		 cmd := exec.Command("ping", "-c", "4", host)
		 output, err := cmd.CombinedOutput()
		 if err != nil {
			 return c.String(http.StatusInternalServerError, fmt.Sprintf("Error: %s", err))
		 }
 
		 // Return the command output to the user
		 return c.HTMLBlob(http.StatusOK, output)
	})
	// Start server
	e.Logger.Fatal(e.Start(":3000"))
}
