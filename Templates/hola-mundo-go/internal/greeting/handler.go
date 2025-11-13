package greeting

import (
	"encoding/json"
	"net/http"
	
	"github.com/gorilla/mux"
)

// Handler maneja las requests HTTP
type Handler struct {
	service Service
}

// NewHandler crea un nuevo handler
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// GreetHandler maneja el endpoint de saludo
func (h *Handler) GreetHandler(w http.ResponseWriter, r *http.Request) {
	// Obtener el nombre del query parameter
	name := r.URL.Query().Get("name")
	
	// Si no viene en query, intentar del JSON body
	if name == "" {
		var req GreetingRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err == nil {
			name = req.Name
		}
	}
	
	// Generar saludo
	response, err := h.service.GenerateGreeting(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	// Responder con JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// RegisterRoutes registra las rutas para este handler
func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/greet", h.GreetHandler).Methods("GET", "POST")
	router.HandleFunc("/greet/{name}", h.GreetHandler).Methods("GET")
}