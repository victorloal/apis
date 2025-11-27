package handlers

import (
    "api_db/internal/models"
    "api_db/internal/services"
    "encoding/json"
    "net/http"
    "strconv"
    "github.com/gorilla/mux"
    "github.com/google/uuid"
)

type BallotsHandler struct {
    service services.BallotsService
}

func NewBallotsHandler(service services.BallotsService) *BallotsHandler {
    return &BallotsHandler{service: service}
}

func (h *BallotsHandler) CreateBallot(w http.ResponseWriter, r *http.Request) {
    var ballot models.Ballot
    if err := json.NewDecoder(r.Body).Decode(&ballot); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    if err := h.service.CreateBallot(&ballot); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(ballot)
}

func (h *BallotsHandler) GetBallot(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]
    electionID, err := uuid.Parse(vars["electionId"])
    if err != nil {
        http.Error(w, "Invalid election ID", http.StatusBadRequest)
        return
    }
    voterID, err := strconv.Atoi(vars["voterId"])
    if err != nil {
        http.Error(w, "Invalid voter ID", http.StatusBadRequest)
        return
    }

    ballot, err := h.service.GetBallot(id, electionID, voterID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(ballot)
}

func (h *BallotsHandler) GetBallotsByElection(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    electionID, err := uuid.Parse(vars["electionId"])
    if err != nil {
        http.Error(w, "Invalid election ID", http.StatusBadRequest)
        return
    }

    ballots, err := h.service.GetBallotsByElection(electionID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(ballots)
}

func (h *BallotsHandler) GetBallotsByVoter(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    voterID, err := strconv.Atoi(vars["voterId"])
    if err != nil {
        http.Error(w, "Invalid voter ID", http.StatusBadRequest)
        return
    }

    ballots, err := h.service.GetBallotsByVoter(voterID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(ballots)
}

func (h *BallotsHandler) UpdateBallot(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]
    electionID, err := uuid.Parse(vars["electionId"])
    if err != nil {
        http.Error(w, "Invalid election ID", http.StatusBadRequest)
        return
    }
    voterID, err := strconv.Atoi(vars["voterId"])
    if err != nil {
        http.Error(w, "Invalid voter ID", http.StatusBadRequest)
        return
    }

    var ballot models.Ballot
    if err := json.NewDecoder(r.Body).Decode(&ballot); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    ballot.ID = id
    ballot.Elections = electionID
    ballot.Voter = voterID

    if err := h.service.UpdateBallot(&ballot); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(ballot)
}

func (h *BallotsHandler) DeleteBallot(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]
    electionID, err := uuid.Parse(vars["electionId"])
    if err != nil {
        http.Error(w, "Invalid election ID", http.StatusBadRequest)
        return
    }
    voterID, err := strconv.Atoi(vars["voterId"])
    if err != nil {
        http.Error(w, "Invalid voter ID", http.StatusBadRequest)
        return
    }

    if err := h.service.DeleteBallot(id, electionID, voterID); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}

func (h *BallotsHandler) GetBallotsWithVoterDetails(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    electionID, err := uuid.Parse(vars["electionId"])
    if err != nil {
        http.Error(w, "Invalid election ID", http.StatusBadRequest)
        return
    }

    ballots, err := h.service.GetBallotsWithVoterDetails(electionID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(ballots)
}