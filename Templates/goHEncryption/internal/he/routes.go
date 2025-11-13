package he

import "github.com/gorilla/mux"

// SetupRoutes configura todas las rutas del dominio he
func SetupRoutes(router *mux.Router) {
	// Crear instancias
	service := NewService()
	handler := NewHandler(service)
	
	// Registrar rutas bajo el path /api
	apiRouter := router.PathPrefix("/he").Subrouter()
	handler.RegisterRoutes(apiRouter)
}

// RegisterRoutes registra las rutas para este handler
func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/greet", h.GreetHandler).Methods("GET", "POST")
	router.HandleFunc("/greet/{name}", h.GreetHandler).Methods("GET")
	router.HandleFunc("/keys/generate", h.GenerateKeys).Methods("POST")
	router.HandleFunc("/keys/public/{name}", h.GreetHandler).Methods("POST")
	router.HandleFunc("/keys/public", h.GreetHandler).Methods("GET")

	router.HandleFunc("/vote", h.NuevoSaludo).Methods("POST")
	router.HandleFunc("/vote/count", h.GreetHandler).Methods("GET")
	router.HandleFunc("/vote/decrypt", h.GreetHandler).Methods("POST")
	router.HandleFunc("/vote/count", h.GreetHandler).Methods("GET")
}