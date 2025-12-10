package election

import (
	"api/internal/models"
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
)



func CreateElectionHandler(w http.ResponseWriter, r *http.Request) {
	// verificar que no esten vacios con service
	var election models.Election
	println(1)
	_ = json.NewDecoder(r.Body).Decode(&election)
		
	//devolver que fue resivido el json y devolver el json
	w.Header().Set("Content-Type", "application/json")
	
	
	//crea uuid en una variable nueva
	ui := strings.ToUpper(strings.ReplaceAll(uuid.New().String(), "-", ""))
	println(1)
	//enviar la elección al blockchain
	url := "http://localhost:3002/newElection"
	println(01)
	data := map[string]interface{}{
		"ui":     ui,
		"nombre": &election.Name,
		"inicio": &election.StartDate,
		"fin":    &election.EndDate,
	}
	println(01)
	
	resp, err := SendJSON(url, data)
	println(01)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		println(01)
		return
	defer resp.Body.Close()
	println(01)

	time.Sleep(5 * time.Second)
	println(2)
	//enviar los votantes al blockchain
	url = "http://localhost:3002/newVoto"
	for _, voter := range election.Voters {
		data := map[string]interface{}{
			"uiVotacion":     ui,
			"IdVoter": &voter.Token,
		}
		
		resp, err := SendJSON(url, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()
	}

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Error al crear votantes  en el blockchain", resp.StatusCode)
		return
	}
	


	println(3)
	//enviar los votantes al blockchain
	url = "http://localhost:3002/newVoto"
	for _, candidate := range election.Candidates {
		data := map[string]interface{}{
			"uiVotacion":     ui,
			"uiCandidato": &candidate.Id,
		}
		
		resp, err := SendJSON(url, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()
	}

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Error al crear candidatos en el blockchain", resp.StatusCode)
		return
	}
	


	
	json.NewEncoder(w).Encode(election)
	// imprimir en la terminal
	
	json.NewEncoder(os.Stdout).Encode(election)
	
	

}
}

func UpdateElectionHandler(w http.ResponseWriter, r *http.Request) {
	// verificar que no esten vacios con service
	var election models.Election
	_ = json.NewDecoder(r.Body).Decode(&election)

	// enviar a

}

func DeleteElectionHandler(w http.ResponseWriter, r *http.Request) {
	// verificar que no esten vacios con service
	var election models.Election
	_ = json.NewDecoder(r.Body).Decode(&election)

	// enviar a

}

func CreateDraftHandler(w http.ResponseWriter, r *http.Request) {
	// verificar que no esten vacios con service
	var election models.Election
	_ = json.NewDecoder(r.Body).Decode(&election)

	// enviar a

}

func UpdateDraftHandler(w http.ResponseWriter, r *http.Request) {
	// verificar que no esten vacios con service
	var election models.Election
	_ = json.NewDecoder(r.Body).Decode(&election)

	// enviar a

}

func DeleteDraftHandler(w http.ResponseWriter, r *http.Request) {
	// verificar que no esten vacios con service
	var election models.Election
	_ = json.NewDecoder(r.Body).Decode(&election)

	// enviar a

}

func SendJSON(url string, data interface{}) (*http.Response, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	return client.Do(req)
}