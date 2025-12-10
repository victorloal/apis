package homomorphicKey

import (
	"api/internal/models"
	"api/pkg/db"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func GetAllHandler(w http.ResponseWriter, r *http.Request) {
	var lits []models.HomomorphicKey
	db.DB.Preload("Election").Find(&lits)
	json.NewEncoder(w).Encode(&lits)
}

func GetIdHandler(w http.ResponseWriter, r *http.Request) {
	var aux models.ElectionAuthority
	param := mux.Vars(r)

	result := db.DB.First(&aux, param["id"])
	result.Preload("Election")


	if result.Error != nil || aux.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("candidate not found"))
		return
	}
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(&aux)
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	var aux models.HomomorphicKey
	_ = json.NewDecoder(r.Body).Decode(&aux)

	createElection := db.DB.Create(&aux)
	err := createElection.Error
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	json.NewEncoder(w).Encode(&aux)
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	var aux models.HomomorphicKey
	params := mux.Vars(r)

	result := db.DB.First(&aux, params["id"])

	if result.Error != nil || aux.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("candidate not found"))
		return
	}
	db.DB.Preload("Election").Delete(&aux, params["id"])
	w.WriteHeader(http.StatusOK)

}

func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	var aux models.HomomorphicKey
	db.DB.Find(&aux)
	json.NewEncoder(w).Encode(&aux)
}
