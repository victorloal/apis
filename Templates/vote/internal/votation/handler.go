package votation

import (
	"api/internal/models"
	"encoding/json"
	"net/http"
)

func CreateElectionHandler(w http.ResponseWriter, r *http.Request) {
	// verificar que no esten vacios con service
	var election models.Election
	_ = json.NewDecoder(r.Body).Decode(&election)

	// enviar a 


}