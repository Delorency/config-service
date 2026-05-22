package schemaHandler

import (
	"net/http"

	. "main/internal/api/http/handlers"

	"github.com/go-chi/chi"
)

// DeleteSchema godoc
// @Summary      Удаление схемы сервиса
// @Description  Удаляет JSON-схему для указанного сервиса по его ID. В случае успеха возвращает информацию об удаленном сервисе.
// @Tags         schemas
// @Accept       json
// @Produce      json
// @Param        service_id   path      string  true  "Уникальный идентификатор сервиса"  example(user-service)
// @Success      200  {object}  map[string]string  "service_id и сообщение об успешном удалении"
// @Failure      400  {object}  map[string]string  "service_id is required"
// @Failure      404  {object}  map[string]string  "схема с указанным service_id не найдена"
// @Router       /schema/{service_id} [delete]
func (h *SchemaHandler) DeleteSchema(w http.ResponseWriter, r *http.Request) {
	serviceID := chi.URLParam(r, "service_id")
	if serviceID == "" {
		WriteError(w, http.StatusBadRequest, "service_id is required")
		return
	}

	if err := h.service.DeleteSchema(serviceID); err != nil {
		WriteError(w, http.StatusNotFound, err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, map[string]string{
		"service_id": serviceID,
		"message":    "schema deleted",
	})
}
