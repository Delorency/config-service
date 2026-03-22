package schemaHandler

import (
	"net/http"

	"github.com/go-chi/chi"
)

func (h *SchemaHandler) DeleteSchema(w http.ResponseWriter, r *http.Request) {
	serviceID := chi.URLParam(r, "service_id")
	if serviceID == "" {
		writeError(w, http.StatusBadRequest, "service_id is required")
		return
	}

	if err := h.service.DeleteSchema(serviceID); err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"service_id": serviceID,
		"message":    "schema deleted",
	})
}
