package site

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"

	"github.com/fortify-presales/insecure-go-api/internal/config"
	"github.com/fortify-presales/insecure-go-api/pkg/log"
)

// SiteHandler is a struct that contains the logger and configuration for the site API
type SiteHandler struct {
	logger log.Logger
	cfg    *config.Config
}

func MakeHTTPHandler(logger log.Logger, cfg *config.Config) http.Handler {

	// Initialize handlers
	siteHandler := &SiteHandler{
		logger: logger,
		cfg:    cfg,
	}

	router := http.NewServeMux()
	router.HandleFunc("GET /api/v1/site/ping", siteHandler.PingSiteByQuery)
	router.HandleFunc("POST /api/v1/site/ping", siteHandler.PingSiteByBody)
	router.HandleFunc("GET /api/v1/site/download/{id}", siteHandler.DownloadFileById)

	return router
}

// GET request with data flow taint source in URL path
//
// @Summary      Ping Site by Query
// @Description  Ping a Site using URL query parameter
// @Tags         site
// @Accept       json
// @Produce      json
// @Param		 hostname	query		string				true	"hostname"	example("localhost")
// @Success      200  {string}  "string"
// @Failure      400  {object}  model.APIError
// @Failure      500  {object}  model.APIError
// @Router       /site/ping [get]
func (s *SiteHandler) PingSiteByQuery(w http.ResponseWriter, r *http.Request) {
	s.logger.Infof("Handling GET at %s\n", r.URL.Path)
	//
	// Get hostname from query parameter
	//
	host := r.URL.Query().Get("hostname")
	if host == "" {
		http.Error(w, "Hostname not provided", http.StatusBadRequest)
		return
	}
	//
	// Command Injection : dataflow
	//
	cmd := exec.Command("ping", "-c", "4", host)
	output, err := cmd.CombinedOutput()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error: %s", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write(output)
}

// POST request with data flow taint source in body
//
// @Summary      Ping Site by Body
// @Description  Ping a Site using JSON Body
// @Description  JSON Body should contain a "hostname" field
// @Description  Example: {"hostname": "localhost"}
// @Description  This is a JSON Injection vulnerability example
// @Tags         site
// @Accept       json
// @Produce      json
// @Param		 Site	body		Site				true	"Site"
// @Success      200  {string}  "string"
// @Failure      400  {object}  model.APIError
// @Failure      500  {object}  model.APIError
// @Router       /site/ping [post]
func (s *SiteHandler) PingSiteByBody(w http.ResponseWriter, r *http.Request) {
	s.logger.Infof("Handling POST at %s\n", r.URL.Path)
	type JsonString struct {
		Hostname string `json:"hostname"`
	}
	var jsonDataToRead JsonString
	//err2 := json.Unmarshal(body, &jsonData)
	err := json.NewDecoder(r.Body).Decode(&jsonDataToRead)
	if err != nil {
		http.Error(w, "Hostname not provided", http.StatusBadRequest)
		return
	}
	//
	// JSON Injection : dataflow
	//
	jsonDataToWrite := map[string]string{
		"command":  "ping",
		"hostname": jsonDataToRead.Hostname,
		"output":   "", // Placeholder for actual output
	}
	s.logger.Infof("Creating file 'command_log.json' with contents: %+v\n", jsonDataToWrite)
	file, _ := os.OpenFile("command_log.json", os.O_CREATE|os.O_WRONLY, os.ModePerm)
	defer file.Close()
	jsonEncoder := json.NewEncoder(file)
	jsonEncoder.SetIndent("", "  ") // Optional: Pretty-print the JSON
	jsonEncoder.Encode(jsonDataToWrite)
	// TODO: actual ping logic can be added here and output placed in "jsonDataToWrite"
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

// GET request with data flow taint source in URL path
//
// @Summary      Download File
// @Description  Download a file by ID
// @Tags         site
// @Accept       json
// @Produce      json
// @Param		 id	path		string				true	"id"	example("12345")
// @Success      200  {string}  "string"
// @Failure      400  {object}  model.APIError
// @Failure      404  {object}  model.APIError
// @Failure      500  {object}  model.APIError
// @Router       /site/download/{id} [get]
func (s *SiteHandler) DownloadFileById(w http.ResponseWriter, r *http.Request) {
	s.logger.Infof("Handling GET at %s\n", r.URL.Path)
	//
	// Get id from URL path
	//
	// PathValue is new in Go 1.22 - Not yet supported by Fortify
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Id not provided", http.StatusBadRequest)
		return
	}
	//
	// Path Manipulation : dataflow
	//
	filename := fmt.Sprintf("%s%c%s%c%s", os.Getenv("PWD"), os.PathSeparator, "downloads", os.PathSeparator, id)
	s.logger.Infof("Retrieving contents of file path: %s\n", filename)
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	data, _ := ioutil.ReadFile(filename)
	w.Header().Set("Content-Type", "text/plain")
	w.Write(data)
}
