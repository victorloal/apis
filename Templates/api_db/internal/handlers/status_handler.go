package handlers

import (
    "api_db/internal/models"
    "api_db/internal/services"
    "encoding/json"
    "net/http"
    "strconv"
    "github.com/gorilla/mux"
)

type StatusHandler struct {
    service services.StatusService
}

func NewStatusHandler(service services.StatusService) *StatusHandler {
    return &StatusHandler{service: service}
}

func (h *StatusHandler) CreateStatus(w http.ResponseWriter, r *http.Request) {
    var status models.Status
    if err := json.NewDecoder(r.Body).Decode(&status); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    if err := h.service.CreateStatus(&status); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(status)
}

func (h *StatusHandler) GetStatus(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid status ID", http.StatusBadRequest)
        return
    }

    status, err := h.service.GetStatus(id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(status)
}

func (h *StatusHandler) GetStatusByName(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    name := vars["name"]

    status, err := h.service.GetStatusByName(name)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(status)
}

func (h *StatusHandler) GetAllStatus(w http.ResponseWriter, r *http.Request) {
    statuses, err := h.service.GetAllStatus()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(statuses)
}

func (h *StatusHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid status ID", http.StatusBadRequest)
        return
    }

    var status models.Status
    if err := json.NewDecoder(r.Body).Decode(&status); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    status.ID = id

    if err := h.service.UpdateStatus(&status); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(status)
}

func (h *StatusHandler) DeleteStatus(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid status ID", http.StatusBadRequest)
        return
    }

    if err := h.service.DeleteStatus(id); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}