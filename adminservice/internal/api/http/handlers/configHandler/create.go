package configHandler

import (
	"encoding/json"
	"io"
	. "main/internal/api/http/handlers"
	"net/http"

	"github.com/go-chi/chi"
)

// CreateConfig godoc
// @Summary      Создание новой конфигурации
// @Description  Создает новую конфигурацию для указанного сервиса. Конфигурация загружается в виде JSON файла через multipart/form-data.
// @Description  Файл должен содержать валидный JSON. После успешного создания автоматически инкрементируется версия конфигурации.
// @Tags         configs
// @Accept       multipart/form-data
// @Produce      json
// @Param        service_id   path      string  true   "ID сервиса"  example(user-service)
// @Param        config       formData  file    true   "JSON файл с конфигурацией (макс. 5MB)"  example(config.json)
// @Success      201  {object}  CreateConfigResponse  "Конфигурация успешно создана"
// @Failure      400  {object}  ErrorResponse  "Невалидный запрос: отсутствует service_id, файл слишком большой, неверный формат или невалидный JSON"
// @Failure      500  {object}  ErrorResponse  "Внутренняя ошибка сервера"
// @Router       /config/{service_id} [post]
func (h *ConfigHandler) CreateConfig(w http.ResponseWriter, r *http.Request) {
	serviceID := chi.URLParam(r, "service_id")
	if serviceID == "" {
		WriteError(w, http.StatusBadRequest, "service_id is required")
		return
	}

	if err := r.ParseMultipartForm(5 << 20); err != nil {
		WriteError(w, http.StatusBadRequest, "file too large or invalid form")
		return
	}

	file, header, err := r.FormFile("config")
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

	config, err := h.service.CreateConfig(r.Context(), serviceID, content)
	if err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	WriteJSON(w, http.StatusCreated, map[string]interface{}{
		"service_id": config.ServiceID,
		"version":    config.Version,
		"filename":   header.Filename,
		"created_at": config.CreatedAt,
		"message":    "config created successfully",
	})
}

type CreateConfigResponse struct {
	ServiceID string `json:"service_id" example:"user-service" description:"ID сервиса"`
	Version   int    `json:"version" example:"7" description:"Номер новой версии конфигурации"`
	Filename  string `json:"filename" example:"config.json" description:"Имя загруженного файла"`
	CreatedAt string `json:"created_at" example:"2024-01-15T10:30:00Z" description:"Дата создания конфигурации"`
	Message   string `json:"message" example:"config created successfully" description:"Сообщение об успешной операции"`
}
