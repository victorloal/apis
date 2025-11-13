package main

import (
	"log"
	"net/http"
	
	"hola-mundo-go/internal/greeting"
	
	"github.com/gorilla/mux"
)

func main() {
	// Crear router
	router := mux.NewRouter()
	
	// Configurar rutas del greeting
	greeting.SetupRoutes(router)
	
	// Ruta de health check
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status": "ok"}`))
	}).Methods("GET")
	
	// Iniciar servidor
	log.Println("Servidor iniciado en :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}