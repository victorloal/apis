package main

import (
	"api/internal/votation/election"
	"net/http"

	"github.com/rs/cors"
	"github.com/gorilla/mux"
)

func main() {
	
	r := mux.NewRouter()
	election.RegisterRoutes(r)
	
	// Configura CORS
	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{
			"*", 
		},
		AllowedMethods: []string{
			"GET", "POST", "OPTIONS", "PUT", "DELETE",
		},
		AllowedHeaders: []string{
			"Content-Type", "Authorization", "X-Requested-With",
		},
		AllowCredentials: true,  // Si necesitas enviar cookies o cabeceras de autenticación
	})

	// Aplica CORS al router
	http.Handle("/", corsHandler.Handler(r))

	println("corriendo en el puerto localhost:3000")
	http.ListenAndServe(":3000", nil)
}
