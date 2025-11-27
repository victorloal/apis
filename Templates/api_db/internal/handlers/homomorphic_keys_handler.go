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

type HomomorphicKeysHandler struct {
    service services.HomomorphicKeysService
}

func NewHomomorphicKeysHandler(service services.HomomorphicKeysService) *HomomorphicKeysHandler {
    return &HomomorphicKeysHandler{service: service}
}

func (h *HomomorphicKeysHandler) CreateKey(w http.ResponseWriter, r *http.Request) {
    var key models.HomomorphicKey
    if err := json.NewDecoder(r.Body).Decode(&key); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    if err := h.service.CreateKey(&key); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(key)
}

func (h *HomomorphicKeysHandler) GetKey(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.ParseUint(vars["id"], 10, 32)
    if err != nil {
        http.Error(w, "Invalid key ID", http.StatusBadRequest)
        return
    }

    key, err := h.service.GetKey(uint(id))
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(key)
}

func (h *HomomorphicKeysHandler) GetKeyByElection(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    electionID, err := uuid.Parse(vars["electionId"])
    if err != nil {
        http.Error(w, "Invalid election ID", http.StatusBadRequest)
        return
    }

    key, err := h.service.GetKeyByElection(electionID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(key)
}

func (h *HomomorphicKeysHandler) UpdateKey(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.ParseUint(vars["id"], 10, 32)
    if err != nil {
        http.Error(w, "Invalid key ID", http.StatusBadRequest)
        return
    }

    var key models.HomomorphicKey
    if err := json.NewDecoder(r.Body).Decode(&key); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    key.ID = uint(id)

    if err := h.service.UpdateKey(&key); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(key)
}

func (h *HomomorphicKeysHandler) DeleteKey(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.ParseUint(vars["id"], 10, 32)
    if err != nil {
        http.Error(w, "Invalid key ID", http.StatusBadRequest)
        return
    }

    if err := h.service.DeleteKey(uint(id)); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}

func (h *HomomorphicKeysHandler) UpdateKeyParams(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    electionID, err := uuid.Parse(vars["electionId"])
    if err != nil {
        http.Error(w, "Invalid election ID", http.StatusBadRequest)
        return
    }

    var request struct {
        Params models.JSONB `json:"params"`
    }
    if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    if err := h.service.UpdateKeyParams(electionID, request.Params); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "message": "Key parameters updated successfully",
        "election_id": electionID,
    })
}