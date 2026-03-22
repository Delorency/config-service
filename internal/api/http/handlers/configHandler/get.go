package configHandler

import (
	. "main/internal/api/http/handlers"
	"net/http"

	"github.com/go-chi/chi"
)

func (h *ConfigHandler) GetConfig(w http.ResponseWriter, r *http.Request) {
	serviceID := chi.URLParam(r, "service_id")
	if serviceID == "" {
		WriteError(w, http.StatusBadRequest, "service_id is required")
		return
	}

	config, err := h.service.GetConfig(r.Context(), serviceID)
	if err != nil {
		WriteError(w, http.StatusNotFound, err.Error())
		return
	}
	WriteJSON(w, 200, config.Data)
}
