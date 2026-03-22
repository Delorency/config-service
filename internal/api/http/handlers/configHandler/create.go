package configHandler

import (
	"encoding/json"
	. "main/internal/api/http/handlers"
	"net/http"

	"github.com/go-chi/chi"
)

func (h *ConfigHandler) CreateConfig(w http.ResponseWriter, r *http.Request) {
	serviceID := chi.URLParam(r, "service_id")
	if serviceID == "" {
		WriteError(w, http.StatusBadRequest, "service_id is required")
		return
	}

	var req struct {
		Config    json.RawMessage `json:"config"`
		CreatedBy string          `json:"created_by"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}

	if len(req.Config) == 0 {
		WriteError(w, http.StatusBadRequest, "config is required")
		return
	}

	if req.CreatedBy == "" {
		req.CreatedBy = "system"
	}

	config, err := h.service.CreateConfig(r.Context(), serviceID, req.Config)
	if err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	WriteJSON(w, http.StatusCreated, map[string]interface{}{
		"service_id": config.ServiceID,
		"version":    config.Version,
		"created_at": config.CreatedAt,
		"message":    "config created successfully",
	})
}
