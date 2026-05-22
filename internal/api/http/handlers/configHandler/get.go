package configHandler

import (
	. "main/internal/api/http/handlers"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

// GetConfigList godoc
// @Summary      Получение списка версий конфигурации
// @Description  Возвращает список всех версий конфигурации для указанного сервиса с пагинацией. Версии отсортированы по убыванию (новые сверху).
// @Tags         configs
// @Accept       json
// @Produce      json
// @Param        service_id   path      string  true   "ID сервиса"  example(user-service)
// @Param        limit        query     int     false  "Количество записей на страницу (макс. 100, по умолчанию 50)"  minimum(1)  maximum(100)  example(20)
// @Param        offset       query     int     false  "Смещение для пагинации"  minimum(0)  example(0)
// @Success      200  {object}  ConfigListResponse  "Список версий конфигурации"
// @Failure      400  {object}  ErrorResponse  "service_id is required"
// @Failure      500  {object}  ErrorResponse
// @Router       /config/{service_id} [get]
func (h *ConfigHandler) GetConfigList(w http.ResponseWriter, r *http.Request) {
	serviceID := chi.URLParam(r, "service_id")
	if serviceID == "" {
		WriteError(w, http.StatusBadRequest, "service_id is required")
		return
	}

	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 50
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	offset := 0
	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	versions, total, err := h.service.ListVersions(r.Context(), serviceID, limit, offset)
	if err != nil {
		WriteError(w, http.StatusNotFound, err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"service_id": serviceID,
		"items":      versions,
		"total":      total,
		"limit":      limit,
		"offset":     offset,
	})
}

// GetActualConfig godoc
// @Summary      Получение актуальной конфигурации сервиса
// @Description  Возвращает последнюю (актуальную) версию конфигурации для указанного сервиса. Включает данные конфигурации и метаинформацию.
// @Tags         configs
// @Accept       json
// @Produce      json
// @Param        service_id   path      string  true  "ID сервиса"  example(user-service)
// @Success      200  {object}  ConfigResponse  "Актуальная конфигурация сервиса"
// @Failure      400  {object}  ErrorResponse  "service_id is required"
// @Failure      404  {object}  ErrorResponse  "Конфигурация для указанного сервиса не найдена"
// @Failure      500  {object}  ErrorResponse  "Внутренняя ошибка сервера"
// @Router       /config/{service_id}/actual [get]
func (h *ConfigHandler) GetActualConfig(w http.ResponseWriter, r *http.Request) {
	serviceID := chi.URLParam(r, "service_id")
	if serviceID == "" {
		WriteError(w, http.StatusBadRequest, "service_id is required")
		return
	}

	config, err := h.service.GetConfig(r.Context(), serviceID)
	if err != nil {
		WriteError(w, http.StatusNotFound, err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"service_id": config.ServiceID,
		"version":    config.Version,
		"data":       config.Data,
		"created_at": config.CreatedAt,
		"updated_at": config.UpdatedAt,
	})
}

// GetVersion godoc
// @Summary      Получение конкретной версии конфигурации
// @Description  Возвращает конфигурацию указанной версии для сервиса. Версии нумеруются с 1 (самая старая) до N (самая новая).
// @Tags         configs
// @Accept       json
// @Produce      json
// @Param        service_id   path      string  true  "ID сервиса"  example(user-service)
// @Param        version      path      int     true  "Номер версии конфигурации"  minimum(1)  example(3)
// @Success      200  {object}  ConfigVersionResponse  "Конфигурация указанной версии"
// @Failure      400  {object}  ErrorResponse  "service_id is required | version must be integer"
// @Failure      404  {object}  ErrorResponse  "Конфигурация с указанной версией не найдена"
// @Failure      500  {object}  ErrorResponse  "Внутренняя ошибка сервера"
// @Router       /config/{service_id}/{version} [get]
func (h *ConfigHandler) GetVersion(w http.ResponseWriter, r *http.Request) {
	serviceID := chi.URLParam(r, "service_id")
	versionStr := chi.URLParam(r, "version")

	if serviceID == "" {
		WriteError(w, http.StatusBadRequest, "service_id is required")
		return
	}

	version, err := strconv.Atoi(versionStr)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "version must be integer")
		return
	}

	history, err := h.service.GetVersion(r.Context(), serviceID, version)
	if err != nil {
		WriteError(w, http.StatusNotFound, err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"service_id": history.ServiceID,
		"version":    history.Version,
		"data":       history.Data,
		"created_at": history.CreatedAt,
	})
}

// ConfigListResponse ответ со списком версий конфигурации
type ConfigListResponse struct {
	ServiceID string       `json:"service_id" example:"user-service" description:"ID сервиса"`
	Items     []ConfigItem `json:"items" description:"Список версий конфигурации"`
	Total     int          `json:"total" example:"15" description:"Общее количество версий"`
	Limit     int          `json:"limit" example:"50" description:"Лимит записей на страницу"`
	Offset    int          `json:"offset" example:"0" description:"Смещение пагинации"`
}

// ConfigItem представляет краткую информацию о версии конфигурации
type ConfigItem struct {
	Version   int    `json:"version" example:"5" description:"Номер версии"`
	CreatedAt string `json:"created_at" example:"2024-01-15T10:30:00Z" description:"Дата создания версии"`
	UpdatedAt string `json:"updated_at" example:"2024-01-15T10:30:00Z" description:"Дата последнего обновления"`
}

// ConfigResponse ответ с актуальной конфигурацией
type ConfigResponse struct {
	ServiceID string                 `json:"service_id" example:"user-service" description:"ID сервиса"`
	Version   int                    `json:"version" example:"5" description:"Номер актуальной версии"`
	Data      map[string]interface{} `json:"data" description:"Данные конфигурации"`
	CreatedAt string                 `json:"created_at" example:"2024-01-15T10:30:00Z" description:"Дата создания"`
	UpdatedAt string                 `json:"updated_at" example:"2024-01-15T15:45:00Z" description:"Дата последнего обновления"`
}

// ConfigVersionResponse ответ с конкретной версией конфигурации
type ConfigVersionResponse struct {
	ServiceID string                 `json:"service_id" example:"user-service" description:"ID сервиса"`
	Version   int                    `json:"version" example:"3" description:"Номер версии"`
	Data      map[string]interface{} `json:"data" description:"Данные конфигурации"`
	CreatedAt string                 `json:"created_at" example:"2024-01-14T09:15:00Z" description:"Дата создания версии"`
}

// ErrorResponse ответ при ошибке
type ErrorResponse struct {
	Error string `json:"error" example:"service_id is required" description:"Описание ошибки"`
}
