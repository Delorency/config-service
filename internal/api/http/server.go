package http

import (
	"fmt"
	"log"
	"main/internal/core/validator"
	"main/internal/repo"
	"main/internal/watcher"
	"net/http"
)

func NewHTTPServer(addr string, port int, logger *log.Logger, schemadir string,
	repo *repo.ConfigRepository,
	watcher *watcher.WatcherManager,
	validator *validator.Validator) *http.Server {
	router := NewRouter(logger, schemadir, repo, watcher, validator)

	server := http.Server{
		Addr:    fmt.Sprintf("%s:%d", addr, port),
		Handler: router,
	}

	return &server
}
