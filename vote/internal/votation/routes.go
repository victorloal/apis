package votation

import "github.com/gorilla/mux"

// RegisterRoutes registra las rutas para este handler
func RegisterRoutes(router *mux.Router) {
	r1 := router.PathPrefix("/election").Subrouter()

	r1.HandleFunc("", CreateElectionHandler).Methods("POST")
	r1.HandleFunc("", UpdateElectionHandler).Methods("PUT")
	r1.HandleFunc("/{id}", DeleteElectionHandler).Methods("DELETE")

	r := r1.PathPrefix("/{id}/candidate").Subrouter()
	r.HandleFunc("", CreateCandidateHandler).Methods("POST")
	r.HandleFunc("/list", CreateCandidateListHandler).Methods("POST")
	r.HandleFunc("/{id}", UpdateCandidateHandler).Methods("PUT")
	r.HandleFunc("/{id}", DeleteCandidateHandler).Methods("DELETE")

	r = r1.PathPrefix("/{id}/voter").Subrouter()
	r.HandleFunc("", CreateVoterHandler).Methods("POST")
	r.HandleFunc("/list", CreateVoterListHandler).Methods("POST")
	r.HandleFunc("/{id}", UpdateVoterHandler).Methods("PUT")
	r.HandleFunc("/{id}", VerifyVoterHandler).Methods("GET")

	r = r1.PathPrefix("/{id}/{voter_id}/ballot").Subrouter()
	r.HandleFunc("", CreateBallotHandler).Methods("POST")
	r.HandleFunc("/{id}", VerifyBallotHandler).Methods("GET")
}
