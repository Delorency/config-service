package schemaHandler

import (
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
