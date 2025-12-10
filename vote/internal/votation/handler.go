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
	// Se configura si es encritada, si tiene autoridades electorales y si tiene la opcion de crearlos,

	// enviar a api PSOT localhost:3001/election
	// Enviar al blocjcahin por la otra api

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

func CreateCandidateHandler(w http.ResponseWriter, r *http.Request) {
	// verificar que no esten vacios con service
	var election models.Election
	_ = json.NewDecoder(r.Body).Decode(&election)

	// enviar a

}

func CreateCandidateListHandler(w http.ResponseWriter, r *http.Request) {
	// verificar que no esten vacios con service
	var election models.Election
	_ = json.NewDecoder(r.Body).Decode(&election)

	// enviar a

}

func UpdateCandidateHandler(w http.ResponseWriter, r *http.Request) {
	// verificar que no esten vacios con service
	var election models.Election
	_ = json.NewDecoder(r.Body).Decode(&election)

	// enviar a

}

func DeleteCandidateHandler(w http.ResponseWriter, r *http.Request) {
	// verificar que no esten vacios con service
	var election models.Election
	_ = json.NewDecoder(r.Body).Decode(&election)

	// enviar a

}

func CreateVoterHandler(w http.ResponseWriter, r *http.Request) {
	// verificar que no esten vacios con service
	var election models.Election
	_ = json.NewDecoder(r.Body).Decode(&election)

	// enviar a

}

func CreateVoterListHandler(w http.ResponseWriter, r *http.Request) {
	// verificar que no esten vacios con service
	var election models.Election
	_ = json.NewDecoder(r.Body).Decode(&election)

	// enviar a

}

func UpdateVoterHandler(w http.ResponseWriter, r *http.Request) {
	// verificar que no esten vacios con service
	var election models.Election
	_ = json.NewDecoder(r.Body).Decode(&election)

	// enviar a

}

func VerifyVoterHandler(w http.ResponseWriter, r *http.Request) {
	// verificar que no esten vacios con service
	var election models.Election
	_ = json.NewDecoder(r.Body).Decode(&election)

	// enviar a

}

func CreateBallotHandler(w http.ResponseWriter, r *http.Request) {
	// verificar que no esten vacios con service
	var election models.Election
	_ = json.NewDecoder(r.Body).Decode(&election)

	// enviar a

}

func VerifyBallotHandler(w http.ResponseWriter, r *http.Request) {
	// verificar que no esten vacios con service
	var election models.Election
	_ = json.NewDecoder(r.Body).Decode(&election)

	// enviar a

}
