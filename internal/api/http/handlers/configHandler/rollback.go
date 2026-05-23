package configHandler

import (
	"net/http"
	"strconv"

	. "main/internal/api/http/handlers"

	"github.com/go-chi/chi"
)

// Rollback godoc
// @Summary      Откат конфигурации к указанной версии
// @Description  Выполняет откат конфигурации сервиса к указанной версии. Создается новая версия с данными из целевой версии.
// @Description  Версия автоматически инкрементируется. Watcher получит уведомление об изменении.
// @Tags         configs
// @Accept       json
// @Produce      json
// @Param        service_id       path      string  true  "ID сервиса"  example(user-service)
// @Param        version   path      int  true  "Номер версии для отката"  minimum(1)  example(3)
// @Success      200  {object}  RollbackResponse  "Успешный откат конфигурации"
// @Failure      400  {object}  ErrorResponse  "Невалидный запрос: отсутствует service_id или target_version"
// @Failure      404  {object}  ErrorResponse  "Сервис не найден или указанная версия не существует"
// @Failure      409  {object}  ErrorResponse  "Невозможно выполнить откат (например, уже на этой версии)"
// @Failure      500  {object}  ErrorResponse  "Внутренняя ошибка сервера"
// @Router       /config/{service_id}/{version}/rollback [post]
func (h *ConfigHandler) Rollback(w http.ResponseWriter, r *http.Request) {
	serviceID := chi.URLParam(r, "service_id")
	if serviceID == "" {
		WriteError(w, http.StatusBadRequest, "service_id is required")
		return
	}

	targetVersionStr := chi.URLParam(r, "version")
	if targetVersionStr == "" {
		WriteError(w, http.StatusBadRequest, "version is required")
		return
	}

	targetVersion, err := strconv.Atoi(targetVersionStr)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "target_version must be integer")
		return
	}

	if targetVersion < 1 {
		WriteError(w, http.StatusBadRequest, "target_version must be greater than 0")
		return
	}

	config, err := h.service.Rollback(r.Context(), serviceID, targetVersion)
	if err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, RollbackResponse{
		ServiceID:     config.ServiceID,
		Version:       config.Version,
		TargetVersion: targetVersion,
		Message:       "successfully rolled back",
		RolledBackAt:  config.UpdatedAt.String(),
	})
}

type RollbackResponse struct {
	ServiceID     string `json:"service_id" example:"user-service" description:"ID сервиса"`
	Version       int    `json:"version" example:"8" description:"Новая версия конфигурации после отката"`
	TargetVersion int    `json:"target_version" example:"3" description:"Версия, к которой выполнен откат"`
	Message       string `json:"message" example:"successfully rolled back" description:"Сообщение об успешной операции"`
	RolledBackAt  string `json:"rolled_back_at" example:"2024-01-15T15:30:00Z" description:"Дата и время отката"`
}
