package configHandler

import (
	"encoding/json"
	"io"
	"net/http"

	. "main/internal/api/http/handlers"

	"github.com/go-chi/chi"
)

func (h *ConfigHandler) UpdateConfig(w http.ResponseWriter, r *http.Request) {
	serviceID := chi.URLParam(r, "service_id")
	if serviceID == "" {
		WriteError(w, http.StatusBadRequest, "service_id is required")
		return
	}

	if err := r.ParseMultipartForm(5 << 20); err != nil {
		WriteError(w, http.StatusBadRequest, "file too large or invalid form")
		return
	}

	file, header, err := r.FormFile("schema")
	if err != nil {
		WriteError(w, http.StatusBadRequest, "config file is required")
		return
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "failed to read file")
		return
	}

	var testJSON map[string]any
	if err := json.Unmarshal(content, &testJSON); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid JSON config: "+err.Error())
		return
	}

	config, err := h.service.UpdateConfig(r.Context(), serviceID, content)
	if err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"service_id": config.ServiceID,
		"version":    config.Version,
		"filename":   header.Filename,
		"updated_at": config.UpdatedAt,
		"message":    "config updated successfully",
	})
}
