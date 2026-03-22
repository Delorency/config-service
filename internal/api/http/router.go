package http

import (
	"log"

	// _ "auth/docs"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	sh "main/internal/api/http/handlers/schemaHandler"
	ss "main/internal/service/schemaService"
)

func AddMiddleware(router *chi.Mux) *chi.Mux {
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	return router
}

func NewRouter(logger *log.Logger, schemadir string) *chi.Mux {
	router := AddMiddleware(chi.NewRouter())
	handler := sh.NewSchemaHandler(ss.NewSchemaService(schemadir))

	// router.Get("/swagger/*", httpSwagger.WrapHandler)

	router.Post("/schema", handler.UploadSchema)
	router.Get("/schema/{id}", handler.GetSchema)
	router.Delete("/schema/{id}", handler.DeleteSchema)

	return router
}
