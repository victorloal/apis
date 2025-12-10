package tallyResult

import (
	"api/internal/models"
	"api/pkg/db"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func GetAllHandler(w http.ResponseWriter, r *http.Request) {
	var tallyResult []models.TallyResult
	db.DB.Find(&tallyResult)
	json.NewEncoder(w).Encode(&tallyResult)
}

func GetIdHandler(w http.ResponseWriter, r *http.Request) {
	var tallyResult models.TallyResult
	param := mux.Vars(r)

	result := db.DB.First(&tallyResult, param["id"])

	if result.Error != nil || tallyResult.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Result not found"))
		return
	}
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(&tallyResult)
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	var tallyResult models.TallyResult
	_ = json.NewDecoder(r.Body).Decode(&tallyResult)

	createElection := db.DB.Create(&tallyResult)
	err := createElection.Error
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	json.NewEncoder(w).Encode(&tallyResult)
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	var tallyResult models.TallyResult
	params := mux.Vars(r)

	result := db.DB.First(&tallyResult, params["id"])

	if result.Error != nil || tallyResult.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Result not found"))
		return
	}
	db.DB.Delete(&tallyResult, params["id"])
	w.WriteHeader(http.StatusOK)

}

func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	var tallyResult models.TallyResult
	db.DB.Find(&tallyResult)
	json.NewEncoder(w).Encode(&tallyResult)
}
