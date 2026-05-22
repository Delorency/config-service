package schemaHandler

import (
	"io"
	. "main/internal/api/http/handlers"
	"net/http"

	"github.com/go-chi/chi"
)

// UploadSchema godoc
// @Summary      Загрузка JSON-схемы из файла
// @Description  Загружает JSON-схему для указанного сервиса через multipart/form-data. Поддерживаются файлы с расширениями .json, .yaml, .yml
// @Tags         schemas
// @Accept       multipart/form-data
// @Produce      json
// @Param        service_id   path      string  true   "Уникальный идентификатор сервиса"  example(user-service)
// @Param        schema       formData  file    true   "JSON-файл схемы"  example(schema.json)
// @Success      201  {object}  UploadSchemaResponse  "Схема успешно загружена"
// @Failure      400  {object}  map[string]string     "Невалидный запрос: отсутствует service_id, файл слишком большой или неверный формат"
// @Failure      500  {object}  map[string]string     "Внутренняя ошибка сервера"
// @Router       /schema/{service_id}/upload [post]
func (h *SchemaHandler) UploadSchema(w http.ResponseWriter, r *http.Request) {
	serviceID := chi.URLParam(r, "service_id")
	if serviceID == "" {
		WriteError(w, http.StatusBadRequest, "service_id is required")
		return
	}

	if err := r.ParseMultipartForm(1 << 20); err != nil {
		WriteError(w, http.StatusBadRequest, "file too large or invalid form")
		return
	}

	file, header, err := r.FormFile("schema")
	if err != nil {
		WriteError(w, http.StatusBadRequest, "schema file is required")
		return
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "failed to read file")
		return
	}

	if err := h.service.CreateSchemaFromFile(serviceID, content, header.Filename); err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	WriteJSON(w, http.StatusCreated, map[string]string{
		"service_id": serviceID,
		"filename":   header.Filename,
		"message":    "schema uploaded successfully",
	})
}

// UploadSchemaResponse представляет успешный ответ при загрузке схемы
type UploadSchemaResponse struct {
	ServiceID string `json:"service_id" example:"user-service" description:"ID сервиса"`
	Filename  string `json:"filename" example:"user-schema.json" description:"Имя загруженного файла"`
	Message   string `json:"message" example:"schema uploaded successfully" description:"Сообщение об успешной операции"`
}
