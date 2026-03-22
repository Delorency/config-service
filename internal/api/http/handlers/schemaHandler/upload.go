package schemaHandler

import (
	"io"
	"net/http"

	"github.com/go-chi/chi"
)

func (h *SchemaHandler) UploadSchema(w http.ResponseWriter, r *http.Request) {
	serviceID := chi.URLParam(r, "service_id")
	if serviceID == "" {
		writeError(w, http.StatusBadRequest, "service_id is required")
		return
	}

	if err := r.ParseMultipartForm(1 << 20); err != nil {
		writeError(w, http.StatusBadRequest, "file too large or invalid form")
		return
	}

	file, header, err := r.FormFile("schema")
	if err != nil {
		writeError(w, http.StatusBadRequest, "schema file is required")
		return
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to read file")
		return
	}

	if err := h.service.CreateSchemaFromFile(serviceID, content, header.Filename); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{
		"service_id": serviceID,
		"filename":   header.Filename,
		"message":    "schema uploaded successfully",
	})
}
