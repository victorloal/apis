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

type CandidatesHandler struct {
    service services.CandidatesService
}

func NewCandidatesHandler(service services.CandidatesService) *CandidatesHandler {
    return &CandidatesHandler{service: service}
}

func (h *CandidatesHandler) CreateCandidate(w http.ResponseWriter, r *http.Request) {
    var candidate models.Candidate
    if err := json.NewDecoder(r.Body).Decode(&candidate); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    if err := h.service.CreateCandidate(&candidate); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(candidate)
}

func (h *CandidatesHandler) GetCandidate(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.ParseInt(vars["id"], 10, 64)
    if err != nil {
        http.Error(w, "Invalid candidate ID", http.StatusBadRequest)
        return
    }

    candidate, err := h.service.GetCandidate(id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(candidate)
}

func (h *CandidatesHandler) GetCandidatesByElection(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    electionID, err := uuid.Parse(vars["electionId"])
    if err != nil {
        http.Error(w, "Invalid election ID", http.StatusBadRequest)
        return
    }

    candidates, err := h.service.GetCandidatesByElection(electionID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(candidates)
}

func (h *CandidatesHandler) UpdateCandidate(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.ParseInt(vars["id"], 10, 64)
    if err != nil {
        http.Error(w, "Invalid candidate ID", http.StatusBadRequest)
        return
    }

    var candidate models.Candidate
    if err := json.NewDecoder(r.Body).Decode(&candidate); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    candidate.ID = id

    if err := h.service.UpdateCandidate(&candidate); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(candidate)
}

func (h *CandidatesHandler) DeleteCandidate(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.ParseInt(vars["id"], 10, 64)
    if err != nil {
        http.Error(w, "Invalid candidate ID", http.StatusBadRequest)
        return
    }

    if err := h.service.DeleteCandidate(id); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}

func (h *CandidatesHandler) GetCandidatesByOrder(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    electionID, err := uuid.Parse(vars["electionId"])
    if err != nil {
        http.Error(w, "Invalid election ID", http.StatusBadRequest)
        return
    }

    candidates, err := h.service.GetCandidatesByOrder(electionID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(candidates)
}