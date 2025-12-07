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

type VotersHandler struct {
    service services.VotersService
}

func NewVotersHandler(service services.VotersService) *VotersHandler {
    return &VotersHandler{service: service}
}

func (h *VotersHandler) CreateVoter(w http.ResponseWriter, r *http.Request) {
    var voter models.Voter
    if err := json.NewDecoder(r.Body).Decode(&voter); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    if err := h.service.CreateVoter(&voter); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(voter)
}

func (h *VotersHandler) GetVoter(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.ParseUint(vars["id"], 10, 32)
    if err != nil {
        http.Error(w, "Invalid voter ID", http.StatusBadRequest)
        return
    }

    voter, err := h.service.GetVoter(uint(id))
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(voter)
}

func (h *VotersHandler) GetVotersByElection(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    electionID, err := uuid.Parse(vars["electionId"])
    if err != nil {
        http.Error(w, "Invalid election ID", http.StatusBadRequest)
        return
    }

    voters, err := h.service.GetVotersByElection(electionID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(voters)
}

func (h *VotersHandler) GetVoterByToken(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    token := vars["token"]

    voter, err := h.service.GetVoterByToken(token)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(voter)
}

func (h *VotersHandler) UpdateVoter(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.ParseUint(vars["id"], 10, 32)
    if err != nil {
        http.Error(w, "Invalid voter ID", http.StatusBadRequest)
        return
    }

    var voter models.Voter
    if err := json.NewDecoder(r.Body).Decode(&voter); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    voter.ID = uint(id)

    if err := h.service.UpdateVoter(&voter); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(voter)
}

func (h *VotersHandler) DeleteVoter(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.ParseUint(vars["id"], 10, 32)
    if err != nil {
        http.Error(w, "Invalid voter ID", http.StatusBadRequest)
        return
    }

    if err := h.service.DeleteVoter(uint(id)); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}

func (h *VotersHandler) UpdateVoteStatus(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.ParseUint(vars["id"], 10, 32)
    if err != nil {
        http.Error(w, "Invalid voter ID", http.StatusBadRequest)
        return
    }

    var request struct {
        Status bool `json:"status"`
    }
    if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    if err := h.service.UpdateVoteStatus(uint(id), request.Status); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "message": "Vote status updated successfully",
        "status":  request.Status,
    })
}