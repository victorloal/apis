package ballot

import (
	"api/internal/models"
	"api/pkg/db"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func GetAllHandler(w http.ResponseWriter, r *http.Request) {
	var ballots []models.Ballot
	db.DB.Preload("Election").Find(&ballots)
	json.NewEncoder(w).Encode(&ballots)
}

func GetIdHandler(w http.ResponseWriter, r *http.Request) {
	var ballot models.Ballot
	param := mux.Vars(r)

	result := db.DB.First(&ballot, param["id"])
	result.Preload("Election")

	if result.Error != nil || ballot.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("ballot not found"))
		return
	}
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(&ballot)
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	var ballot models.Ballot
	_ = json.NewDecoder(r.Body).Decode(&ballot)

	createElection := db.DB.Create(&ballot)
	err := createElection.Error
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	json.NewEncoder(w).Encode(&ballot)
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	var ballot models.Ballot
	params := mux.Vars(r)

	result := db.DB.First(&ballot, params["id"])

	if result.Error != nil || ballot.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("ballot not found"))
		return
	}
	db.DB.Preload("Election").Delete(&ballot, params["id"])
	w.WriteHeader(http.StatusOK)

}

func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	var ballot []models.Ballot
	db.DB.Find(&ballot)
	json.NewEncoder(w).Encode(&ballot)
}
