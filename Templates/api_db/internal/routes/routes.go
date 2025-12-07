package routes

import (
	"api_db/internal/handlers"
	"api_db/internal/repositories"
	"api_db/internal/services"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func SetupRoutes(db *gorm.DB) *mux.Router {
    router := mux.NewRouter()

    // Initialize all repositories
    electionRepo := repositories.NewElectionRepository(db)
    authoritiesRepo := repositories.NewAuthoritiesRepository(db)
    votersRepo := repositories.NewVotersRepository(db)
    candidatesRepo := repositories.NewCandidatesRepository(db)
    ballotsRepo := repositories.NewBallotsRepository(db)
    homomorphicKeysRepo := repositories.NewHomomorphicKeysRepository(db)
    tallyResultsRepo := repositories.NewTallyResultsRepository(db)
    statusRepo := repositories.NewStatusRepository(db)
    auditConfigRepo := repositories.NewAuditConfigRepository(db)
    auditLogsRepo := repositories.NewAuditLogsRepository(db)

    // Initialize all services
    electionService := services.NewElectionService(electionRepo)
    authoritiesService := services.NewAuthoritiesService(authoritiesRepo)
    votersService := services.NewVotersService(votersRepo)
    candidatesService := services.NewCandidatesService(candidatesRepo)
    ballotsService := services.NewBallotsService(ballotsRepo)
    homomorphicKeysService := services.NewHomomorphicKeysService(homomorphicKeysRepo)
    tallyResultsService := services.NewTallyResultsService(tallyResultsRepo)
    statusService := services.NewStatusService(statusRepo)
    auditConfigService := services.NewAuditConfigService(auditConfigRepo)
    auditLogsService := services.NewAuditLogsService(auditLogsRepo)

    // Initialize all handlers
    electionHandler := handlers.NewElectionHandler(electionService)
    authoritiesHandler := handlers.NewAuthoritiesHandler(authoritiesService)
    votersHandler := handlers.NewVotersHandler(votersService)
    candidatesHandler := handlers.NewCandidatesHandler(candidatesService)
    ballotsHandler := handlers.NewBallotsHandler(ballotsService)
    homomorphicKeysHandler := handlers.NewHomomorphicKeysHandler(homomorphicKeysService)
    tallyResultsHandler := handlers.NewTallyResultsHandler(tallyResultsService)
    statusHandler := handlers.NewStatusHandler(statusService)
    auditConfigHandler := handlers.NewAuditConfigHandler(auditConfigService)
    auditLogsHandler := handlers.NewAuditLogsHandler(auditLogsService)

    // Election routes
    electionRouter := router.PathPrefix("/elections").Subrouter()
    electionRouter.HandleFunc("", electionHandler.CreateElection).Methods("POST")
    electionRouter.HandleFunc("", electionHandler.GetAllElections).Methods("GET")
    electionRouter.HandleFunc("/{id}", electionHandler.GetElection).Methods("GET")
    electionRouter.HandleFunc("/{id}", electionHandler.UpdateElection).Methods("PUT")
    electionRouter.HandleFunc("/{id}", electionHandler.DeleteElection).Methods("DELETE")

    // Authorities routes
    authoritiesRouter := router.PathPrefix("/authorities").Subrouter()
    authoritiesRouter.HandleFunc("", authoritiesHandler.CreateAuthority).Methods("POST")
    authoritiesRouter.HandleFunc("/{id}", authoritiesHandler.GetAuthority).Methods("GET")
    authoritiesRouter.HandleFunc("/{id}", authoritiesHandler.UpdateAuthority).Methods("PUT")
    authoritiesRouter.HandleFunc("/{id}", authoritiesHandler.UpdateAuthorityPartial).Methods("PATCH")
    authoritiesRouter.HandleFunc("/{id}", authoritiesHandler.DeleteAuthority).Methods("DELETE")
    authoritiesRouter.HandleFunc("/{id}/password", authoritiesHandler.UpdateAuthorityPassword).Methods("PUT")
    authoritiesRouter.HandleFunc("/{id}/secret-key", authoritiesHandler.UpdateAuthoritySecretKey).Methods("PUT")
    authoritiesRouter.HandleFunc("/email/{email}", authoritiesHandler.GetAuthorityByEmail).Methods("GET")
    authoritiesRouter.HandleFunc("/election/{electionId}", authoritiesHandler.GetAuthoritiesByElection).Methods("GET")

    // Voters routes
    votersRouter := router.PathPrefix("/voters").Subrouter()
    votersRouter.HandleFunc("", votersHandler.CreateVoter).Methods("POST")
    votersRouter.HandleFunc("/{id}", votersHandler.GetVoter).Methods("GET")
    votersRouter.HandleFunc("/{id}", votersHandler.UpdateVoter).Methods("PUT")
    votersRouter.HandleFunc("/{id}", votersHandler.DeleteVoter).Methods("DELETE")
    votersRouter.HandleFunc("/{id}/vote-status", votersHandler.UpdateVoteStatus).Methods("PUT")
    votersRouter.HandleFunc("/election/{electionId}", votersHandler.GetVotersByElection).Methods("GET")
    votersRouter.HandleFunc("/token/{token}", votersHandler.GetVoterByToken).Methods("GET")

    // Candidates routes
    candidatesRouter := router.PathPrefix("/candidates").Subrouter()
    candidatesRouter.HandleFunc("", candidatesHandler.CreateCandidate).Methods("POST")
    candidatesRouter.HandleFunc("/{id}", candidatesHandler.GetCandidate).Methods("GET")
    candidatesRouter.HandleFunc("/{id}", candidatesHandler.UpdateCandidate).Methods("PUT")
    candidatesRouter.HandleFunc("/{id}", candidatesHandler.DeleteCandidate).Methods("DELETE")
    candidatesRouter.HandleFunc("/election/{electionId}", candidatesHandler.GetCandidatesByElection).Methods("GET")
    candidatesRouter.HandleFunc("/election/{electionId}/order", candidatesHandler.GetCandidatesByOrder).Methods("GET")

    // Ballots routes
    ballotsRouter := router.PathPrefix("/ballots").Subrouter()
    ballotsRouter.HandleFunc("", ballotsHandler.CreateBallot).Methods("POST")
    ballotsRouter.HandleFunc("/election/{electionId}/voter/{voterId}/id/{id}", ballotsHandler.GetBallot).Methods("GET")
    ballotsRouter.HandleFunc("/election/{electionId}/voter/{voterId}/id/{id}", ballotsHandler.UpdateBallot).Methods("PUT")
    ballotsRouter.HandleFunc("/election/{electionId}/voter/{voterId}/id/{id}", ballotsHandler.DeleteBallot).Methods("DELETE")
    ballotsRouter.HandleFunc("/election/{electionId}", ballotsHandler.GetBallotsByElection).Methods("GET")
    ballotsRouter.HandleFunc("/voter/{voterId}", ballotsHandler.GetBallotsByVoter).Methods("GET")
    ballotsRouter.HandleFunc("/election/{electionId}/with-details", ballotsHandler.GetBallotsWithVoterDetails).Methods("GET")

    // Homomorphic Keys routes
    keysRouter := router.PathPrefix("/keys").Subrouter()
    keysRouter.HandleFunc("", homomorphicKeysHandler.CreateKey).Methods("POST")
    keysRouter.HandleFunc("/{id}", homomorphicKeysHandler.GetKey).Methods("GET")
    keysRouter.HandleFunc("/{id}", homomorphicKeysHandler.UpdateKey).Methods("PUT")
    keysRouter.HandleFunc("/{id}", homomorphicKeysHandler.DeleteKey).Methods("DELETE")
    keysRouter.HandleFunc("/election/{electionId}", homomorphicKeysHandler.GetKeyByElection).Methods("GET")
    keysRouter.HandleFunc("/election/{electionId}/params", homomorphicKeysHandler.UpdateKeyParams).Methods("PUT")

    // Tally Results routes
    tallyRouter := router.PathPrefix("/tally-results").Subrouter()
    tallyRouter.HandleFunc("", tallyResultsHandler.CreateTallyResult).Methods("POST")
    tallyRouter.HandleFunc("/{id}", tallyResultsHandler.GetTallyResult).Methods("GET")
    tallyRouter.HandleFunc("/{id}", tallyResultsHandler.UpdateTallyResult).Methods("PUT")
    tallyRouter.HandleFunc("/{id}", tallyResultsHandler.DeleteTallyResult).Methods("DELETE")
    tallyRouter.HandleFunc("/election/{electionId}", tallyResultsHandler.GetTallyResultByElection).Methods("GET")
    tallyRouter.HandleFunc("/election/{electionId}/compute", tallyResultsHandler.ComputeTallyResult).Methods("POST")
    tallyRouter.HandleFunc("/with-details", tallyResultsHandler.GetTallyResultsWithElection).Methods("GET")

    // Status routes
    statusRouter := router.PathPrefix("/status").Subrouter()
    statusRouter.HandleFunc("", statusHandler.CreateStatus).Methods("POST")
    statusRouter.HandleFunc("", statusHandler.GetAllStatus).Methods("GET")
    statusRouter.HandleFunc("/{id}", statusHandler.GetStatus).Methods("GET")
    statusRouter.HandleFunc("/{id}", statusHandler.UpdateStatus).Methods("PUT")
    statusRouter.HandleFunc("/{id}", statusHandler.DeleteStatus).Methods("DELETE")
    statusRouter.HandleFunc("/name/{name}", statusHandler.GetStatusByName).Methods("GET")

    // Audit Config routes
    auditConfigRouter := router.PathPrefix("/audit-config").Subrouter()
    auditConfigRouter.HandleFunc("", auditConfigHandler.CreateAuditConfig).Methods("POST")
    auditConfigRouter.HandleFunc("/{id}", auditConfigHandler.GetAuditConfig).Methods("GET")
    auditConfigRouter.HandleFunc("/{id}", auditConfigHandler.UpdateAuditConfig).Methods("PUT")
    auditConfigRouter.HandleFunc("/{id}", auditConfigHandler.DeleteAuditConfig).Methods("DELETE")
    auditConfigRouter.HandleFunc("/election/{electionId}", auditConfigHandler.GetAuditConfigByElection).Methods("GET")
    auditConfigRouter.HandleFunc("/election/{electionId}/ballot-audit", auditConfigHandler.EnableBallotAudit).Methods("PUT")
    auditConfigRouter.HandleFunc("/election/{electionId}/access-logs", auditConfigHandler.EnableAccessLogs).Methods("PUT")

    // Audit Logs routes
    auditLogsRouter := router.PathPrefix("/audit-logs").Subrouter()
    auditLogsRouter.HandleFunc("", auditLogsHandler.CreateAuditLog).Methods("POST")
    auditLogsRouter.HandleFunc("/{id}", auditLogsHandler.GetAuditLog).Methods("GET")
    auditLogsRouter.HandleFunc("/{id}", auditLogsHandler.DeleteAuditLog).Methods("DELETE")
    auditLogsRouter.HandleFunc("/election/{electionId}", auditLogsHandler.GetAuditLogsByElection).Methods("GET")
    auditLogsRouter.HandleFunc("/action/{action}", auditLogsHandler.GetAuditLogsByAction).Methods("GET")
    auditLogsRouter.HandleFunc("/user/{userType}/{userId}", auditLogsHandler.GetAuditLogsByUser).Methods("GET")
    auditLogsRouter.HandleFunc("/date-range", auditLogsHandler.GetAuditLogsByDateRange).Methods("POST")
    auditLogsRouter.HandleFunc("/election/{electionId}/vote", auditLogsHandler.LogVoteAction).Methods("POST")
    auditLogsRouter.HandleFunc("/election/{electionId}/authority", auditLogsHandler.LogAuthorityAction).Methods("POST")

    // Health check
    router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        w.Write([]byte(`{"status": "ok"}`))
    }).Methods("GET")

    return router
}