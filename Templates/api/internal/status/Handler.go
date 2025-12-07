package status

import (
	"api/internal/models"
	"api/pkg/db"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func GetAllHandler(w http.ResponseWriter, r *http.Request) {
	var status []models.Status
	db.DB.Find(&status)
	json.NewEncoder(w).Encode(&status)
}

func GetIdHandler(w http.ResponseWriter, r *http.Request) {
	var state models.Status
	param := mux.Vars(r)
	
	result := db.DB.First(&state, param["id"])


	if result.Error != nil || state.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Status not found"))
		return
	}
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(&state)
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	var state models.Status
	_ = json.NewDecoder(r.Body).Decode(&state)

	createElection := db.DB.Create(&state)
	err := createElection.Error
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	json.NewEncoder(w).Encode(&state)
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	var state models.Status
	params := mux.Vars(r)

	result := db.DB.First(&state, params["id"])

	if result.Error != nil || state.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Status not found"))
		return
	}
	db.DB.Preload("Tasks").Delete(&state, params["id"])
	w.WriteHeader(http.StatusOK)

}

func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	var state models.Status
	db.DB.Find(&state)
	json.NewEncoder(w).Encode(&state)
}
