package election

import (
	"github.com/gorilla/mux"
)

// RegisterRoutes registra las rutas para este handler
func RegisterRoutes(router *mux.Router) {
	r := router.PathPrefix("/election").Subrouter()

	r.HandleFunc("", CreateElectionHandler).Methods("POST")
	r.HandleFunc("", UpdateElectionHandler).Methods("PUT")
	r.HandleFunc("/{id}", DeleteElectionHandler).Methods("DELETE")
	
	r1 := r.PathPrefix("/drafts").Subrouter()
	r1.HandleFunc("", CreateDraftHandler).Methods("POST")
	r1.HandleFunc("", UpdateDraftHandler).Methods("PUT")
	r1.HandleFunc("/{id}", DeleteDraftHandler).Methods("DELETE")



}

// 
// func CreateElectionHandler(w http.ResponseWriter, r *http.Request) {
// 	// 1. Decodificar JSON
// 	var election models.Election
// 	println(1)
// 	if err := json.NewDecoder(r.Body).Decode(&election); err != nil {
// 		http.Error(w, "JSON inválido", http.StatusBadRequest)
// 		return
// 	}
// 
// 	println(2)
// 	w.Header().Set("Content-Type", "application/json")
// 
// 	// 2. Crear UI para la votación
// 	ui := uuid.New().String()
// 
// 	// 3. Formatear fechas (Fabric requiere RFC3339)
// 	inicio := election.StartDate
// 	fin := election.EndDate
// 
// 	//-----------------------------------------------------------
// 	// A) CREAR LA ELECCIÓN
// 	//-----------------------------------------------------------
// 	print(3)
// 	url := "http://localhost:3002/newElection"
// 
// 	data := map[string]interface{}{
// 		"ui":     ui,
// 		"nombre": election.Name,      // SIN &
// 		"inicio": inicio,             // formato string RFC3339
// 		"fin":    fin,
// 	}
// 
// 	print(4)
// 	resp, err := SendJSON(url, data)
// 	if err != nil {
// 		http.Error(w, "Error creando votación en blockchain: "+err.Error(), 500)
// 		return
// 	}
// 	defer resp.Body.Close()
// 
// 	if resp.StatusCode != 200 {
// 		http.Error(w, "Blockchain rechazó la creación de la votación", resp.StatusCode)
// 		return
// 	}
// 	print(5)
// 
// 	//-----------------------------------------------------------
// 	// B) AGREGAR VOTANTES
// 	//-----------------------------------------------------------
// 
// 	url = "http://localhost:3002/newVoto"
// print(6)
// 	for _, voter := range election.Voters {
// 		data := map[string]interface{}{
// 			"uiVotacion": ui,
// 			"idVoter":    voter.Token, // SIN &
// 		}
// 
// 		resp, err := SendJSON(url, data)
// 		if err != nil {
// 			http.Error(w, "Error agregando votante: "+err.Error(), 500)
// 			return
// 		}
// 		resp.Body.Close()
// 
// 		if resp.StatusCode != 200 {
// 			http.Error(w, "Blockchain rechazó un votante", resp.StatusCode)
// 			return
// 		}
// 	}
// print(7)
// 	//-----------------------------------------------------------
// 	// C) AGREGAR CANDIDATOS
// 	//-----------------------------------------------------------
// 
// 	url = "http://localhost:3002/newCandidate"
// 
// 	for _, candidate := range election.Candidates {
// 		data := map[string]interface{}{
// 			"uiVotacion": ui,
// 			"uiCandidato": candidate.Id, // SIN &
// 		}
// 
// 		resp, err := SendJSON(url, data)
// 		if err != nil {
// 			http.Error(w, "Error agregando candidato: "+err.Error(), 500)
// 			return
// 		}
// 		resp.Body.Close()
// 
// 		if resp.StatusCode != 200 {
// 			http.Error(w, "Blockchain rechazó un candidato", resp.StatusCode)
// 			return
// 		}
// 	}
// print(8)
// 	//-----------------------------------------------------------
// 	// D) Respuesta final al frontend
// 	//-----------------------------------------------------------
// 
// 	json.NewEncoder(w).Encode(map[string]interface{}{
// 		"ok":     true,
// 		"ui":     ui,
// 		"status": "Election created",
// 	})
// print(9)
// 	// Log en servidor
// 	fmt.Println("Elección creada correctamente:", ui)
// }

