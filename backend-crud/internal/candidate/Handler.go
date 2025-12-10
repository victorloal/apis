package candidate

import (
	"api/internal/models"
	"api/pkg/db"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func GetAllHandler(w http.ResponseWriter, r *http.Request) {
	var candidates []models.Candidate
	db.DB.Preload("Election").Find(&candidates)
	json.NewEncoder(w).Encode(&candidates)
}

func GetIdHandler(w http.ResponseWriter, r *http.Request) {
	var candidate models.Candidate
	param := mux.Vars(r)

	result := db.DB.First(&candidate, param["id"])
	result.Preload("Election")



	if result.Error != nil || candidate.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("candidate not found"))
		return
	}
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(&candidate)
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	var candidate models.Candidate
	_ = json.NewDecoder(r.Body).Decode(&candidate)

	createElection := db.DB.Create(&candidate)
	err := createElection.Error
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	json.NewEncoder(w).Encode(&candidate)
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	var candidate models.Candidate
	params := mux.Vars(r)

	result := db.DB.First(&candidate, params["id"])

	if result.Error != nil || candidate.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("candidate not found"))
		return
	}
	db.DB.Preload("Election").Delete(&candidate, params["id"])
	w.WriteHeader(http.StatusOK)

}

func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	var candidate []models.Candidate
	db.DB.Find(&candidate)
	json.NewEncoder(w).Encode(&candidate)
}
