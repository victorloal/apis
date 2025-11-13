package greeting

import "github.com/gorilla/mux"

// SetupRoutes configura todas las rutas del dominio greeting
func SetupRoutes(router *mux.Router) {
	// Crear instancias
	service := NewService()
	handler := NewHandler(service)
	
	// Registrar rutas bajo el path /api
	apiRouter := router.PathPrefix("/api").Subrouter()
	handler.RegisterRoutes(apiRouter)
}