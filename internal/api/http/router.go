package http

import (
	"log"

	// _ "auth/docs"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

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

	// router.Get("/swagger/*", httpSwagger.WrapHandler)

	router.Post("/schema", schemaH.UploadSchema)
	router.Get("/schema/{id}", schemaH.GetSchema)
	router.Delete("/schema/{id}", schemaH.DeleteSchema)

	router.Post("/schema", configH.CreateConfig)
	router.Get("/schema/{id}", configH.GetConfig)

	return router
}
