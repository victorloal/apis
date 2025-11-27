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

type TallyResultsHandler struct {
    service services.TallyResultsService
}

func NewTallyResultsHandler(service services.TallyResultsService) *TallyResultsHandler {
    return &TallyResultsHandler{service: service}
}

func (h *TallyResultsHandler) CreateTallyResult(w http.ResponseWriter, r *http.Request) {
    var result models.TallyResult
    if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    if err := h.service.CreateTallyResult(&result); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(result)
}

func (h *TallyResultsHandler) GetTallyResult(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.ParseUint(vars["id"], 10, 32)
    if err != nil {
        http.Error(w, "Invalid tally result ID", http.StatusBadRequest)
        return
    }

    result, err := h.service.GetTallyResult(uint(id))
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(result)
}

func (h *TallyResultsHandler) GetTallyResultByElection(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    electionID, err := uuid.Parse(vars["electionId"])
    if err != nil {
        http.Error(w, "Invalid election ID", http.StatusBadRequest)
        return
    }

    result, err := h.service.GetTallyResultByElection(electionID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(result)
}

func (h *TallyResultsHandler) UpdateTallyResult(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.ParseUint(vars["id"], 10, 32)
    if err != nil {
        http.Error(w, "Invalid tally result ID", http.StatusBadRequest)
        return
    }

    var result models.TallyResult
    if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    result.ID = uint(id)

    if err := h.service.UpdateTallyResult(&result); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(result)
}

func (h *TallyResultsHandler) DeleteTallyResult(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.ParseUint(vars["id"], 10, 32)
    if err != nil {
        http.Error(w, "Invalid tally result ID", http.StatusBadRequest)
        return
    }

    if err := h.service.DeleteTallyResult(uint(id)); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}

func (h *TallyResultsHandler) ComputeTallyResult(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    electionID, err := uuid.Parse(vars["electionId"])
    if err != nil {
        http.Error(w, "Invalid election ID", http.StatusBadRequest)
        return
    }

    var request struct {
        ComputedBy string `json:"computed_by"`
    }
    if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    if err := h.service.ComputeTallyResult(electionID, request.ComputedBy); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "message": "Tally result computed successfully",
        "election_id": electionID,
    })
}

func (h *TallyResultsHandler) GetTallyResultsWithElection(w http.ResponseWriter, r *http.Request) {
    results, err := h.service.GetTallyResultsWithElection()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(results)
}