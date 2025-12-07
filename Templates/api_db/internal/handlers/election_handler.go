package handlers

import (
    "api_db/internal/models"
    "api_db/internal/services"
    "encoding/json"
    "net/http"
    "github.com/gorilla/mux"
    "github.com/google/uuid"
)

type ElectionHandler struct {
    service services.ElectionService
}

func NewElectionHandler(service services.ElectionService) *ElectionHandler {
    return &ElectionHandler{service: service}
}

func (h *ElectionHandler) CreateElection(w http.ResponseWriter, r *http.Request) {
    var election models.Election
    if err := json.NewDecoder(r.Body).Decode(&election); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    if err := h.service.CreateElection(&election); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(election)
}

func (h *ElectionHandler) GetElection(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := uuid.Parse(vars["id"])
    if err != nil {
        http.Error(w, "Invalid election ID", http.StatusBadRequest)
        return
    }

    election, err := h.service.GetElection(id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(election)
}

func (h *ElectionHandler) GetAllElections(w http.ResponseWriter, r *http.Request) {
    elections, err := h.service.GetAllElections()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(elections)
}

func (h *ElectionHandler) UpdateElection(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := uuid.Parse(vars["id"])
    if err != nil {
        http.Error(w, "Invalid election ID", http.StatusBadRequest)
        return
    }

    var election models.Election
    if err := json.NewDecoder(r.Body).Decode(&election); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    election.ID = id

    if err := h.service.UpdateElection(&election); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(election)
}

func (h *ElectionHandler) DeleteElection(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := uuid.Parse(vars["id"])
    if err != nil {
        http.Error(w, "Invalid election ID", http.StatusBadRequest)
        return
    }

    if err := h.service.DeleteElection(id); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}