package election

import (
	"api/internal/models"
	"api/pkg/db"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func GetAllHandler(w http.ResponseWriter, r *http.Request) {
	var election []models.Election
	db.DB.Find(&election)
	json.NewEncoder(w).Encode(&election)
}

func GetIdHandler(w http.ResponseWriter, r *http.Request) {
	var election models.Election
	param := mux.Vars(r)
	
	result := db.DB.First(&election, param["id"])
	result.Preload("ElectionAuthorities","Candidates","Voters","Status","AuditConfig","HomomorphicKey","TallyResult")


	if result.Error != nil || election.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("User not found"))
		return
	}
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(&election)
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	var election models.Election
	_ = json.NewDecoder(r.Body).Decode(&election)

	createElection := db.DB.Create(&election)
	err := createElection.Error
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	json.NewEncoder(w).Encode(&election)
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	var election models.Election
	params := mux.Vars(r)

	result := db.DB.Preload("ElectionAuthorities","Candidates","Voters","Status","AuditConfig","HomomorphicKey","TallyResult").First(&election, params["id"])

	if result.Error != nil || election.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("election not found"))
		return
	}
	db.DB.Preload("Tasks").Delete(&election, params["id"])
	w.WriteHeader(http.StatusOK)

}

func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	var election models.Election
	db.DB.Find(&election)
	json.NewEncoder(w).Encode(&election)
}
