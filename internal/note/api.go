package note

import (
	"encoding/json"
	"errors"
	"net/http"

	model "github.com/fortify-presales/insecure-go-api/internal/models"
)

// NoteHandler organizes HTTP handler functions for CRUD on Note entity
type NoteHandler struct {
	Repository Repository // interface for persistence
}

func MakeHTTPHandler(repo Repository) http.Handler {

	// Iniitialize handlers
	noteHandler := &NoteHandler{
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

// Post handles HTTP Post
//
// @Summary      Create Note
// @Description  Create a new Note
// @Tags         notes
// @Accept       json
// @Produce      json
// @Param		 Note	body		Note			true	"Note"
// @Success      200  {object}  Note
// @Failure      400  {object}  model.APIError
// @Failure      404  {object}  model.APIError
// @Failure      500  {object}  model.APIError
// @Router       /notes/ [post]
func (h *NoteHandler) Post(w http.ResponseWriter, r *http.Request) {
	var note Note
	// Decode the incoming note json
	err := json.NewDecoder(r.Body).Decode(&note)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create note
	if _, err := h.Repository.Create(note); err != nil {
		if errors.Is(err, ErrNoteExists) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// GetAll handles HTTP Get with no Id
//
// @Summary      Get Notes
// @Description  Get all Notes
// @Tags         notes
// @Accept       json
// @Produce      json
// @Param        keywords    query     string  false  "search by keywords"  example(alphadex)
// @Success      200  {array}  	Note
// @Failure      400  {object}  model.APIError
// @Failure      404  {object}  model.APIError
// @Failure      500  {object}  model.APIError
// @Router       /notes [get]
func (h *NoteHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	keywords := r.URL.Query().Get("keywords")
	// Get all
	if notes, err := h.Repository.GetAll(keywords); err != nil {
		if errors.Is(err, model.ErrNotFound) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return

	} else {
		j, err := json.Marshal(notes)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(j)
	}
}

// Get handles HTTP Get with Id
//
// @Summary      Get Note
// @Description  Get a Note
// @Tags         notes
// @Accept       json
// @Produce      json
// @Param		 id	path		string				true	"Note ID"
// @Success      200  {object}  Note
// @Failure      400  {object}  model.APIError
// @Failure      404  {object}  model.APIError	"Could not find Note Id"
// @Failure      500  {object}  model.APIError
// @Router       /notes/{id} [get]
func (h *NoteHandler) Get(w http.ResponseWriter, r *http.Request) {
	// Getting route parameter id
	id := r.PathValue("id")
	// Get by id
	if note, err := h.Repository.GetById(id); err != nil {
		if errors.Is(err, model.ErrNotFound) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		w.Header().Set("Content-Type", "application/json")
		j, err := json.Marshal(note)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusOK)
		w.Write(j)
	}
}

// Put handles HTTP Put with Id
//
// @Summary      Update Note
// @Description  Update an existing Note
// @Tags         notes
// @Accept       json
// @Produce      json
// @Param		 id		path	string				true	"Note ID"
// @Param		 Note	body	Note			true	"Note"
// @Success      200  {object}  Note
// @Failure      400  {object}  model.APIError
// @Failure      404  {object}  model.APIError	"Could not find Note Id"
// @Failure      500  {object}  model.APIError
// @Router       /notes/{id} [put]
func (h *NoteHandler) Put(w http.ResponseWriter, r *http.Request) {
	// Getting route parameter id
	id := r.PathValue("id")
	var note Note
	// Decode the incoming note json
	err := json.NewDecoder(r.Body).Decode(&note)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Update
	if err := h.Repository.Update(id, note); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Delete handles HTTP Delete with Id
//
// @Summary      Delete Note
// @Description  Delete a Note
// @Tags         notes
// @Accept       json
// @Produce      json
// @Param		 id		path	string				true	"Note ID"
// @Success      200  {object}  model.APIMessage
// @Failure      400  {object}  model.APIError
// @Failure      404  {object}  model.APIError	"Could not find Note Id"
// @Failure      500  {object}  model.APIError
// @Router       /notes/{id} [delete]
func (h *NoteHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// Getting route parameter id
	id := r.PathValue("id")
	// delete
	if err := h.Repository.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
