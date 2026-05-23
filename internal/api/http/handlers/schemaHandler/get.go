package schemaHandler

import (
	. "main/internal/api/http/handlers"
	"net/http"

	"github.com/go-chi/chi"
)

// GetSchema godoc
// @Summary      Получение схемы сервиса
// @Description  Возвращает JSON-схему для указанного сервиса по его ID. Схема используется для валидации конфигураций сервиса.
// @Tags         schemas
// @Produce      json
// @Param        service_id   path      string  true  "Уникальный идентификатор сервиса"  example(user-service)
// @Success      200  {object}  map[string]interface{}  "JSON-схема сервиса"
// @Failure      400  {object}  map[string]string  "service_id is required"
// @Failure      404  {object}  map[string]string  "схема с указанным service_id не найдена"
// @Router       /schema/{service_id} [get]
func (h *SchemaHandler) GetSchema(w http.ResponseWriter, r *http.Request) {
	serviceID := chi.URLParam(r, "service_id")
	if serviceID == "" {
		WriteError(w, http.StatusBadRequest, "service_id is required")
		return
	}

	data, err := h.service.GetSchema(serviceID)
	if err != nil {
		WriteError(w, http.StatusNotFound, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
