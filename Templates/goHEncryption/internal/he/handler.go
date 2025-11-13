package he

import (
	"encoding/json"
	"fmt"
	"net/http"
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

func (h *Handler) GenerateKeys(w http.ResponseWriter, r *http.Request) {
	// 1. Verificar que el content-type es JSON
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type debe ser application/json", http.StatusBadRequest)
		return
	}

	// 2. Leer el JSON del body como estructura Init
	var initData Init
	if err := json.NewDecoder(r.Body).Decode(&initData); err != nil {
		http.Error(w, "JSON inválido: "+err.Error(), http.StatusBadRequest)
		return
	}

	// 3. Validar los datos
	if initData.T <= 0 {
		http.Error(w, "T debe ser mayor a 0", http.StatusBadRequest)
		return
	}

	if len(initData.N) == 0 {
		http.Error(w, "La lista N no puede estar vacía", http.StatusBadRequest)
		return
	}

	// 4. Validar cada elemento Ae
	for i, ae := range initData.N {
		if ae.Name == "" {
			http.Error(w, fmt.Sprintf("El campo 'name' es requerido en el elemento %d", i), http.StatusBadRequest)
			return
		}
		if ae.Mail == "" {
			http.Error(w, fmt.Sprintf("El campo 'mail' es requerido en el elemento %d", i), http.StatusBadRequest)
			return
		}
	}

	// 5. Procesar (pasamos toda la estructura Init al servicio)
	response, err := h.service.GenerateKeys(initData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 6. Responder
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) NuevoSaludo(w http.ResponseWriter, r *http.Request) {

	// Generar saludo
	response, err := h.service.Nuevo()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Responder con JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
