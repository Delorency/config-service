package http

import (
	"log"

	_ "main/docs"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	httpSwagger "github.com/swaggo/http-swagger"

	ch "main/internal/api/http/handlers/configHandler"
	sh "main/internal/api/http/handlers/schemaHandler"
	"main/internal/core/validator"
	"main/internal/repo"
	cs "main/internal/service/configService"
	ss "main/internal/service/schemaService"
	"main/internal/watcher"
)

func AddMiddleware(router *chi.Mux) *chi.Mux {
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	return router
}

func NewRouter(logger *log.Logger, schemadir string,
	repo *repo.ConfigRepository,
	watcher *watcher.WatcherManager,
	validator *validator.Validator) *chi.Mux {

	router := AddMiddleware(chi.NewRouter())
	schemaH := sh.NewSchemaHandler(ss.NewSchemaService(schemadir))
	configH := ch.NewConfigHandler(cs.NewConfigService(repo, watcher, validator))

	router.Get("/swagger/*", httpSwagger.WrapHandler)

	router.Post("/schema/{service_id}/upload", schemaH.UploadSchema)
	router.Get("/schema/{service_id}", schemaH.GetSchema)
	router.Delete("/schema/{service_id}", schemaH.DeleteSchema)

	router.Post("/config/{service_id}", configH.CreateConfig)
	router.Get("/config/{service_id}", configH.GetConfigList)
	router.Get("/config/{service_id}/actual", configH.GetActualConfig)
	router.Get("/config/{service_id}/{version}", configH.GetVersion)
	router.Post("/config/{service_id}/{version}/rollback", configH.Rollback)

	return router
}
