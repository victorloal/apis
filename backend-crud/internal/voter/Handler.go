package voter

import (
	"api/internal/models"
	"api/pkg/db"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func GetAllHandler(w http.ResponseWriter, r *http.Request) {
	var voters []models.Voter
	db.DB.Find(&voters)
	json.NewEncoder(w).Encode(&voters)
}

func GetIdHandler(w http.ResponseWriter, r *http.Request) {
	var voter models.Voter
	param := mux.Vars(r)

	result := db.DB.First(&voter, param["id"])
	result.Preload("Election","Ballot")

	if result.Error != nil || voter.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Voter not found"))
		return
	}
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(&voter)
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	var voter models.Voter
	_ = json.NewDecoder(r.Body).Decode(&voter)

	createElection := db.DB.Create(&voter)
	err := createElection.Error
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	json.NewEncoder(w).Encode(&voter)
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	var voter models.Voter
	params := mux.Vars(r)

	result := db.DB.First(&voter, params["id"])

	if result.Error != nil || voter.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Voter not found"))
		return
	}
	db.DB.Preload("Election","Ballot").Delete(&voter, params["id"])
	w.WriteHeader(http.StatusOK)

}

func UpdateHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    // Ver ID
    params := mux.Vars(r)
    var voter models.Voter

    if err := db.DB.First(&voter, params["id"]).Error; err != nil {
        w.WriteHeader(http.StatusNotFound)
        json.NewEncoder(w).Encode(map[string]string{"error": "Voter not found"})
        return
    }

    // Decodificar JSON a un mapa dinámico
    var updates map[string]interface{}
    if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{"error": "Invalid JSON"})
        return
    }

    // Actualizar valores dinámicos
    if err := db.DB.Model(&voter).Updates(updates).Error; err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
        return
    }

    // Respuesta final
    json.NewEncoder(w).Encode(voter)
}
