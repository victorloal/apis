package main

import (
	"api/internal/auditConfig"
	"api/internal/ballot"
	"api/internal/candidate"
	"api/internal/election"
	"api/internal/electionAuthority"
	"api/internal/homomorphicKey"
	"api/internal/status"
	"api/internal/statusBallot"
	"api/internal/tallyResult"
	"api/internal/voter"
	"api/pkg/db"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	db.DBConetion()
	// corregir
	db.Migrate(db.DB)
	
	db.Initialize(db.DB)
	r := mux.NewRouter()
	auditConfig.RegisterRoutes(r)
	ballot.RegisterRoutes(r)
	candidate.RegisterRoutes(r)
	election.RegisterRoutes(r)
	electionAuthority.RegisterRoutes(r)
	homomorphicKey.RegisterRoutes(r)
	status.RegisterRoutes(r)
	statusBallot.RegisterRoutes(r)
	tallyResult.RegisterRoutes(r)
	voter.RegisterRoutes(r)


	// routes.RegisterRoutes(r)

	http.ListenAndServe(":3001", r)
}
