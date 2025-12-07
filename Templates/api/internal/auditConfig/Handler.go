package auditConfig

import (
	"api/internal/models"
	"api/pkg/db"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func GetAllHandler(w http.ResponseWriter, r *http.Request) {
	var auditConfigs []models.AuditConfig
	db.DB.Preload("Election").Find(&auditConfigs)
	json.NewEncoder(w).Encode(&auditConfigs)
}

func GetIdHandler(w http.ResponseWriter, r *http.Request) {
	var auditConfig models.AuditConfig
	param := mux.Vars(r)

	result := db.DB.First(&auditConfig, param["id"])
	result.Preload("Election")

	if result.Error != nil || auditConfig.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("auditConfig not found"))
		return
	}
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(&auditConfig)
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	var auditConfig models.AuditConfig
	_ = json.NewDecoder(r.Body).Decode(&auditConfig)

	createElection := db.DB.Create(&auditConfig)
	err := createElection.Error
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	json.NewEncoder(w).Encode(&auditConfig)
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	var auditConfig models.AuditConfig
	params := mux.Vars(r)

	result := db.DB.First(&auditConfig, params["id"])

	if result.Error != nil || auditConfig.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("auditConfig not found"))
		return
	}
	db.DB.Preload("Election").Delete(&auditConfig, params["id"])
	w.WriteHeader(http.StatusOK)

}

func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	var auditConfig []models.AuditConfig
	db.DB.Find(&auditConfig)
	json.NewEncoder(w).Encode(&auditConfig)
}
