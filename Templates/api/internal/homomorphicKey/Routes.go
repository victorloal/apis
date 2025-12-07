package homomorphicKey

import (
	"github.com/gorilla/mux"
)

func RegisterRoutes(router *mux.Router) {
	r := router.PathPrefix("/homomorphicKey").Subrouter()

	r.HandleFunc("", GetAllHandler).Methods("GET")
	r.HandleFunc("/{id}", GetIdHandler).Methods("GET")
	r.HandleFunc("", PostHandler).Methods("POST")
	r.HandleFunc("/{id}", DeleteHandler).Methods("DELETE")
	r.HandleFunc("/{id}", UpdateHandler).Methods("UPDATE")
}
