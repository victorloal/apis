package handlers

import (
    "api_db/internal/models"
    "api_db/internal/services"
    "encoding/json"
    "net/http"
    "strconv"
    "time"
    "github.com/gorilla/mux"
    "github.com/google/uuid"
)

type AuditLogsHandler struct {
    service services.AuditLogsService
}

func NewAuditLogsHandler(service services.AuditLogsService) *AuditLogsHandler {
    return &AuditLogsHandler{service: service}
}

func (h *AuditLogsHandler) CreateAuditLog(w http.ResponseWriter, r *http.Request) {
    var log models.AuditLog
    if err := json.NewDecoder(r.Body).Decode(&log); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    if err := h.service.CreateAuditLog(&log); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(log)
}

func (h *AuditLogsHandler) GetAuditLog(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.ParseUint(vars["id"], 10, 32)
    if err != nil {
        http.Error(w, "Invalid audit log ID", http.StatusBadRequest)
        return
    }

    log, err := h.service.GetAuditLog(uint(id))
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(log)
}

func (h *AuditLogsHandler) GetAuditLogsByElection(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    electionID, err := uuid.Parse(vars["electionId"])
    if err != nil {
        http.Error(w, "Invalid election ID", http.StatusBadRequest)
        return
    }

    logs, err := h.service.GetAuditLogsByElection(electionID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(logs)
}

func (h *AuditLogsHandler) GetAuditLogsByAction(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    action := vars["action"]

    logs, err := h.service.GetAuditLogsByAction(action)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(logs)
}

func (h *AuditLogsHandler) GetAuditLogsByUser(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    userType := vars["userType"]
    userID := vars["userId"]

    logs, err := h.service.GetAuditLogsByUser(userType, userID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(logs)
}

func (h *AuditLogsHandler) GetAuditLogsByDateRange(w http.ResponseWriter, r *http.Request) {
    var request struct {
        Start time.Time `json:"start"`
        End   time.Time `json:"end"`
    }
    if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    logs, err := h.service.GetAuditLogsByDateRange(request.Start, request.End)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(logs)
}

func (h *AuditLogsHandler) DeleteAuditLog(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.ParseUint(vars["id"], 10, 32)
    if err != nil {
        http.Error(w, "Invalid audit log ID", http.StatusBadRequest)
        return
    }

    if err := h.service.DeleteAuditLog(uint(id)); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}

func (h *AuditLogsHandler) LogVoteAction(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    electionID, err := uuid.Parse(vars["electionId"])
    if err != nil {
        http.Error(w, "Invalid election ID", http.StatusBadRequest)
        return
    }

    var request struct {
        VoterID   string `json:"voter_id"`
        IPAddress string `json:"ip_address"`
        UserAgent string `json:"user_agent"`
    }
    if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    if err := h.service.LogVoteAction(electionID, request.VoterID, request.IPAddress, request.UserAgent); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "message": "Vote action logged successfully",
        "election_id": electionID,
    })
}

func (h *AuditLogsHandler) LogAuthorityAction(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    electionID, err := uuid.Parse(vars["electionId"])
    if err != nil {
        http.Error(w, "Invalid election ID", http.StatusBadRequest)
        return
    }

    var request struct {
        AuthorityID string         `json:"authority_id"`
        Action      string         `json:"action"`
        Details     models.JSONB   `json:"details"`
    }
    if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    if err := h.service.LogAuthorityAction(electionID, request.AuthorityID, request.Action, request.Details); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "message": "Authority action logged successfully",
        "election_id": electionID,
    })
}