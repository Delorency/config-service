package schemaHandler

import (
	"net/http"

	"github.com/go-chi/chi"
)

func (h *SchemaHandler) GetSchema(w http.ResponseWriter, r *http.Request) {
	serviceID := chi.URLParam(r, "service_id")
	if serviceID == "" {
		writeError(w, http.StatusBadRequest, "service_id is required")
		return
	}

	data, err := h.service.GetSchema(serviceID)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
