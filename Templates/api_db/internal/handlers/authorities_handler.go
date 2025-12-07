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

type AuthoritiesHandler struct {
    service services.AuthoritiesService
}

func NewAuthoritiesHandler(service services.AuthoritiesService) *AuthoritiesHandler {
    return &AuthoritiesHandler{service: service}
}

func (h *AuthoritiesHandler) CreateAuthority(w http.ResponseWriter, r *http.Request) {
    var authority models.ElectionAuthority
    if err := json.NewDecoder(r.Body).Decode(&authority); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    if err := h.service.CreateAuthority(&authority); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(authority)
}

func (h *AuthoritiesHandler) GetAuthority(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.ParseUint(vars["id"], 10, 32)
    if err != nil {
        http.Error(w, "Invalid authority ID", http.StatusBadRequest)
        return
    }

    authority, err := h.service.GetAuthority(uint(id))
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(authority)
}

func (h *AuthoritiesHandler) GetAuthoritiesByElection(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    electionID, err := uuid.Parse(vars["electionId"])
    if err != nil {
        http.Error(w, "Invalid election ID", http.StatusBadRequest)
        return
    }

    authorities, err := h.service.GetAuthoritiesByElection(electionID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(authorities)
}

func (h *AuthoritiesHandler) GetAuthorityByEmail(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    email := vars["email"]

    authority, err := h.service.GetAuthorityByEmail(email)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(authority)
}

func (h *AuthoritiesHandler) UpdateAuthority(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.ParseUint(vars["id"], 10, 32)
    if err != nil {
        http.Error(w, "Invalid authority ID", http.StatusBadRequest)
        return
    }

    // Primero obtenemos la autoridad existente para preservar algunos campos si es necesario
    existingAuthority, err := h.service.GetAuthority(uint(id))
    if err != nil {
        http.Error(w, "Authority not found", http.StatusNotFound)
        return
    }

    var authority models.ElectionAuthority
    if err := json.NewDecoder(r.Body).Decode(&authority); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Aseguramos que el ID sea el correcto
    authority.ID = uint(id)
    
    // Preservamos campos sensibles que no deberían actualizarse desde la API
    if authority.Password == "" {
        authority.Password = existingAuthority.Password
    }
    if len(authority.SKey) == 0 {
        authority.SKey = existingAuthority.SKey
    }
    if authority.Election.String() == "00000000-0000-0000-0000-000000000000" {
        authority.Election = existingAuthority.Election
    }

    if err := h.service.UpdateAuthority(&authority); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(authority)
}

func (h *AuthoritiesHandler) UpdateAuthorityPartial(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.ParseUint(vars["id"], 10, 32)
    if err != nil {
        http.Error(w, "Invalid authority ID", http.StatusBadRequest)
        return
    }

    // Obtenemos la autoridad existente
    existingAuthority, err := h.service.GetAuthority(uint(id))
    if err != nil {
        http.Error(w, "Authority not found", http.StatusNotFound)
        return
    }

    var updates map[string]interface{}
    if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Actualizamos solo los campos permitidos
    if name, ok := updates["name"].(string); ok {
        existingAuthority.Name = name
    }
    if email, ok := updates["email"].(string); ok {
        existingAuthority.Email = email
    }
    if cc, ok := updates["cc"].(float64); ok {
        existingAuthority.CC = int(cc)
    }

    if err := h.service.UpdateAuthority(existingAuthority); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(existingAuthority)
}

func (h *AuthoritiesHandler) DeleteAuthority(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.ParseUint(vars["id"], 10, 32)
    if err != nil {
        http.Error(w, "Invalid authority ID", http.StatusBadRequest)
        return
    }

    // Verificamos que la autoridad exista antes de eliminar
    _, err = h.service.GetAuthority(uint(id))
    if err != nil {
        http.Error(w, "Authority not found", http.StatusNotFound)
        return
    }

    if err := h.service.DeleteAuthority(uint(id)); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "message": "Authority deleted successfully",
        "id":      id,
    })
}

func (h *AuthoritiesHandler) UpdateAuthorityPassword(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.ParseUint(vars["id"], 10, 32)
    if err != nil {
        http.Error(w, "Invalid authority ID", http.StatusBadRequest)
        return
    }

    var request struct {
        Password string `json:"password"`
    }
    if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    if request.Password == "" {
        http.Error(w, "Password cannot be empty", http.StatusBadRequest)
        return
    }

    authority, err := h.service.GetAuthority(uint(id))
    if err != nil {
        http.Error(w, "Authority not found", http.StatusNotFound)
        return
    }

    authority.Password = request.Password

    if err := h.service.UpdateAuthority(authority); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "message": "Password updated successfully",
        "id":      id,
    })
}

func (h *AuthoritiesHandler) UpdateAuthoritySecretKey(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.ParseUint(vars["id"], 10, 32)
    if err != nil {
        http.Error(w, "Invalid authority ID", http.StatusBadRequest)
        return
    }

    var request struct {
        SKey []byte `json:"s_key"`
    }
    if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    if len(request.SKey) == 0 {
        http.Error(w, "Secret key cannot be empty", http.StatusBadRequest)
        return
    }

    authority, err := h.service.GetAuthority(uint(id))
    if err != nil {
        http.Error(w, "Authority not found", http.StatusNotFound)
        return
    }

    authority.SKey = request.SKey

    if err := h.service.UpdateAuthority(authority); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "message": "Secret key updated successfully",
        "id":      id,
    })
}