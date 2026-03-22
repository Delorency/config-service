package schemaHandler

import (
	"encoding/json"
	ss "main/internal/service/schemaService"
	"net/http"
)

type SchemaHandler struct {
	service ss.SchemaServiceI
}

type SchemaHandlerI interface {
	UploadSchema(http.ResponseWriter, *http.Request)
	GetSchema(http.ResponseWriter, *http.Request)
	DeleteSchema(http.ResponseWriter, *http.Request)
}

func NewSchemaHandler(service ss.SchemaServiceI) SchemaHandlerI {
	return &SchemaHandler{service}
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
